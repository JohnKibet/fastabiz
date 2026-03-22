using frontend.Models;
using frontend.Models.Storefront;

namespace frontend.Services;

public class OrderService
{
    private readonly ApiService   _api;
    private readonly ToastService _toast;

    private readonly CacheEntry<List<Order>> _cache = new(TimeSpan.FromMinutes(5));

    public OrderService(ApiService api, ToastService toast)
    {
        _api   = api;
        _toast = toast;
    }

    // ── CRUD ──────────────────────────────────────────────────────────

    public Task<ApiResult<CreateOrderResponse>> AddOrder(CreateOrderRequest req)
        => _api.PostAsync<CreateOrderRequest, CreateOrderResponse>("orders/create", req);

    public Task<ApiResult<CreateOrderResponse>> CreatePendingOrder(CreateOrderRequest req)
        => _api.PostAsync<CreateOrderRequest, CreateOrderResponse>("orders/pending", req);

    public Task<ApiResult<Order>> GetOrderByID(Guid id)
        => _api.GetAsync<Order>($"orders/by-id/{id}");

    public Task<ApiResult<List<Order>>> GetOrdersByCustomer(Guid customerId)
        => _api.GetAsync<List<Order>>($"orders/by-customer/{customerId}");

    public Task<ApiResult<ApiMessageResponse>> UpdateOrder(Guid orderId, UpdateOrderRequest req)
        => _api.PutAsync<UpdateOrderRequest, ApiMessageResponse>($"orders/{orderId}/update", req);

    public Task<ApiResult<List<Order>>> GetAllOrders()
        => _api.GetAsync<List<Order>>("orders/all_orders");

    public Task<ApiResult<bool>> DeleteOrder(Guid orderId)
        => _api.DeleteAsync($"orders/{orderId}");

    // ── Cached list ───────────────────────────────────────────────────

    /// <summary>
    /// Returns all orders, serving from cache when still valid.
    /// Silent on success and cache hits — only toasts on empty or error.
    /// </summary>
    public async Task<ApiResult<List<Order>>> GetAllCachedOrders(bool forceRefresh = false)
    {
        if (!forceRefresh && _cache.TryGet(out var cached))
            return ApiResult<List<Order>>.Ok(cached!, fromCache: true);

        var result = await GetAllOrders();

        if (result.Success && result.Data != null)
        {
            _cache.Set(result.Data);

            if (!result.Data.Any())
                _toast.ShowToast("No orders found.", string.Empty, ToastService.ToastLevel.Warning);
        }
        else
        {
            _toast.ShowToast(
                result.Error?.Message ?? "Failed to load orders.",
                string.Empty,
                ToastService.ToastLevel.Error);
        }

        return result;
    }

    /// <summary>Forces the next GetAllCachedOrders call to hit the network.</summary>
    public void InvalidateCache() => _cache.Invalidate();
}
