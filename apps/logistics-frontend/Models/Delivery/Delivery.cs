namespace logistics_frontend.Models.Delivery;

public class Delivery
{
    public Guid ID { get; set; }
    public Guid OrderID { get; set; }
    public Guid DriverId { get; set; }
    public string AssignedAt { get; set; } = string.Empty;
    public string PickedUpAt { get; set; } = string.Empty;
    public string DeliveredAt { get; set; } = string.Empty;
    public string Status { get; set; } = string.Empty;
}

public class CreateDelivery
{
    public Guid OrderID { get; set; }
    public Guid DriverID { get; set; }
}