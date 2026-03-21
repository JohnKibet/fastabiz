using System.Collections.Generic;
using System.Linq;
using System.Text.Json.Serialization;

namespace frontend.Models.Storefront;

/// <summary>
/// DTO shapes used for API responses — separate from the domain models
/// above so serialization concerns don't leak into business logic.
/// </summary>
public static class ProductMapper
{
    // ── DTOs ──────────────────────────────────────────────────────────

    public class ProductDto
    {
        [JsonPropertyName("id")]          public System.Guid   Id          { get; set; }
        [JsonPropertyName("store_id")]    public System.Guid   StoreId     { get; set; }
        [JsonPropertyName("name")]        public string        Name        { get; set; } = string.Empty;
        [JsonPropertyName("description")] public string        Description { get; set; } = string.Empty;
        [JsonPropertyName("category")]    public string        Category    { get; set; } = string.Empty;
        [JsonPropertyName("images")]      public List<string>? Images      { get; set; } = new();
        [JsonPropertyName("has_variants")] public bool         HasVariants { get; set; }
        [JsonPropertyName("created_at")]  public System.DateTime CreatedAt { get; set; }
        [JsonPropertyName("updated_at")]  public System.DateTime UpdatedAt { get; set; }

        public List<OptionDto>  Options  { get; set; } = new();
        public List<VariantDto> Variants { get; set; } = new();
    }

    public class VariantDto
    {
        [JsonPropertyName("id")]          public System.Guid   Id        { get; set; }
        [JsonPropertyName("product_id")]  public System.Guid   ProductId { get; set; }
        [JsonPropertyName("sku")]         public string        Sku       { get; set; } = string.Empty;
        [JsonPropertyName("price")]       public double        Price     { get; set; }
        [JsonPropertyName("stock")]       public int           Stock     { get; set; }
        [JsonPropertyName("image_url")]   public string?       ImageUrl  { get; set; }
        [JsonPropertyName("created_at")]  public System.DateTime CreatedAt { get; set; }
        [JsonPropertyName("updated_at")]  public System.DateTime UpdatedAt { get; set; }
        [JsonPropertyName("options")]     public Dictionary<string, string> Options { get; set; } = new();
    }

    public class OptionDto
    {
        [JsonPropertyName("name")]   public string       Name   { get; set; } = string.Empty;
        [JsonPropertyName("values")] public List<string> Values { get; set; } = new();
    }

    // ── Mapping ───────────────────────────────────────────────────────

    public static ProductDto ToDto(this Product product)
    {
        var variants = product.Variants?.Select(v => new VariantDto
        {
            Id        = v.Id,
            ProductId = v.ProductId,
            Sku       = v.SKU,
            Price     = v.Price,
            Stock     = v.Stock,
            ImageUrl  = v.ImageUrl,
            Options   = v.Options,
            CreatedAt = v.CreatedAt,
            UpdatedAt = v.UpdatedAt
        }).ToList() ?? new();

        return new ProductDto
        {
            Id          = product.Id,
            StoreId     = product.StoreId,
            Name        = product.Name,
            Description = product.Description,
            Category    = product.Category,
            Images      = product.Images,
            HasVariants = product.HasVariants,
            Variants    = variants,
            Options     = product.HasVariants ? DeriveOptions(variants) : new(),
            CreatedAt   = product.CreatedAt,
            UpdatedAt   = product.UpdatedAt
        };
    }

    private static List<OptionDto> DeriveOptions(List<VariantDto> variants) =>
        variants
            .SelectMany(v => v.Options)
            .GroupBy(o => o.Key)
            .Select(g => new OptionDto
            {
                Name   = g.Key,
                Values = g.Select(x => x.Value).Distinct().OrderBy(v => v).ToList()
            })
            .ToList();
}
