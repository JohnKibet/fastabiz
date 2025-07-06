using System.Net.Http.Json;
using logistics_frontend.Models.Order;
public class OrderService
{
    private readonly HttpClient _http;
    public OrderService(IHttpClientFactory httpClientFactory)
    {
        _http = httpClientFactory.CreateClient("AuthenticatedApi");
    }

    public async Task AddOrder(CreateOrderRequest order)
    {
        var response = await _http.PostAsJsonAsync("orders/create", order);
        response.EnsureSuccessStatusCode();
    }

    public async Task<List<Order>> GetOrderByID(Guid id)
    {
        var order = await _http.GetFromJsonAsync<List<Order>>($"orders/{id}");
        return order ?? new List<Order>();
    }


    public async Task<List<Order>> GetOrdersByCustomer(Guid customerId)
    {
        var orders = await _http.GetFromJsonAsync<List<Order>>($"orders/{customerId}");
        return orders ?? new List<Order>();
    }

    public async Task<List<Order>> GetAllOrders()
    {
        var orders = await _http.GetFromJsonAsync<List<Order>>("orders/all_orders");
        return orders ?? new List<Order>();
    }
}
