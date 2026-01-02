using System.ComponentModel.DataAnnotations;
using System.Text.Json.Serialization;

namespace frontend.Models;

public class Order
{
    public Guid Id { get; set; }
    public Guid StoreId { get; set; }
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

public enum OrderStatus
    {
        Pending,
        Assigned,
        InTransit,
        Delivered,
        Cancelled
    }

public sealed class CreateOrderRequest
{
    [Required]
    [JsonPropertyName("store_id")]
    public Guid StoreId { get; set; }

    [Required]
    [JsonPropertyName("customer_id")]
    public Guid CustomerId { get; set; }

    [Required]
    [JsonPropertyName("payment_method")]
    public string PaymentMethod { get; set; } = "cash";

    [Required]
    [JsonPropertyName("delivery")]
    public LocationDto Delivery { get; set; } = new();

    [Required]
    [JsonPropertyName("pickup")]
    public LocationDto Pickup { get; set; } = new();

    [Required]
    [JsonPropertyName("items")]
    public List<CreateOrderItemDto> Items { get; set; } = new();
}

public sealed class CreateOrderItemDto
{
    [Required]
    [JsonPropertyName("product_id")]
    public Guid ProductId { get; set; }

    [JsonPropertyName("variant_id")]
    public Guid? VariantId { get; set; }

    [Required]
    [JsonPropertyName("quantity")]
    public int Quantity { get; set; }
}

public sealed class LocationDto
{
    public string Address { get; set; } = string.Empty;
    public double Lat { get; set; }
    public double Lng { get; set; }
}
