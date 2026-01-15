using System;
using System.Text.Json;
using System.Net.Http.Json;
using frontend.Models;
public class StoreService
{
    private readonly HttpClient _http;
    private readonly ToastService _toastService;
    public StoreService(IHttpClientFactory httpClientFactory, ToastService toastService)
    {
        _http = httpClientFactory.CreateClient("AuthenticatedApi");
        _toastService = toastService;
    }

    public async Task<ServiceResult2<HttpResponseMessage>> CreateStore(CreateStoreRequest store)
    {
        try
        {
            var response = await _http.PostAsJsonAsync("stores/create", store);
            if (response.IsSuccessStatusCode)
            {
                return ServiceResult2<HttpResponseMessage>.Ok(response);
            }

            var error = await ParseError(response);
            return ServiceResult2<HttpResponseMessage>.Fail(error);
        }
        catch (HttpRequestException ex)
        {
            return ServiceResult2<HttpResponseMessage>.Fail($"Network error: {ex.Message}");
        }
        catch (Exception ex)
        {
            return ServiceResult2<HttpResponseMessage>.Fail($"Unexpected error: {ex.Message}");
        }
    }

    public async Task<ServiceResult2<List<Store>>> GetAllStores()
    {
        return await GetFromJsonSafe<List<Store>>("stores/all_stores");
    }

    public async Task<ServiceResult2<List<Store>>> GetStoresByOwner()
    {
        return await GetFromJsonSafe<List<Store>>("stores/me");
    }

    public async Task<ServiceResult2<List<Store>>> ListStoresPaginated()
    {
        return await GetFromJsonSafe<List<Store>>("stores/me/paged");
    }

    public async Task<ServiceResult2<Store>> GetStoreById(Guid storeId)
    {
        return await GetFromJsonSafe<Store>($"stores/by-id/{storeId}");
    }

    public async Task<ServiceResult2<Store>> GetStoreSummary(Guid storeId)
    {
        return await GetFromJsonSafe<Store>($"stores/{storeId}/summary");
    }

    public async Task<ServiceResult2<List<ProductX>>> UpdateStore(Guid StoreId, UpdateStoreRequest updateRequest)
    {
        try
        {
            var response = await _http.PutAsJsonAsync($"stores/{StoreId}/update", updateRequest);
            if (response.IsSuccessStatusCode)
            {
                var updatedStore = await response.Content.ReadFromJsonAsync<Store>();
                return ServiceResult2<List<ProductX>>.Ok(updatedStore?.Products ?? new List<ProductX>());
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

    public async Task<ServiceResult2<HttpResponseMessage>> DeleteStore(Guid storeId)
    {
        try
        {
            var response = await _http.DeleteAsync($"stores/{storeId}/delete");
            if (response.IsSuccessStatusCode)
            {
                return ServiceResult2<HttpResponseMessage>.Ok(response);
            }

            var error = await ParseError(response);
            return ServiceResult2<HttpResponseMessage>.Fail(error);
        }
        catch (HttpRequestException ex)
        {
            return ServiceResult2<HttpResponseMessage>.Fail($"Network error: {ex.Message}");
        }
        catch (Exception ex)
        {
            return ServiceResult2<HttpResponseMessage>.Fail($"Unexpected error: {ex.Message}");
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

            return error?.Detail ?? "Unknown error occurred.";
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