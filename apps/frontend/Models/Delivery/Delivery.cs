using System.ComponentModel.DataAnnotations;
using System.Text.Json.Serialization;

namespace frontend.Models;

public class Delivery
{
    [JsonPropertyName("id")]
    public Guid ID { get; set; }

    [JsonPropertyName("order_id")]
    public Guid OrderID { get; set; }

    [JsonPropertyName("driver_id")]
    public Guid DriverId { get; set; }

    [JsonPropertyName("assigned_at")]
    public DateTime AssignedAt { get; set; }

    [JsonPropertyName("picked_up_at")]
    public DateTime PickedUpAt { get; set; }

    [JsonPropertyName("delivered_at")]
    public DateTime DeliveredAt { get; set; }

    [JsonPropertyName("status")]
    [JsonConverter(typeof(JsonStringEnumConverter))]
    public DeliveryStatus Status { get; set; }

    [JsonPropertyName("pickup_address")]
    public string PickupAddress { get; set; } = string.Empty;

    [JsonPropertyName("delivery_address")]
    public string DeliveryAddress { get; set; } = string.Empty;
}

public class CreateDeliveryRequest
{
    [Required]
    public Guid OrderID { get; set; }

    [Required]
    public Guid DriverID { get; set; }
}

public class CreateDeliveryResponse
{

}

public enum DeliveryStatus
{
    Assigned,
    PickedUp,
    Delivered,
    Failed

}