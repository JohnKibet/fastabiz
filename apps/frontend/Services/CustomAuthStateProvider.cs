using System.Security.Claims;
using System.Text.Json;
using Microsoft.AspNetCore.Components;
using Microsoft.AspNetCore.Components.Authorization;
using Microsoft.JSInterop;
using frontend.Models;

namespace frontend.Services;

public class CustomAuthStateProvider : AuthenticationStateProvider
{
    private readonly IJSRuntime _jSRuntime;
    private readonly NavigationManager _navigationManager;
    private const string TokenStorageKey = "auth_token";
    private User? _cachedUser = null;
    public CustomAuthStateProvider(IJSRuntime jSRuntime, NavigationManager navigationManager)
    {
        _jSRuntime = jSRuntime;
        _navigationManager = navigationManager;
    }

    public override async Task<AuthenticationState> GetAuthenticationStateAsync()
    {
        try
        {
            if (_cachedUser == null)
            {
                var token = await _jSRuntime.InvokeAsync<string>("localStorage.getItem", TokenStorageKey);

                if (string.IsNullOrWhiteSpace(token))
                    return new AuthenticationState(new ClaimsPrincipal(new ClaimsIdentity()));

                _cachedUser = ParseToken(token);
            }

            if (_cachedUser == null)
                return new AuthenticationState(new ClaimsPrincipal(new ClaimsIdentity()));

            var identity = await CreateIdentity(_cachedUser);

            if (identity == null)
            {
                return new AuthenticationState(
                    new ClaimsPrincipal(new ClaimsIdentity())
                );
            }

            return new AuthenticationState(new ClaimsPrincipal(identity));
        }
        catch
        {
            return new AuthenticationState(new ClaimsPrincipal(new ClaimsIdentity()));
        }
    }

    public async Task SignInAsync(string token)
    {
        await _jSRuntime.InvokeVoidAsync("localStorage.setItem", TokenStorageKey, token);
        _cachedUser = ParseToken(token);
        NotifyAuthenticationStateChanged(GetAuthenticationStateAsync());
    }

    public async Task SignOutAsync()
    {
        await _jSRuntime.InvokeVoidAsync("localStorage.removeItem", TokenStorageKey);
        await _jSRuntime.InvokeVoidAsync("localStorage.removeItem", "cart");
        _cachedUser = null;

        NotifyAuthenticationStateChanged(GetAuthenticationStateAsync());
        await Task.Delay(50);

        _navigationManager.NavigateTo("/auth/login", forceLoad: true);
    }

    private async Task<ClaimsIdentity?> CreateIdentity(User user)
    {
        if (_cachedUser?.ExpiresAt != null && _cachedUser.ExpiresAt <= DateTime.UtcNow)
        {
            await ForceLogoutAsync();
            return null;
        }

        var claims = new List<Claim>
        {
            new Claim(ClaimTypes.NameIdentifier, user.ID.ToString()),
            new Claim(ClaimTypes.Name, user.FullName),
            new Claim(ClaimTypes.Email, user.Email),
            new Claim(ClaimTypes.Role, user.Role),
        };

        // Phone (custom claim)
        if (!string.IsNullOrWhiteSpace(user.Phone))
            claims.Add(new Claim("phone", user.Phone));

        // Slug (custom claim)
        if (!string.IsNullOrWhiteSpace(user.Slug))
            claims.Add(new Claim("slug", user.Slug));

        // Status (enum, nullable)
        if (user.Status != null)
            claims.Add(new Claim("status", user.Status.Value.ToString()));

        // Last Login (nullable datetime)
        if (user.LastLogin != null)
            claims.Add(new Claim("last_login", user.LastLogin.Value.ToString("o")));
        // ISO-8601 format

        // Created At (nullable datetime)
        if (user.CreatedAt != null)
            claims.Add(new Claim("created_at", user.CreatedAt.Value.ToString("o")));

        return new ClaimsIdentity(claims, "jwtAuth");
    }

    public User? ParseToken(string jwt)
    {
        try
        {
            var parts = jwt.Split('.');
            if (parts.Length != 3) return null;

            var payload = parts[1];
            var jsonBytes = Convert.FromBase64String(PadBase64(payload));
            var json = System.Text.Encoding.UTF8.GetString(jsonBytes);
            var claimsDict = JsonSerializer.Deserialize<Dictionary<string, object>>(json);

            var user = new User();

            // ID
            if (claimsDict.TryGetValue("sub", out var idVal))
                user.ID = Guid.TryParse(idVal?.ToString(), out var id) ? id : Guid.Empty;

            // Email
            if (claimsDict.TryGetValue("email", out var emailVal))
                user.Email = emailVal?.ToString() ?? "";

            // Name
            if (claimsDict.TryGetValue("name", out var nameVal))
                user.FullName = nameVal?.ToString() ?? "";

            // Role
            if (claimsDict.TryGetValue("role", out var roleVal))
                user.Role = roleVal?.ToString() ?? "";

            // Phone (nullable)
            if (claimsDict.TryGetValue("phone", out var phoneVal))
                user.Phone = phoneVal?.ToString() ?? "";

            // Slug (nullable)
            if (claimsDict.TryGetValue("slug", out var slugVal))
                user.Slug = slugVal?.ToString() ?? "";

            // Status (enum)
            if (claimsDict.TryGetValue("status", out var statusVal) && statusVal != null)
            {
                try
                {
                    user.Status = Enum.Parse<UserStatus>(statusVal.ToString(), ignoreCase: true);
                }
                catch
                {
                    user.Status = null; // invalid or missing
                }
            }

            // lastLogin (Unix timestamp or null)
            if (claimsDict.TryGetValue("last_login", out var lastLoginVal) &&
                lastLoginVal != null &&
                long.TryParse(lastLoginVal.ToString(), out var unix))
            {
                user.LastLogin = DateTimeOffset.FromUnixTimeSeconds(unix).UtcDateTime;
            }
            else
            {
                user.LastLogin = null;
            }

            // created_at (Unix timestamp)
            if (claimsDict.TryGetValue("created_at", out var createdVal) &&
                createdVal != null &&
                long.TryParse(createdVal.ToString(), out var createdUnix))
            {
                user.CreatedAt = DateTimeOffset.FromUnixTimeSeconds(createdUnix).UtcDateTime;
            }
            else
            {
                user.CreatedAt = null;
            }

            // exp (Unix timestamp)
            if (claimsDict.TryGetValue("exp", out var expVal) &&
                expVal != null &&
                long.TryParse(expVal.ToString(), out var expUnix))
            {
                user.ExpiresAt = DateTimeOffset.FromUnixTimeSeconds(expUnix).UtcDateTime;
            }
            else
            {
                user.ExpiresAt = null;
            }

            return user;
        }
        catch
        {
            return null;
        }
    }

    private string PadBase64(string base64)
    {
        return base64.PadRight(base64.Length + (4 - base64.Length % 4) % 4, '=');
    }

    private async Task ForceLogoutAsync()
    {
        await _jSRuntime.InvokeVoidAsync("localStorage.removeItem", TokenStorageKey);
        await _jSRuntime.InvokeVoidAsync("localStorage.removeItem", "cart");
        _cachedUser = null;

        NotifyAuthenticationStateChanged(GetAuthenticationStateAsync());

        _navigationManager.NavigateTo("/auth/login", forceLoad: true);
    }
}

public class AuthExpiredHandler : DelegatingHandler
{
    private readonly CustomAuthStateProvider _auth;
    public AuthExpiredHandler(CustomAuthStateProvider auth)
    {
        _auth = auth;
    }

    protected override async Task<HttpResponseMessage>SendAsync(
        HttpRequestMessage request, 
        CancellationToken cancellationToken)
    {
        var response = await base.SendAsync(request, cancellationToken);

        if (response.StatusCode == System.Net.HttpStatusCode.Unauthorized)
        {
            await _auth.SignOutAsync();
        }

        return response;
    }
}

// proactive logout with a timer
// private Timer? _expiryTimer;

// private void ScheduleExpiry(User user)
// {
//     if (user.ExpiresAt == null) return;

//     var delay = user.ExpiresAt.Value - DateTime.UtcNow;
//     if (delay <= TimeSpan.Zero) return;

//     _expiryTimer?.Dispose();
//     _expiryTimer = new Timer(async _ =>
//     {
//         await ForceLogoutAsync();
//     }, null, delay, Timeout.InfiniteTimeSpan);
// }
// Call it in:
// SignInAsync
// GetAuthenticationStateAsync
