using System.Security.Claims;
using Microsoft.AspNetCore.Components.Authorization;

public class MockAuthProvider : AuthenticationStateProvider
{
    public override Task<AuthenticationState> GetAuthenticationStateAsync()
    {
        // simulated a logged-in admin user
        var identity = new ClaimsIdentity(new[]
        {
            new Claim(ClaimTypes.Name, "adminuser"),
            new Claim(ClaimTypes.Role, "admin")
        }, "mock");

        var user = new ClaimsPrincipal(identity);
        return Task.FromResult(new AuthenticationState(user));
    }
}