using System.Net.Http.Json;
using logistics_frontend.Models.Notification;

public class NotificationService
{
    private readonly HttpClient _http;

    public NotificationService(HttpClient http)
    {
        _http = http;
    }

    public async Task CreateNotification(CreateNotificationRequest notification)
    {
        var response = await _http.PostAsJsonAsync("notifications/create", notification);
        response.EnsureSuccessStatusCode();
    }

    public async Task<List<Notification>> GetNotificationByUser(Guid userId)
    {
        var notifications = await _http.GetFromJsonAsync<List<Notification>>($"notifications/{userId}");
        return notifications ?? new List<Notification>();
    }

    public async Task<List<Notification>> GetAllNotifications()
    {
        var notifications = await _http.GetFromJsonAsync<List<Notification>>("notifications/all_notifications");
        return notifications ?? new List<Notification>();
    }
}