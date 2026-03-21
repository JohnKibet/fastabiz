using System.Net.Http.Json;
using System.Text.Json;
using frontend.Models;
using frontend.Models.Storefront;

namespace frontend.Services;

public class OrderService
{
    private readonly ApiService _api;
    private readonly HttpClient _cloudinaryHttp;
    private readonly ToastService _toast;

    private List<Order>? _cachedOrders;
    private DateTime _lastFetchTime;
    private readonly TimeSpan _cacheDuration = TimeSpan.FromMinutes(5);

    public OrderService(ApiService api, IHttpClientFactory factory, ToastService toast)
    {
        _api = api;
        _cloudinaryHttp = factory.CreateClient("CloudinaryClient");
        _toast = toast;
    }

    public Task<ApiResult<CreateOrderResponse>> AddOrder(CreateOrderRequest req) => _api.PostAsync<CreateOrderRequest, CreateOrderResponse>("orders/create", req);

    public Task<ApiResult<CreateOrderResponse>> CreatePendingOrder(CreateOrderRequest req) => _api.PostAsync<CreateOrderRequest, CreateOrderResponse>("orders/pending", req);

    public Task<ApiResult<Order>> GetOrderByID(Guid id) => _api.GetAsync<Order>($"orders/by-id/{id}");

    public Task<ApiResult<List<Order>>> GetOrdersByCustomer(Guid customerId) => _api.GetAsync<List<Order>>($"orders/by-customer/{customerId}");

    public Task<ApiResult<ApiMessageResponse>> UpdateOrder(Guid orderId, UpdateOrderRequest req)

        => _api.PutAsync<UpdateOrderRequest, ApiMessageResponse>($"orders/{orderId}/update", req);

    public Task<ApiResult<List<Order>>> GetAllOrders() => _api.GetAsync<List<Order>>("orders/all_orders");

    // cache orders
    // public async Task<ServiceResult2<List<Order>>> GetAllCachedOrders(bool forceRefresh = false)
    // {
    //     if (!forceRefresh && _cachedOrders != null && DateTime.UtcNow - _lastFetchTime < _cacheDuration)
    //     {
    //         return ServiceResult2<List<Order>>.Ok(_cachedOrders, fromCache: true);
    //     }

    //     var result = await GetAllOrders();
    //     if (result.Success && result.Data != null && result.Data.Any())
    //     {
    //         _cachedOrders = result.Data;
    //         _lastFetchTime = DateTime.UtcNow;

    //         _toastService.ShowToast("Orders fetched successfully.", string.Empty, ToastService.ToastLevel.Success);
    //     }
    //     else if (result.Success && (result.Data == null || !result.Data.Any()))
    //     {
    //         _toastService.ShowToast("No orders.", string.Empty, ToastService.ToastLevel.Warning);
    //     }
    //     else
    //     {
    //         _toastService.ShowToast("Failed to load orders.", string.Empty, ToastService.ToastLevel.Error);
    //     }

    //     return result;
    // }

    public void InvalidateCache()
    {
        _cachedOrders = null;
    }

    public Task<ApiResult<bool>> DeleteOrder(Guid orderId) => _api.DeleteAsync($"orders/{orderId}");
}

