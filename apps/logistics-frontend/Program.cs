using Microsoft.AspNetCore.Components.Web;
using Microsoft.AspNetCore.Components.WebAssembly.Hosting;
using Microsoft.AspNetCore.Components.Authorization;
using logistics_frontend;
using logistics_frontend.Services.CustomAuthStateProvider;
using System.Text.Json;

var builder = WebAssemblyHostBuilder.CreateDefault(args);
builder.RootComponents.Add<App>("#app");
builder.RootComponents.Add<HeadOutlet>("head::after");

// Use the host environment's base address to fetch appsettings.json safely
var tempHttpClient = new HttpClient
{
    BaseAddress = new Uri(builder.HostEnvironment.BaseAddress)
};

using var settingsResponse = await tempHttpClient.GetAsync("appsettings.json");
if (!settingsResponse.IsSuccessStatusCode)
{
    throw new Exception("Failed to load appsettings.json");
}

var json = await settingsResponse.Content.ReadAsStringAsync();

Dictionary<string, string>? config;
try
{
    config = JsonSerializer.Deserialize<Dictionary<string, string>>(json);
}
catch (Exception ex)
{
    throw new Exception("Failed to parse appsettings.json", ex);
}

if (config == null || !config.TryGetValue("ApiBaseUrl", out var apiBaseUrl) || string.IsNullOrWhiteSpace(apiBaseUrl))
{
    throw new Exception("ApiBaseUrl not found or invalid in appsettings.json");
}


// Now use this as the main HttpClient for API calls
builder.Services.AddScoped(sp => new HttpClient
{
    BaseAddress = new Uri(apiBaseUrl)
});

Console.WriteLine($"API Base URL set to: {apiBaseUrl}");

// Register services
builder.Services.AddScoped<UserService>();
builder.Services.AddScoped<OrderService>();
builder.Services.AddScoped<DriverService>();
builder.Services.AddScoped<PaymentService>();
builder.Services.AddScoped<DeliveryService>();
builder.Services.AddScoped<FeedbackService>();
builder.Services.AddScoped<UserSessionService>();
builder.Services.AddScoped<NotificationService>();

// Authorization
builder.Services.AddOptions();
builder.Services.AddAuthorizationCore();
builder.Services.AddScoped<AuthenticationStateProvider, CustomAuthStateProvider>();
builder.Services.AddScoped<CustomAuthStateProvider>();

await builder.Build().RunAsync();   
