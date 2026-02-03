using System.ComponentModel.DataAnnotations;
using System.Text.Json.Serialization;

namespace frontend.Models;

public class Order
{
    [JsonPropertyName("id")]
    public Guid Id { get; set; }

    [JsonPropertyName("store_id")]
    public Guid StoreId { get; set; }

    [JsonPropertyName("customer_id")]
    public Guid CustomerId { get; set; }

    [JsonPropertyName("product_id")]
    public Guid ProductId { get; set; }

    [JsonPropertyName("variant_id")]
    public Guid? VariantId { get; set; }

    [JsonPropertyName("quantity")]
    public int Quantity { get; set; }

    [JsonPropertyName("unit_price")]
    public double UnitPrice { get; set; }

    [JsonPropertyName("currency")]
    public string Currency { get; set; } = "KES";

    [JsonPropertyName("total")]
    public double Total { get; set; }

    [JsonPropertyName("product_name")]
    public string ProductName { get; set; } = string.Empty;

    [JsonPropertyName("variant_name")]
    public string? VariantName { get; set; }

    [JsonPropertyName("image_url")]
    public string? ImageUrl { get; set; }

    [JsonPropertyName("pickup_address")]
    public string PickupAddress { get; set; } = string.Empty;

    [JsonPropertyName("delivery_address")]
    public string DeliveryAddress { get; set; } = string.Empty;

    [JsonPropertyName("status")]
    public OrderStatus Status { get; set; }

    [JsonPropertyName("created_at")]
    public DateTime CreatedAt { get; set; }
}

public enum OrderStatus
{
    Pending,
    Assigned,
    InTransit,
    Delivered,
    Cancelled
}

// mapping CartItem → CreateOrderItemDto.
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
    public string? PaymentMethod { get; set; }

    [Required]
    [JsonPropertyName("delivery")]
    public LocationDto Delivery { get; set; } = new();

    [Required]
    [JsonPropertyName("pickup")]
    public LocationDto Pickup { get; set; } = new();

    [Required]
    [JsonPropertyName("items")]
    public List<CreateOrderItemDto> Items { get; set; } = new();

    public Dictionary<string, string>? ExtraData { get; set; }
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
