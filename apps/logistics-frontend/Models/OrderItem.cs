using System.Text;

public class OrderItem 
{
    public int Id {get; set;}
    public string? CustomerName {get; set;}
    public string? ProductName {get; set;}
    public int Quantity {get; set;}
    public string? Address {get; set;}
    public DateTime OrderDate {get; set;} = DateTime.Now;
    public string Status {get; set;} = "Pending";
}