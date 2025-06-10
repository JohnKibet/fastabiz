using System.Security.Claims;

namespace logistics_frontend.Helpers
{
    public static class ClaimsPrincipalExtensions
    {
        public static bool IsInRole(this ClaimsPrincipal user, string role)
        {
            return user.Identity?.IsAuthenticated == true &&
                user.FindFirst(ClaimTypes.Role)?.Value == role;
        }

        public static string? GetUserId(this ClaimsPrincipal user)
        {
            return user.FindFirst(ClaimTypes.NameIdentifier)?.Value;
        }
    }   
}