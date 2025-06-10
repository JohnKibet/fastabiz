using System.Net.Http.Json;
using logistics_frontend.Models.Delivery;

public class DeliveryService
{
    private readonly HttpClient _http;

    public DeliveryService(HttpClient http)
    {
        _http = http;
    }

    public async Task CreateDelivery(CreateDelivery delivery)
    {
        var response = await _http.PostAsJsonAsync("deliveries/create", delivery);
        response.EnsureSuccessStatusCode();
    }

    public async Task<Delivery> GetDeliveryById(Guid Id)
    {
        var delivery = await _http.GetFromJsonAsync<Delivery>($"deliveries/id/{Id}");
        return delivery ?? throw new Exception("No delivery found");
    }
    public async Task<List<Delivery>> GetDeliveries()
    {
        var deliveries = await _http.GetFromJsonAsync<List<Delivery>>("deliveries/all_deliveries");
        return deliveries ?? new List<Delivery>();
    }
}