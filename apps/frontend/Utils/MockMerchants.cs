using System.Collections.Generic;

namespace frontend.Models;

public static class MockData
{
    public static List<Store> Stores = new()
    {
        new Store
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
                    StoreId = "m001",
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
                    StoreId = "m001",
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
                    StoreId = "m001",
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
                    StoreId = "m002",
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
                    StoreId = "m002",
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
                    StoreId = "m003",
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
                    StoreId = "m003",
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
                    StoreId = "m004",
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
                    StoreId = "m004",
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
                    StoreId = "m004",
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
                    StoreId = "m005",
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
                    StoreId = "m005",
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
                    StoreId = "m006",
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
                    StoreId = "m006",
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
                    StoreId = "m007",
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
                    StoreId = "m007",
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
            Id = "m008",
            Name = "TechWorld Electronics",
            Logo = "techworld.jpg",
            Rating = 4.4,
            TotalProducts = 3,
            Products = new List<ProductX> {
                new ProductX {
                    Id = "p701",
                    StoreId = "m008",
                    Name = "Wireless Headphones",
                    Category = "Audio",
                    Description = "Bluetooth noise-canceling headphones.",
                    Images = new List<string>{ "headphones_main.jpg" },
                    HasVariants = true,
                    Options = new List<Option>{
                        new Option { Name = "Color", Values = new List<string>{ "Black", "White" } }
                    },
                    Variants = new List<Variant>{
                        new Variant { Id = "v701", SKU = "HD-BLK-701", Price = 59.99, Stock = 40, Options = new Dictionary<string,string>{{"Color","Black"}} },
                        new Variant { Id = "v702", SKU = "HD-WHT-701", Price = 59.99, Stock = 30, Options = new Dictionary<string,string>{{"Color","White"}} }
                    }
                },
                new ProductX {
                    Id = "p702",
                    StoreId = "m008",
                    Name = "Smartphone Charger",
                    Category = "Accessories",
                    Description = "Fast-charging USB-C charger.",
                    Images = new List<string>{ "charger.jpg" },
                    HasVariants = false,
                    Price = 15.99,
                    Stock = 120
                },
                new ProductX {
                    Id = "p703",
                    StoreId = "m008",
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
            Id = "m009",
            Name = "UrbanStyle Clothing",
            Logo = "urbanstyle.jpg",
            Rating = 4.2,
            TotalProducts = 2,
            Products = new List<ProductX> {
                new ProductX {
                    Id = "p801",
                    StoreId = "m009",
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
                        new Variant { Id = "v801", SKU="TS-M-BLK", Price=12.99, Stock=100, Options = new(){{"Size","M"},{"Color","Black"}}},
                        new Variant { Id = "v802", SKU="TS-L-GRY", Price=12.99, Stock=80, Options = new(){{"Size","L"},{"Color","Gray"}}}
                    }
                },
                new ProductX {
                    Id = "p802",
                    StoreId = "m009",
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
            Id = "m010",
            Name = "HomeLiving Furniture",
            Logo = "homeliving.jpg",
            Rating = 4.7,
            TotalProducts = 2,
            Products = new List<ProductX>{
                new ProductX {
                    Id = "p901",
                    StoreId = "m010",
                    Name = "Wooden Coffee Table",
                    Category = "Furniture",
                    Description = "Solid oak handcrafted table.",
                    Images = new List<string>{ "coffee_table.jpg" },
                    HasVariants = false,
                    Price = 130.00,
                    Stock = 20
                },
                new ProductX {
                    Id = "p902",
                    StoreId = "m010",
                    Name = "Office Chair",
                    Category = "Furniture",
                    Description = "Ergonomic office chair with lumbar support.",
                    Images = new List<string>{ "office_chair.jpg" },
                    HasVariants = true,
                    Options = new List<Option>{
                        new Option { Name="Color", Values = new(){"Black","Brown"} }
                    },
                    Variants = new List<Variant>{
                        new Variant { Id="v901", SKU="CH-BLK-01", Price=89.99, Stock=40, Options=new(){{"Color","Black"}}},
                        new Variant { Id="v902", SKU="CH-BRN-01", Price=89.99, Stock=30, Options=new(){{"Color","Brown"}}}
                    }
                }
            }
        },

        new Store 
        {
            Id = "m011",
            Name = "Glow Beauty Essentials",
            Logo = "glowbeauty.jpg",
            Rating = 4.5,
            TotalProducts = 3,
            Products = new List<ProductX>{
                new ProductX {
                    Id = "p1001",
                    StoreId = "m011",
                    Name = "Face Serum",
                    Category = "Skincare",
                    Description = "Vitamin C brightening serum.",
                    Images = new List<string>{ "serum.jpg" },
                    HasVariants = false,
                    Price = 18.99,
                    Stock = 150
                },
                new ProductX {
                    Id = "p1002",
                    StoreId = "m011",
                    Name = "Lipstick",
                    Category = "Cosmetics",
                    Description = "Matte finish lipstick.",
                    Images = new List<string>{ "lipstick.jpg" },
                    HasVariants = true,
                    Options = new List<Option>{
                        new Option { Name="Shade", Values=new(){"Red","Nude","Rose"}}
                    },
                    Variants = new List<Variant>{
                        new Variant{ Id="v1001", SKU="LS-RED-01", Price=7.99, Stock=60, Options=new(){{"Shade","Red"}}},
                        new Variant{ Id="v1002", SKU="LS-NUD-01", Price=7.99, Stock=70, Options=new(){{"Shade","Nude"}}}
                    }
                },
                new ProductX {
                    Id = "p1003",
                    StoreId = "m011",
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
            Id = "m012",
            Name = "PowerFit Sports",
            Logo = "powerfit.jpg",
            Rating = 4.6,
            TotalProducts = 2,
            Products = new List<ProductX>{
                new ProductX {
                    Id = "p1101",
                    StoreId = "m012",
                    Name = "Yoga Mat",
                    Category = "Fitness",
                    Description = "Non-slip thick yoga mat.",
                    Images = new List<string>{ "yogamat.jpg" },
                    HasVariants = true,
                    Options = new(){ new Option{Name="Color", Values=new(){"Purple","Blue","Black"}} },
                    Variants = new List<Variant>{
                        new Variant{ Id="v1101", SKU="YM-PUR-01", Price=22.99, Stock=80, Options=new(){{"Color","Purple"}}},
                        new Variant{ Id="v1102", SKU="YM-BLU-01", Price=22.99, Stock=90, Options=new(){{"Color","Blue"}}}
                    }
                },
                new ProductX {
                    Id = "p1102",
                    StoreId = "m012",
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
            Id = "m013",
            Name = "PetCare Supplies",
            Logo = "petcare.jpg",
            Rating = 4.4,
            TotalProducts = 2,
            Products = new List<ProductX>{
                new ProductX{
                    Id="p1201",
                    StoreId="m013",
                    Name="Dog Food - Chicken Flavor",
                    Category="Pet Food",
                    Description="Nutritious dry dog food.",
                    Images=new(){"dog_food.jpg"},
                    HasVariants=false,
                    Price=24.50,
                    Stock=90
                },
                new ProductX{
                    Id="p1202",
                    StoreId="m013",
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
            Id="m014",
            Name="AutoGear Shop",
            Logo="autogear.jpg",
            Rating=4.3,
            TotalProducts=2,
            Products=new(){
                new ProductX{
                    Id="p1301",
                    StoreId="m014",
                    Name="Car Vacuum Cleaner",
                    Category="Car Accessories",
                    Description="Portable car cleaner.",
                    Images=new(){"car_vacuum.jpg"},
                    HasVariants=false,
                    Price=35.99,
                    Stock=60
                },
                new ProductX{
                    Id="p1302",
                    StoreId="m014",
                    Name="Car Seat Covers",
                    Category="Car Accessories",
                    Description="Waterproof leather covers.",
                    Images=new(){"seat_covers.jpg"},
                    HasVariants=true,
                    Options=new(){ new Option { Name="Color", Values=new(){"Black","Beige"} } },
                    Variants=new(){
                        new Variant{Id="v1301", SKU="SC-BLK-01", Price=59.99, Stock=20, Options=new(){{"Color","Black"}}},
                        new Variant{Id="v1302", SKU="SC-BGE-01", Price=59.99, Stock=30, Options=new(){{"Color","Beige"}}}
                    }
                }
            }
        },

        new Store 
        {
            Id="m015",
            Name="PaperTrail Books",
            Logo="papertrail.jpg",
            Rating=4.8,
            TotalProducts=2,
            Products=new(){
                new ProductX{
                    Id="p1401",
                    StoreId="m015",
                    Name="Mystery Novel",
                    Category="Books",
                    Description="Thrilling detective story.",
                    Images=new(){"mystery_book.jpg"},
                    HasVariants=false,
                    Price=14.99,
                    Stock=150
                },
                new ProductX{
                    Id="p1402",
                    StoreId="m015",
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
            Id="m016",
            Name="Organic Health Store",
            Logo="organic.jpg",
            Rating=4.7,
            TotalProducts=2,
            Products=new(){
                new ProductX{
                    Id="p1501",
                    StoreId="m016",
                    Name="Organic Honey",
                    Category="Food",
                    Description="Raw unprocessed honey.",
                    Images=new(){"honey.jpg"},
                    HasVariants=false,
                    Price=12.00,
                    Stock=80
                },
                new ProductX{
                    Id="p1502",
                    StoreId="m016",
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
            Id="m017",
            Name="KitchenPro Supplies",
            Logo="kitchenpro.jpg",
            Rating=4.4,
            TotalProducts=2,
            Products=new(){
                new ProductX{
                    Id="p1601",
                    StoreId="m017",
                    Name="Stainless Steel Knife Set",
                    Category="Kitchen",
                    Description="Professional chef knives.",
                    Images=new(){"knife_set.jpg"},
                    HasVariants=false,
                    Price=39.99,
                    Stock=50
                },
                new ProductX{
                    Id="p1602",
                    StoreId="m017",
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
            Id="m018",
            Name="HandyTools Hardware",
            Logo="handytools.jpg",
            Rating=4.5,
            TotalProducts=2,
            Products=new(){
                new ProductX{
                    Id="p1701",
                    StoreId="m018",
                    Name="Electric Drill",
                    Category="Tools",
                    Description="Cordless drill set.",
                    Images=new(){"drill.jpg"},
                    HasVariants=false,
                    Price=49.99,
                    Stock=40
                },
                new ProductX{
                    Id="p1702",
                    StoreId="m018",
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
            Id="m019",
            Name="FunLand Toys",
            Logo="funland.jpg",
            Rating=4.3,
            TotalProducts=2,
            Products=new(){
                new ProductX{
                    Id="p1801",
                    StoreId="m019",
                    Name="Remote Control Car",
                    Category="Toys",
                    Description="High-speed RC car.",
                    Images=new(){"rc_car.jpg"},
                    HasVariants=false,
                    Price=29.99,
                    Stock=70
                },
                new ProductX{
                    Id="p1802",
                    StoreId="m019",
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
            Id="m020",
            Name="SweetTreats Bakery",
            Logo="sweettreats.jpg",
            Rating=4.9,
            TotalProducts=2,
            Products=new(){
                new ProductX{
                    Id="p1901",
                    StoreId="m020",
                    Name="Cupcake Box",
                    Category="Bakery",
                    Description="Box of 6 assorted cupcakes.",
                    Images=new(){"cupcakes.jpg"},
                    HasVariants=false,
                    Price=15.00,
                    Stock=60
                },
                new ProductX{
                    Id="p1902",
                    StoreId="m020",
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
