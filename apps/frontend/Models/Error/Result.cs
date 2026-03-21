using System.Collections.Generic;
using System.Text.Json.Serialization;

namespace frontend.Models;

// ── ApiError ───────────────────────────────────────────────────────────
//
// Mirrors the Go backend error shape exactly:
//
//   type ErrorResponse struct {
//       Error  string `json:"error"`            // always present
//       Detail string `json:"detail,omitempty"` // only on 5xx
//   }
//
// writeJSONError() always produces:
//   { "error": "user-friendly message", "detail": "internal error (5xx only)" }
//
// Renamed from ErrorResponse → ApiError to avoid confusion with the
// Blazor/ASP.NET ErrorResponse class that exists in some project templates.

public class ApiError
{
    [JsonPropertyName("error")]
    public string  Error  { get; set; } = string.Empty;

    // Only populated by the backend on 5xx responses
    [JsonPropertyName("detail")]
    public string? Detail { get; set; }

    // Future: if you add a Go validation middleware that emits field errors
    // in the shape { "errors": { "fieldName": ["msg1", "msg2"] } }
    [JsonPropertyName("errors")]
    public Dictionary<string, string[]>? Errors { get; set; }

    // Convenience: the most specific message available
    public string Message =>
        !string.IsNullOrWhiteSpace(Detail)  ? Detail  :
        !string.IsNullOrWhiteSpace(Error)   ? Error   :
        "An unexpected error occurred.";
}

// ── ApiResult<T> ───────────────────────────────────────────────────────
//
// Single result wrapper used by ALL services and components.
// Replaces both ServiceResult<T> and ServiceResult2<T>.
//
// RULES:
//   • Success = true  → Data is populated, Error is null
//   • Success = false → Error is populated, Data is default
//   • Never check Data without checking Success first
//
// FACTORY METHODS:
//   ApiResult<T>.Ok(data)         — success with data
//   ApiResult<T>.Ok(data, true)   — success, data came from cache
//   ApiResult<T>.Fail(apiError)   — failure with full structured error
//   ApiResult<T>.Fail("message")  — failure with plain string (wraps into ApiError)
//
// USAGE IN SERVICES (via ApiService):
//   public Task<ApiResult<List<ProductListItem>>> ListProducts(Guid storeId)
//       => _api.GetAsync<List<ProductListItem>>($"products/{storeId}/all_products");
//
// USAGE IN COMPONENTS:
//   var result = await ProductService.ListProducts(StoreCtx.StoreId);
//   if (!result.Success)
//   {
//       _error = result.Error?.Message;   // always safe — never null on failure
//       return;
//   }
//   _products = result.Data!;             // safe to ! after Success check

public class ApiResult<T>
{
    public bool     Success   { get; private set; }
    public T?       Data      { get; private set; }
    public ApiError? Error    { get; private set; }
    public bool     FromCache { get; private set; }

    // Convenience guards
    public bool HasData             => Success && Data != null;
    public bool HasValidationErrors => Error?.Errors != null && Error.Errors.Count > 0;

    // ── Factories ────────────────────────────────────────────────────

    public static ApiResult<T> Ok(T data, bool fromCache = false) => new()
    {
        Success   = true,
        Data      = data,
        FromCache = fromCache
    };

    public static ApiResult<T> Fail(ApiError error) => new()
    {
        Success = false,
        Error   = error
    };

    /// Convenience overload — wraps a plain string in an ApiError so
    /// callers never have to construct ApiError manually for simple cases.
    public static ApiResult<T> Fail(string message) => new()
    {
        Success = false,
        Error   = new ApiError { Error = message }
    };
}
