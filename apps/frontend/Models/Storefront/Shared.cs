using System.Text.Json.Serialization;

namespace frontend.Models.Storefront;

/// <summary>
/// Generic API message envelope — used for operations that return
/// only a confirmation string (delete, update, etc.)
/// </summary>
public class ApiMessageResponse
{
    [JsonPropertyName("message")] public string Message { get; set; } = string.Empty;
}
