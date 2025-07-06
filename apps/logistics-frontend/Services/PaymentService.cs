using System.Net.Http.Json;
using logistics_frontend.Models.Payment;

public class PaymentService
{
    private readonly HttpClient _http;
    public PaymentService(IHttpClientFactory httpClientFactory)
    {
        _http = httpClientFactory.CreateClient("AuthenticatedApi");
    }

    public async Task MakePayment(PaymentRequest payment)
    {
        var response = await _http.PostAsJsonAsync("payments/create", payment);
        response.EnsureSuccessStatusCode();
    }

    public async Task<Payment> GetPaymentById(Guid paymentId)
    {
        var payment = await _http.GetFromJsonAsync<Payment>($"payments/{paymentId}");
        return payment ?? throw new Exception("Payment not found");
    }

    public async Task<List<Payment>> GetPaymentsByOrderId(Guid orderId)
    {
        var payments = await _http.GetFromJsonAsync<List<Payment>>($"payments/{orderId}");
        return payments ?? new List<Payment>();
    }

    public async Task<List<Payment>> GetAllPayments()
    {
        var payments = await _http.GetFromJsonAsync<List<Payment>>("payments/all_payments");
        return payments ?? new List<Payment>();
    }
}