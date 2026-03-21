using System.Text.Json;
using Microsoft.AspNetCore.Components.Web;
using Microsoft.AspNetCore.Components.WebAssembly.Hosting;
using Microsoft.AspNetCore.Components.Authorization;
using frontend;
using frontend.Services;
using MudBlazor.Services;

var builder = WebAssemblyHostBuilder.CreateDefault(args);
builder.RootComponents.Add<App>("#app");
builder.RootComponents.Add<HeadOutlet>("head::after");

// ── 1. Load runtime config ────────────────────────────────────────────
// Must run before any service registration that depends on the API URL.
// Uses a raw HttpClient (not DI) because the container isn't built yet.

var tempHttpClient = new HttpClient
{
    BaseAddress = new Uri(builder.HostEnvironment.BaseAddress)
};

using var settingsResponse = await tempHttpClient.GetAsync("appsettings.json");
if (!settingsResponse.IsSuccessStatusCode)
    throw new Exception("Failed to load appsettings.json");

var json   = await settingsResponse.Content.ReadAsStringAsync();
var config = JsonSerializer.Deserialize<Dictionary<string, string>>(json)
             ?? throw new Exception("Failed to parse appsettings.json");

if (!config.TryGetValue("ApiBaseUrl", out var apiBaseUrl) || string.IsNullOrWhiteSpace(apiBaseUrl))
    throw new Exception("ApiBaseUrl not found or invalid in appsettings.json");

Console.WriteLine($"API Base URL: {apiBaseUrl}");

// ── 2. Authorization & authentication ────────────────────────────────
// AddOptions() and AddAuthorizationCore() first — they are prerequisites
// for the auth state provider registrations below.
//
// CustomAuthStateProvider is registered ONCE as itself, then a factory
// resolves AuthenticationStateProvider from that same scoped instance.
// Two separate AddScoped<> calls would create two instances per scope —
// SignInAsync would update one but AuthorizeView would read the other,
// causing the login redirect loop.

builder.Services.AddOptions();
builder.Services.AddAuthorizationCore();

builder.Services.AddScoped<CustomAuthStateProvider>();
builder.Services.AddScoped<AuthenticationStateProvider>(sp =>
    sp.GetRequiredService<CustomAuthStateProvider>());

// ── 3. HTTP message handlers — before the named clients that use them ─
// DelegatingHandlers must be registered before the HttpClient pipelines
// that reference them via AddHttpMessageHandler<T>().

builder.Services.AddScoped<AuthHeaderHandler>();
builder.Services.AddScoped<AuthExpiredHandler>();

// ── 4. Named HTTP clients ─────────────────────────────────────────────
// AnonymousApi    — public endpoints (login, register). No auth header.
// AuthenticatedApi — protected endpoints. Handler order matters:
//                    AuthHeaderHandler attaches the Bearer token first,
//                    AuthExpiredHandler then intercepts 401 responses.
// CloudinaryClient — direct Cloudinary upload. No base address needed;
//                    ProductService constructs the full URL at call time.

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

// ── 5. MudBlazor ─────────────────────────────────────────────────────

builder.Services.AddMudServices();

// ── 6. Infrastructure services ────────────────────────────────────────
// Generic layers consumed by all domain services below.
// ApiService and ToastService have no domain dependencies so they
// come before the services that inject them.

builder.Services.AddScoped<ApiService>();
builder.Services.AddScoped<ToastService>();

// ── 7. Domain services ────────────────────────────────────────────────
// Scoped  — one instance per user session. Correct for services that
//           hold user-specific state or call authenticated endpoints.
// Singleton — one instance for the entire app lifetime. Used for
//             services with shared global state (cart persists across
//             page navigations) or stateless utilities (map, geo).
//
// Alphabetical within each lifetime group for easy scanning.

builder.Services.AddScoped<DeliveryService>();
builder.Services.AddScoped<DriverService>();
builder.Services.AddScoped<FeedbackService>();
builder.Services.AddScoped<NotificationService>();
builder.Services.AddScoped<OrderService>();
// builder.Services.AddScoped<PaymentService>();
builder.Services.AddScoped<ProductService>();
builder.Services.AddScoped<StoreService>();
builder.Services.AddScoped<UserService>();
builder.Services.AddScoped<UserSessionService>();

builder.Services.AddSingleton<CartService>();
builder.Services.AddSingleton<GeoService>();
builder.Services.AddSingleton<MapService>();

await builder.Build().RunAsync();
