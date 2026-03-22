using System.Text.Json;
using frontend.Models.Storefront;
using Microsoft.JSInterop;

namespace frontend.Services;

/// <summary>
/// Singleton cart service — one instance for the entire app lifetime so
/// cart contents survive Blazor page navigations without re-fetching.
///
/// Persistence: serialised to localStorage under the key "cart".
/// Call LoadCartAsync() once on app startup (e.g. in MainLayout or App.razor)
/// to rehydrate from any previous session.
/// </summary>
public class CartService
{
    private readonly IJSRuntime _js;
    private const string StorageKey = "cart";

    private List<CartItem> _items = new();

    public CartService(IJSRuntime js) => _js = js;

    // ── State ──────────────────────────────────────────────────────────

    /// Raised whenever the cart changes — subscribe in components
    /// with CartService.OnChange += StateHasChanged.
    public event Action? OnChange;

    public IReadOnlyList<CartItem> Items => _items.AsReadOnly();

    /// Number of distinct line-items (not total quantity)
    public int Count => _items.Count;

    /// Sum of all line-item quantities
    public int TotalQuantity => _items.Sum(i => i.Quantity);

    /// Total price across all items
    public decimal Total => (decimal)_items.Sum(i => i.Price * i.Quantity);

    // ── Mutations ─────────────────────────────────────────────────────

    public async Task AddItem(ProductMapper.ProductDto product, ProductMapper.VariantDto? variant, int qty = 1)
    {
        Guid? variantId = variant?.Id;

        var existing = _items.FirstOrDefault(i =>
            i.ProductId == product.Id && i.VariantId == variantId);

        if (existing is not null)
        {
            existing.Quantity += qty;
        }
        else
        {
            _items.Add(new CartItem
            {
                StoreId     = product.StoreId,
                ProductId   = product.Id,
                VariantId   = variantId,
                Name        = product.Name,
                Description = product.Description,
                VariantName = variant?.Sku ?? string.Empty,
                SKU         = variant?.Sku,
                Price       = variant?.Price ?? 0,
                Thumbnail   = !string.IsNullOrEmpty(variant?.ImageUrl)
                                  ? variant.ImageUrl
                                  : product.Images?.FirstOrDefault() ?? "placeholder.jpg",
                Quantity    = qty,
            });
        }

        await PersistAndNotify();
    }

    public async Task RemoveItem(Guid productId, Guid? variantId)
    {
        var item = _items.FirstOrDefault(i =>
            i.ProductId == productId && i.VariantId == variantId);

        if (item is not null)
            _items.Remove(item);

        await PersistAndNotify();
    }

    public async Task IncrementItem(CartItem item)
    {
        var existing = _items.FirstOrDefault(i => i == item);
        if (existing is null) return;

        existing.Quantity++;
        await PersistAndNotify();
    }

    public async Task DecrementItem(CartItem item)
    {
        var existing = _items.FirstOrDefault(i => i == item);
        if (existing is null) return;

        existing.Quantity--;
        if (existing.Quantity <= 0)
            _items.Remove(existing);

        await PersistAndNotify();
    }

    public async Task UpdateQuantity(Guid productId, Guid? variantId, int qty)
    {
        var item = _items.FirstOrDefault(i =>
            i.ProductId == productId && i.VariantId == variantId);

        if (item is null) return;

        if (qty <= 0)
            _items.Remove(item);
        else
            item.Quantity = qty;

        await PersistAndNotify();
    }

    public void ClearCart()
    {
        _items.Clear();
        _ = PersistAndNotify(); // fire-and-forget; void method can't await
    }

    // ── Queries ───────────────────────────────────────────────────────

    public bool IsInCart(Guid productId, Guid? variantId) =>
        _items.Any(i => i.ProductId == productId && i.VariantId == variantId);

    public CartItem? GetItem(Guid productId, Guid? variantId) =>
        _items.FirstOrDefault(i => i.ProductId == productId && i.VariantId == variantId);

    // ── Persistence ───────────────────────────────────────────────────

    /// Call once during app initialisation to rehydrate from localStorage.
    public async Task LoadCartAsync()
    {
        try
        {
            var json = await _js.InvokeAsync<string?>("localStorage.getItem", StorageKey);
            if (!string.IsNullOrWhiteSpace(json))
                _items = JsonSerializer.Deserialize<List<CartItem>>(json) ?? new();
        }
        catch
        {
            // localStorage unavailable (SSR, private mode, etc.) — start with empty cart
            _items = new();
        }

        Notify();
    }

    public async Task SaveCartAsync()
    {
        try
        {
            await _js.InvokeVoidAsync(
                "localStorage.setItem", StorageKey,
                JsonSerializer.Serialize(_items));
        }
        catch
        {
            // Swallow — persistence failure should not crash the UI
        }
    }

    // ── Internals ─────────────────────────────────────────────────────

    private async Task PersistAndNotify()
    {
        await SaveCartAsync();
        Notify();
    }

    private void Notify() => OnChange?.Invoke();
}

/// <summary>
/// A single line-item in the cart.
/// Kept in the same file — it only exists to support CartService.
/// Move to Models/Cart.cs if it grows (e.g. discount fields, metadata).
/// </summary>
public class CartItem
{
    public Guid    StoreId     { get; set; }
    public Guid    ProductId   { get; set; }
    public Guid?   VariantId   { get; set; }
    public string  Name        { get; set; } = string.Empty;
    public string  VariantName { get; set; } = string.Empty;
    public string  Description { get; set; } = string.Empty;
    public string  Thumbnail   { get; set; } = string.Empty;
    public string? SKU         { get; set; }
    public double  Price       { get; set; }
    public int     Quantity    { get; set; }
}
