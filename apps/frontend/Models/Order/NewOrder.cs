namespace frontend.Models;

public class NewOrderModel
{
    public Guid Id { get; set; }
    public Guid MerchantId { get; set; }
    public Guid CustomerId { get; set; }

    // temp: string instead of Guid
    public string ProductId { get; set; } = string.Empty;
    public string VariantId { get; set; } = string.Empty;

    public int Quantity { get; set; }

    public double UnitPrice { get; set; }
    public string Currency { get; set; } = "KES";
    public int Total { get; set; }

    public string PickupAddress { get; set; } = string.Empty;
    public double PickupLat { get; set; }
    public double PickupLng { get; set; }

    public string DeliveryAddress { get; set; } = string.Empty;
    public double DeliveryLat { get; set; }
    public double DeliveryLng { get; set; }

    public OrderStatus Status { get; set; }
    public DateTime CreatedAt { get; set; }
    public DateTime UpdatedAt { get; set; }

}

public class NewCreateOrderRequest
{
    public Guid CustomerId { get; set; }
    public string DeliveryAddress { get; set; } = string.Empty;
    public double DeliveryLatitude { get; set; }
    public double DeliveryLongitude { get; set; }
    public string PickupAddress { get; set; } = string.Empty;
    public double PickupLatitude { get; set; }
    public double PickupLongitude { get; set; }
    public List<OrderLineItem> Items { get; set; } = new();
    public string PaymentMethod { get; set; } = "cash";
}

public class OrderLineItem
{
    public string ProductId { get; set; } = string.Empty;
    public string VariantId { get; set; } = string.Empty;
    public int Quantity { get; set; }
    public double Price { get; set; }
}
