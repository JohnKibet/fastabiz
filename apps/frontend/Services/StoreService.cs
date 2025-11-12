using System.Net.Http.Json;
using frontend.Models;
using System.Text.Json;

public class StoreService
{
    private readonly HttpClient _http;

    public StoreService(HttpClient http)
    {
        _http = http;
    }

    public async Task<ServiceResult<HttpResponseMessage>> CreateStore(CreateStoreRequest store)
    {
        try
        {
            var response = await _http.PostAsJsonAsync("stores/create", store);
            if (response.IsSuccessStatusCode)
            {
                return ServiceResult<HttpResponseMessage>.Ok(response);
            }

            var error = await ParseError(response);
            return ServiceResult<HttpResponseMessage>.Fail(error);
        }
        catch (HttpRequestException ex)
        {
            return ServiceResult<HttpResponseMessage>.Fail($"Network error: {ex.Message}");
        }
        catch (Exception ex)
        {
            return ServiceResult<HttpResponseMessage>.Fail($"Unexpected error: {ex.Message}");
        }
    }

    public async Task<ServiceResult<Store>> GetStoreByID(Guid id)
    {
        return await GetFromJsonSafe<Store>($"stores/by-id/{id}");
    }

    public async Task<Store?> GetStoreBySlug(string slug)
    {
        return await _http.GetFromJsonAsync<Store>($"stores/by-slug?slug={slug}");
    }

    public async Task<Store?> GetStoreByOwner(Guid ownerId)
    {
        return await _http.GetFromJsonAsync<Store>($"stores/owner/{ownerId}");
    }

    public async Task<Store?> UpdateStore(Guid storeId, string column, object value)
    {
        var requestBody = new
        {
            column,
            value
        };

        var response = await _http.PutAsJsonAsync($"stores/{storeId}/update", requestBody);
        if (response.IsSuccessStatusCode)
        {
            return await response.Content.ReadFromJsonAsync<Store>() ?? new Store();
        }

        return null;
    }

    public async Task<ServiceResult<List<Store>>> GetAllPublicStores()
    {
        return await GetFromJsonSafe<List<Store>>("stores/public");
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
    
    private async Task<ServiceResult<T>> GetFromJsonSafe<T>(string url)
    {
        try
        {
            var response = await _http.GetAsync(url);

            if (response.IsSuccessStatusCode)
            {
                var result = await response.Content.ReadFromJsonAsync<T>();
                return ServiceResult<T>.Ok(result ?? Activator.CreateInstance<T>());
            }

            var error = await ParseError(response);
            return ServiceResult<T>.Fail(error);
        }
        catch (HttpRequestException ex)
        {
            return ServiceResult<T>.Fail($"Network error: {ex.Message}");
        }
        catch (Exception ex)
        {
            return ServiceResult<T>.Fail($"Unexpected error: {ex.Message}");
        }
    }
}
