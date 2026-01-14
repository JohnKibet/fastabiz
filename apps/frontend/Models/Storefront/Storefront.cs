using System;
using System.Collections.Generic;
using System.ComponentModel.DataAnnotations;
using System.Linq;
using System.Text.Json.Serialization;
namespace frontend.Models;

// Represents a Store
public class Store
{
    public Guid Id { get; set; }
    public string Name { get; set; } = string.Empty;
    public string Logo { get; set; } = string.Empty;
    public double Rating { get; set; }
    public int TotalProducts { get; set; }

    // Optional: Products found in Store (for Store detail page)
    public List<ProductX> Products { get; set; } = new();
}

public class CreateStoreRequest
{
    [Required]
    [JsonPropertyName("id")]
    public Guid Id { get; set; }

    [Required]
    [JsonPropertyName("name")]
    public string Name { get; set; } = string.Empty;

    [JsonPropertyName("logo_url")]
    public string Logo { get; set; } = string.Empty;

    [Required]
    [JsonPropertyName("location")]
    public string Location { get; set; } = string.Empty;
}

// Represents a product (variant or simple)
public class ProductX
{
    public Guid Id { get; set; }
    public Guid StoreId { get; set; }
    public string Name { get; set; } = string.Empty;
    public string Description { get; set; } = string.Empty;
    public string Category { get; set; } = string.Empty;

    public List<string> Images { get; set; } = new();

    public bool HasVariants { get; set; }

    // Simple products (HasVariants == false)
    public double Price { get; set; }
    public int Stock { get; set; }

    // Variant-enabled products
    public List<Option> Options { get; set; } = new();
    public List<Variant> Variants { get; set; } = new();

    // Optional frontend helpers
    public double? MinPrice => HasVariants && Variants.Any() ? Variants.Min(v => v.Price) : Price;
    public double? MaxPrice => HasVariants && Variants.Any() ? Variants.Max(v => v.Price) : Price;
}

public class CreateProductRequest
{
    [Required]
    [JsonPropertyName("store_id")]
    public Guid StoreId { get; set; }

    [Required]
    [JsonPropertyName("name")]
    public string Name { get; set; } = string.Empty;

    [Required]
    [JsonPropertyName("description")]
    public string Description { get; set; } = string.Empty;

    [Required]
    [JsonPropertyName("category")]
    public string Category { get; set; } = string.Empty;
}

public class UpdateProductDetailsRequest
{
    [Required]
    [JsonPropertyName("product_id")]
    public Guid ProductId { get; set; }

    [JsonPropertyName("name")]
    public string? Name { get; set; }

    [JsonPropertyName("description")]
    public string? Description { get; set; }

    [JsonPropertyName("category")]
    public string? Category { get; set; }
}

// Product option, e.g., Size, Weight, Type
public class Option
{
    public string Name { get; set; } = string.Empty;
    public List<string> Values { get; set; } = new();
}

public class CreateOptionNameRequest
{
    [Required]
    [JsonPropertyName("product_id")]
    public Guid ProductId { get; set; }

    [Required]
    [JsonPropertyName("name")]
    public string Name { get; set; } = string.Empty;
}

public class AddOptionValuesRequest
{
    [Required]
    [JsonPropertyName("product_id")]
    public Guid ProductId { get; set; }

    [Required]
    [JsonPropertyName("option_id")]
    public Guid OptionId { get; set; }

    [Required]
    [JsonPropertyName("values")]
    public List<string> Values { get; set; } = new();
}

public class Image
{
    public string URL { get; set; } = string.Empty;
    public bool IsPrimary { get; set; }
}

public class AddImageRequest
{
    [Required]
    [JsonPropertyName("product_id")]
    public Guid ProductId { get; set; }

    [Required]
    [JsonPropertyName("images")]
    public List<Image> Images { get; set; } = new();
}

public class ReorderImagesRequest
{
    [Required]
    [JsonPropertyName("product_id")]
    public Guid ProductId { get; set; }

    [Required]
    [JsonPropertyName("image_ids")]
    public List<Guid> ImageIDs { get; set; } = new();
}

// Variant of a product
public class Variant
{
    public Guid Id { get; set; }
    public string SKU { get; set; } = string.Empty;

    public double Price { get; set; }
    public int Stock { get; set; }

    public string Image { get; set; } = string.Empty;

    // Dictionary of option name → selected value
    public Dictionary<string, string> Options { get; set; } = new();
}

public class CreateVariantRequest
{
    [Required]
    [JsonPropertyName("product_id")]
    public Guid ProductId { get; set; }

    [Required]
    [JsonPropertyName("sku")]
    public string SKU { get; set; } = string.Empty;

    [Required]
    [JsonPropertyName("price")]
    public double Price { get; set; }

    [Required]
    [JsonPropertyName("stock")]
    public int Stock { get; set; }

    [JsonPropertyName("image_url")]
    public string Image { get; set; } = string.Empty;

    [Required]
    [JsonPropertyName("options")]
    public Dictionary<string, string> Options { get; set; } = new();
}

public class UpdateVariantStockRequest
{
    [Required]
    [JsonPropertyName("variant_id")]
    public Guid VariantId { get; set; }

    [Required]
    [JsonPropertyName("stock")]
    public int Stock { get; set; }
}

public class UpdateVariantPriceRequest
{
    [Required]
    [JsonPropertyName("variant_id")]
    public Guid VariantId { get; set; }

    [Required]
    [JsonPropertyName("price")]
    public double Price { get; set; }
}

// Lightweight version for list/grid views
public class ProductListItem
{
    public Guid Id { get; set; }
    public string Name { get; set; } = string.Empty;
    public string Category { get; set; } = string.Empty;
    public string Description { get; set; } = string.Empty;
    public string Thumbnail { get; set; } = string.Empty;

    public bool HasVariants { get; set; }
    public double? Price { get; set; }
    public double? MinPrice { get; set; }
    public double? MaxPrice { get; set; }

    public Guid StoreId { get; set; }
    public string StoreName { get; set; } = string.Empty;
    public double Rating { get; set; }
    public int TotalSold { get; set; }
}
