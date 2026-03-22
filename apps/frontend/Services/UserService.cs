using frontend.Models;
using frontend.Models.Storefront;

namespace frontend.Services;

public class UserService
{
    private readonly ApiService   _api;
    private readonly ToastService _toast;

    // Cache: one instance per service lifetime (Scoped), 5-minute TTL
    private readonly CacheEntry<List<User>> _cache = new(TimeSpan.FromMinutes(5));

    public UserService(ApiService api, ToastService toast)
    {
        _api   = api;
        _toast = toast;
    }

    // ── CRUD ──────────────────────────────────────────────────────────

    public Task<ApiResult<CreateUserResponse>> AddUser(CreateUserRequest req)
        => _api.PostAsync<CreateUserRequest, CreateUserResponse>("public/create", req);

    public Task<ApiResult<User>> GetUserByID(Guid id)
        => _api.GetAsync<User>($"users/by-id/{id}");

    public Task<ApiResult<List<User>>> GetAllUsers()
        => _api.GetAsync<List<User>>("users/all_users");

    public Task<ApiResult<ApiMessageResponse>> UpdateUser(Guid userId, UpdateOrderRequest req)
        => _api.PutAsync<UpdateOrderRequest, ApiMessageResponse>($"users/{userId}/update", req);

    public Task<ApiResult<ApiMessageResponse>> ChangePassword(Guid userId, ChangePasswordDto req)
        => _api.PutAsync<ChangePasswordDto, ApiMessageResponse>($"users/{userId}/password", req);

    public Task<ApiResult<ApiMessageResponse>> UpdateProfile(Guid userId, EditUserModel model)
        => _api.PutAsync<EditUserModel, ApiMessageResponse>($"users/{userId}/profile", model);

    public Task<ApiResult<ApiMessageResponse>> UpdateStatus(Guid userId, UserStatus newStatus)
        => _api.PutAsync<UserStatus, ApiMessageResponse>($"users/{userId}/status", newStatus);

    public Task<ApiResult<bool>> DeleteUser(Guid userId)
        => _api.DeleteAsync($"users/{userId}");

    // ── Cached list ───────────────────────────────────────────────────

    /// <summary>
    /// Returns the full user list, serving from cache when still valid.
    ///
    /// Toast policy:
    ///   • Cache hit  → silent (no toast — user doesn't need to know)
    ///   • Empty list → warning toast (actionable: go create a user)
    ///   • Error      → error toast
    ///   • Success    → silent (the data appearing on screen is feedback enough)
    /// </summary>
    public async Task<ApiResult<List<User>>> GetAllCachedUsers(bool forceRefresh = false)
    {
        if (!forceRefresh && _cache.TryGet(out var cached))
            return ApiResult<List<User>>.Ok(cached!, fromCache: true);

        var result = await GetAllUsers();

        if (result.Success && result.Data != null)
        {
            _cache.Set(result.Data);

            if (!result.Data.Any())
                _toast.ShowToast("No users found.", string.Empty, ToastService.ToastLevel.Warning);
        }
        else
        {
            _toast.ShowToast(
                result.Error?.Message ?? "Failed to load users.",
                string.Empty,
                ToastService.ToastLevel.Error);
        }

        return result;
    }

    /// <summary>Forces the next GetAllCachedUsers call to hit the network.</summary>
    public void InvalidateCache() => _cache.Invalidate();
}
