using System;
using System.Text.Json;
using System.Net.Http.Json;
using frontend.Models;
using frontend.Models.Storefront;
namespace frontend.Services;

public class StoreService
{
    private readonly ApiService _api;
    private readonly HttpClient _cloudinaryHttp;
    private readonly ToastService _toast;

    public StoreService(ApiService api, IHttpClientFactory factory, ToastService toast)
    {
        _api = api;
        _cloudinaryHttp = factory.CreateClient("CloudinaryClient");
        _toast = toast;
    }

    public Task<ApiResult<CreateStoreResponse>> CreateStore(CreateStoreRequest req)
        => _api.PostAsync<CreateStoreRequest, CreateStoreResponse>("stores/create", req);

    public Task<ApiResult<List<StoreDto>>> GetAllStores() => _api.GetAsync<List<StoreDto>>("stores/all_stores");

    public Task<ApiResult<List<ActiveStoreContext>>> GetStoresByOwner() => _api.GetAsync<List<ActiveStoreContext>>("stores/me");

    public Task<ApiResult<List<Store>>> ListStoresPaginated() => _api.GetAsync<List<Store>>("stores/me/paged");

    public Task<ApiResult<Store>> GetStoreById(Guid storeId) => _api.GetAsync<Store>($"stores/by-id/{storeId}");

    public Task<ApiResult<Store>> GetStoreSummary(Guid storeId) => _api.GetAsync<Store>($"stores/{storeId}/summary");

    public Task<ApiResult<ApiMessageResponse>> UpdateStore(Guid StoreId, UpdateStoreRequest req)
        => _api.PutAsync<UpdateStoreRequest, ApiMessageResponse>($"stores{StoreId}/update", req);

    public Task<ApiResult<bool>> DeleteStore(Guid storeId) => _api.DeleteAsync($"stores/{storeId}/delete");
}