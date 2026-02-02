using System.Text.Json.Serialization;

public class CloudinarySignatureResponse
{
  [JsonPropertyName("cloud_name")]
  public string Cloud_Name { get; set; } = string.Empty;

  [JsonPropertyName("api_key")]
  public string Api_Key { get; set; } = string.Empty;

  [JsonPropertyName("timestamp")]
  public long Timestamp { get; set; }

  [JsonPropertyName("signature")]
  public string Signature { get; set; } = string.Empty;

  [JsonPropertyName("folder")]
  public string Folder { get; set; } = string.Empty;
}
