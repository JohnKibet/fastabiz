using frontend.Models;
using System.Net.Http.Json;
using System.Text.Json;
public class UserService
{
    private readonly HttpClient _http;
    private readonly DropdownDataService _dropdownService;
    private readonly ToastService _toastService;
    private List<User>? _cachedUsers;
    private DateTime _lastFetchTime;
    private readonly TimeSpan _cacheDuration = TimeSpan.FromMinutes(5);
    public UserService(IHttpClientFactory httpClientFactory, DropdownDataService dropdownService, ToastService toastService)
    {
        _http = httpClientFactory.CreateClient("AuthenticatedApi");;
        _dropdownService = dropdownService;
        _toastService = toastService;
    }

    public async Task<ServiceResult2<HttpResponseMessage>> AddUser(CreateUserRequest user)
    {
        try
        {
            var response = await _http.PostAsJsonAsync("public/create", user);
            if (response.IsSuccessStatusCode)
            {
                InvalidateCache();
                _dropdownService.InvalidateCache();
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

    public async Task<User> GetUserByID(Guid id)
    {
        var user = await _http.GetFromJsonAsync<User>($"users/by-id/{id}");
        return user ?? new User();
    }
    public async Task<ServiceResult2<List<User>>> GetAllUsers()
    {
        return await GetFromJsonSafe<List<User>>("users/all_users");
    }

    public async Task<ServiceResult2<HttpResponseMessage>> UpdateUser(Guid userId, string column, object value)
    {
        try
        {
            var requestBody = new
            {
                column,
                value
            };

            var response = await _http.PutAsJsonAsync($"users/{userId}/update", requestBody);
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

    public async Task<ServiceResult2<HttpResponseMessage>> ChangePassword(Guid userId, string currentPassword, string newPassword)
    {
        var body = new ChangePasswordDto
        {
            CurrentPassword = currentPassword,
            NewPassword = newPassword
        };

        var response = await _http.PutAsJsonAsync($"users/{userId}/password", body);

        if (response.IsSuccessStatusCode)
            return ServiceResult2<HttpResponseMessage>.Ok(response);

        var error = await ParseError(response);
        return ServiceResult2<HttpResponseMessage>.Fail(error);
    }
    public async Task<ServiceResult2<HttpResponseMessage>> UpdateProfile(Guid userId, EditUserModel model)
    {
        var response = await _http.PatchAsJsonAsync($"users/{userId}/profile", model);

        if (response.IsSuccessStatusCode)
            return ServiceResult2<HttpResponseMessage>.Ok(response);

        // Parse structured backend error (ErrorResponse)
        var error = await ParseError(response);
        return ServiceResult2<HttpResponseMessage>.Fail(error);
    }

    public async Task<ServiceResult2<List<User>>> GetAllCachedUsers(bool forceRefresh = false)
    {
        if (!forceRefresh && _cachedUsers != null && DateTime.UtcNow - _lastFetchTime < _cacheDuration)
        {
            return ServiceResult2<List<User>>.Ok(_cachedUsers, fromCache: true);
        }

        var result = await GetAllUsers(); 
        if (result.Success && result.Data != null && result.Data.Any())
        {
            _cachedUsers = result.Data;
            _lastFetchTime = DateTime.UtcNow;

            _toastService.ShowToast("Users fetched successfully.", ToastService.ToastLevel.Success);
        }
        else if (result.Success && (result.Data == null || !result.Data.Any()))
        {
            _toastService.ShowToast("No users.", ToastService.ToastLevel.Warning);
        }
        else
        {
            _toastService.ShowToast("Failed to load users.", ToastService.ToastLevel.Error);
        }

        return result;
    }

    public void InvalidateCache()
    {
        _cachedUsers = null;
    }

    public async Task<bool> DeleteUser(Guid id)
    {
        var res = await _http.DeleteAsync($"users/{id}");
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
