using System;
using System.Collections.Generic;
using System.ComponentModel.DataAnnotations;
using System.Text.Json.Serialization;

namespace frontend.Models.Storefront;

public class Store
{
    public Guid Id { get; set; }
    public string Name { get; set; } = string.Empty;
    public string LogoUrl { get; set; } = string.Empty;
    public string Location { get; set; } = string.Empty;
    public double Rating { get; set; }
    public int? TotalProducts { get; set; }
    public List<Product> Products { get; set; } = new();
    public DateTime CreatedAt { get; set; }
    public DateTime UpdatedAt { get; set; }
}

public class StoreDto
{
    [JsonPropertyName("id")] public Guid Id { get; set; }
    [JsonPropertyName("merchant_id")] public Guid OwnerId { get; set; }
    [JsonPropertyName("name")] public string Name { get; set; } = string.Empty;
    [JsonPropertyName("logo_url")] public string LogoUrl { get; set; } = string.Empty;
    [JsonPropertyName("location")] public string Location { get; set; } = string.Empty;
    [JsonPropertyName("created_at")] public string CreatedAt { get; set; } = string.Empty;
    [JsonPropertyName("updated_at")] public string UpdatedAt { get; set; } = string.Empty;
    public double? Rating { get; set; }
    public int? TotalProducts { get; set; }
}

public class CreateStoreRequest
{
    [Required]
    [JsonPropertyName("id")] public Guid Id { get; set; }

    [Required]
    [JsonPropertyName("name")] public string Name { get; set; } = string.Empty;

    [JsonPropertyName("logo_url")] public string Logo { get; set; } = string.Empty;

    [Required]
    [JsonPropertyName("location")] public string Location { get; set; } = string.Empty;
}

public class CreateStoreResponse
{
    [JsonPropertyName("id")] public Guid Id { get; set; }
    [JsonPropertyName("store_id")] public Guid OwnerId { get; set; }
    [JsonPropertyName("name")] public string Name { get; set; } = string.Empty;
    [JsonPropertyName("logo_url")] public string LogoURL { get; set; } = string.Empty;
    [JsonPropertyName("location")] public string Location { get; set; } = string.Empty;
    [JsonPropertyName("created_at")] public DateTime CreatedAt { get; set; }
    [JsonPropertyName("updated_at")] public DateTime UpdatedAt { get; set; }

}

public class UpdateStoreRequest
{
    [JsonPropertyName("name")] public string? Name { get; set; }
    [JsonPropertyName("logo_url")] public string? LogoUrl { get; set; }
    [JsonPropertyName("location")] public string? Location { get; set; }
}

/// <summary>
/// Cascaded from MerchantLayout down to any component that needs the
/// currently active store (sidebar, product forms, etc.)
/// </summary>
public sealed class ActiveStoreContext
{
    [JsonPropertyName("id")] public Guid StoreId { get; set; }
    [JsonPropertyName("name")] public string StoreName { get; set; } = string.Empty;
}
