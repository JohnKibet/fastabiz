using logistics_frontend.Models;

public class OrderService 
{
    private List<OrderItem> _orders = new();

    public List<OrderItem> GetOrders()
    {
        return _orders;
    }

    public void AddOrder(OrderItem order) 
    {
        order.Id = _orders.Count + 1;
        _orders.Add(order);
    }
}
