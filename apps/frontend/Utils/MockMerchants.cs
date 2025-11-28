using System.Collections.Generic;

namespace frontend.Models;

public static class MockData
{
    public static List<Merchant> Merchants = new()
    {
        new Merchant
        {
            Id = "m001",
            Name = "Green Valley Farms",
            Logo = "green_valley_logo.jpg",
            Rating = 4.7,
            TotalProducts = 3,
            Products = new List<ProductX>
            {
                new ProductX
                {
                    Id = "p001",
                    MerchantId = "m001",
                    Name = "Fresh Apples",
                    Category = "Fruits",
                    Description = "Crisp, sweet organic apples harvested this season.",
                    Images = new List<string> { "apple_main.jpg", "apple_2.jpg" },
                    HasVariants = true,
                    Options = new List<Option>
                    {
                        new Option { Name = "Size", Values = new List<string> { "Small", "Medium", "Large" } }
                    },
                    Variants = new List<Variant>
                    {
                        new Variant { Id = "v001", SKU = "APP-SM-001", Price = 1.20, Stock = 120, Image = "apple_small.jpg", Options = new Dictionary<string,string>{{"Size","Small"}} },
                        new Variant { Id = "v002", SKU = "APP-MD-001", Price = 1.50, Stock = 200, Image = "apple_medium.jpg", Options = new Dictionary<string,string>{{"Size","Medium"}} },
                        new Variant { Id = "v003", SKU = "APP-LG-001", Price = 1.80, Stock = 80, Image = "apple_large.jpg", Options = new Dictionary<string,string>{{"Size","Large"}} },
                    }
                },
                new ProductX
                {
                    Id = "p002",
                    MerchantId = "m001",
                    Name = "Bananas",
                    Category = "Fruits",
                    Description = "Sweet ripe bananas grown locally.",
                    Images = new List<string> { "banana.jpg" },
                    HasVariants = false,
                    Price = 0.80,
                    Stock = 300
                },
                new ProductX
                {
                    Id = "p003",
                    MerchantId = "m001",
                    Name = "Carrots",
                    Category = "Vegetables",
                    Description = "Fresh organic carrots.",
                    Images = new List<string> { "carrot.jpg" },
                    HasVariants = false,
                    Price = 0.50,
                    Stock = 150
                }
            }
        },

        new Merchant
        {
            Id = "m002",
            Name = "Sunrise Orchard",
            Logo = "sunrise_logo.jpg",
            Rating = 4.5,
            TotalProducts = 2,
            Products = new List<ProductX>
            {
                new ProductX
                {
                    Id = "p101",
                    MerchantId = "m002",
                    Name = "Oranges",
                    Category = "Fruits",
                    Description = "Juicy oranges rich in vitamin C.",
                    Images = new List<string> { "orange_main.jpg" },
                    HasVariants = true,
                    Options = new List<Option>
                    {
                        new Option { Name = "Type", Values = new List<string> { "Sweet", "Tangy" } }
                    },
                    Variants = new List<Variant>
                    {
                        new Variant { Id = "v101", SKU = "ORG-SW-101", Price = 0.95, Stock = 180, Image = "orange_sweet.jpg", Options = new Dictionary<string,string>{{"Type","Sweet"}} },
                        new Variant { Id = "v102", SKU = "ORG-TG-101", Price = 1.00, Stock = 150, Image = "orange_tangy.jpg", Options = new Dictionary<string,string>{{"Type","Tangy"}} },
                    }
                },
                new ProductX
                {
                    Id = "p102",
                    MerchantId = "m002",
                    Name = "Limes",
                    Category = "Fruits",
                    Description = "Fresh and zesty limes.",
                    Images = new List<string> { "lime.jpg" },
                    HasVariants = false,
                    Price = 0.60,
                    Stock = 200
                }
            }
        },

        new Merchant
        {
            Id = "m003",
            Name = "Highland Fresh Produce",
            Logo = "highland_logo.jpg",
            Rating = 4.6,
            TotalProducts = 2,
            Products = new List<ProductX>
            {
                new ProductX
                {
                    Id = "p201",
                    MerchantId = "m003",
                    Name = "Strawberries",
                    Category = "Berries",
                    Description = "Fresh, sweet strawberries delivered daily.",
                    Images = new List<string> { "strawberry.jpg" },
                    HasVariants = false,
                    Price = 3.50,
                    Stock = 60
                },
                new ProductX
                {
                    Id = "p202",
                    MerchantId = "m003",
                    Name = "Blueberries",
                    Category = "Berries",
                    Description = "Antioxidant-rich fresh blueberries.",
                    Images = new List<string> { "blueberry.jpg" },
                    HasVariants = false,
                    Price = 4.50,
                    Stock = 40
                }
            }
        },

        new Merchant
        {
            Id = "m004",
            Name = "Coastal Citrus Co.",
            Logo = "coastal_logo.jpg",
            Rating = 4.8,
            TotalProducts = 3,
            Products = new List<ProductX>
            {
                new ProductX
                {
                    Id = "p301",
                    MerchantId = "m004",
                    Name = "Lemons",
                    Category = "Citrus",
                    Description = "Fresh lemons with strong aroma.",
                    Images = new List<string> { "lemon.jpg" },
                    HasVariants = false,
                    Price = 0.60,
                    Stock = 250
                },
                new ProductX
                {
                    Id = "p302",
                    MerchantId = "m004",
                    Name = "Grapefruits",
                    Category = "Citrus",
                    Description = "Juicy grapefruits.",
                    Images = new List<string> { "grapefruit.jpg" },
                    HasVariants = false,
                    Price = 1.20,
                    Stock = 180
                },
                new ProductX
                {
                    Id = "p303",
                    MerchantId = "m004",
                    Name = "Oranges",
                    Category = "Fruits",
                    Description = "Tangy oranges from coastal farms.",
                    Images = new List<string> { "orange_coastal.jpg" },
                    HasVariants = false,
                    Price = 0.95,
                    Stock = 150
                }
            }
        },

        new Merchant
        {
            Id = "m005",
            Name = "Happy Harvest Market",
            Logo = "happy_harvest.jpg",
            Rating = 4.3,
            TotalProducts = 2,
            Products = new List<ProductX>
            {
                new ProductX
                {
                    Id = "p401",
                    MerchantId = "m005",
                    Name = "Pineapples",
                    Category = "Fruits",
                    Description = "Sweet tropical pineapples.",
                    Images = new List<string> { "pineapple.jpg" },
                    HasVariants = false,
                    Price = 2.80,
                    Stock = 90
                },
                new ProductX
                {
                    Id = "p402",
                    MerchantId = "m005",
                    Name = "Papayas",
                    Category = "Fruits",
                    Description = "Sweet and soft papayas.",
                    Images = new List<string> { "papaya.jpg" },
                    HasVariants = false,
                    Price = 2.20,
                    Stock = 70
                }
            }
        },

        new Merchant
        {
            Id = "m006",
            Name = "Farm Direct Foods",
            Logo = "farm_direct.jpg",
            Rating = 4.6,
            TotalProducts = 2,
            Products = new List<ProductX>
            {
                new ProductX
                {
                    Id = "p501",
                    MerchantId = "m006",
                    Name = "Grapes",
                    Category = "Fruits",
                    Description = "Seedless grapes, red or green.",
                    Images = new List<string> { "grapes.jpg" },
                    HasVariants = true,
                    Options = new List<Option>
                    {
                        new Option { Name = "Color", Values = new List<string> { "Red", "Green" } }
                    },
                    Variants = new List<Variant>
                    {
                        new Variant { Id = "v501", SKU = "GRP-RD-501", Price = 2.40, Stock = 110, Options = new Dictionary<string,string>{{"Color","Red"}} },
                        new Variant { Id = "v502", SKU = "GRP-GR-501", Price = 2.30, Stock = 95, Options = new Dictionary<string,string>{{"Color","Green"}} },
                    }
                },
                new ProductX
                {
                    Id = "p502",
                    MerchantId = "m006",
                    Name = "Avocados",
                    Category = "Fruits",
                    Description = "Fresh organic avocados.",
                    Images = new List<string> { "avocado.jpg" },
                    HasVariants = false,
                    Price = 1.60,
                    Stock = 180
                }
            }
        },

        new Merchant
        {
            Id = "m007",
            Name = "Tropical Delight Farm",
            Logo = "tropical_delight.jpg",
            Rating = 4.9,
            TotalProducts = 2,
            Products = new List<ProductX>
            {
                new ProductX
                {
                    Id = "p601",
                    MerchantId = "m007",
                    Name = "Mangoes",
                    Category = "Fruits",
                    Description = "Juicy ripe mangoes from coastal farms.",
                    Images = new List<string> { "mango.jpg" },
                    HasVariants = false,
                    Price = 1.90,
                    Stock = 200
                },
                new ProductX
                {
                    Id = "p602",
                    MerchantId = "m007",
                    Name = "Pineapples",
                    Category = "Fruits",
                    Description = "Sweet pineapples.",
                    Images = new List<string> { "pineapple.jpg" },
                    HasVariants = false,
                    Price = 2.80,
                    Stock = 100
                }
            }
        }
    };

    /// <summary>
    /// Flatten all products across all merchants into a ProductListItem view
    /// </summary>
    public static List<ProductListItem> GetAllProducts()
    {
        var products = new List<ProductListItem>();

        foreach (var merchant in Merchants)
        {
            foreach (var p in merchant.Products)
            {
                products.Add(new ProductListItem
                {
                    Id = p.Id,
                    Name = p.Name,
                    Category = p.Category,
                    Thumbnail = p.Images.FirstOrDefault() ?? "https://via.placeholder.com/150",
                    HasVariants = p.HasVariants,
                    Price = p.Price,
                    MinPrice = p.MinPrice,
                    MaxPrice = p.MaxPrice,
                    MerchantId = merchant.Id,
                    MerchantName = merchant.Name,
                    Rating = merchant.Rating,
                    TotalSold = new Random().Next(10, 200) // mock sold count
                });
            }
        }

        return products;
    }

    /// <summary>
    /// Get all unique categories from all products
    /// </summary>
    public static List<string> GetAllCategories()
    {
        return GetAllProducts()
            .Select(p => p.Category)
            .Distinct()
            .OrderBy(c => c)
            .ToList();
    }
}
