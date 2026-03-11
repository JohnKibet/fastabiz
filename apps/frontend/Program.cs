using System.Text.Json;
using Microsoft.Extensions.DependencyInjection;
using Microsoft.AspNetCore.Components.Web;
using Microsoft.AspNetCore.Components.WebAssembly.Hosting;
using Microsoft.AspNetCore.Components.Authorization;
using frontend;
using frontend.Services;
using MudBlazor.Services;

var builder = WebAssemblyHostBuilder.CreateDefault(args);
builder.RootComponents.Add<App>("#app");
builder.RootComponents.Add<HeadOutlet>("head::after");

var tempHttpClient = new HttpClient
{
    BaseAddress = new Uri(builder.HostEnvironment.BaseAddress)
};

using var settingsResponse = await tempHttpClient.GetAsync("appsettings.json");
if (!settingsResponse.IsSuccessStatusCode)
    throw new Exception("Failed to load appsettings.json");

var json = await settingsResponse.Content.ReadAsStringAsync();
var config = JsonSerializer.Deserialize<Dictionary<string, string>>(json)
    ?? throw new Exception("Failed to parse appsettings.json");

if (!config.TryGetValue("ApiBaseUrl", out var apiBaseUrl) || string.IsNullOrWhiteSpace(apiBaseUrl))
    throw new Exception("ApiBaseUrl not found or invalid in appsettings.json");

Console.WriteLine($"API Base URL set to: {apiBaseUrl}");

builder.Services.AddMudServices();

// This guarantees SignInAsync/SignOutAsync and
// AuthorizeView all talk to the same object in memory.
builder.Services.AddScoped<CustomAuthStateProvider>();
builder.Services.AddScoped<AuthenticationStateProvider>(sp =>
    sp.GetRequiredService<CustomAuthStateProvider>());

builder.Services.AddScoped<AuthHeaderHandler>();
builder.Services.AddScoped<AuthExpiredHandler>();

builder.Services.AddHttpClient("AnonymousApi", client =>
{
    client.BaseAddress = new Uri(apiBaseUrl);
});

builder.Services.AddHttpClient("AuthenticatedApi", client =>
{
    client.BaseAddress = new Uri(apiBaseUrl);
})
.AddHttpMessageHandler<AuthHeaderHandler>()
.AddHttpMessageHandler<AuthExpiredHandler>();

builder.Services.AddHttpClient("CloudinaryClient");

builder.Services.AddScoped<UserService>();
builder.Services.AddScoped<OrderService>();
builder.Services.AddScoped<DriverService>();
builder.Services.AddScoped<PaymentService>();
builder.Services.AddScoped<DeliveryService>();
builder.Services.AddScoped<FeedbackService>();
builder.Services.AddScoped<UserSessionService>();
builder.Services.AddScoped<NotificationService>();
builder.Services.AddScoped<ToastService>();
builder.Services.AddScoped<StoreService>();
builder.Services.AddScoped<ProductService>();
builder.Services.AddSingleton<CartService>();
builder.Services.AddSingleton<MapService>();
builder.Services.AddSingleton<GeoService>();

builder.Services.AddOptions();
builder.Services.AddAuthorizationCore();

await builder.Build().RunAsync();
