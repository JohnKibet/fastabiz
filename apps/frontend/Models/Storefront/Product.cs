using System;
using System.Collections.Generic;
using System.ComponentModel.DataAnnotations;
using System.Linq;
using System.Text.Json.Serialization;

namespace frontend.Models.Storefront;

public class Product
{
    public Guid   Id          { get; set; }
    public Guid   StoreId     { get; set; }
    public string Name        { get; set; } = string.Empty;
    public string Description { get; set; } = string.Empty;
    public string Category    { get; set; } = string.Empty;

    public List<string>  Images   { get; set; } = new();
    public List<Option>  Options  { get; set; } = new();
    public List<Variant> Variants { get; set; } = new();

    public bool HasVariants { get; set; }

    // Simple product fields (HasVariants == false)
    public double Price { get; set; }
    public int    Stock { get; set; }

    // Computed helpers
    public double? MinPrice => HasVariants && Variants.Any() ? Variants.Min(v => v.Price) : Price;
    public double? MaxPrice => HasVariants && Variants.Any() ? Variants.Max(v => v.Price) : Price;

    public DateTime CreatedAt { get; set; }
    public DateTime UpdatedAt { get; set; }
}

/// <summary>
/// Lightweight projection used in list/grid views — avoids fetching
/// full variant/option trees for every row.
/// </summary>
public class ProductListItem
{
    [JsonPropertyName("id")]            public Guid   Id           { get; set; }
    [JsonPropertyName("store_id")]      public Guid   StoreId      { get; set; }
    [JsonPropertyName("name")]          public string Name         { get; set; } = string.Empty;
    [JsonPropertyName("category")]      public string Category     { get; set; } = string.Empty;
    [JsonPropertyName("description")]   public string Description  { get; set; } = string.Empty;
    [JsonPropertyName("image_url")]     public string ImageUrl     { get; set; } = string.Empty;
    [JsonPropertyName("has_variants")]  public bool   HasVariants  { get; set; }
    [JsonPropertyName("variant_count")] public int    VariantCount { get; set; }
    [JsonPropertyName("min_price")]     public double? MinPrice    { get; set; }
    [JsonPropertyName("max_price")]     public double? MaxPrice    { get; set; }
    [JsonPropertyName("price")]         public double? Price       { get; set; }
    [JsonPropertyName("stock")]         public int    Stock        { get; set; }
    [JsonPropertyName("created_at")]    public DateTime CreatedAt  { get; set; }

    public double Rating    { get; set; }
    public int    TotalSold { get; set; }
}

public class CreateProductRequest
{
    [Required]
    [JsonPropertyName("store_id")]    public Guid   StoreId     { get; set; }

    [Required]
    [JsonPropertyName("name")]        public string Name        { get; set; } = string.Empty;

    [Required]
    [JsonPropertyName("description")] public string Description { get; set; } = string.Empty;

    [Required]
    [JsonPropertyName("category")]    public string Category    { get; set; } = string.Empty;

    public bool   HasVariants { get; set; }
    public double Price       { get; set; }
    public int    Stock       { get; set; }
}

public class CreateProductResponse
{
    [JsonPropertyName("id")]         public Guid     Id        { get; set; }
    [JsonPropertyName("store_id")]   public Guid     StoreId   { get; set; }
    [JsonPropertyName("name")]       public string   Name      { get; set; } = string.Empty;
    [JsonPropertyName("created_at")] public DateTime CreatedAt { get; set; }
}

public class UpdateProductDetailsRequest
{
    [Required]
    [JsonPropertyName("product_id")]  public Guid    ProductId   { get; set; }
    [JsonPropertyName("name")]        public string? Name        { get; set; }
    [JsonPropertyName("description")] public string? Description { get; set; }
    [JsonPropertyName("category")]    public string? Category    { get; set; }
}

public class UpdateProductInventoryRequest
{
    [Required]
    [JsonPropertyName("product_id")] public Guid   ProductId { get; set; }
    [Required]
    [JsonPropertyName("price")]      public double Price     { get; set; }
    [Required]
    [JsonPropertyName("stock")]      public int    Stock     { get; set; }
}

/// <summary>
/// Cascaded from MerchantLayout so the sidebar can build the dynamic
/// "Manage Product" link when the user is editing a specific product.
/// </summary>
public class ActiveProductContext
{
    public Guid?  ProductId   { get; set; }
    public string? ProductName { get; set; }
    public bool   HasValue    => ProductId.HasValue;
}
