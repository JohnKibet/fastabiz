using System.Net.Http.Json;
using System.Text.Json;
using frontend.Models;


public class DeliveryService
{
    private readonly HttpClient _http;
    public DeliveryService(IHttpClientFactory httpClientFactory)
    {
        _http = httpClientFactory.CreateClient("AuthenticatedApi");
    }

    public async Task<ServiceResult2<HttpResponseMessage>> CreateDelivery(CreateDelivery delivery)
    {
        try
        {
            var response = await _http.PostAsJsonAsync("deliveries/create", delivery);
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

    public async Task<ServiceResult2<List<Delivery>>> GetDeliveryById(Guid Id)
    {
        return await GetFromJsonSafe<List<Delivery>>($"deliveries/by-id/{Id}"); // backend handler changes
    }
    public async Task<ServiceResult2<List<Delivery>>> GetDeliveries()
    {
        return await GetFromJsonSafe<List<Delivery>>("deliveries/all_deliveries");
    }

    public async Task<bool> DeleteDelivery(Guid id)
    {
        var res = await _http.DeleteAsync($"deliveries/{id}");
        
        return res.IsSuccessStatusCode;
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
