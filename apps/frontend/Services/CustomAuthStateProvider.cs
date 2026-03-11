using System.Security.Claims;
using System.Text.Json;
using Microsoft.AspNetCore.Components;
using Microsoft.AspNetCore.Components.Authorization;
using Microsoft.JSInterop;
using frontend.Models;

namespace frontend.Services;

public class CustomAuthStateProvider : AuthenticationStateProvider
{
    private readonly IJSRuntime _jsRuntime;
    private readonly NavigationManager _navigationManager;
    private const string TokenStorageKey = "auth_token";
    private User? _cachedUser = null;
    private AuthenticationState? _cachedAuthState = null;
    private string? _cachedToken = null; // raw JWT — read by AuthHeaderHandler synchronously

    public CustomAuthStateProvider(IJSRuntime jsRuntime, NavigationManager navigationManager)
    {
        _jsRuntime = jsRuntime;
        _navigationManager = navigationManager;
    }

    /// Returns the raw JWT from memory. Used by AuthHeaderHandler to attach
    /// the Bearer token without hitting JS interop on every request.
    public string? GetCachedToken() => _cachedToken;

    public override async Task<AuthenticationState> GetAuthenticationStateAsync()
    {
        try
        {
            // Return cached state immediately — no JS interop, no async gap, no race condition.
            if (_cachedAuthState != null)
                return _cachedAuthState;

            if (_cachedUser == null)
            {
                var token = await _jsRuntime.InvokeAsync<string>("localStorage.getItem", TokenStorageKey);

                if (string.IsNullOrWhiteSpace(token))
                    return Unauthenticated();

                _cachedToken = token; // restore raw token from localStorage on cold load
                _cachedUser = ParseToken(token);
            }

            if (_cachedUser == null)
                return Unauthenticated();

            var identity = await CreateIdentity(_cachedUser);
            if (identity == null)
                return Unauthenticated();

            _cachedAuthState = new AuthenticationState(new ClaimsPrincipal(identity));
            return _cachedAuthState;
        }
        catch
        {
            return Unauthenticated();
        }
    }

    public async Task SignInAsync(string token)
    {
        // 1. Persist token to localStorage
        await _jsRuntime.InvokeVoidAsync("localStorage.setItem", TokenStorageKey, token);

        // 2. Cache the raw token — AuthHeaderHandler reads this synchronously
        _cachedToken = token;

        // 3. Parse and cache user
        _cachedUser = ParseToken(token);

        if (_cachedUser == null)
        {
            _cachedToken     = null;
            _cachedAuthState = null;
            NotifyAuthenticationStateChanged(Task.FromResult(Unauthenticated()));
            return;
        }

        // 4. Build and cache auth state
        var identity     = await CreateIdentity(_cachedUser);
        var principal    = new ClaimsPrincipal(identity ?? new ClaimsIdentity());
        _cachedAuthState = new AuthenticationState(principal);

        // 5. Notify — returns already-resolved Task, no async gap
        NotifyAuthenticationStateChanged(Task.FromResult(_cachedAuthState));
    }

    public async Task SignOutAsync()
    {
        await _jsRuntime.InvokeVoidAsync("localStorage.removeItem", TokenStorageKey);
        await _jsRuntime.InvokeVoidAsync("localStorage.removeItem", "cart");

        _cachedToken     = null;
        _cachedUser      = null;
        _cachedAuthState = null;

        NotifyAuthenticationStateChanged(Task.FromResult(Unauthenticated()));

        // Small delay gives AuthorizeView time to re-render before navigation
        await Task.Delay(50);
        _navigationManager.NavigateTo("/auth/login", forceLoad: true);
    }

    private static AuthenticationState Unauthenticated() =>
        new(new ClaimsPrincipal(new ClaimsIdentity()));

    private async Task<ClaimsIdentity?> CreateIdentity(User user)
    {
        if (user.ExpiresAt != null && user.ExpiresAt <= DateTime.UtcNow)
        {
            await ForceLogoutAsync();
            return null;
        }

        var claims = new List<Claim>
        {
            new(ClaimTypes.NameIdentifier, user.ID.ToString()),
            new(ClaimTypes.Name,           user.FullName),
            new(ClaimTypes.Email,          user.Email),
            new(ClaimTypes.Role,           user.Role),
        };

        if (!string.IsNullOrWhiteSpace(user.Phone))  claims.Add(new("phone",      user.Phone));
        if (!string.IsNullOrWhiteSpace(user.Slug))   claims.Add(new("slug",       user.Slug));
        if (user.Status    != null)                  claims.Add(new("status",     user.Status.Value.ToString()));
        if (user.LastLogin != null)                  claims.Add(new("last_login", user.LastLogin.Value.ToString("o")));
        if (user.CreatedAt != null)                  claims.Add(new("created_at", user.CreatedAt.Value.ToString("o")));

        return new ClaimsIdentity(claims, "jwtAuth");
    }

    private async Task ForceLogoutAsync()
    {
        await _jsRuntime.InvokeVoidAsync("localStorage.removeItem", TokenStorageKey);
        await _jsRuntime.InvokeVoidAsync("localStorage.removeItem", "cart");

        _cachedToken     = null;
        _cachedUser      = null;
        _cachedAuthState = null;

        NotifyAuthenticationStateChanged(Task.FromResult(Unauthenticated()));
        _navigationManager.NavigateTo("/auth/login", forceLoad: true);
    }

    public User? ParseToken(string jwt)
    {
        try
        {
            var parts = jwt.Split('.');
            if (parts.Length != 3) return null;

            var jsonBytes   = Convert.FromBase64String(PadBase64(parts[1]));
            var json        = System.Text.Encoding.UTF8.GetString(jsonBytes);
            var claimsDict  = JsonSerializer.Deserialize<Dictionary<string, object>>(json);
            if (claimsDict == null) return null;

            var user = new User();

            if (claimsDict.TryGetValue("sub",        out var idVal))    user.ID       = Guid.TryParse(idVal?.ToString(), out var id) ? id : Guid.Empty;
            if (claimsDict.TryGetValue("email",      out var emailVal)) user.Email    = emailVal?.ToString() ?? "";
            if (claimsDict.TryGetValue("name",       out var nameVal))  user.FullName = nameVal?.ToString()  ?? "";
            if (claimsDict.TryGetValue("role",       out var roleVal))  user.Role     = roleVal?.ToString()  ?? "";
            if (claimsDict.TryGetValue("phone",      out var phoneVal)) user.Phone    = phoneVal?.ToString() ?? "";
            if (claimsDict.TryGetValue("slug",       out var slugVal))  user.Slug     = slugVal?.ToString()  ?? "";

            if (claimsDict.TryGetValue("status", out var statusVal) && statusVal != null)
            {
                try { user.Status = Enum.Parse<UserStatus>(statusVal.ToString()!, ignoreCase: true); }
                catch { user.Status = null; }
            }

            if (claimsDict.TryGetValue("last_login", out var llVal) && long.TryParse(llVal?.ToString(), out var llUnix))
                user.LastLogin = DateTimeOffset.FromUnixTimeSeconds(llUnix).UtcDateTime;

            if (claimsDict.TryGetValue("created_at", out var caVal) && long.TryParse(caVal?.ToString(), out var caUnix))
                user.CreatedAt = DateTimeOffset.FromUnixTimeSeconds(caUnix).UtcDateTime;

            if (claimsDict.TryGetValue("exp", out var expVal) && long.TryParse(expVal?.ToString(), out var expUnix))
                user.ExpiresAt = DateTimeOffset.FromUnixTimeSeconds(expUnix).UtcDateTime;

            return user;
        }
        catch
        {
            return null;
        }
    }

    private static string PadBase64(string base64) =>
        base64.PadRight(base64.Length + (4 - base64.Length % 4) % 4, '=');
}

// ---------------------------------------------------------------------------
// AuthHeaderHandler
// Attaches the Bearer token from memory — never touches JS interop directly.
// Reading _cachedToken is synchronous so the token is always ready the
// instant after SignInAsync completes, eliminating the 401 race on first load.
// ---------------------------------------------------------------------------
public class AuthHeaderHandler : DelegatingHandler
{
    private readonly CustomAuthStateProvider _auth;

    public AuthHeaderHandler(CustomAuthStateProvider auth) => _auth = auth;

    protected override async Task<HttpResponseMessage> SendAsync(
        HttpRequestMessage request, CancellationToken cancellationToken)
    {
        var token = _auth.GetCachedToken();

        if (!string.IsNullOrEmpty(token))
            request.Headers.Authorization =
                new System.Net.Http.Headers.AuthenticationHeaderValue("Bearer", token);

        return await base.SendAsync(request, cancellationToken);
    }
}


// ---------------------------------------------------------------------------
// AuthExpiredHandler
//
// PURPOSE: Silently log out users whose token has expired mid-session
// (e.g. they left a tab open overnight and resume activity).
//
// GUARD: Only triggers SignOutAsync when _cachedToken is present — meaning
// we were genuinely authenticated and the server rejected a real session.
// Without this guard, a racing HTTP call fired before SignInAsync completes
// could receive a 401 (token not yet attached) and incorrectly sign the
// user out immediately after a successful login.
//
// TODO (WIP — Session Tokens):
// When server-side session tokens are implemented, replace this simple 401
// check with a token-refresh flow:
//   1. On 401, attempt POST /auth/refresh with the refresh token.
//   2. If refresh succeeds, update _cachedToken + localStorage and retry
//      the original request transparently.
//   3. Only call SignOutAsync if the refresh itself returns 401/403.
// This will allow seamless session continuation without forcing re-login.
// ---------------------------------------------------------------------------
public class AuthExpiredHandler : DelegatingHandler
{
    private readonly CustomAuthStateProvider _auth;

    public AuthExpiredHandler(CustomAuthStateProvider auth) => _auth = auth;

    protected override async Task<HttpResponseMessage> SendAsync(
        HttpRequestMessage request, CancellationToken cancellationToken)
    {
        var response = await base.SendAsync(request, cancellationToken);

        if (response.StatusCode == System.Net.HttpStatusCode.Unauthorized)
        {
            // Only sign out if we have a cached token — i.e. we were genuinely
            // authenticated and the server rejected a real (likely expired) session.
            // If _cachedToken is null, the request raced ahead of SignInAsync and
            // the 401 is a timing artifact, not an expired session — ignore it.
            if (!string.IsNullOrEmpty(_auth.GetCachedToken()))
            {
                await _auth.SignOutAsync();
            }
        }

        return response;
    }
}
