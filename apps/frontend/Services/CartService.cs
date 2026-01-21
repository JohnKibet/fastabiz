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
    public int Count => items.Count();

    public async Task AddItem(ProductX product, Variant? variant, int qty = 1 )
    {
        var variantId = variant?.Id ?? null;

        var existing = items.FirstOrDefault(i =>
            i.ProductId == product.Id &&
            i.VariantId == variantId
        );

        double price = (variant != null && variant.Price > 0) ? variant.Price : product.Price;
        string thumbnail = (variant != null && !string.IsNullOrEmpty(variant.ImageUrl)) 
                        ? variant.ImageUrl 
                        : product.Images.FirstOrDefault() ?? "placeholder.jpg";

        if (existing != null)
        {
            existing.Quantity += qty;
        }
        else
        {
            items.Add(new CartItem
            {
                ProductId = product.Id,
                VariantId = variantId!.Value,
                Name = product.Name,
                VariantName = variant != null 
                    ? string.Join(", ", variant.Options.Select(o => $"{o.Key}: {o.Value}")) 
                    : null,
                Price = price,
                Thumbnail = thumbnail,
                Quantity = qty,
                SKU = variant?.SKU ?? null,
                StoreId = product.StoreId,
                Description = product.Description,
            });
        }

        await SaveCartAsync();
        NotifyStateChanged();
    }

    public bool IsInCart(Guid productId, Guid? variantId)
    {
        return Items.Any(i =>
            i.ProductId == productId &&
            i.VariantId == variantId);
    }

    public async Task RemoveItem(Guid productId, Guid? variantId)
    {
        var item = Items.FirstOrDefault(i =>
            i.ProductId == productId &&
            i.VariantId == variantId);

        if (item != null)
            items.Remove(item);

        await SaveCartAsync();
        NotifyStateChanged();
    }

    public async Task IncrementItem(CartItem item)
    {
        var existing = items.FirstOrDefault(i => i == item);
        if (existing != null)
        {
            existing.Quantity++;
            await SaveCartAsync();
            NotifyStateChanged();
        }
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
            items = JsonSerializer.Deserialize<List<CartItem>>(json)!;

        NotifyStateChanged();
    }
}

public class CartItem
{
    public Guid ProductId { get; set; }
    public Guid VariantId { get; set; }
    public string Name { get; set; } = string.Empty;
    public string? VariantName { get; set; }
    public string Thumbnail { get; set; } = string.Empty;
    public double Price { get; set; }
    public int Quantity { get; set; }
    public string Description { get; set; } = string.Empty;

    // Optional but recommended:
    public Guid StoreId { get; set; }
    public string? SKU { get; set; }
}
