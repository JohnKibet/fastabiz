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

            var identity = CreateIdentity(_cachedUser);
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
        _cachedUser = null;

        NotifyAuthenticationStateChanged(GetAuthenticationStateAsync());
        await Task.Delay(50);

        _navigationManager.NavigateTo("/auth/login", forceLoad: true);
    }

    private ClaimsIdentity CreateIdentity(User user)
    {
        return new ClaimsIdentity(new[]
        {
            new Claim(ClaimTypes.NameIdentifier, user.ID.ToString()),
            new Claim(ClaimTypes.Name, user.FullName),
            new Claim(ClaimTypes.Email, user.Email),
            new Claim(ClaimTypes.Role, user.Role),
        }, "jwtAuth");
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

            return new User
            {
                ID = Guid.TryParse(claimsDict["sub"]?.ToString(), out var id) ? id : Guid.Empty,
                Email = claimsDict["email"]?.ToString() ?? "",
                Role = claimsDict["role"]?.ToString() ?? "",
                FullName = claimsDict.ContainsKey("name") ? claimsDict["name"]?.ToString() ?? "" : ""
            };
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
}