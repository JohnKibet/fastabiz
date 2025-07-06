using System.ComponentModel.DataAnnotations;

namespace logistics_frontend.Models.Inventory
{
    public class Inventory
    {
        public Guid ID { get; set; }
        public Guid AdminID { get; set; }
        public string Name { get; set; } = string.Empty;
        public string? Slug { get; set; }
        public string Category { get; set; } = string.Empty;
        public int Stock { get; set; }
        public float Price { get; set; }
        public string Images { get; set; } = string.Empty;
        public string Unit { get; set; } = string.Empty;
        public string Packaging { get; set; } = string.Empty;
        public string Description { get; set; } = string.Empty;
        public string Location { get; set; } = string.Empty;
        public DateTime CreatedAt { get; set; }
        public DateTime UpdatedAt { get; set; }

    }

    public class CreateInventoryRequest
    {
        [Required]
        public Guid AdminID { get; set; }

        [Required]
        public string Name { get; set; } = string.Empty;

        [Required]
        public string? Slug { get; set; }

        [Required]
        public string Category { get; set; } = string.Empty;

        [Required]
        public int Stock { get; set; }

        [Required]
        public float Price { get; set; }

        [Required]
        public string Images { get; set; } = string.Empty;

        [Required]
        public string Unit { get; set; } = string.Empty;

        [Required]
        public string Packaging { get; set; } = string.Empty;

        [Required]
        public string Description { get; set; } = string.Empty;

        [Required]
        public string Location { get; set; } = string.Empty;

    }

    public class StorePublicView
    {
        public string AdminName { get; set; } = string.Empty;
        public List<Inventory> Products { get; set; } = new();
    }
}
