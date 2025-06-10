namespace logistics_frontend.Models.Payment
{
    public class Payment
    {
        public Guid ID { get; set; }
        public Guid OrderID { get; set; }
        public double Amount { get; set; }
        public string Method { get; set; } = string.Empty;
        public string Status { get; set; } = string.Empty;
        public DateTime PaidAt { get; set; }
    }

    public class PaymentRequest
    {
        public Guid OrderID { get; set; }
        public double Amount { get; set; }
        public string Method { get; set; } = string.Empty;
    }
}
