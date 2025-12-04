using System;
using System.Collections.Generic;
using System.Linq;
using System.Text.Json;
using System.Threading.Tasks;
using frontend.Models;
using Microsoft.JSInterop;

public class CartService
{
    private readonly IJSRuntime _jSRuntime;
    public CartService(IJSRuntime jSRuntime)
    {
        _jSRuntime = jSRuntime;
    }
    public event Action? OnChange;

    private List<CartItem> items = new();

    public IReadOnlyList<CartItem> Items => items.AsReadOnly();
    public int Count => items.Sum(i => i.Quantity); // total number of products in cart

    public async Task AddItem(ProductX product, Variant? variant = null)
    {
        var existing = items.FirstOrDefault(i =>
            i.ProductId == product.Id &&
            i.VariantId == variant?.Id
        );

        if (existing != null)
        {
            existing.Quantity++;
        }
        else
        {
            items.Add(new CartItem
            {
                ProductId = product.Id,
                VariantId = variant?.Id,
                Name = product.Name,
                VariantName = variant != null ? string.Join(", ", variant.Options.Select(o => $"{o.Key}: {o.Value}")) : null,
                Price = variant?.Price ?? 0,
                Quantity = 1
            });
        }

        await SaveCartAsync();
        NotifyStateChanged();
    }

    public bool IsInCart(string productId, string? variantId)
    {
        return Items.Any(i =>
            i.ProductId == productId &&
            i.VariantId == variantId);
    }

    public async Task RemoveItem(string productId, string? variantId)
    {
        var item = Items.FirstOrDefault(i =>
            i.ProductId == productId &&
            i.VariantId == variantId);

        if (item != null)
            items.Remove(item);

        await SaveCartAsync();
        NotifyStateChanged();
    }

    public async Task DecrementItem(CartItem item)
    {
        var existing = items.FirstOrDefault(i => i == item);
        if (existing != null)
        {
            existing.Quantity--;
            if (existing.Quantity <= 0) items.Remove(existing);
            await SaveCartAsync();
            NotifyStateChanged();
        }
    }

    private void NotifyStateChanged() => OnChange?.Invoke();

    public async Task SaveCartAsync()
    {
        await _jSRuntime.InvokeVoidAsync("localStorage.setItem", "cart", JsonSerializer.Serialize(items));
    }

    public async Task LoadCartAsync()
    {
        var json = await _jSRuntime.InvokeAsync<string>("localStorage.getItem", "cart");
        if (!string.IsNullOrEmpty(json))
        {
            items = JsonSerializer.Deserialize<List<CartItem>>(json)!;
            OnChange?.Invoke();
        }
    }
}

public class CartItem
{
    public string ProductId { get; set; } = string.Empty;
    public string? VariantId { get; set; }
    public string Name { get; set; } = string.Empty;
    public string? VariantName { get; set; } // e.g., Size: M, Color: Red
    public double Price { get; set; }
    public int Quantity { get; set; }
    public string Thumbnail { get; set; } = string.Empty;
}
