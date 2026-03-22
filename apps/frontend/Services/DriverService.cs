using frontend.Models;
using frontend.Models.Storefront;

namespace frontend.Services;

public class DriverService
{
    private readonly ApiService _api;

    public DriverService(ApiService api)
    {
        _api = api;
    }

    public Task<ApiResult<ApiMessageResponse>> RegisterDriver(CreateDriverRequest req)
        => _api.PostAsync<CreateDriverRequest, ApiMessageResponse>("drivers/create", req);

    public Task<ApiResult<Driver>> GetDriverById(Guid driverId)
        => _api.GetAsync<Driver>($"drivers/{driverId}");

    public Task<ApiResult<Driver>> GetDriverByEmail(string email)
        => _api.GetAsync<Driver>($"drivers/{email}");

    public Task<ApiResult<List<Driver>>> GetAllDrivers()
        => _api.GetAsync<List<Driver>>("drivers/all_drivers");
}
