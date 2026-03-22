using frontend.Models;
using frontend.Models.Storefront;

namespace frontend.Services;

public class NotificationService
{
    private readonly ApiService   _api;
    private readonly ToastService _toast;

    public NotificationService(ApiService api, ToastService toast)
    {
        _api   = api;
        _toast = toast;
    }

    /// <summary>Send a notification (user-triggered: feedback, message, etc.)</summary>
    public Task<ApiResult<ApiMessageResponse>> CreateNotification(CreateNotificationRequest req)
        => _api.PostAsync<CreateNotificationRequest, ApiMessageResponse>("notifications/create", req);

    /// <summary>Fetch all notifications for a given user.</summary>
    public Task<ApiResult<List<Notification>>> GetNotificationByUser(Guid userId)
        => _api.GetAsync<List<Notification>>($"notifications/all_my_notifications/{userId}");

    /// <summary>Mark a single notification as read/unread.</summary>
    public Task<ApiResult<Notification>> UpdateNotificationStatus(Guid id, NotificationStatus status)
        => _api.PutAsync<NotificationStatus, Notification>($"notifications/{id}/status", status);

    /// <summary>Mark all notifications for a user as read.</summary>
    public async Task<ApiResult<ApiMessageResponse>> MarkAllAsRead(Guid userId)
    {
        var result = await _api.PatchAsync<object, ApiMessageResponse>(
            $"notifications/mark_all_as_read/{userId}", new { });

        if (result.Success)
            _toast.ShowToast("All notifications marked as read.", string.Empty, ToastService.ToastLevel.Success);
        else
            _toast.ShowToast(result.Error?.Message ?? "Failed to mark notifications as read.", string.Empty, ToastService.ToastLevel.Error);

        return result;
    }
}
