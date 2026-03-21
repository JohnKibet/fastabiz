using System;
using System.Collections.Generic;
using System.Net.Http;
using System.Text.Json;
using System.Threading.Tasks;
using Microsoft.AspNetCore.Components.Forms;
using frontend.Models;
using frontend.Models.Storefront;

namespace frontend.Services;

/// <summary>
/// Product domain service.
/// All HTTP plumbing is handled by ApiService — this class only knows
/// about product-specific endpoints and request/response shapes.
/// </summary>
public class ProductService
{
    private readonly ApiService      _api;
    private readonly HttpClient      _cloudinaryHttp;
    private readonly ToastService    _toast;

    public ProductService(ApiService api, IHttpClientFactory factory, ToastService toast)
    {
        _api            = api;
        _cloudinaryHttp = factory.CreateClient("CloudinaryClient");
        _toast          = toast;
    }

    // ── Products ──────────────────────────────────────────────────────

    public Task<ApiResult<CreateProductResponse>> AddProduct(CreateProductRequest req)
        => _api.PostAsync<CreateProductRequest, CreateProductResponse>("products/create", req);

    public Task<ApiResult<List<ProductListItem>>> ListAllStoreProducts(Guid storeId)
        => _api.GetAsync<List<ProductListItem>>($"products/{storeId}/all_products");

    public Task<ApiResult<ProductMapper.ProductDto>> GetProductById(Guid productId)
        => _api.GetAsync<ProductMapper.ProductDto>($"products/by-id/{productId}");

    public Task<ApiResult<ProductMapper.ProductDto>> UpdateProductDetails(UpdateProductDetailsRequest req)
        => _api.PostAsync<UpdateProductDetailsRequest, ProductMapper.ProductDto>(
               $"products/{req.ProductId}/product_details", req);

    public Task<ApiResult<ApiMessageResponse>> UpdateProductInventory(UpdateProductInventoryRequest req)
        => _api.PatchAsync<UpdateProductInventoryRequest, ApiMessageResponse>("products/inventory", req);

    public Task<ApiResult<bool>> DeleteProduct(Guid productId)
        => _api.DeleteAsync($"products/{productId}/delete");

    // ── Options ───────────────────────────────────────────────────────

    public Task<ApiResult<List<Option>>> ListProductOptions(Guid productId)
        => _api.GetAsync<List<Option>>($"products/{productId}/options");

    public Task<ApiResult<CreateOptionNameResponse>> CreateOptionName(CreateOptionNameRequest req)
        => _api.PostAsync<CreateOptionNameRequest, CreateOptionNameResponse>("products/options/add", req);

    public Task<ApiResult<Product>> AddOptionValues(AddOptionValuesRequest req)
        => _api.PostAsync<AddOptionValuesRequest, Product>("products/options/values/add", req);

    public Task<ApiResult<bool>> DeleteOptionName(Guid optionId)
        => _api.DeleteAsync($"products/options/{optionId}/delete");

    public Task<ApiResult<bool>> DeleteOptionValue(Guid valueId)
        => _api.DeleteAsync($"products/options/values/{valueId}/delete");

    // ── Variants ──────────────────────────────────────────────────────

    public Task<ApiResult<CreateVariantResponse>> CreateVariant(CreateVariantRequest req)
        => _api.PostAsync<CreateVariantRequest, CreateVariantResponse>("products/variants/add", req);

    public Task<ApiResult<List<Product>>> UpdateVariantStock(UpdateVariantStockRequest req)
        => _api.PostAsync<UpdateVariantStockRequest, List<Product>>("products/variants/stock/update", req);

    public Task<ApiResult<Product>> UpdateVariantPrice(UpdateVariantPriceRequest req)
        => _api.PostAsync<UpdateVariantPriceRequest, Product>("products/variants/price/update", req);

    public Task<ApiResult<bool>> DeleteVariant(Guid variantId)
        => _api.DeleteAsync($"products/variants/{variantId}/delete");

    // ── Images ────────────────────────────────────────────────────────

    public Task<ApiResult<bool>> AddProductImages(AddImageRequest req)
        => _api.PostAsync<AddImageRequest, bool>("products/images/add", req);

    public Task<ApiResult<List<Product>>> ReorderProductImages(ReorderImagesRequest req)
        => _api.PostAsync<ReorderImagesRequest, List<Product>>("products/images/reorder", req);

    public Task<ApiResult<bool>> DeleteImage(Guid imageId)
        => _api.DeleteAsync($"products/images/{imageId}/delete");

    // ── Cloudinary (separate HTTP client — not routed through ApiService) ──

    public Task<ApiResult<CloudinarySignatureResponse>> GetCloudinarySignature()
        => _api.PostAsync<CloudinarySignatureResponse>("products/cloudinary/signature");

    public async Task<string?> UploadToCloudinary(IBrowserFile file, CloudinarySignatureResponse sig)
    {
        using var content = new MultipartFormDataContent();
        content.Add(new StreamContent(file.OpenReadStream(10_000_000)), "file", file.Name);
        content.Add(new StringContent(sig.Api_Key),              "api_key");
        content.Add(new StringContent(sig.Timestamp.ToString()), "timestamp");
        content.Add(new StringContent(sig.Signature),            "signature");
        content.Add(new StringContent(sig.Folder),               "folder");

        var url      = $"https://api.cloudinary.com/v1_1/{sig.Cloud_Name}/image/upload";
        var response = await _cloudinaryHttp.PostAsync(url, content);

        if (!response.IsSuccessStatusCode) return null;

        using var doc = JsonDocument.Parse(await response.Content.ReadAsStringAsync());
        return doc.RootElement.GetProperty("secure_url").GetString();
    }
}
