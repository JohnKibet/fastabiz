using System;
using System.Text.Json;
using System.Net.Http.Json;
using frontend.Models;
using Microsoft.AspNetCore.Components.Forms;
public class ProductService
{
  private readonly HttpClient _http;
  private readonly HttpClient _cloudinaryHttp;
  private readonly ToastService _toastService;
  public ProductService(IHttpClientFactory httpClientFactory, ToastService toastService)
  {
    _http = httpClientFactory.CreateClient("AuthenticatedApi");
    _cloudinaryHttp = httpClientFactory.CreateClient("CloudinaryClient");
    _toastService = toastService;
  }

  public async Task<ServiceResult2<CreateProductResponse>> AddProduct(CreateProductRequest product)
  {
    try
    {
      var response = await _http.PostAsJsonAsync("products/create", product);

      if (!response.IsSuccessStatusCode)
      {
        var error = await ParseError(response);
        return ServiceResult2<CreateProductResponse>.Fail(error);
      }

      var createdProduct = await response.Content.ReadFromJsonAsync<CreateProductResponse>();

      if (createdProduct is null)
      {
        return ServiceResult2<CreateProductResponse>.Fail("Invalid server response");
      }

      return ServiceResult2<CreateProductResponse>.Ok(createdProduct);
    }
    catch (HttpRequestException ex)
    {
      return ServiceResult2<CreateProductResponse>.Fail($"Network error: {ex.Message}");
    }
    catch (Exception ex)
    {
      return ServiceResult2<CreateProductResponse>.Fail($"Unexpected error: {ex.Message}");
    }
  }

  public async Task<ServiceResult2<List<ProductListItem>>> ListAllStoreProducts(Guid StoreId)
  {
    return await GetFromJsonSafe<List<ProductListItem>>($"products/{StoreId}/all_products");
  }

  // get signature
  public async Task<ServiceResult2<CloudinarySignatureResponse>> GetCloudinarySignature()
  {
    try
    {
      var response = await _http.PostAsync("products/cloudinary/signature", null);

      if (!response.IsSuccessStatusCode)
        return ServiceResult2<CloudinarySignatureResponse>.Fail(
            await ParseError(response)
        );

      var data = await response.Content.ReadFromJsonAsync<CloudinarySignatureResponse>();
      return ServiceResult2<CloudinarySignatureResponse>.Ok(data!);
    }
    catch (Exception ex)
    {
      return ServiceResult2<CloudinarySignatureResponse>.Fail(ex.Message);
    }
  }

  // upload image to cloudinary
  public async Task<string?> UploadToCloudinary(IBrowserFile file, CloudinarySignatureResponse sig)
  {
    using var content = new MultipartFormDataContent();

    content.Add(new StreamContent(file.OpenReadStream(10_000_000)), "file", file.Name);

    content.Add(new StringContent(sig.Api_Key), "api_key");
    content.Add(new StringContent(sig.Timestamp.ToString()), "timestamp");
    content.Add(new StringContent(sig.Signature), "signature");
    content.Add(new StringContent(sig.Folder), "folder");

    var url = $"https://api.cloudinary.com/v1_1/{sig.Cloud_Name}/image/upload";

    var response = await _cloudinaryHttp.PostAsync(url, content);

    if (!response.IsSuccessStatusCode)
      return null;

    using var doc = JsonDocument.Parse(await response.Content.ReadAsStringAsync());
    return doc.RootElement.GetProperty("secure_url").GetString();
  }


  public async Task<ServiceResult2<bool>> AddProductImages(AddImageRequest request)
  {
    try
    {
      var response = await _http.PostAsJsonAsync("products/images/add", request);
      if (response.IsSuccessStatusCode)
      {
        var result = await response.Content.ReadFromJsonAsync<bool>();
        return ServiceResult2<bool>.Ok(true);
      }

      var error = await ParseError(response);
      return ServiceResult2<bool>.Fail(error);
    }
    catch (HttpRequestException ex)
    {
      return ServiceResult2<bool>.Fail($"Network error: {ex.Message}");
    }
    catch (Exception ex)
    {
      return ServiceResult2<bool>.Fail($"Unexpected error: {ex.Message}");
    }
  }

  public async Task<ServiceResult2<CreateOptionNameResponse>> CreateOptionName(CreateOptionNameRequest request)
  {
    try
    {
      var response = await _http.PostAsJsonAsync("products/options/add", request);
      if (!response.IsSuccessStatusCode)
      {
        var error = await ParseError(response);
        return ServiceResult2<CreateOptionNameResponse>.Fail(error);
      }

      var optionId = await response.Content.ReadFromJsonAsync<CreateOptionNameResponse>();

      if (optionId is null)
      {
        return ServiceResult2<CreateOptionNameResponse>.Fail("Invalid server response");
      }

      return ServiceResult2<CreateOptionNameResponse>.Ok(optionId);
    }
    catch (HttpRequestException ex)
    {
      return ServiceResult2<CreateOptionNameResponse>.Fail($"Network error: {ex.Message}");
    }
    catch (Exception ex)
    {
      return ServiceResult2<CreateOptionNameResponse>.Fail($"Unexpected error: {ex.Message}");
    }
  }

  public async Task<ServiceResult2<ProductX>> AddOptionValues(AddOptionValuesRequest request)
  {
    try
    {
      var response = await _http.PostAsJsonAsync("products/options/values/add", request);
      if (response.IsSuccessStatusCode)
      {
        var result = await response.Content.ReadFromJsonAsync<ProductX>();
        return ServiceResult2<ProductX>.Ok(result ?? new ProductX());
      }

      var error = await ParseError(response);
      return ServiceResult2<ProductX>.Fail(error);
    }
    catch (HttpRequestException ex)
    {
      return ServiceResult2<ProductX>.Fail($"Network error: {ex.Message}");
    }
    catch (Exception ex)
    {
      return ServiceResult2<ProductX>.Fail($"Unexpected error: {ex.Message}");
    }
  }

  public async Task<ServiceResult2<CreateVariantResponse>> CreateVariant(CreateVariantRequest request)
  {
    try
    {
      var response = await _http.PostAsJsonAsync("products/variants/add", request);
      if (response.IsSuccessStatusCode)
      {
        var result = await response.Content.ReadFromJsonAsync<CreateVariantResponse>();
        return ServiceResult2<CreateVariantResponse>.Ok(result ?? new CreateVariantResponse());
      }

      var error = await ParseError(response);
      return ServiceResult2<CreateVariantResponse>.Fail(error);
    }
    catch (HttpRequestException ex)
    {
      return ServiceResult2<CreateVariantResponse>.Fail($"Network error: {ex.Message}");
    }
    catch (Exception ex)
    {
      return ServiceResult2<CreateVariantResponse>.Fail($"Unexpected error: {ex.Message}");
    }
  }

  public async Task<ServiceResult2<List<ProductX>>> UpdateVariantStock(UpdateVariantStockRequest request)
  {
    try
    {
      var response = await _http.PostAsJsonAsync("products/variants/stock/update", request);
      if (response.IsSuccessStatusCode)
      {
        var result = await response.Content.ReadFromJsonAsync<List<ProductX>>();
        return ServiceResult2<List<ProductX>>.Ok(result ?? new List<ProductX>());
      }

      var error = await ParseError(response);
      return ServiceResult2<List<ProductX>>.Fail(error);
    }
    catch (HttpRequestException ex)
    {
      return ServiceResult2<List<ProductX>>.Fail($"Network error: {ex.Message}");
    }
    catch (Exception ex)
    {
      return ServiceResult2<List<ProductX>>.Fail($"Unexpected error: {ex.Message}");
    }
  }

  public async Task<ServiceResult2<ProductX>> UpdateVariantPrice(UpdateVariantPriceRequest request)
  {
    try
    {
      var response = await _http.PostAsJsonAsync("products/variants/price/update", request);
      if (response.IsSuccessStatusCode)
      {
        var result = await response.Content.ReadFromJsonAsync<ProductX>();
        return ServiceResult2<ProductX>.Ok(result ?? new ProductX());
      }

      var error = await ParseError(response);
      return ServiceResult2<ProductX>.Fail(error);
    }
    catch (HttpRequestException ex)
    {
      return ServiceResult2<ProductX>.Fail($"Network error: {ex.Message}");
    }
    catch (Exception ex)
    {
      return ServiceResult2<ProductX>.Fail($"Unexpected error: {ex.Message}");
    }
  }

  public async Task<ServiceResult2<List<ProductX>>> ReorderProductImages(ReorderImagesRequest request)
  {
    try
    {
      var response = await _http.PostAsJsonAsync("products/images/reorder", request);
      if (response.IsSuccessStatusCode)
      {
        var result = await response.Content.ReadFromJsonAsync<List<ProductX>>();
        return ServiceResult2<List<ProductX>>.Ok(result ?? new List<ProductX>());
      }

      var error = await ParseError(response);
      return ServiceResult2<List<ProductX>>.Fail(error);
    }
    catch (HttpRequestException ex)
    {
      return ServiceResult2<List<ProductX>>.Fail($"Network error: {ex.Message}");
    }
    catch (Exception ex)
    {
      return ServiceResult2<List<ProductX>>.Fail($"Unexpected error: {ex.Message}");
    }
  }

  public async Task<ServiceResult2<ProductMapper.ProductDto>> UpdateProductDetails(UpdateProductDetailsRequest request)
  {
    try
    {
      var response = await _http.PostAsJsonAsync($"products/{request.ProductId}/product_details", request);
      if (response.IsSuccessStatusCode)
      {
        var result = await response.Content.ReadFromJsonAsync<ProductMapper.ProductDto>();
        return ServiceResult2<ProductMapper.ProductDto>.Ok(result ?? new ProductMapper.ProductDto());
      }

      var error = await ParseError(response);
      return ServiceResult2<ProductMapper.ProductDto>.Fail(error);
    }
    catch (HttpRequestException ex)
    {
      return ServiceResult2<ProductMapper.ProductDto>.Fail($"Network error: {ex.Message}");
    }
    catch (Exception ex)
    {
      return ServiceResult2<ProductMapper.ProductDto>.Fail($"Unexpected error: {ex.Message}");
    }
  }

  public async Task<ServiceResult2<List<Option>>> ListProductOptions(Guid productId)
  {
    return await GetFromJsonSafe<List<Option>>($"products/{productId}/options");
  }

  public async Task<ServiceResult2<ProductMapper.ProductDto>> GetProductByID(Guid productId)
  {
    return await GetFromJsonSafe<ProductMapper.ProductDto>($"products/by-id/{productId}");
  }

  public async Task<ServiceResult2<bool>> DeleteProduct(Guid productId)
  {
    try
    {
      var response = await _http.DeleteAsync($"products/{productId}/delete");
      if (response.IsSuccessStatusCode)
      {
        return ServiceResult2<bool>.Ok(true);
      }

      var error = await ParseError(response);
      return ServiceResult2<bool>.Fail(error);
    }
    catch (HttpRequestException ex)
    {
      return ServiceResult2<bool>.Fail($"Network error: {ex.Message}");
    }
    catch (Exception ex)
    {
      return ServiceResult2<bool>.Fail($"Unexpected error: {ex.Message}");
    }
  }

  public async Task<ServiceResult2<bool>> DeleteImage(Guid imageId)
  {
    try
    {
      var response = await _http.DeleteAsync($"products/images/{imageId}/delete");
      if (response.IsSuccessStatusCode)
      {
        return ServiceResult2<bool>.Ok(true);
      }

      var error = await ParseError(response);
      return ServiceResult2<bool>.Fail(error);
    }
    catch (HttpRequestException ex)
    {
      return ServiceResult2<bool>.Fail($"Network error: {ex.Message}");
    }
    catch (Exception ex)
    {
      return ServiceResult2<bool>.Fail($"Unexpected error: {ex.Message}");
    }
  }

  public async Task<ServiceResult2<bool>> DeleteOptionName(Guid optionId)
  {
    try
    {
      var response = await _http.DeleteAsync($"products/options/{optionId}/delete");
      if (response.IsSuccessStatusCode)
      {
        return ServiceResult2<bool>.Ok(true);
      }

      var error = await ParseError(response);
      return ServiceResult2<bool>.Fail(error);
    }
    catch (HttpRequestException ex)
    {
      return ServiceResult2<bool>.Fail($"Network error: {ex.Message}");
    }
    catch (Exception ex)
    {
      return ServiceResult2<bool>.Fail($"Unexpected error: {ex.Message}");
    }
  }

  public async Task<ServiceResult2<bool>> DeleteOptionValue(Guid valueId)
  {
    try
    {
      var response = await _http.DeleteAsync($"products/options/values/{valueId}/delete");
      if (response.IsSuccessStatusCode)
      {
        return ServiceResult2<bool>.Ok(true);
      }

      var error = await ParseError(response);
      return ServiceResult2<bool>.Fail(error);
    }
    catch (HttpRequestException ex)
    {
      return ServiceResult2<bool>.Fail($"Network error: {ex.Message}");
    }
    catch (Exception ex)
    {
      return ServiceResult2<bool>.Fail($"Unexpected error: {ex.Message}");
    }
  }

  public async Task<ServiceResult2<bool>> DeleteVariant(Guid variantId)
  {
    try
    {
      var response = await _http.DeleteAsync($"products/variants/{variantId}/delete");
      if (response.IsSuccessStatusCode)
      {
        return ServiceResult2<bool>.Ok(true);
      }

      var error = await ParseError(response);
      return ServiceResult2<bool>.Fail(error);
    }
    catch (HttpRequestException ex)
    {
      return ServiceResult2<bool>.Fail($"Network error: {ex.Message}");
    }
    catch (Exception ex)
    {
      return ServiceResult2<bool>.Fail($"Unexpected error: {ex.Message}");
    }
  }

  public async Task<ServiceResult2<ApiMessageResponse>> UpdateProductInventory(UpdateProductInventoryRequest request)
  {
    try
    {
      var response = await _http.PatchAsJsonAsync("products/inventory", request);

      if (response.IsSuccessStatusCode)
      {
        var result = await response.Content.ReadFromJsonAsync<ApiMessageResponse>();
        return ServiceResult2<ApiMessageResponse>.Ok(result ?? new ApiMessageResponse());
      }

      var error = await ParseError(response);
      return ServiceResult2<ApiMessageResponse>.Fail(error);
    }
    catch (HttpRequestException ex)
    {
      return ServiceResult2<ApiMessageResponse>.Fail($"Network error: {ex.Message}");
    }
    catch (Exception ex)
    {
      return ServiceResult2<ApiMessageResponse>.Fail($"Unexpected error: {ex.Message}");
    }
  }

  public async Task<string> ParseError(HttpResponseMessage response)
  {
    try
    {
      var json = await response.Content.ReadAsStringAsync();
      var error = JsonSerializer.Deserialize<ErrorResponse>(json, new JsonSerializerOptions
      {
        PropertyNameCaseInsensitive = true
      });

      if (error == null)
        return $"HTTP {(int)response.StatusCode} - {response.ReasonPhrase}";

      if (error.Errors != null && error.Errors.Any())
      {
        // Flatten field-level errors: "PickupAddress: Required"
        var fieldErrors = error.Errors
          .SelectMany(kvp => kvp.Value.Select(v => $"{kvp.Key}: {v}"));
        return string.Join("; ", fieldErrors);
      }

      // Fall back to detail or generic error
      return !string.IsNullOrWhiteSpace(error.Detail)
        ? error.Detail
        : error.Error ?? $"HTTP {(int)response.StatusCode} - {response.ReasonPhrase}";
    }
    catch
    {
      return $"HTTP {(int)response.StatusCode} - {response.ReasonPhrase}";
    }
  }

  private async Task<ServiceResult2<T>> GetFromJsonSafe<T>(string url)
  {
    try
    {
      var response = await _http.GetAsync(url);

      if (response.IsSuccessStatusCode)
      {
        var result = await response.Content.ReadFromJsonAsync<T>();
        return ServiceResult2<T>.Ok(result ?? Activator.CreateInstance<T>());
      }

      var error = await ParseError(response);
      return ServiceResult2<T>.Fail(error);
    }
    catch (HttpRequestException ex)
    {
      return ServiceResult2<T>.Fail($"Network error: {ex.Message}");
    }
    catch (Exception ex)
    {
      return ServiceResult2<T>.Fail($"Unexpected error: {ex.Message}");
    }
  }
}