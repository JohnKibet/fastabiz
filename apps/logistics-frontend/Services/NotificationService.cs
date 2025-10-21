using System.Net.Http.Json;
using logistics_frontend.Models.Notification;

public class NotificationService
{
    private readonly HttpClient _http;
    private readonly ToastService _toastService;
    public NotificationService(IHttpClientFactory httpClientFactory, ToastService toastService)
    {
        _http = httpClientFactory.CreateClient("AuthenticatedApi");
        _toastService = toastService;
    }

    // User-triggered notifications (send feedback, message another user, etc.)
    public async Task CreateNotification(CreateNotificationRequest notification)
    {
        var response = await _http.PostAsJsonAsync("notifications/create", notification);
        response.EnsureSuccessStatusCode();
    }

    public async Task<List<Notification>> GetNotificationByUser(Guid userId)
    {
        var notifications = await _http.GetFromJsonAsync<List<Notification>>($"notifications/all_my_notifications/{userId}");
        return notifications ?? new List<Notification>();
    }

    // mark single read
    public async Task<Notification> UpdateNotificationStatus(Guid id, NotificationStatus status)
    {
        var response = await _http.PutAsJsonAsync($"notifications/{id}/status", status);
        if (response.IsSuccessStatusCode)
        {
            return await response.Content.ReadFromJsonAsync<Notification>() ?? new Notification();
        }

        return null;
    }

    public async Task MarkAllAsRead(Guid userId)
    {
        var response = await _http.PatchAsync($"notifications/mark_all_as_read/{userId}", null);
        if (response.IsSuccessStatusCode)
        {
            _toastService.ShowToast("Marked Read !", ToastService.ToastLevel.Success);
        }
        else
        {
            _toastService.ShowToast("Mark All Read Failed !", ToastService.ToastLevel.Error);
        }
    }
}