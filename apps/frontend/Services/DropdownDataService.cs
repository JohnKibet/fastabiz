using frontend.Models;

public class DropdownDataService
{
    private readonly OrderService _orderService;
    private readonly ToastService _toastService;

    public List<Customer> Customers { get; private set; } = new();
    public List<AllInventory> Inventories { get; private set; } = new();

    private DateTime _lastFetchTime;
    private readonly TimeSpan _cacheDuration = TimeSpan.FromMinutes(10);
    public DropdownDataService(OrderService orderService, ToastService toastService)
    {
        _orderService = orderService;
        _toastService = toastService;
    }

    public async Task<ServiceResult2<bool>> LoadCachedDropdownData(bool forceRefresh = false)
    {
        if (!forceRefresh && Customers.Count > 0 && Inventories.Count > 0
            && DateTime.UtcNow - _lastFetchTime < _cacheDuration)
        {
            return ServiceResult2<bool>.Ok(true);
        }

        var result = await _orderService.GetDropdownMenuData();
        if (result.Success && result.Data != null)
        {
            Customers = result.Data.Customers ?? new();
            Inventories = result.Data.Inventories ?? new();
            _lastFetchTime = DateTime.UtcNow;            
            return ServiceResult2<bool>.Ok(true);
        }
        else if (result.Error != null )
        {    
            _toastService.ShowToast(result.Error.Detail, ToastService.ToastLevel.Warning);
            return ServiceResult2<bool>.Fail(result.Error);
        }
        else
        {
            return ServiceResult2<bool>.Fail("Failed to load dropdown data (unknown error).");
        }
    }

    public void InvalidateCache()
    {
        Customers.Clear();
        Inventories.Clear();
        _lastFetchTime = DateTime.MinValue;
    }
}