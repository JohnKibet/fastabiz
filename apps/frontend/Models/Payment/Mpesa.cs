using System.Text.Json.Serialization;

public class MpesaExpressRequest
{
    [JsonPropertyName("phone")]
    public string Phone { get; set; } = string.Empty;

    [JsonPropertyName("amount")]
    public string Amount { get; set; } = string.Empty;
}

public class MpesaExpressResponse
{
    [JsonPropertyName("MerchantRequestID")]
    public string? MerchantRequestID { get; set; }

    [JsonPropertyName("CheckoutRequestID")]
    public string? CheckoutRequestID { get; set; }

    [JsonPropertyName("ResponseCode")]
    public string? ResponseCode { get; set; }

    [JsonPropertyName("ResponseDescription")]
    public string? ResponseDescription { get; set; }

    [JsonPropertyName("CustomerMessage")]
    public string? CustomerMessage { get; set; }
}
