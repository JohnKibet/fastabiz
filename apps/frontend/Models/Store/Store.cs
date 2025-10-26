using System.ComponentModel.DataAnnotations;
using System.Text.Json.Serialization;
using InvModel = frontend.Models.Inventory.Inventory;

namespace frontend.Models.Store
{
    public class Store
    {
        [JsonPropertyName("id")]
        public Guid ID { get; set; }

        [JsonPropertyName("owner_id")]
        public Guid OwnerID { get; set; }

        [JsonPropertyName("name")]
        public string Name { get; set; } = string.Empty;

        [JsonPropertyName("slug")]
        public string Slug { get; set; } = string.Empty;

        [JsonPropertyName("description")]
        public string Description { get; set; } = string.Empty;

        [JsonPropertyName("logo_url")]
        public string? LogoURL { get; set; }

        [JsonPropertyName("banner_url")]
        public string? BannerURL { get; set; }

        [JsonPropertyName("is_public")]
        public bool IsPublic { get; set; }

        [JsonPropertyName("location")]
        public string Location { get; set; } = string.Empty;

        [JsonPropertyName("created_at")]
        public DateTime CreatedAt { get; set; }

        [JsonPropertyName("Updated_at")]
        public DateTime UpdatedAt { get; set; }

        // Optional: if API returns inventories for this store
        [JsonPropertyName("inventories")]
        public List<InvModel>? Inventories { get; set; }
    }

    public class CreateStoreRequest
    {
        [Required(ErrorMessage = "Owner ID is required")]
        [JsonPropertyName("owner_id")]
        public Guid OwnerID { get; set; }

        [Required(ErrorMessage = "Name is required")]
        [JsonPropertyName("name")]
        public string Name { get; set; } = string.Empty;

        [Required(ErrorMessage = "Slug is required")]
        [JsonPropertyName("slug")]
        public string Slug { get; set; } = string.Empty;

        [JsonPropertyName("description")]
        public string Description { get; set; } = string.Empty;

        [JsonPropertyName("logo_url")]
        public string? LogoURL { get; set; }

        [JsonPropertyName("banner_url")]
        public string? BannerURL { get; set; }

        [JsonPropertyName("is_public")]
        public bool IsPublic { get; set; } = true;

        [JsonPropertyName("location")]
        public string Location { get; set; } = string.Empty;
    }

    // Optional for lightweight previews
    public class StoreSummary
    {
        [JsonPropertyName("id")]
        public Guid ID { get; set; }

        [JsonPropertyName("name")]
        public string Name { get; set; } = string.Empty;

        [JsonPropertyName("slug")]
        public string Slug { get; set; } = string.Empty;

        [JsonPropertyName("logo_url")]
        public string? LogoURL { get; set; }
    }

    public class Product
    {
        public Guid ID { get; set; }
        public string Name { get; set; } = string.Empty;
        public string Price { get; set; } = string.Empty;
        public string Image { get; set; } = string.Empty;
        public string Stock { get; set; } = string.Empty;
        public string Description { get; set; } = string.Empty;
    }

}