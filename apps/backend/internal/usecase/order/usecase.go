package order

import (
	"backend/internal/domain/driver"
	"backend/internal/domain/notification"
	"backend/internal/domain/order"
	prod "backend/internal/domain/product"
	"backend/internal/usecase/common"
	"context"
	"fmt"
	"strings"

	"github.com/cridenour/go-postgis"
	"github.com/google/uuid"
)

// UseCase encapsulates order business logic and dependencies.
type UseCase struct {
	repo      order.Repository
	usrRepo   order.CustomerReader
	drvRepo   order.DriverReader
	txManager common.TxManager
	notfRepo  order.NotificationReader
	prodvrt   order.ProductOrVariantReader
	storeRepo order.StoreReader
}

// NewUseCase creates a new order UseCase.
func NewUseCase(
	repo order.Repository,
	usrRepo order.CustomerReader,
	drvRepo order.DriverReader,
	txm common.TxManager,
	notf order.NotificationReader,
	prodvrt order.ProductOrVariantReader,
	strRepo order.StoreReader,
) *UseCase {
	return &UseCase{
		repo:      repo,
		usrRepo:   usrRepo,
		drvRepo:   drvRepo,
		txManager: txm,
		notfRepo:  notf,
		prodvrt:   prodvrt,
		storeRepo: strRepo,
	}
}

// CreateOrder creates new confirmed orders and returns them.
func (uc *UseCase) CreateOrder(ctx context.Context, customerID uuid.UUID, req *order.CreateOrderRequest) ([]*order.Order, error) {
	orders, err := uc.createOrders(ctx, customerID, req, order.Pending)
	if err != nil {
		return nil, err
	}

	// Fire notifications asynchronously
	for _, o := range orders {
		go uc.fireOrderNotifications(ctx, o)
	}

	return orders, nil
}

// CreatePendingOrders creates new orders marked as pending (used for carts or pre-orders)
func (uc *UseCase) CreatePendingOrders(ctx context.Context, customerID uuid.UUID, req *order.CreateOrderRequest) ([]*order.Order, error) {
	return uc.createOrders(ctx, customerID, req, order.Pending)
}

// createOrders is a unified helper that handles both pending and confirmed orders.
func (uc *UseCase) createOrders(
	ctx context.Context,
	customerID uuid.UUID,
	req *order.CreateOrderRequest,
	status order.OrderStatus,
) ([]*order.Order, error) {
	store, err := uc.storeRepo.GetByID(ctx, req.StoreID)
	if err != nil {
		return nil, fmt.Errorf("store not found: %w", err)
	}

	var createdOrders []*order.Order

	for _, item := range req.Items {
		// Load product
		p, err := uc.prodvrt.GetProductByID(ctx, item.ProductID)
		if err != nil {
			return nil, fmt.Errorf("product not found: %w", err)
		}

		// Handle variants if present
		var variant *prod.Variant
		if p.HasVariants {
			if item.VariantID == nil {
				return nil, order.ErrorVariantRequired
			}
			variant = findVariant(p, item.VariantID)
			if variant == nil {
				return nil, fmt.Errorf("variant %s not found in product %s", item.VariantID, p.ID)
			}
			if variant.Stock < item.Quantity {
				return nil, order.ErrorOutOfStock
			}
		} else {
			if item.VariantID != nil {
				return nil, order.ErrorVariantNotAllowed
			}
			if p.Stock < item.Quantity {
				return nil, order.ErrorOutOfStock
			}
		}

		// Build order snapshot
		o := req.ToOrder()
		o.CustomerID = customerID
		o.MerchantID = store.OwnerID
		o.ProductID = p.ID
		o.VariantID = item.VariantID
		o.Quantity = item.Quantity
		o.Status = status

		if p.HasVariants {
			o.UnitPrice = int64(variant.Price * 100)
			o.Total = o.UnitPrice * int64(item.Quantity)
			o.ProductName = p.Name
			o.VariantName = variant.SKU
			o.ImageURL = variant.ImageURL

			// Decrement stock
			newStock := variant.Stock - item.Quantity
			if err := uc.prodvrt.UpdateVariantStock(ctx, variant.ID, newStock); err != nil {
				return nil, fmt.Errorf("update variant stock failed: %w", err)
			}
		} else {
			o.UnitPrice = int64(p.Price * 100)
			o.Total = o.UnitPrice * int64(item.Quantity)
			o.ProductName = p.Name
			o.ImageURL = firstImage(p.Images)

			// Decrement stock
			newStock := p.Stock - item.Quantity
			if err := uc.prodvrt.UpdateProductStock(ctx, p.ID, newStock); err != nil {
				return nil, fmt.Errorf("update product stock failed: %w", err)
			}
		}

		// Persist
		if err := uc.repo.Create(ctx, o); err != nil {
			return nil, fmt.Errorf("create order failed: %w", err)
		}

		createdOrders = append(createdOrders, o)
	}

	return createdOrders, nil
}

// GetOrder fetches a single order by its ID
func (uc *UseCase) GetOrder(ctx context.Context, id uuid.UUID) (*order.Order, error) {
	return uc.repo.GetByID(ctx, id)
}

// GetOrderPickupPoint fetches the pickup location for an order
func (uc *UseCase) GetOrderPickupPoint(ctx context.Context, orderID uuid.UUID) (postgis.PointS, error) {
	return uc.repo.GetPickupPoint(ctx, orderID)
}

// GetOrderDeliveryPoint fetches the delivery location for an order
func (uc *UseCase) GetOrderDeliveryPoint(ctx context.Context, orderID uuid.UUID) (postgis.PointS, error) {
	return uc.repo.GetDeliveryPoint(ctx, orderID)
}

// GetOrderByCustomer fetches all orders for a given customer
func (uc *UseCase) GetOrderByCustomer(ctx context.Context, customerID uuid.UUID) ([]*order.Order, error) {
	return uc.repo.ListByCustomer(ctx, customerID)
}

// UpdateOrder updates a single column value in an order row
func (uc *UseCase) UpdateOrder(ctx context.Context, orderID uuid.UUID, column string, value any) error {
	return uc.txManager.Do(ctx, func(txCtx context.Context) error {
		if err := uc.repo.Update(txCtx, orderID, column, value); err != nil {
			return fmt.Errorf("update order failed: %w", err)
		}
		return nil
	})
}

// ListOrders returns all orders
func (uc *UseCase) ListOrders(ctx context.Context) ([]*order.Order, error) {
	return uc.repo.List(ctx)
}

// DeleteOrder removes an order by ID
func (uc *UseCase) DeleteOrder(ctx context.Context, id uuid.UUID) error {
	return uc.txManager.Do(ctx, func(txCtx context.Context) error {
		if err := uc.repo.Delete(txCtx, id); err != nil {
			return fmt.Errorf("delete order failed: %w", err)
		}
		return nil
	})
}

// AssignOrderToDriver assigns a driver to an order
func (uc *UseCase) AssignOrderToDriver(ctx context.Context, orderID, driverID uuid.UUID, maxDistance float64) (*driver.Driver, error) {
	o, err := uc.GetOrder(ctx, orderID)
	if err != nil {
		return nil, fmt.Errorf("fetch order: %w", err)
	}
	if o.Status != order.Pending {
		return nil, fmt.Errorf("order not pending")
	}

	pickupPoint, err := uc.repo.GetPickupPoint(ctx, orderID)
	if err != nil {
		return nil, fmt.Errorf("get pickup point: %w", err)
	}

	nearestDriver, err := uc.drvRepo.GetNearestDriver(ctx, pickupPoint, maxDistance)
	if err != nil || nearestDriver == nil {
		return nil, fmt.Errorf("no available driver within %.2f meters", maxDistance)
	}

	if err := uc.UpdateOrder(ctx, orderID, "status", order.Assigned); err != nil {
		return nil, fmt.Errorf("update order status: %w", err)
	}

	return nearestDriver, nil
}

// notify sends a notification to a user
func (uc *UseCase) notify(ctx context.Context, userID uuid.UUID, message string) error {
	n := &notification.Notification{
		UserID:  userID,
		Message: message,
		Type:    notification.System,
		Status:  notification.Pending,
	}
	return uc.notfRepo.Create(ctx, n)
}

// fireOrderNotifications triggers customer notifications
func (uc *UseCase) fireOrderNotifications(ctx context.Context, o *order.Order) {
	msg := fmt.Sprintf("Your order %s has been placed successfully.", o.ID)
	_ = uc.notify(ctx, o.CustomerID, msg)
}

// firstImage returns the first image from the list or empty
func firstImage(images []string) string {
	if len(images) == 0 {
		return ""
	}
	return images[0]
}

// findVariant searches for a variant by ID
func findVariant(product *prod.Product, variantID *uuid.UUID) *prod.Variant {
	if variantID == nil {
		return nil
	}
	for _, v := range product.Variants {
		if v.ID == *variantID {
			return &v
		}
	}
	return nil
}

// resolvePrice returns the correct unit price (product or variant)
func resolvePrice(product *prod.Product, variant *prod.Variant) int64 {
	if variant != nil {
		return int64(variant.Price * 100)
	}
	return int64(product.Price * 100)
}

// variantName returns a human-readable name for a variant
func variantName(v *prod.Variant) string {
	if v == nil {
		return ""
	}
	if len(v.Options) > 0 {
		nameParts := []string{}
		for k, val := range v.Options {
			nameParts = append(nameParts, fmt.Sprintf("%s:%s", k, val))
		}
		return fmt.Sprintf("%s (%s)", v.SKU, strings.Join(nameParts, ", "))
	}
	return v.SKU
}

// primaryImage returns the image URL for a variant or product
func primaryImage(images []string, variant *prod.Variant) string {
	if variant != nil && variant.ImageURL != "" {
		return variant.ImageURL
	}
	if len(images) > 0 {
		return images[0]
	}
	return ""
}
