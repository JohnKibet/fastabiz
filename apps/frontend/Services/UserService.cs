using frontend.Models;
using frontend.Models.Storefront;
using System.Net.Http.Json;
using System.Text.Json;

namespace frontend.Services;

public class UserService
{
    private readonly ApiService _api;
    private readonly ToastService _toast;

    private List<User>? _cachedUsers;
    private DateTime _lastFetchTime;
    private readonly TimeSpan _cacheDuration = TimeSpan.FromMinutes(5);

    public UserService(ApiService api, ToastService toast)
    {
        _api = api;
        _toast = toast;
    }

    public Task<ApiResult<CreateUserResponse>> AddUser(CreateUserRequest req) => _api.PostAsync<CreateUserRequest, CreateUserResponse>("public/create", req);

    public Task<ApiResult<User>> GetUserByID(Guid id) => _api.GetAsync<User>($"users/by-id/{id}");

    public Task<ApiResult<List<User>>> GetAllUsers() => _api.GetAsync<List<User>>("users/all_users");

    // temp use of req order-type
    public Task<ApiResult<ApiMessageResponse>> UpdateUser(Guid userId, UpdateOrderRequest req)

        => _api.PutAsync<UpdateOrderRequest, ApiMessageResponse>($"users/{userId}/update", req);


    public Task<ApiResult<ApiMessageResponse>> ChangePassword(Guid userId, ChangePasswordDto req)

        => _api.PutAsync<ChangePasswordDto, ApiMessageResponse>($"users/{userId}/password", req);


    public Task<ApiResult<ApiMessageResponse>> UpdateProfile(Guid userId, EditUserModel model)
        => _api.PutAsync<EditUserModel, ApiMessageResponse>($"users/{userId}/profile", model);

    public Task<ApiResult<ApiMessageResponse>> UpdateStatus(Guid userId, UserStatus newStatus)
        => _api.PutAsync<UserStatus, ApiMessageResponse>($"users/{userId}/status", newStatus);

    // public async Task<ServiceResult2<List<User>>> GetAllCachedUsers(bool forceRefresh = false)
    // {
    //     if (!forceRefresh && _cachedUsers != null && DateTime.UtcNow - _lastFetchTime < _cacheDuration)
    //     {
    //         return ServiceResult2<List<User>>.Ok(_cachedUsers, fromCache: true);
    //     }

    //     var result = await GetAllUsers(); 
    //     if (result.Success && result.Data != null && result.Data.Any())
    //     {
    //         _cachedUsers = result.Data;
    //         _lastFetchTime = DateTime.UtcNow;

    //         _toastService.ShowToast("Users fetched successfully.", string.Empty, ToastService.ToastLevel.Success);
    //     }
    //     else if (result.Success && (result.Data == null || !result.Data.Any()))
    //     {
    //         _toastService.ShowToast("No users.", string.Empty, ToastService.ToastLevel.Warning);
    //     }
    //     else
    //     {
    //         _toastService.ShowToast("Failed to load users.", string.Empty, ToastService.ToastLevel.Error);
    //     }

    //     return result;
    // }

    public void InvalidateCache()
    {
        _cachedUsers = null;
    }

    public Task<ApiResult<bool>> DeleteUser(Guid userId) => _api.DeleteAsync($"users/{userId}");
}
