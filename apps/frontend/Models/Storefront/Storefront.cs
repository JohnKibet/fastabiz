namespace frontend.Models;

// Represents a Store
public class Store
{
    public string Id { get; set; } = string.Empty;
    public string Name { get; set; } = string.Empty;
    public string Logo { get; set; } = string.Empty;
    public double Rating { get; set; }
    public int TotalProducts { get; set; }

    // Optional: Products found in Store (for Store detail page)
    public List<ProductX> Products { get; set; } = new();
}

// Represents a product (variant or simple)
public class ProductX
{
    public string Id { get; set; } = string.Empty;
    public string StoreId { get; set; } = string.Empty;
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

// Product option, e.g., Size, Weight, Type
public class Option
{
    public string Name { get; set; } = string.Empty;
    public List<string> Values { get; set; } = new();
}

// Variant of a product
public class Variant
{
    public string Id { get; set; } = string.Empty;
    public string SKU { get; set; } = string.Empty;

    public double Price { get; set; }
    public int Stock { get; set; }

    public string Image { get; set; } = string.Empty;

    // Dictionary of option name → selected value
    public Dictionary<string, string> Options { get; set; } = new();
}

// Lightweight version for list/grid views
public class ProductListItem
{
    public string Id { get; set; } = string.Empty;
    public string Name { get; set; } = string.Empty;
    public string Category { get; set; } = string.Empty;
    public string Description { get; set; } = string.Empty;
    public string Thumbnail { get; set; } = string.Empty;

    public bool HasVariants { get; set; }
    public double? Price { get; set; }
    public double? MinPrice { get; set; }
    public double? MaxPrice { get; set; }

    public string StoreId { get; set; } = string.Empty;
    public string StoreName { get; set; } = string.Empty;
    public double Rating { get; set; }
    public int TotalSold { get; set; }
}
