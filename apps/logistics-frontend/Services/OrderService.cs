using System.Net.Http.Json;
using logistics_frontend.Models.Order;
public class OrderService 
{
    private readonly HttpClient _http;

    public OrderService(HttpClient http)
    {
        _http = http;
    }

    public async Task AddOrder(CreateOrderRequest order)
    {
        var response = await _http.PostAsJsonAsync("http://192.168.1.18:8080/orders/create", order);
        response.EnsureSuccessStatusCode();
    }

    public async Task<List<Order>> GetAllOrders()
    {
        var orders = await _http.GetFromJsonAsync<List<Order>>("http://192.168.1.18:8080/orders/all_orders");
        if (orders == null)
            return new List<Order>(); // fallback in case of null

        return orders;
    }
}
