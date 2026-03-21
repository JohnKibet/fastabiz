using System;
using System.Collections.Generic;
using System.ComponentModel.DataAnnotations;
using System.Text.Json.Serialization;

namespace frontend.Models.Storefront;

public class ProductImage
{
    public string URL       { get; set; } = string.Empty;
    public bool   IsPrimary { get; set; }
}

public class AddImageRequest
{
    [Required]
    [JsonPropertyName("product_id")] public Guid              ProductId { get; set; }
    [Required]
    [JsonPropertyName("images")]     public List<ProductImage> Images   { get; set; } = new();
}

public class ReorderImagesRequest
{
    [Required]
    [JsonPropertyName("product_id")] public Guid       ProductId { get; set; }
    [Required]
    [JsonPropertyName("image_ids")]  public List<Guid> ImageIDs  { get; set; } = new();
}
