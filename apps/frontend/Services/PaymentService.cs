using System.Net.Http.Json;
using frontend.Models;

public class PaymentService
{
  private readonly HttpClient _http;

  public PaymentService(IHttpClientFactory httpClientFactory)
  {
    _http = httpClientFactory.CreateClient("AuthenticatedApi");
  }

  public async Task<ServiceResult<MpesaExpressResponse>> StartMpesaExpressAsync(string phone, string amount)
  {
    try
    {
      var payload = new MpesaExpressRequest
      {
        Phone = phone,
        Amount = amount
      };

      var response = await _http.PostAsJsonAsync("payments/mpesa-express", payload);

      if (!response.IsSuccessStatusCode)
      {
        var error = await response.Content.ReadAsStringAsync();
        return ServiceResult<MpesaExpressResponse>.Fail(
            string.IsNullOrWhiteSpace(error)
                ? "Failed to initiate MPesa payment"
                : error
        );
      }

      // SAFE JSON read
      var data = await response.Content
          .ReadFromJsonAsync<MpesaExpressResponse>();

      if (data is null)
      {
        return ServiceResult<MpesaExpressResponse>.Fail(
            "Invalid server response"
        );
      }

      return ServiceResult<MpesaExpressResponse>.Ok(data);
    }
    catch (HttpRequestException ex)
    {
      // Browser / CORS / network issues
      return ServiceResult<MpesaExpressResponse>.Fail(
          $"Network error: {ex.Message}"
      );
    }
    catch (Exception ex)
    {
      return ServiceResult<MpesaExpressResponse>.Fail(
          $"Unexpected error: {ex.Message}"
      );
    }
  }

}
