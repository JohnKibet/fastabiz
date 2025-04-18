namespace logistics_frontend.Models {
    public class DeliveryItem
    {
        public int Id {get; set;}
        public required string CustomerName {get; set;}
        public required string ProductName {get; set;}
        public DateTime DeliveryDate {get; set;}
        public required string Status {get; set;}
    }

}