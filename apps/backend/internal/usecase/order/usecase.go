package order

import (
	"backend/internal/domain/driver"
	"backend/internal/domain/notification"
	"backend/internal/domain/order"
	"backend/internal/usecase/common"
	"context"
	"fmt"

	"github.com/cridenour/go-postgis"
	"github.com/google/uuid"
)

type UseCase struct {
	repo      order.Repository
	invRepo   order.InventoryReader // X
	usrRepo   order.CustomerReader
	drvRepo   order.DriverReader
	txManager common.TxManager
	notfRepo  order.NotificationReader
	prodvrt   order.ProductOrVariantReader
}

func NewUseCase(repo order.Repository, invRepo order.InventoryReader, usrRepo order.CustomerReader, drvRepo order.DriverReader, txm common.TxManager, notf order.NotificationReader, prodvrt order.ProductOrVariantReader) *UseCase {
	return &UseCase{repo: repo, invRepo: invRepo, usrRepo: usrRepo, drvRepo: drvRepo, txManager: txm, notfRepo: notf, prodvrt: prodvrt}
}

func (uc *UseCase) CreateOrder(ctx context.Context, o *order.Order) (err error) {
	return uc.txManager.Do(ctx, func(txCtx context.Context) error {

		if o.Quantity <= 0 {
			return order.ErrorInvalidQuantity
		}

		// 1. Load product
		p, err := uc.prodvrt.GetProductByID(txCtx, o.ProductID)
		if err != nil {
			return fmt.Errorf("product not found: %w", err)
		}

		// 2. Branch by variant mode
		switch p.HasVariants {

		case true:
			if o.VariantID == nil {
				return order.ErrorVariantRequired
			}

			v, err := uc.prodvrt.GetVariantByID(txCtx, *o.VariantID)
			if err != nil {
				return fmt.Errorf("variant not found: %w", err)
			}

			if v.Stock < o.Quantity {
				return order.ErrorOutOfStock
			}

			// Snapshot
			o.UnitPrice = int64(v.Price * 100)
			o.Total = o.UnitPrice * int64(o.Quantity)
			o.ProductName = p.Name
			o.VariantName = v.SKU
			o.ImageURL = v.ImageURL

			// Decrement stock
			newStock := v.Stock - o.Quantity
			if err := uc.prodvrt.UpdateVariantStock(txCtx, v.ID, newStock); err != nil {
				return fmt.Errorf("update variant stock failed: %w", err)
			}

		case false:
			if o.VariantID != nil {
				return order.ErrorVariantNotAllowed
			}

			if p.Stock < o.Quantity {
				return order.ErrorOutOfStock
			}

			// Snapshot
			o.UnitPrice = int64(p.Price * 100)
			o.Total = o.UnitPrice * int64(o.Quantity)
			o.ProductName = p.Name
			o.ImageURL = firstImage(p.Images)

			newStock := p.Stock - o.Quantity
			if err := uc.prodvrt.UpdateProductStock(txCtx, p.ID, newStock); err != nil {
				return fmt.Errorf("update product stock failed: %w", err)
			}
		}

		// 3. Create order
		o.Status = order.Pending
		if err := uc.repo.Create(txCtx, o); err != nil {
			return fmt.Errorf("create order failed: %w", err)
		}

		// 4. Notifications AFTER commit
		go uc.fireOrderNotifications(ctx, o)

		return nil
	})
}

func (uc *UseCase) GetOrder(ctx context.Context, id uuid.UUID) (*order.Order, error) {
	return uc.repo.GetByID(ctx, id)
}

func (uc *UseCase) GetOrderByCustomer(ctx context.Context, customerID uuid.UUID) ([]*order.Order, error) {
	return uc.repo.ListByCustomer(ctx, customerID)
}

func (uc *UseCase) UpdateOrder(ctx context.Context, orderID uuid.UUID, column string, value any) error {
	return uc.txManager.Do(ctx, func(txCtx context.Context) error {
		if err := uc.repo.Update(txCtx, orderID, column, value); err != nil {
			return fmt.Errorf("update order failed: %w", err)
		}

		return nil
	})
}

func (uc *UseCase) ListOrders(ctx context.Context) ([]*order.Order, error) {
	return uc.repo.List(ctx)
}

func (uc *UseCase) DeleteOrder(ctx context.Context, id uuid.UUID) error {
	return uc.txManager.Do(ctx, func(txCtx context.Context) error {
		if err := uc.repo.Delete(txCtx, id); err != nil {
			return fmt.Errorf("delete order failed: %w", err)
		}

		return nil
	})
}

func (uc *UseCase) GetAllInventories(ctx context.Context) ([]order.Inventory, error) {
	return uc.invRepo.GetAllInventories(ctx)
}

func (uc *UseCase) GetAllCustomers(ctx context.Context) ([]order.Customer, error) {
	return uc.usrRepo.GetAllCustomers(ctx)
}

func (uc *UseCase) GetOrderPickupPoint(ctx context.Context, orderID uuid.UUID) (postgis.PointS, error) {
	return uc.repo.GetPickupPoint(ctx, orderID)
}

func (uc *UseCase) GetOrderDeliveryPoint(ctx context.Context, orderID uuid.UUID) (postgis.PointS, error) {
	return uc.repo.GetDeliveryPoint(ctx, orderID)
}

func (uc *UseCase) AssignOrderToDriver(ctx context.Context, orderID, driverID uuid.UUID, maxDistance float64) (*driver.Driver, error) {
	// 1. Fetch the order
	o, err := uc.GetOrder(ctx, orderID)
	if err != nil {
		return nil, fmt.Errorf("fetch order: %w", err)
	}

	if o.Status != order.Pending {
		return nil, fmt.Errorf("order not pending")
	}

	// 2. Get pickup point
	pickupPoint, err := uc.GetOrderPickupPoint(ctx, orderID)
	if err != nil {
		return nil, fmt.Errorf("get pickup point: %w", err)
	}

	// 3. Find nearest available driver within maxDistance
	nearestDriver, err := uc.drvRepo.GetNearestDriver(ctx, pickupPoint, maxDistance)
	if err != nil || nearestDriver == nil {
		return nil, fmt.Errorf("no available driver within %.2f meters", maxDistance)
	}

	// 4. Update order status to assigned
	if err := uc.UpdateOrder(ctx, orderID, "status", order.Assigned); err != nil {
		return nil, fmt.Errorf("update order status: %w", err)
	}

	// 5. Optionally: create a pending assignment record
	// This could be a lightweight table: order_id, driver_id, assigned_at
	// You can also push a websocket / notification here to the driver
	// if err := uc.repo.CreatePendingAssignment(ctx, orderID, nearestDriver.ID); err != nil {
	//     return nil, fmt.Errorf("create pending assignment: %w", err)
	// }

	// 6. Return assigned driver for confirmation / logging
	return nearestDriver, nil
}

func (uc *UseCase) notify(ctx context.Context, userID uuid.UUID, message string) error {
	n := &notification.Notification{
		UserID:  userID,
		Message: message,
		Type:    notification.System,
		Status:  notification.Pending,
	}
	return uc.notfRepo.Create(ctx, n)
}

func (uc *UseCase) fireOrderNotifications(ctx context.Context, o *order.Order) {
	msg := fmt.Sprintf("Your order %s has been placed successfully.", o.ID)
	_ = uc.notify(ctx, o.CustomerID, msg)
}

func firstImage(images []string) string {
	if len(images) == 0 {
		return ""
	}
	return images[0]
}
