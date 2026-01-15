using System.Collections.Generic;

namespace frontend.Models;

public static class MockData
{
    public static List<Store> Stores = new()
    {
        new Store
        {
            Id = Guid.NewGuid(),
            Name = "Green Valley Farms",
            LogoUrl = "green_valley_logo.jpg",
            Rating = 4.7,
            TotalProducts = 3,
            Products = new List<ProductX>
            {
                new ProductX
                {
                    Id = Guid.NewGuid(),
                    StoreId = Guid.NewGuid(),
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
                        new Variant { Id = Guid.NewGuid(), SKU = "APP-SM-001", Price = 1.20, Stock = 120, Image = "apple_small.jpg", Options = new Dictionary<string,string>{{"Size","Small"}} },
                        new Variant { Id = Guid.NewGuid(), SKU = "APP-MD-001", Price = 1.50, Stock = 200, Image = "apple_medium.jpg", Options = new Dictionary<string,string>{{"Size","Medium"}} },
                        new Variant { Id = Guid.NewGuid(), SKU = "APP-LG-001", Price = 1.80, Stock = 80, Image = "apple_large.jpg", Options = new Dictionary<string,string>{{"Size","Large"}} },
                    }
                },
                new ProductX
                {
                    Id = Guid.NewGuid(),
                    StoreId = Guid.NewGuid(),
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
                    Id = Guid.NewGuid(),
                    StoreId = Guid.NewGuid(),
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

        new Store
        {
            Id = Guid.NewGuid(),
            Name = "Sunrise Orchard",
            LogoUrl = "sunrise_logo.jpg",
            Rating = 4.5,
            TotalProducts = 2,
            Products = new List<ProductX>
            {
                new ProductX
                {
                    Id = Guid.NewGuid(),
                    StoreId = Guid.NewGuid(),
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
                        new Variant { Id = Guid.NewGuid(), SKU = "ORG-SW-101", Price = 0.95, Stock = 180, Image = "orange_sweet.jpg", Options = new Dictionary<string,string>{{"Type","Sweet"}} },
                        new Variant { Id = Guid.NewGuid(), SKU = "ORG-TG-101", Price = 1.00, Stock = 150, Image = "orange_tangy.jpg", Options = new Dictionary<string,string>{{"Type","Tangy"}} },
                    }
                },
                new ProductX
                {
                    Id = Guid.NewGuid(),
                    StoreId = Guid.NewGuid(),
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

        new Store
        {
            Id = Guid.NewGuid(),
            Name = "Highland Fresh Produce",
            LogoUrl = "highland_logo.jpg",
            Rating = 4.6,
            TotalProducts = 2,
            Products = new List<ProductX>
            {
                new ProductX
                {
                    Id = Guid.NewGuid(),
                    StoreId = Guid.NewGuid(),
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
                    Id = Guid.NewGuid(),
                    StoreId = Guid.NewGuid(),
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

        new Store
        {
            Id = Guid.NewGuid(),
            Name = "Coastal Citrus Co.",
            LogoUrl = "coastal_logo.jpg",
            Rating = 4.8,
            TotalProducts = 3,
            Products = new List<ProductX>
            {
                new ProductX
                {
                    Id = Guid.NewGuid(),
                    StoreId = Guid.NewGuid(),
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
                    Id = Guid.NewGuid(),
                    StoreId = Guid.NewGuid(),
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
                    Id = Guid.NewGuid(),
                    StoreId = Guid.NewGuid(),
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

        new Store
        {
            Id = Guid.NewGuid(),
            Name = "Happy Harvest Market",
            LogoUrl = "happy_harvest.jpg",
            Rating = 4.3,
            TotalProducts = 2,
            Products = new List<ProductX>
            {
                new ProductX
                {
                    Id = Guid.NewGuid(),
                    StoreId = Guid.NewGuid(),
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
                    Id = Guid.NewGuid(),
                    StoreId = Guid.NewGuid(),
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

        new Store
        {
            Id = Guid.NewGuid(),
            Name = "Farm Direct Foods",
            LogoUrl = "farm_direct.jpg",
            Rating = 4.6,
            TotalProducts = 2,
            Products = new List<ProductX>
            {
                new ProductX
                {
                    Id = Guid.NewGuid(),
                    StoreId = Guid.NewGuid(),
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
                        new Variant { Id = Guid.NewGuid(), SKU = "GRP-RD-501", Price = 2.40, Stock = 110, Options = new Dictionary<string,string>{{"Color","Red"}} },
                        new Variant { Id = Guid.NewGuid(), SKU = "GRP-GR-501", Price = 2.30, Stock = 95, Options = new Dictionary<string,string>{{"Color","Green"}} },
                    }
                },
                new ProductX
                {
                    Id = Guid.NewGuid(),
                    StoreId = Guid.NewGuid(),
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

        new Store
        {
            Id = Guid.NewGuid(),
            Name = "Tropical Delight Farm",
            LogoUrl = "tropical_delight.jpg",
            Rating = 4.9,
            TotalProducts = 2,
            Products = new List<ProductX>
            {
                new ProductX
                {
                    Id = Guid.NewGuid(),
                    StoreId = Guid.NewGuid(),
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
                    Id = Guid.NewGuid(),
                    StoreId = Guid.NewGuid(),
                    Name = "Pineapples",
                    Category = "Fruits",
                    Description = "Sweet pineapples.",
                    Images = new List<string> { "pineapple.jpg" },
                    HasVariants = false,
                    Price = 2.80,
                    Stock = 100
                }
            }
        },

        new Store
        {
            Id = Guid.NewGuid(),
            Name = "TechWorld Electronics",
            LogoUrl = "techworld.jpg",
            Rating = 4.4,
            TotalProducts = 3,
            Products = new List<ProductX> {
                new ProductX {
                    Id = Guid.NewGuid(),
                    StoreId = Guid.NewGuid(),
                    Name = "Wireless Headphones",
                    Category = "Audio",
                    Description = "Bluetooth noise-canceling headphones.",
                    Images = new List<string>{ "headphones_main.jpg" },
                    HasVariants = true,
                    Options = new List<Option>{
                        new Option { Name = "Color", Values = new List<string>{ "Black", "White" } }
                    },
                    Variants = new List<Variant>{
                        new Variant { Id = Guid.NewGuid(), SKU = "HD-BLK-701", Price = 59.99, Stock = 40, Options = new Dictionary<string,string>{{"Color","Black"}} },
                        new Variant { Id = Guid.NewGuid(), SKU = "HD-WHT-701", Price = 59.99, Stock = 30, Options = new Dictionary<string,string>{{"Color","White"}} }
                    }
                },
                new ProductX {
                    Id = Guid.NewGuid(),
                    StoreId = Guid.NewGuid(),
                    Name = "Smartphone Charger",
                    Category = "Accessories",
                    Description = "Fast-charging USB-C charger.",
                    Images = new List<string>{ "charger.jpg" },
                    HasVariants = false,
                    Price = 15.99,
                    Stock = 120
                },
                new ProductX {
                    Id = Guid.NewGuid(),
                    StoreId = Guid.NewGuid(),
                    Name = "4K LED TV",
                    Category = "TV & Home",
                    Description = "Ultra HD 50-inch smart TV.",
                    Images = new List<string>{ "tv_4k.jpg" },
                    HasVariants = false,
                    Price = 450.00,
                    Stock = 15
                }
            }
        },

        new Store
        {
            Id = Guid.NewGuid(),
            Name = "UrbanStyle Clothing",
            LogoUrl = "urbanstyle.jpg",
            Rating = 4.2,
            TotalProducts = 2,
            Products = new List<ProductX> {
                new ProductX {
                    Id = Guid.NewGuid(),
                    StoreId = Guid.NewGuid(),
                    Name = "Men's T-Shirts",
                    Category = "Apparel",
                    Description = "100% cotton premium tees.",
                    Images = new List<string>{ "tshirt_main.jpg" },
                    HasVariants = true,
                    Options = new List<Option> {
                        new Option { Name = "Size", Values = new List<string>{ "S","M","L","XL" } },
                        new Option { Name = "Color", Values = new List<string>{ "Black","Blue","Gray" } }
                    },
                    Variants = new List<Variant>{
                        new Variant { Id = Guid.NewGuid(), SKU="TS-M-BLK", Price=12.99, Stock=100, Options = new(){{"Size","M"},{"Color","Black"}}},
                        new Variant { Id = Guid.NewGuid(), SKU="TS-L-GRY", Price=12.99, Stock=80, Options = new(){{"Size","L"},{"Color","Gray"}}}
                    }
                },
                new ProductX {
                    Id = Guid.NewGuid(),
                    StoreId = Guid.NewGuid(),
                    Name = "Women's Jeans",
                    Category = "Apparel",
                    Description = "Stretch denim, premium quality.",
                    Images = new List<string>{ "jeans.jpg" },
                    HasVariants = false,
                    Price = 25.50,
                    Stock = 75
                }
            }
        },

        new Store
        {
            Id = Guid.NewGuid(),
            Name = "HomeLiving Furniture",
            LogoUrl = "homeliving.jpg",
            Rating = 4.7,
            TotalProducts = 2,
            Products = new List<ProductX>{
                new ProductX {
                    Id = Guid.NewGuid(),
                    StoreId = Guid.NewGuid(),
                    Name = "Wooden Coffee Table",
                    Category = "Furniture",
                    Description = "Solid oak handcrafted table.",
                    Images = new List<string>{ "coffee_table.jpg" },
                    HasVariants = false,
                    Price = 130.00,
                    Stock = 20
                },
                new ProductX {
                    Id = Guid.NewGuid(),
                    StoreId = Guid.NewGuid(),
                    Name = "Office Chair",
                    Category = "Furniture",
                    Description = "Ergonomic office chair with lumbar support.",
                    Images = new List<string>{ "office_chair.jpg" },
                    HasVariants = true,
                    Options = new List<Option>{
                        new Option { Name="Color", Values = new(){"Black","Brown"} }
                    },
                    Variants = new List<Variant>{
                        new Variant { Id=Guid.NewGuid(), SKU="CH-BLK-01", Price=89.99, Stock=40, Options=new(){{"Color","Black"}}},
                        new Variant { Id=Guid.NewGuid(), SKU="CH-BRN-01", Price=89.99, Stock=30, Options=new(){{"Color","Brown"}}}
                    }
                }
            }
        },

        new Store
        {
            Id = Guid.NewGuid(),
            Name = "Glow Beauty Essentials",
            LogoUrl = "glowbeauty.jpg",
            Rating = 4.5,
            TotalProducts = 3,
            Products = new List<ProductX>{
                new ProductX {
                    Id = Guid.NewGuid(),
                    StoreId = Guid.NewGuid(),
                    Name = "Face Serum",
                    Category = "Skincare",
                    Description = "Vitamin C brightening serum.",
                    Images = new List<string>{ "serum.jpg" },
                    HasVariants = false,
                    Price = 18.99,
                    Stock = 150
                },
                new ProductX {
                    Id = Guid.NewGuid(),
                    StoreId = Guid.NewGuid(),
                    Name = "Lipstick",
                    Category = "Cosmetics",
                    Description = "Matte finish lipstick.",
                    Images = new List<string>{ "lipstick.jpg" },
                    HasVariants = true,
                    Options = new List<Option>{
                        new Option { Name="Shade", Values=new(){"Red","Nude","Rose"}}
                    },
                    Variants = new List<Variant>{
                        new Variant{ Id=Guid.NewGuid(), SKU="LS-RED-01", Price=7.99, Stock=60, Options=new(){{"Shade","Red"}}},
                        new Variant{ Id=Guid.NewGuid(), SKU="LS-NUD-01", Price=7.99, Stock=70, Options=new(){{"Shade","Nude"}}}
                    }
                },
                new ProductX {
                    Id = Guid.NewGuid(),
                    StoreId = Guid.NewGuid(),
                    Name = "Face Cleanser",
                    Category = "Skincare",
                    Description = "Gentle daily cleansing foam.",
                    Images = new List<string>{ "cleanser.jpg" },
                    HasVariants = false,
                    Price = 9.99,
                    Stock = 110
                }
            }
        },

        new Store
        {
            Id = Guid.NewGuid(),
            Name = "PowerFit Sports",
            LogoUrl = "powerfit.jpg",
            Rating = 4.6,
            TotalProducts = 2,
            Products = new List<ProductX>{
                new ProductX {
                    Id = Guid.NewGuid(),
                    StoreId = Guid.NewGuid(),
                    Name = "Yoga Mat",
                    Category = "Fitness",
                    Description = "Non-slip thick yoga mat.",
                    Images = new List<string>{ "yogamat.jpg" },
                    HasVariants = true,
                    Options = new(){ new Option{Name="Color", Values=new(){"Purple","Blue","Black"}} },
                    Variants = new List<Variant>{
                        new Variant{ Id=Guid.NewGuid(), SKU="YM-PUR-01", Price=22.99, Stock=80, Options=new(){{"Color","Purple"}}},
                        new Variant{ Id=Guid.NewGuid(), SKU="YM-BLU-01", Price=22.99, Stock=90, Options=new(){{"Color","Blue"}}}
                    }
                },
                new ProductX {
                    Id = Guid.NewGuid(),
                    StoreId = Guid.NewGuid(),
                    Name = "Dumbbell Set",
                    Category = "Fitness",
                    Description = "20kg adjustable dumbbells.",
                    Images = new List<string>{ "dumbbells.jpg" },
                    HasVariants = false,
                    Price = 45.00,
                    Stock = 50
                }
            }
        },

        new Store
        {
            Id = Guid.NewGuid(),
            Name = "PetCare Supplies",
            LogoUrl = "petcare.jpg",
            Rating = 4.4,
            TotalProducts = 2,
            Products = new List<ProductX>{
                new ProductX{
                    Id=Guid.NewGuid(),
                    StoreId=Guid.NewGuid(),
                    Name="Dog Food - Chicken Flavor",
                    Category="Pet Food",
                    Description="Nutritious dry dog food.",
                    Images=new(){"dog_food.jpg"},
                    HasVariants=false,
                    Price=24.50,
                    Stock=90
                },
                new ProductX{
                    Id=Guid.NewGuid(),
                    StoreId=Guid.NewGuid(),
                    Name="Cat Toy Set",
                    Category="Pet Accessories",
                    Description="Interactive toys for cats.",
                    Images=new(){"cat_toy.jpg"},
                    HasVariants=false,
                    Price=12.00,
                    Stock=140
                }
            }
        },

        new Store
        {
            Id=Guid.NewGuid(),
            Name="AutoGear Shop",
            LogoUrl="autogear.jpg",
            Rating=4.3,
            TotalProducts=2,
            Products=new(){
                new ProductX{
                    Id=Guid.NewGuid(),
                    StoreId=Guid.NewGuid(),
                    Name="Car Vacuum Cleaner",
                    Category="Car Accessories",
                    Description="Portable car cleaner.",
                    Images=new(){"car_vacuum.jpg"},
                    HasVariants=false,
                    Price=35.99,
                    Stock=60
                },
                new ProductX{
                    Id=Guid.NewGuid(),
                    StoreId=Guid.NewGuid(),
                    Name="Car Seat Covers",
                    Category="Car Accessories",
                    Description="Waterproof leather covers.",
                    Images=new(){"seat_covers.jpg"},
                    HasVariants=true,
                    Options=new(){ new Option { Name="Color", Values=new(){"Black","Beige"} } },
                    Variants=new(){
                        new Variant{Id=Guid.NewGuid(), SKU="SC-BLK-01", Price=59.99, Stock=20, Options=new(){{"Color","Black"}}},
                        new Variant{Id=Guid.NewGuid(), SKU="SC-BGE-01", Price=59.99, Stock=30, Options=new(){{"Color","Beige"}}}
                    }
                }
            }
        },

        new Store
        {
            Id=Guid.NewGuid(),
            Name="PaperTrail Books",
            LogoUrl="papertrail.jpg",
            Rating=4.8,
            TotalProducts=2,
            Products=new(){
                new ProductX{
                    Id=Guid.NewGuid(),
                    StoreId=Guid.NewGuid(),
                    Name="Mystery Novel",
                    Category="Books",
                    Description="Thrilling detective story.",
                    Images=new(){"mystery_book.jpg"},
                    HasVariants=false,
                    Price=14.99,
                    Stock=150
                },
                new ProductX{
                    Id=Guid.NewGuid(),
                    StoreId=Guid.NewGuid(),
                    Name="Notebook Set",
                    Category="Stationery",
                    Description="Pack of 3 premium notebooks.",
                    Images=new(){"notebooks.jpg"},
                    HasVariants=false,
                    Price=9.99,
                    Stock=200
                }
            }
        },

        new Store
        {
            Id=Guid.NewGuid(),
            Name="Organic Health Store",
            LogoUrl="organic.jpg",
            Rating=4.7,
            TotalProducts=2,
            Products=new(){
                new ProductX{
                    Id=Guid.NewGuid(),
                    StoreId=Guid.NewGuid(),
                    Name="Organic Honey",
                    Category="Food",
                    Description="Raw unprocessed honey.",
                    Images=new(){"honey.jpg"},
                    HasVariants=false,
                    Price=12.00,
                    Stock=80
                },
                new ProductX{
                    Id=Guid.NewGuid(),
                    StoreId=Guid.NewGuid(),
                    Name="Herbal Tea Mix",
                    Category="Beverages",
                    Description="Relaxing herbal infusion.",
                    Images=new(){"tea.jpg"},
                    HasVariants=false,
                    Price=7.50,
                    Stock=100
                }
            }
        },

        new Store
        {
            Id=Guid.NewGuid(),
            Name="KitchenPro Supplies",
            LogoUrl="kitchenpro.jpg",
            Rating=4.4,
            TotalProducts=2,
            Products=new(){
                new ProductX{
                    Id=Guid.NewGuid(),
                    StoreId=Guid.NewGuid(),
                    Name="Stainless Steel Knife Set",
                    Category="Kitchen",
                    Description="Professional chef knives.",
                    Images=new(){"knife_set.jpg"},
                    HasVariants=false,
                    Price=39.99,
                    Stock=50
                },
                new ProductX{
                    Id=Guid.NewGuid(),
                    StoreId=Guid.NewGuid(),
                    Name="Cutting Board",
                    Category="Kitchen",
                    Description="Bamboo cutting board.",
                    Images=new(){"cutting_board.jpg"},
                    HasVariants=false,
                    Price=12.99,
                    Stock=120
                }
            }
        },

        new Store
        {
            Id=Guid.NewGuid(),
            Name="HandyTools Hardware",
            LogoUrl="handytools.jpg",
            Rating=4.5,
            TotalProducts=2,
            Products=new(){
                new ProductX{
                    Id=Guid.NewGuid(),
                    StoreId=Guid.NewGuid(),
                    Name="Electric Drill",
                    Category="Tools",
                    Description="Cordless drill set.",
                    Images=new(){"drill.jpg"},
                    HasVariants=false,
                    Price=49.99,
                    Stock=40
                },
                new ProductX{
                    Id=Guid.NewGuid(),
                    StoreId=Guid.NewGuid(),
                    Name="Hammer",
                    Category="Tools",
                    Description="Steel claw hammer.",
                    Images=new(){"hammer.jpg"},
                    HasVariants=false,
                    Price=9.99,
                    Stock=150
                }
            }
        },

        new Store
        {
            Id=Guid.NewGuid(),
            Name="FunLand Toys",
            LogoUrl="funland.jpg",
            Rating=4.3,
            TotalProducts=2,
            Products=new(){
                new ProductX{
                    Id=Guid.NewGuid(),
                    StoreId=Guid.NewGuid(),
                    Name="Remote Control Car",
                    Category="Toys",
                    Description="High-speed RC car.",
                    Images=new(){"rc_car.jpg"},
                    HasVariants=false,
                    Price=29.99,
                    Stock=70
                },
                new ProductX{
                    Id=Guid.NewGuid(),
                    StoreId=Guid.NewGuid(),
                    Name="Puzzle Set",
                    Category="Kids",
                    Description="500-piece puzzle.",
                    Images=new(){"puzzle.jpg"},
                    HasVariants=false,
                    Price=8.99,
                    Stock=200
                }
            }
        },

        new Store
        {
            Id=Guid.NewGuid(),
            Name="SweetTreats Bakery",
            LogoUrl="sweettreats.jpg",
            Rating=4.9,
            TotalProducts=2,
            Products=new(){
                new ProductX{
                    Id=Guid.NewGuid(),
                    StoreId=Guid.NewGuid(),
                    Name="Cupcake Box",
                    Category="Bakery",
                    Description="Box of 6 assorted cupcakes.",
                    Images=new(){"cupcakes.jpg"},
                    HasVariants=false,
                    Price=15.00,
                    Stock=60
                },
                new ProductX{
                    Id=Guid.NewGuid(),
                    StoreId=Guid.NewGuid(),
                    Name="Chocolate Cake",
                    Category="Bakery",
                    Description="Rich chocolate cake.",
                    Images=new(){"chocolate_cake.jpg"},
                    HasVariants=false,
                    Price=25.00,
                    Stock=40
                }
            }
        },
    };

    /// <summary>
    /// Flatten all products across all Stores into a ProductListItem view
    /// </summary>
    public static List<ProductListItem> GetAllProducts()
    {
        var products = new List<ProductListItem>();

        foreach (var Store in Stores)
        {
            foreach (var p in Store.Products)
            {
                products.Add(new ProductListItem
                {
                    Id = p.Id,
                    Name = p.Name,
                    Category = p.Category,
                    Description = p.Description,
                    Thumbnail = p.Images.FirstOrDefault() ?? "https://via.placeholder.com/150",
                    HasVariants = p.HasVariants,
                    Price = p.Price,
                    MinPrice = p.MinPrice,
                    MaxPrice = p.MaxPrice,
                    StoreId = Store.Id,
                    StoreName = Store.Name,
                    Rating = Store.Rating,
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
