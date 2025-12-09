namespace frontend.Models;

public class NewOrder
{
    public Guid Id { get; set; }
    public Guid MerchantId { get; set; }
    public Guid CustomerId { get; set; }

    public Guid ProductId { get; set; }
    public Guid VariantId { get; set; }

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