using System.Text.Json.Serialization;

namespace frontend.Models;

public class Invoice
{
    // Razor uses inv.Id — kept lowercase to match C# convention and page usage
    [JsonPropertyName("id")]
    public Guid Id { get; set; }

    [JsonPropertyName("invoice_number")]
    public string InvoiceNumber { get; set; } = string.Empty;

    [JsonPropertyName("customer_name")]
    public string CustomerName { get; set; } = string.Empty;

    [JsonPropertyName("customer_email")]
    public string CustomerEmail { get; set; } = string.Empty;

    [JsonPropertyName("amount")]
    public decimal Amount { get; set; }

    // "Pending" | "Paid" | "Overdue" — matches SelectField options in Invoices.razor
    [JsonPropertyName("status")]
    public string Status { get; set; } = "Pending";

    [JsonPropertyName("due_date")]
    public DateTime DueDate { get; set; }

    [JsonPropertyName("created_at")]
    public DateTime CreatedAt { get; set; }

    [JsonPropertyName("items")]
    public List<InvoiceItem> Items { get; set; } = new();
}

public class InvoiceItem
{
    [JsonPropertyName("item_name")]
    public string ItemName { get; set; } = string.Empty;

    [JsonPropertyName("unit_price")]
    public decimal UnitPrice { get; set; }

    [JsonPropertyName("quantity")]
    public int Quantity { get; set; }

    // Computed — not serialised
    [JsonIgnore]
    public decimal Total => UnitPrice * Quantity;
}

public class InvoiceSummary
{
    public int     Total       { get; set; }
    public int     Pending     { get; set; }
    public int     Paid        { get; set; }
    public int     Overdue     { get; set; }
    public decimal TotalAmount { get; set; }
}
