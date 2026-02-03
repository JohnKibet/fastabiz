using System.Net.Http.Json;
using System.Text.Json;
using frontend.Models;
public class OrderService
{
    private readonly HttpClient _http;
    private readonly ToastService _toastService;
    private List<Order>? _cachedOrders;
    private DateTime _lastFetchTime;
    private readonly TimeSpan _cacheDuration = TimeSpan.FromMinutes(5);
    public OrderService(IHttpClientFactory httpClientFactory, ToastService toastService)
    {
        _http = httpClientFactory.CreateClient("AuthenticatedApi");
        _toastService = toastService;
    }

    public async Task<ServiceResult2<HttpResponseMessage>> AddOrder(CreateOrderRequest order)
    {
        try
        {
            var response = await _http.PostAsJsonAsync("orders/create", order);
            if (response.IsSuccessStatusCode)
            {
                InvalidateCache();
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

    public async Task<ServiceResult2<List<Guid>>> CreatePendingOrder(CreateOrderRequest order)
    {
        try
        {
            var response = await _http.PostAsJsonAsync("orders/pending", order);
            if (!response.IsSuccessStatusCode)
                return ServiceResult2<List<Guid>>.Fail(await ParseError(response));

            var ids = await response.Content.ReadFromJsonAsync<List<Guid>>();
            return ServiceResult2<List<Guid>>.Ok(ids!);
        }
        catch (Exception ex)
        {
            return ServiceResult2<List<Guid>>.Fail(ex.Message);
        }
    }


    public async Task<ServiceResult2<List<Order>>> GetOrderByID(Guid id)
    {
        return await GetFromJsonSafe<List<Order>>($"orders/by-id/{id}");
    }


    public async Task<ServiceResult2<List<Order>>> GetOrdersByCustomer(Guid customerId)
    {
        return await GetFromJsonSafe<List<Order>>($"orders/by-customer/{customerId}");
    }

    public async Task<Order?> UpdateOrder(Guid orderId, string column, object value)
    {
        var requestBody = new
        {
            column,
            value
        };

        var response = await _http.PutAsJsonAsync($"orders/{orderId}/update", requestBody);
        if (response.IsSuccessStatusCode)
        {
            InvalidateCache();
            return await response.Content.ReadFromJsonAsync<Order>() ?? new Order();
        }

        return null;

    }

    public async Task<ServiceResult2<List<Order>>> GetAllOrders()
    {
        return await GetFromJsonSafe<List<Order>>("orders/all_orders");
    }

    // cache orders
    public async Task<ServiceResult2<List<Order>>> GetAllCachedOrders(bool forceRefresh = false)
    {
        if (!forceRefresh && _cachedOrders != null && DateTime.UtcNow - _lastFetchTime < _cacheDuration)
        {
            return ServiceResult2<List<Order>>.Ok(_cachedOrders, fromCache: true);
        }

        var result = await GetAllOrders();
        if (result.Success && result.Data != null && result.Data.Any())
        {
            _cachedOrders = result.Data;
            _lastFetchTime = DateTime.UtcNow;

            _toastService.ShowToast("Orders fetched successfully.", string.Empty, ToastService.ToastLevel.Success);
        }
        else if (result.Success && (result.Data == null || !result.Data.Any()))
        {
            _toastService.ShowToast("No orders.", string.Empty, ToastService.ToastLevel.Warning);
        }
        else
        {
            _toastService.ShowToast("Failed to load orders.", string.Empty, ToastService.ToastLevel.Error);
        }

        return result;
    }

    public void InvalidateCache()
    {
        _cachedOrders = null;
    }

    public async Task<bool> DeleteOrder(Guid id)
    {
        var res = await _http.DeleteAsync($"orders/{id}");
        if (res.IsSuccessStatusCode)
        {
            InvalidateCache();
        }
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

