using System.Net.Http.Json;
using System.Text.Json;
using frontend.Models;
using frontend.Models.Storefront;

namespace frontend.Services;

public class DeliveryService
{
    private readonly ApiService _api;
    private readonly ToastService _toast;

    public DeliveryService(ApiService api, ToastService toast)
    {
        _api = api;
        _toast = toast;
    }

    public Task<ApiResult<CreateDeliveryResponse>> CreateDelivery(CreateDeliveryRequest req) => _api.PostAsync<CreateDeliveryRequest, CreateDeliveryResponse>("deliveries/create", req);

    public Task<ApiResult<Delivery>> GetDeliveryById(Guid Id) => _api.GetAsync<Delivery>($"deliveries/by-id/{Id}");

    public Task<ApiResult<List<Delivery>>> GetDeliveries() => _api.GetAsync<List<Delivery>>("deliveries/all_deliveries");

    public Task<ApiResult<ApiMessageResponse>> UpdateDeliveryStatus(Guid deliveryId, DeliveryStatus newStatus) => _api.PutAsync<DeliveryStatus, ApiMessageResponse>($"deliveries/{deliveryId}/update", newStatus);

    public Task<ApiResult<bool>> DeleteDelivery(Guid storeId) => _api.DeleteAsync($"deliveries/{storeId}");
}
