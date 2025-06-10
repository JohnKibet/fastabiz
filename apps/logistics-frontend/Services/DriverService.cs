using System.Net.Http.Json;
using logistics_frontend.Models.Driver;

public class DriverService
{
    private readonly HttpClient _http;

    public DriverService(HttpClient http)
    {
        _http = http;
    }

    public async Task RegisterDriver(CreateDriverRequest driver)
    {
        var response = await _http.PostAsJsonAsync("drivers/create", driver);
        response.EnsureSuccessStatusCode();
    }

    public async Task<Driver> GetDriverById(Guid driverId)
    {
        var driver = await _http.GetFromJsonAsync<Driver>($"drivers/id/{driverId}");
        return driver ?? throw new Exception("No driver found");
    }

    public async Task<List<Driver>> GetAllDrivers()
    {
        var drivers = await _http.GetFromJsonAsync<List<Driver>>("drivers/all_drivers");
        return drivers ?? new List<Driver>();
    }
}