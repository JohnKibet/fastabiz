using System;
using System.Net.Http;
using System.Net.Http.Json;
using System.Text.Json;
using System.Threading.Tasks;
using System.Linq;
using frontend.Models;

namespace frontend.Services;

/// <summary>
/// Generic HTTP service — all CRUD operations in one place.
/// Specific services (ProductService, StoreService, etc.) delegate here.
/// Uses ApiResult<T> throughout — no ServiceResult<T> or ServiceResult2<T>.
/// </summary>
public class ApiService
{
    private readonly HttpClient _http;
    private readonly JsonSerializerOptions _jsonOpts = new()
    {
        PropertyNameCaseInsensitive = true
    };

    public ApiService(IHttpClientFactory factory)
    {
        _http = factory.CreateClient("AuthenticatedApi");
    }

    // ── GET ──────────────────────────────────────────────────────────

    public async Task<ApiResult<T>> GetAsync<T>(string url)
    {
        try
        {
            var response = await _http.GetAsync(url);
            return await ParseResponse<T>(response);
        }
        catch (Exception ex) { return Fail<T>(ex); }
    }

    // ── POST ─────────────────────────────────────────────────────────

    public async Task<ApiResult<TResponse>> PostAsync<TRequest, TResponse>(string url, TRequest body)
    {
        try
        {
            var response = await _http.PostAsJsonAsync(url, body);
            return await ParseResponse<TResponse>(response);
        }
        catch (Exception ex) { return Fail<TResponse>(ex); }
    }

    /// POST with no body (trigger endpoints, e.g. cloudinary signature)
    public async Task<ApiResult<TResponse>> PostAsync<TResponse>(string url)
    {
        try
        {
            var response = await _http.PostAsync(url, null);
            return await ParseResponse<TResponse>(response);
        }
        catch (Exception ex) { return Fail<TResponse>(ex); }
    }

    // ── PUT ──────────────────────────────────────────────────────────

    public async Task<ApiResult<TResponse>> PutAsync<TRequest, TResponse>(string url, TRequest body)
    {
        try
        {
            var response = await _http.PutAsJsonAsync(url, body);
            return await ParseResponse<TResponse>(response);
        }
        catch (Exception ex) { return Fail<TResponse>(ex); }
    }

    // ── PATCH ────────────────────────────────────────────────────────

    public async Task<ApiResult<TResponse>> PatchAsync<TRequest, TResponse>(string url, TRequest body)
    {
        try
        {
            var response = await _http.PatchAsJsonAsync(url, body);
            return await ParseResponse<TResponse>(response);
        }
        catch (Exception ex) { return Fail<TResponse>(ex); }
    }

    // ── DELETE ───────────────────────────────────────────────────────

    public async Task<ApiResult<bool>> DeleteAsync(string url)
    {
        try
        {
            var response = await _http.DeleteAsync(url);
            if (response.IsSuccessStatusCode) return ApiResult<bool>.Ok(true);
            return ApiResult<bool>.Fail(await ParseApiError(response));
        }
        catch (Exception ex) { return Fail<bool>(ex); }
    }

    // ── Shared helpers ────────────────────────────────────────────────

    private async Task<ApiResult<T>> ParseResponse<T>(HttpResponseMessage response)
    {
        if (!response.IsSuccessStatusCode)
            return ApiResult<T>.Fail(await ParseApiError(response));

        if (response.StatusCode == System.Net.HttpStatusCode.NoContent)
            return ApiResult<T>.Ok(default!);

        var result = await response.Content.ReadFromJsonAsync<T>(_jsonOpts);
        return ApiResult<T>.Ok(result ?? Activator.CreateInstance<T>());
    }

    /// Deserialises the Go backend { "error": "...", "detail": "..." } shape.
    public async Task<ApiError> ParseApiError(HttpResponseMessage response)
    {
        try
        {
            var json  = await response.Content.ReadAsStringAsync();
            var error = JsonSerializer.Deserialize<ApiError>(json, _jsonOpts);

            return error ?? new ApiError
            {
                Error = $"HTTP {(int)response.StatusCode} – {response.ReasonPhrase}"
            };
        }
        catch
        {
            return new ApiError
            {
                Error = $"HTTP {(int)response.StatusCode} – {response.ReasonPhrase}"
            };
        }
    }

    private static ApiResult<T> Fail<T>(Exception ex) =>
        ApiResult<T>.Fail(
            ex is HttpRequestException
                ? new ApiError { Error = "Network error", Detail = ex.Message }
                : new ApiError { Error = "Unexpected error", Detail = ex.Message }
        );
}
