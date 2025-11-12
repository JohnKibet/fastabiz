// using System.Net.Http.Json;
// using frontend.Models.Payment;

// public class PaymentService
// {
//     private readonly HttpClient _http;
//     public PaymentService(IHttpClientFactory httpClientFactory)
//     {
//         _http = httpClientFactory.CreateClient("AuthenticatedApi");
//     }

//     public async Task MakePayment(PaymentRequest payment)
//     {
//         var response = await _http.PostAsJsonAsync("payments/create", payment);
//         response.EnsureSuccessStatusCode();
//     }

//     public async Task<Payment> GetPaymentById(Guid paymentId)
//     {
//         var payment = await _http.GetFromJsonAsync<Payment>($"payments/{paymentId}");
//         return payment ?? throw new Exception("Payment not found");
//     }

//     public async Task<List<Payment>> GetPaymentsByOrderId(Guid orderId)
//     {
//         var payments = await _http.GetFromJsonAsync<List<Payment>>($"payments/{orderId}");
//         return payments ?? new List<Payment>();
//     }

//     public async Task<List<Payment>> GetAllPayments()
//     {
//         var payments = await _http.GetFromJsonAsync<List<Payment>>("payments/all_payments");
//         return payments ?? new List<Payment>();
//     }
// }


using System;
using System.Collections.Generic;
using System.Linq;
using System.Threading.Tasks;
using frontend.Models;

namespace frontend.Services
{
    public class PaymentService
    {
        private readonly InvoiceService _invoiceService;

        public PaymentService(InvoiceService invoiceService)
        {
            _invoiceService = invoiceService;
        }

        public Task<Invoice?> GetInvoiceAsync(Guid invoiceId)
        {
            return _invoiceService.GetByIdAsync(invoiceId);
        }

        public async Task<string> ProcessPaymentAsync(Guid invoiceId, string? method)
        {
            var invoice = await _invoiceService.GetByIdAsync(invoiceId);
            if (invoice == null)
                return "failed";

            // Simulate processing
            await Task.Delay(1000);

            // Mock logic: mark invoice as paid for demonstration
            await _invoiceService.MarkAsPaidAsync(invoiceId);

            return "success";
        }
    }
}