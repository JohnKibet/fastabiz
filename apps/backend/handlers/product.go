package handlers

import (
	"backend/internal/application"
	"backend/internal/domain/product"
	"encoding/json"
	"errors"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
)

type ProductHandler struct {
	UC *application.OrderService
}

func NewProductHandler(uc *application.OrderService) *ProductHandler {
	return &ProductHandler{UC: uc}
}

// CreateProduct godoc
// @Summary Create a new product
// @Description Creates a new product for a merchant, including optional images, options, and variants.
// @Tags products
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param product body product.CreateProductRequest true "Product creation payload"
// @Success 201 {object} product.Product
// @Failure 400 {object} handlers.ErrorResponse "Invalid request body"
// @Failure 401 {object} handlers.ErrorResponse "Unauthorized"
// @Failure 500 {object} handlers.ErrorResponse "Internal server error"
// @Router /products/create [post]
func (h *ProductHandler) CreateProduct(w http.ResponseWriter, r *http.Request) {
	var req product.CreateProductRequest

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeJSONError(w, http.StatusBadRequest, "Invalid request body", nil)
		return
	}

	p := req.ToProduct()

	if err := h.UC.Products.UseCase.CreateProduct(r.Context(), p); err != nil {
		switch {
		case errors.Is(err, product.ErrProductInvalidStore):
			writeJSONError(w, http.StatusBadRequest, "Invalid store", nil)

		case errors.Is(err, product.ErrProductInvalidInput):
			writeJSONError(w, http.StatusBadRequest, "Invalid product data", nil)

		case errors.Is(err, product.ErrProductAlreadyExists):
			writeJSONError(w, http.StatusConflict, "Product already exists", nil)

		default:
			writeJSONError(w, http.StatusInternalServerError, "Failed to create product", nil)
		}

		return
	}

	writeJSON(w, http.StatusCreated, map[string]any{
		"id":           p.ID,
		"store_id":     p.StoreID,
		"name":         p.Name,
		"description":  p.Description,
		"category":     p.Category,
		"images":       p.Images,
		"has_variants": p.HasVariants,
		"price":        p.Price,
		"stock":        p.Stock,
		"options":      p.Options,
		"variants":     p.Variants,
		"created_at":   p.CreatedAt,
		"updated_at":   p.UpdatedAt,
	})
}

// GetProductByID godoc
// @Summary Get product by ID
// @Description Retrieves a product by its ID, including images, options, and variants.
// @Tags products
// @Accept json
// @Produce json
// @Param id path string true "Product ID"
// @Success 200 {object} product.Product
// @Failure 400 {object} handlers.ErrorResponse "Invalid product ID"
// @Failure 404 {object} handlers.ErrorResponse "Product not found"
// @Failure 500 {object} handlers.ErrorResponse "Internal server error"
// @Router /products/by-id/{id} [get]
func (h *ProductHandler) GetProductByID(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	id, err := uuid.Parse(idStr)

	if err != nil {
		writeJSONError(w, http.StatusBadRequest, "Invalid ID", nil)
		return
	}

	p, err := h.UC.Products.UseCase.GetProductByID(r.Context(), id)
	if err != nil {
		writeJSONError(w, http.StatusNotFound, "Product not found", err)
		return
	}

	writeJSON(w, http.StatusOK, p)
}

// UpdateProductDetails godoc
// @Summary Update product details
// @Description Updates the core details of a product (name, description, category).
// @Tags products
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param product body product.UpdateProductDetailsRequest true "Product update payload"
// @Success 200 {object} map[string]string "Product updated successfully"
// @Failure 400 {object} handlers.ErrorResponse "Invalid request body"
// @Failure 401 {object} handlers.ErrorResponse "Unauthorized"
// @Failure 404 {object} handlers.ErrorResponse "Product not found"
// @Failure 500 {object} handlers.ErrorResponse "Internal server error"
// @Router /products/{id}/product_details [patch]
func (h *ProductHandler) UpdateProductDetails(w http.ResponseWriter, r *http.Request) {
	var req product.UpdateProductDetailsRequest

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeJSONError(w, http.StatusBadRequest, "Invalid request body", nil)
		return
	}

	if err := h.UC.Products.UseCase.UpdateProductDetails(r.Context(), &req); err != nil {
		writeJSONError(w, http.StatusInternalServerError, "Failed to update product details", err)
		return
	}

	writeJSON(w, http.StatusOK, map[string]any{
		"message": "Product updated successfully",
	})
}

// ListProductsByStore godoc
// @Summary List products by store
// @Description Retrieves all products belonging to a specific store.
// @Tags products
// @Security BearerAuth
// @Produce json
// @Param storeId path string true "Store ID"
// @Success 200 {array} product.ProductListItem "List of products"
// @Failure 400 {object} handlers.ErrorResponse "Invalid store ID"
// @Failure 401 {object} handlers.ErrorResponse "Unauthorized"
// @Failure 500 {object} handlers.ErrorResponse "Internal server error"
// @Router /products/{store_id}/all_products [get]
func (h *ProductHandler) ListProductsByStore(w http.ResponseWriter, r *http.Request) {
	storeIDStr := chi.URLParam(r, "store_id")
	storeID, err := uuid.Parse(storeIDStr)

	if err != nil {
		writeJSONError(w, http.StatusBadRequest, "Invalid store ID", err)
		return
	}

	products, err := h.UC.Products.UseCase.GetAllProducts(r.Context(), storeID)
	if err != nil {
		writeJSONError(w, http.StatusInternalServerError, "Could not fetch products", err)
		return
	}

	writeJSON(w, http.StatusOK, products)
}

// DeleteProduct godoc
// @Summary Delete a product
// @Description Permanently deletes a product by its ID.
// @Tags products
// @Security BearerAuth
// @Produce json
// @Param id path string true "Product ID"
// @Success 200 {object} map[string]string "Product deleted successfully"
// @Failure 400 {object} handlers.ErrorResponse "Invalid product ID"
// @Failure 401 {object} handlers.ErrorResponse "Unauthorized"
// @Failure 404 {object} handlers.ErrorResponse "Product not found"
// @Failure 500 {object} handlers.ErrorResponse "Internal server error"
// @Router /products/{id}/delete [delete]
func (h *ProductHandler) DeleteProduct(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	productID, err := uuid.Parse(idStr)

	if err != nil {
		writeJSONError(w, http.StatusBadRequest, "Invalid ID", nil)
		return
	}

	if err := h.UC.Products.UseCase.DeleteProduct(r.Context(), productID); err != nil {
		writeJSONError(w, http.StatusInternalServerError, "Failed to delete product", err)
		return
	}

	writeJSON(w, http.StatusOK, map[string]any{
		"message": "Product deleted",
	})
}

// AddImage godoc
// @Summary Add product image
// @Description Adds an image to a product. Optionally marks it as the primary image.
// @Tags products
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param image body product.AddImageRequest true "Add product image payload"
// @Success 201 {object} map[string]string "Image added successfully"
// @Failure 400 {object} handlers.ErrorResponse "Invalid request body"
// @Failure 401 {object} handlers.ErrorResponse "Unauthorized"
// @Failure 404 {object} handlers.ErrorResponse "Product not found"
// @Failure 500 {object} handlers.ErrorResponse "Internal server error"
// @Router /products/images/add [post]
func (h *ProductHandler) AddImage(w http.ResponseWriter, r *http.Request) {
	var req product.AddImageRequest

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeJSONError(w, http.StatusBadRequest, "Invalid request body", nil)
		return
	}

	images := make([]product.Image, 0, len(req.Images))
	for _, img := range req.Images {
		images = append(images, product.Image{
			URL:       img.URL,
			IsPrimary: img.IsPrimary,
		})
	}

	if err := h.UC.Products.UseCase.AddImage(r.Context(), req.ProductID, images); err != nil {
		writeJSONError(w, http.StatusInternalServerError, "Failed to add image", err)
		return
	}

	writeJSON(w, http.StatusOK, map[string]any{
		"message": "Image added",
	})
}

// DeleteImage godoc
// @Summary Delete product image
// @Description Deletes a product image by its ID.
// @Tags products
// @Security BearerAuth
// @Produce json
// @Param imageId path string true "Image ID"
// @Success 200 {object} map[string]string "Image deleted successfully"
// @Failure 400 {object} handlers.ErrorResponse "Invalid image ID"
// @Failure 401 {object} handlers.ErrorResponse "Unauthorized"
// @Failure 404 {object} handlers.ErrorResponse "Image not found"
// @Failure 500 {object} handlers.ErrorResponse "Internal server error"
// @Router /products/images/{imageId}/delete [delete]
func (h *ProductHandler) DeleteImage(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "imageId")
	imageID, err := uuid.Parse(idStr)

	if err != nil {
		writeJSONError(w, http.StatusBadRequest, "Invalid image ID", nil)
		return
	}

	if err := h.UC.Products.UseCase.DeleteImage(r.Context(), imageID); err != nil {
		writeJSONError(w, http.StatusInternalServerError, "Failed to delete image", err)
		return
	}

	writeJSON(w, http.StatusOK, map[string]string{"message": "Image deleted"})
}

// ReorderImages godoc
// @Summary Reorder product images
// @Description Updates the display order of product images using an ordered list of image IDs.
// @Tags products
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param reorder body product.ReorderImagesRequest true "Reorder images payload"
// @Success 200 {object} map[string]string "Images reordered successfully"
// @Failure 400 {object} handlers.ErrorResponse "Invalid request body"
// @Failure 401 {object} handlers.ErrorResponse "Unauthorized"
// @Failure 404 {object} handlers.ErrorResponse "Product or images not found"
// @Failure 500 {object} handlers.ErrorResponse "Internal server error"
// @Router /products/images/reorder [patch]
func (h *ProductHandler) ReorderImages(w http.ResponseWriter, r *http.Request) {
	var req product.ReorderImagesRequest

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeJSONError(w, http.StatusBadRequest, "Invalid request body", nil)
		return
	}

	if err := h.UC.Products.UseCase.ReorderImages(r.Context(), req.ProductID, req.ImageIDs); err != nil {
		writeJSONError(w, http.StatusInternalServerError, "Failed to reorder images", err)
		return
	}

	writeJSON(w, http.StatusOK, map[string]string{"message": "Images reordered"})
}

// AddOptionName godoc
// @Summary Add product option
// @Description Adds a new option name (e.g. Size, Color) to a product.
// @Tags products
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param option body product.AddOptionNameRequest true "Add option name payload"
// @Success 201 {object} map[string]uuid.UUID "Option created"
// @Failure 400 {object} handlers.ErrorResponse "Invalid request body"
// @Failure 401 {object} handlers.ErrorResponse "Unauthorized"
// @Failure 404 {object} handlers.ErrorResponse "Product not found"
// @Failure 500 {object} handlers.ErrorResponse "Internal server error"
// @Router /products/options/add [post]
func (h *ProductHandler) AddOptionName(w http.ResponseWriter, r *http.Request) {
	var req product.AddOptionNameRequest

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeJSONError(w, http.StatusBadRequest, "Invalid request body", nil)
		return
	}

	optionID, err := h.UC.Products.UseCase.AddOptionName(r.Context(), req.ProductID, req.Name)
	if err != nil {
		writeJSONError(w, http.StatusInternalServerError, "Failed to add option", err)
		return
	}

	writeJSON(w, http.StatusOK, map[string]any{
		"option_id": optionID,
	})
}

// DeleteOptionName godoc
// @Summary Delete product option
// @Description Deletes an option name from a product.
// @Tags products
// @Security BearerAuth
// @Produce json
// @Param optionId path string true "Option ID"
// @Success 200 {object} map[string]string "Option deleted successfully"
// @Failure 400 {object} handlers.ErrorResponse "Invalid option ID"
// @Failure 401 {object} handlers.ErrorResponse "Unauthorized"
// @Failure 404 {object} handlers.ErrorResponse "Option not found"
// @Failure 500 {object} handlers.ErrorResponse "Internal server error"
// @Router /products/options/{optionId}/delete [delete]
func (h *ProductHandler) DeleteOptionName(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "optionId")
	optionID, err := uuid.Parse(idStr)

	if err != nil {
		writeJSONError(w, http.StatusBadRequest, "Invalid option ID", nil)
		return
	}

	if err := h.UC.Products.UseCase.DeleteOptionName(r.Context(), optionID); err != nil {
		writeJSONError(w, http.StatusInternalServerError, "Failed to delete option", err)
		return
	}

	writeJSON(w, http.StatusOK, map[string]string{"message": "Option deleted"})
}

// AddOptionValue godoc
// @Summary Add option values
// @Description Adds values to an existing product option (e.g. Red, Large).
// @Tags products
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param value body product.AddOptionValuesRequest true "Add option values payload"
// @Success 201 {object} map[string]string "Option values added successfully"
// @Failure 400 {object} handlers.ErrorResponse "Invalid request body"
// @Failure 401 {object} handlers.ErrorResponse "Unauthorized"
// @Failure 404 {object} handlers.ErrorResponse "Option not found"
// @Failure 500 {object} handlers.ErrorResponse "Internal server error"
// @Router /products/options/values/add [post]
func (h *ProductHandler) AddOptionValue(w http.ResponseWriter, r *http.Request) {
	var req product.AddOptionValuesRequest

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeJSONError(w, http.StatusBadRequest, "Invalid request body", nil)
		return
	}

	if len(req.Values) == 0 {
		writeJSONError(w, http.StatusBadRequest, "Values array cannot be empty", nil)
		return
	}

	if err := h.UC.Products.UseCase.AddOptionValue(r.Context(), req.ProductID, req.OptionID, req.Values); err != nil {
		writeJSONError(w, http.StatusInternalServerError, "Failed to add option value", err)
		return
	}

	writeJSON(w, http.StatusOK, map[string]string{"message": "Option value added"})
}

// DeleteOptionValue godoc
// @Summary Delete option value
// @Description Deletes a value from a product option.
// @Tags products
// @Security BearerAuth
// @Produce json
// @Param valueId path string true "Option Value ID"
// @Success 200 {object} map[string]string "Option value deleted successfully"
// @Failure 400 {object} handlers.ErrorResponse "Invalid option value ID"
// @Failure 401 {object} handlers.ErrorResponse "Unauthorized"
// @Failure 404 {object} handlers.ErrorResponse "Option value not found"
// @Failure 500 {object} handlers.ErrorResponse "Internal server error"
// @Router /products/options/values/{valueId}/delete [delete]
func (h *ProductHandler) DeleteOptionValue(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "valueId")
	valueID, err := uuid.Parse(idStr)

	if err != nil {
		writeJSONError(w, http.StatusBadRequest, "Invalid option value ID", nil)
		return
	}

	if err := h.UC.Products.UseCase.DeleteOptionValue(r.Context(), valueID); err != nil {
		switch {
		case errors.Is(err, product.ErrOptionValueInUse):
			writeJSONError(w, http.StatusConflict, "Option value is used by existing variants", err)
		default:
			writeJSONError(w, http.StatusInternalServerError, "Failed to delete option value", err)
		}

		return
	}

	writeJSON(w, http.StatusOK, map[string]string{"message": "Option value deleted"})
}

// ListOptions godoc
// @Summary List product options
// @Description Returns all options and their values for a given product
// @Tags products
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param productId path string true "Product ID (UUID)"
// @Success 200 {array} product.Option
// @Failure 400 {object} handlers.ErrorResponse "Invalid product ID"
// @Failure 401 {object} handlers.ErrorResponse "Unauthorized"
// @Failure 500 {object} handlers.ErrorResponse "Failed to list options"
// @Router /products/{productId}/options [get]
func (h *ProductHandler) ListOptions(w http.ResponseWriter, r *http.Request) {
	productID, err := uuid.Parse(chi.URLParam(r, "productId"))
	if err != nil {
		writeJSONError(w, http.StatusBadRequest, "Invalid product ID", nil)
		return
	}

	options, err := h.UC.Products.UseCase.ListProductOptions(r.Context(), productID)
	if err != nil {
		writeJSONError(w, http.StatusInternalServerError, "Failed to list options", err)
		return
	}

	writeJSON(w, http.StatusOK, options)
}

// CreateVariant godoc
// @Summary Create product variant
// @Description Creates a purchasable variant for a product based on option values.
// @Tags products
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param variant body product.CreateVariantRequest true "Create variant payload"
// @Success 201 {object} product.Variant
// @Failure 400 {object} handlers.ErrorResponse "Invalid request body"
// @Failure 401 {object} handlers.ErrorResponse "Unauthorized"
// @Failure 404 {object} handlers.ErrorResponse "Product or option values not found"
// @Failure 500 {object} handlers.ErrorResponse "Internal server error"
// @Router /products/variants/add [post]
func (h *ProductHandler) CreateVariant(w http.ResponseWriter, r *http.Request) {
	var req product.CreateVariantRequest

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeJSONError(w, http.StatusBadRequest, "Invalid request body", nil)
		return
	}

	variant, err := h.UC.Products.UseCase.CreateVariant(r.Context(), req)
	if err != nil {
		switch {
		case errors.Is(err, product.ErrProductNotFound):
			writeJSONError(w, http.StatusNotFound, "Product not found", err)
		case errors.Is(err, product.ErrOptionNotFound):
			writeJSONError(w, http.StatusNotFound, "Option not found", err)
		case errors.Is(err, product.ErrOptionValueNotFound):
			writeJSONError(w, http.StatusNotFound, "Option value not found", err)
		default:
			writeJSONError(w, http.StatusInternalServerError, "Failed to create variant", err)
		}
		return
	}

	writeJSON(w, http.StatusCreated, variant)
}

// UpdateVariantStock godoc
// @Summary Update variant stock
// @Description Updates the stock quantity of a product variant.
// @Tags products
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param stock body product.UpdateVariantStockRequest true "Update variant stock payload"
// @Success 200 {object} map[string]string "Variant stock updated successfully"
// @Failure 400 {object} handlers.ErrorResponse "Invalid request body"
// @Failure 401 {object} handlers.ErrorResponse "Unauthorized"
// @Failure 404 {object} handlers.ErrorResponse "Variant not found"
// @Failure 500 {object} handlers.ErrorResponse "Internal server error"
// @Router /products/variants/stock/update [patch]
func (h *ProductHandler) UpdateVariantStock(w http.ResponseWriter, r *http.Request) {
	var req product.UpdateVariantStockRequest

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeJSONError(w, http.StatusBadRequest, "Invalid request body", nil)
		return
	}

	if err := h.UC.Products.UseCase.UpdateVariantStock(r.Context(), req.VariantID, req.Stock); err != nil {
		writeJSONError(w, http.StatusInternalServerError, "Failed to update stock", err)
		return
	}

	writeJSON(w, http.StatusOK, map[string]string{"message": "Variant stock updated"})
}

// UpdateVariantPrice godoc
// @Summary Update variant price
// @Description Updates the price of a product variant.
// @Tags products
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param price body product.UpdateVariantPriceRequest true "Update variant price payload"
// @Success 200 {object} map[string]string "Variant price updated successfully"
// @Failure 400 {object} handlers.ErrorResponse "Invalid request body"
// @Failure 401 {object} handlers.ErrorResponse "Unauthorized"
// @Failure 404 {object} handlers.ErrorResponse "Variant not found"
// @Failure 500 {object} handlers.ErrorResponse "Internal server error"
// @Router /products/variants/price/update [patch]
func (h *ProductHandler) UpdateVariantPrice(w http.ResponseWriter, r *http.Request) {
	var req product.UpdateVariantPriceRequest

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeJSONError(w, http.StatusBadRequest, "Invalid request body", nil)
		return
	}

	if err := h.UC.Products.UseCase.UpdateVariantPrice(r.Context(), req.VariantID, req.Price); err != nil {
		writeJSONError(w, http.StatusInternalServerError, "Failed to update price", err)
		return
	}

	writeJSON(w, http.StatusOK, map[string]string{"message": "Variant price updated"})
}

// DeleteVariant godoc
// @Summary Delete product variant
// @Description Permanently deletes a product variant.
// @Tags products
// @Security BearerAuth
// @Produce json
// @Param variantId path string true "Variant ID"
// @Success 200 {object} map[string]string "Variant deleted successfully"
// @Failure 400 {object} handlers.ErrorResponse "Invalid variant ID"
// @Failure 401 {object} handlers.ErrorResponse "Unauthorized"
// @Failure 404 {object} handlers.ErrorResponse "Variant not found"
// @Failure 500 {object} handlers.ErrorResponse "Internal server error"
// @Router /products/variants/{variantId}/delete [delete]
func (h *ProductHandler) DeleteVariant(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "variantId")
	variantID, err := uuid.Parse(idStr)

	if err != nil {
		writeJSONError(w, http.StatusBadRequest, "Invalid variant ID", nil)
		return
	}

	if err := h.UC.Products.UseCase.DeleteVariant(r.Context(), variantID); err != nil {
		writeJSONError(w, http.StatusInternalServerError, "Failed to delete variant", err)
		return
	}

	writeJSON(w, http.StatusOK, map[string]string{"message": "Variant deleted"})
}
