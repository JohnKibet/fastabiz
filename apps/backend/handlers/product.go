package handlers

import (
	"backend/internal/application"
	"backend/internal/domain/product"
	"encoding/json"
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
		writeJSONError(w, http.StatusInternalServerError, "Failed to create product", err)
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

// ListProducts godoc
// @Summary List all products
// @Description Retrieves all products available to the merchant/admin.
// @Tags products
// @Security BearerAuth
// @Produce json
// @Success 200 {array} product.Product "List of products"
// @Failure 401 {object} handlers.ErrorResponse "Unauthorized"
// @Failure 500 {object} handlers.ErrorResponse "Internal server error"
// @Router /products/all_products [get]
func (h *ProductHandler) ListProducts(w http.ResponseWriter, r *http.Request) {
	products, err := h.UC.Products.UseCase.GetAllProducts(r.Context())
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
// @Router /products/{id}/product [delete]
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

	if err := h.UC.Products.UseCase.AddImage(r.Context(), req.ProductID, req.URL, req.IsPrimary); err != nil {
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
// @Router /products/images/{imageId} [delete]
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
// @Router /products/options/{optionId} [delete]
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
// @Summary Add option value
// @Description Adds a value to an existing product option (e.g. Red, Large).
// @Tags products
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param value body product.AddOptionValueRequest true "Add option value payload"
// @Success 201 {object} map[string]string "Option value added successfully"
// @Failure 400 {object} handlers.ErrorResponse "Invalid request body"
// @Failure 401 {object} handlers.ErrorResponse "Unauthorized"
// @Failure 404 {object} handlers.ErrorResponse "Option not found"
// @Failure 500 {object} handlers.ErrorResponse "Internal server error"
// @Router /products/options/values/add [post]
func (h *ProductHandler) AddOptionValue(w http.ResponseWriter, r *http.Request) {
	var req product.AddOptionValueRequest

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeJSONError(w, http.StatusBadRequest, "Invalid request body", nil)
		return
	}

	if err := h.UC.Products.UseCase.AddOptionValue(r.Context(), req.OptionID, req.Value); err != nil {
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
// @Router /products/options/values/{valueId} [delete]
func (h *ProductHandler) DeleteOptionValue(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "valueId")
	valueID, err := uuid.Parse(idStr)

	if err != nil {
		writeJSONError(w, http.StatusBadRequest, "Invalid option value ID", nil)
		return
	}

	if err := h.UC.Products.UseCase.DeleteOptionValue(r.Context(), valueID); err != nil {
		writeJSONError(w, http.StatusInternalServerError, "Failed to delete option value", err)
		return
	}

	writeJSON(w, http.StatusOK, map[string]string{"message": "Option value deleted"})
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

	v := req.ToVariant()

	if err := h.UC.Products.UseCase.CreateVariant(r.Context(), v); err != nil {
		writeJSONError(w, http.StatusInternalServerError, "Failed to create variant", err)
		return
	}

	writeJSON(w, http.StatusCreated, v)
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
// @Router /products/variants/{variantId} [delete]
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
