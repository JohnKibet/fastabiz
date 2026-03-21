using System;
using System.Collections.Generic;
using System.ComponentModel.DataAnnotations;
using System.Text.Json.Serialization;

namespace frontend.Models.Storefront;

// ── Option ────────────────────────────────────────────────────────────
// Represents a product dimension: e.g. "Size" with values ["S","M","L"]

public class Option
{
    public string       Name   { get; set; } = string.Empty;
    public List<string> Values { get; set; } = new();
}

public class CreateOptionNameRequest
{
    [Required]
    [JsonPropertyName("product_id")] public Guid   ProductId { get; set; }
    [Required]
    [JsonPropertyName("name")]       public string Name      { get; set; } = string.Empty;
}

public class CreateOptionNameResponse
{
    [JsonPropertyName("option_id")] public Guid OptionId { get; set; }
}

public class AddOptionValuesRequest
{
    [Required]
    [JsonPropertyName("product_id")] public Guid        ProductId { get; set; }
    [Required]
    [JsonPropertyName("option_id")]  public Guid        OptionId  { get; set; }
    [Required]
    [JsonPropertyName("values")]     public List<string> Values   { get; set; } = new();
}

// ── Variant ───────────────────────────────────────────────────────────
// A concrete sellable unit combining a set of option values (Size=S, Color=Red)

public class Variant
{
    public Guid   Id        { get; set; }
    public Guid   ProductId { get; set; }
    public string SKU       { get; set; } = string.Empty;
    public double Price     { get; set; }
    public int    Stock     { get; set; }
    public string ImageUrl  { get; set; } = string.Empty;

    // e.g. { "Size": "Small", "Color": "Red" }
    public Dictionary<string, string> Options { get; set; } = new();

    public DateTime CreatedAt { get; set; }
    public DateTime UpdatedAt { get; set; }
}

public class CreateVariantRequest
{
    [Required]
    [JsonPropertyName("product_id")] public Guid   ProductId { get; set; }
    [Required]
    [JsonPropertyName("sku")]        public string SKU       { get; set; } = string.Empty;
    [Required]
    [JsonPropertyName("price")]      public double Price     { get; set; }
    [Required]
    [JsonPropertyName("stock")]      public int    Stock     { get; set; }

    [JsonPropertyName("image_url")]  public string Image     { get; set; } = string.Empty;

    [Required]
    [JsonPropertyName("options")]    public Dictionary<string, string> Options { get; set; } = new();
}

public class CreateVariantResponse
{
    [JsonPropertyName("id")]         public Guid     Id        { get; set; }
    [JsonPropertyName("product_id")] public Guid     ProductId { get; set; }
    [JsonPropertyName("sku")]        public string   SKU       { get; set; } = string.Empty;
    [JsonPropertyName("price")]      public double   Price     { get; set; }
    [JsonPropertyName("stock")]      public int      Stock     { get; set; }
    [JsonPropertyName("image_url")]  public string   ImageUrl  { get; set; } = string.Empty;
    [JsonPropertyName("options")]    public Dictionary<string, string> Options { get; set; } = new();
    [JsonPropertyName("created_at")] public DateTime CreatedAt { get; set; }
    [JsonPropertyName("updated_at")] public DateTime UpdatedAt { get; set; }
}

public class UpdateVariantStockRequest
{
    [Required]
    [JsonPropertyName("variant_id")] public Guid VariantId { get; set; }
    [Required]
    [JsonPropertyName("stock")]      public int  Stock     { get; set; }
}

public class UpdateVariantPriceRequest
{
    [Required]
    [JsonPropertyName("variant_id")] public Guid   VariantId { get; set; }
    [Required]
    [JsonPropertyName("price")]      public double Price     { get; set; }
}
