using logistics_frontend.Models.User;
using logistics_frontend.Services.AuthHeaderHandler;
using System.Net.Http.Json;

public class UserService
{
    private readonly HttpClient _http;
    private readonly ILogger<UserService> _logger;
    public UserService(HttpClient http, ILogger<UserService> logger)
    {
        _http = http;
        _logger = logger;

        _logger.LogInformation("UserService constructed.");
    }
    public async Task<List<User>> GetAllUsers()
    {
        _logger.LogInformation("UserService: GetAllUsers() called.");
        var users = await _http.GetFromJsonAsync<List<User>>("users/all_users");
        return users ?? new List<User>();
    }
}
