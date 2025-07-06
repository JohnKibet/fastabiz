using System.Net.Http.Json;
using logistics_frontend.Models.Inventory;
using Microsoft.JSInterop;

public class InventoryService
{
    private readonly HttpClient _http;
    public InventoryService(IHttpClientFactory httpClientFactory)
    {
        _http = httpClientFactory.CreateClient("AuthenticatedApi");
    }

    public async Task AddInventory(CreateInventoryRequest inventory)
    {
        var response = await _http.PostAsJsonAsync("inventories/create", inventory);
        response.EnsureSuccessStatusCode();
    }

    public async Task<List<Inventory>> GetInventoriesByID(Guid inventory_id)
    {
        var inventories = await _http.GetFromJsonAsync<List<Inventory>>($"inventories/by-id?id={inventory_id}");
        return inventories ?? new List<Inventory>();
    }

    public async Task<List<Inventory>> GetInventoriesByName(string name)
    {
        var encodedName = Uri.EscapeDataString(name);
        var inventories = await _http.GetFromJsonAsync<List<Inventory>>($"inventories/by-name?name={encodedName}");
        return inventories ?? new List<Inventory>();
    }

    public async Task<List<Inventory>> GetAllInventories()
    {
        var inventories = await _http.GetFromJsonAsync<List<Inventory>>("inventories/all_inventories");
        return inventories ?? new List<Inventory>();
    }

    public async Task<List<Inventory>> GetInventoriesByCategory(string category)
    {
        var encodedCategory = Uri.EscapeDataString(category);
        var inventories = await _http.GetFromJsonAsync<List<Inventory>>($"inventories/by-category?category={encodedCategory}");
        return inventories ?? new List<Inventory>();
    }

    public async Task<List<string>> GetCategories()
    {
        var categories = await _http.GetFromJsonAsync<List<string>>("inventories/categories");
        return categories ?? new List<string>();
    }
}