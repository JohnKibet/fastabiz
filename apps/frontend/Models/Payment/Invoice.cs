using System.Dynamic;
using System.Net.NetworkInformation;

namespace frontend.Models;

public class Invoice
{
    public Guid ID { get; set; }
    public string InvoiceNumber { get; set; } = string.Empty;
    public string CustomerName { get; set; } = string.Empty;
    public string CustomerEmail { get; set; } = string.Empty;
    public decimal Amount { get; set; } 
    public string Status { get; set; } = "Pending";
    public DateTime CreatedAt { get; set; }
    public List<InvoiceItem> Items { get; set; }
}

public class InvoiceItem
{
    public string ItemName { get; set; } = string.Empty;
    public decimal UnitPrice { get; set; }
    public int Quantity { get; set; }
    public decimal Total => UnitPrice * Quantity;
}


public class InvoiceSummary
{   
    public int Total { get; set; }
    public int Pending { get; set; }
    public int Paid { get; set; }
    public int Overdue { get; set; }
    public decimal TotalAmount { get; set; }
}