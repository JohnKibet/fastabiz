using frontend.Models;

namespace frontend.Services;

public class PaymentService
{
    private readonly ApiService _api;

    public PaymentService(ApiService api)
    {
        _api = api;
    }

    /// <summary>
    /// Initiates an M-Pesa STK push (Daraja Express).
    /// Returns the CheckoutRequestID and other metadata on success.
    /// </summary>
    public Task<ApiResult<MpesaExpressResponse>> StartMpesaExpressAsync(string phone, string amount)
        => _api.PostAsync<MpesaExpressRequest, MpesaExpressResponse>(
               "payments/mpesa-express",
               new MpesaExpressRequest { Phone = phone, Amount = amount });
}
