namespace frontend.Models
{
    public class Payment
    {
        public Guid ID { get; set; }
        public Guid InvoiceID { get; set; }
        public string Method { get; set; } = string.Empty;
        public string Status { get; set; } = "Pending";
        public string? TransactionRef { get; set; }
        public DateTime CreatedAt { get; set; } = DateTime.UtcNow;
    }
}
