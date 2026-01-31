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
    public string LogoUrl { get; set; } = string.Empty;
    public string Location { get; set; } = string.Empty;
    public double Rating { get; set; }
    public int? TotalProducts { get; set; }

    // Optional: Products found in Store (for Store detail page)
    public List<ProductX> Products { get; set; } = new();
    public DateTime CreatedAt { get; set; }
    public DateTime UpdatedAt { get; set; }
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

public class UpdateStoreRequest
{
    [JsonPropertyName("name")]
    public string? Name { get; set; }

    [JsonPropertyName("logo_url")]
    public string? LogoUrl { get; set; }

    [JsonPropertyName("location")]
    public string? Location { get; set; }
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

    public DateTime CreatedAt { get; set; }
    public DateTime UpdatedAt { get; set; }
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

    public bool HasVariants { get; set; }
}

// createproduct response
public class CreateProductResponse
{
    [JsonPropertyName("id")]
    public Guid Id { get; set; }

    [JsonPropertyName("store_id")]
    public Guid StoreId { get; set; }

    [JsonPropertyName("name")]
    public string Name { get; set; } = string.Empty;

    [JsonPropertyName("created_at")]
    public DateTime CreatedAt { get; set; }
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

public class CreateOptionNameResponse
{
    [JsonPropertyName("option_id")]
    public Guid OptionId { get; set; }
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
    public Guid ProductId { get; set; }
    public string SKU { get; set; } = string.Empty;
    public double Price { get; set; }
    public int Stock { get; set; }
    public string ImageUrl { get; set; } = string.Empty;
    public DateTime CreatedAt { get; set; }
    public DateTime UpdatedAt { get; set; }

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
public class CreateVariantResponse
{
    [JsonPropertyName("id")]
    public Guid Id { get; set; }
    [JsonPropertyName("product_id")]
    public Guid ProductId { get; set; }
    [JsonPropertyName("sku")]
    public string SKU { get; set; } = string.Empty;
    [JsonPropertyName("price")]
    public double Price { get; set; }
    [JsonPropertyName("stock")]
    public int stock { get; set; }
    [JsonPropertyName("image_url")]
    public string ImageURL { get; set; } = string.Empty;
    [JsonPropertyName("options")]
    public Dictionary<string, string> Options { get; set; } = new(); // Size → Small
    [JsonPropertyName("created_at")]
    public DateTime CreatedAt { get; set; }
    [JsonPropertyName("updated_at")]
    public DateTime UpdatedAt { get; set; }
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
    [JsonPropertyName("id")]
    public Guid Id { get; set; }

    [JsonPropertyName("store_id")]
    public Guid StoreId { get; set; }

    [JsonPropertyName("name")]
    public string Name { get; set; } = string.Empty;

    [JsonPropertyName("category")]
    public string Category { get; set; } = string.Empty;

    [JsonPropertyName("description")]
    public string Description { get; set; } = string.Empty;

    [JsonPropertyName("image_url")]
    public string ImageUrl { get; set; } = string.Empty;

    [JsonPropertyName("has_variants")]
    public bool HasVariants { get; set; }

    // simple variants
    [JsonPropertyName("variant_count")]
    public int VariantCount { get; set; }

    [JsonPropertyName("min_price")]
    public double? MinPrice { get; set; }

    [JsonPropertyName("max_price")]
    public double? MaxPrice { get; set; }

    // products
    [JsonPropertyName("price")]
    public double? Price { get; set; }

    [JsonPropertyName("stock")]
    public int Stock { get; set; }

    [JsonPropertyName("created_at")]
    public DateTime CreatedAt { get; set; }

    public double Rating { get; set; }
    public int TotalSold { get; set; }
}

// getstoresbyowner json mapper
public sealed class ActiveStoreContext
{
    [JsonPropertyName("id")]
    public Guid StoreId { get; set; }

    [JsonPropertyName("name")]
    public string StoreName { get; set; } = string.Empty;
    // public bool HasValue => StoreId.HasValue && StoreId != Guid.Empty;
}
public class ActiveProductContext
{
    public Guid? ProductId { get; set; }
    public string? ProductName { get; set; }
    public bool HasValue => ProductId.HasValue;
}

// getproductbyid json mapper
public static class ProductMapper
{
    public class ProductDto
    {
        [JsonPropertyName("id")]
        public Guid Id { get; set; }

        [JsonPropertyName("store_id")]
        public Guid StoreId { get; set; }

        [JsonPropertyName("name")]
        public string Name { get; set; } = string.Empty;

        [JsonPropertyName("description")]
        public string Description { get; set; } = string.Empty;

        [JsonPropertyName("category")]
        public string Category { get; set; } = string.Empty;

        [JsonPropertyName("images")]
        public List<string>? Images { get; set; } = new();

        [JsonPropertyName("has_variants")]
        public bool HasVariants { get; set; }

        public List<VariantDto> Variants { get; set; } = new();

        [JsonPropertyName("created_at")]
        public DateTime CreatedAt { get; set; }

        [JsonPropertyName("updated_at")]
        public DateTime UpdatedAt { get; set; }
    }
    public class VariantDto
    {
        [JsonPropertyName("id")]
        public Guid Id { get; set; }

        [JsonPropertyName("product_id")]
        public Guid ProductId { get; set; }

        [JsonPropertyName("sku")]
        public string Sku { get; set; } = string.Empty;

        [JsonPropertyName("price")]
        public double Price { get; set; }

        [JsonPropertyName("stock")]
        public int Stock { get; set; }

        [JsonPropertyName("image_url")]
        public string? ImageUrl { get; set; }

        [JsonPropertyName("created_at")]
        public DateTime CreatedAt { get; set; }

        [JsonPropertyName("updated_at")]
        public DateTime UpdatedAt { get; set; }
    }
    public static ProductDto ToDto(this ProductX product)
    {
        return new ProductDto
        {
            Id = product.Id,
            StoreId = product.StoreId,
            Name = product.Name,
            Description = product.Description,
            Category = product.Category,
            Images = product.Images, // assuming it's List<string> or null
            HasVariants = product.HasVariants,
            Variants = product.Variants?.Select(v => new VariantDto
            {
                Id = v.Id,
                ProductId = v.ProductId,
                Sku = v.SKU,
                Price = v.Price,
                Stock = v.Stock,
                ImageUrl = v.ImageUrl,
                CreatedAt = v.CreatedAt,
                UpdatedAt = v.UpdatedAt
            }).ToList() ?? new List<VariantDto>(),
            CreatedAt = product.CreatedAt,
            UpdatedAt = product.UpdatedAt
        };
    }
}
