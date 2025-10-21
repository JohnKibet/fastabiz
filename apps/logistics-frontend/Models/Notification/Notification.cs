using System.ComponentModel.DataAnnotations;
using System.Net;
using System.Text.Json;
using System.Text.Json.Serialization;

namespace logistics_frontend.Models.Notification
{
    public class Notification
    {

        [JsonPropertyName("id")]
        public Guid ID { get; set; }

        [JsonPropertyName("user_id")]
        public Guid UserID { get; set; }

        [JsonPropertyName("message")]
        public string Message { get; set; } = string.Empty;

        [JsonPropertyName("type")]
        public NotificationType Type { get; set; }

        [JsonPropertyName("status")]
        public NotificationStatus Status { get; set; }

        [JsonPropertyName("sent_at")]
        public DateTime SentAt { get; set; }

        [JsonPropertyName("updated_at")]
        public DateTime UpdatedAt { get; set; }

        [JsonPropertyName("created_at")]
        public DateTime CreatedAt { get; set; }
    }

    public enum NotificationStatus
    {
        Pending,
        Sent,
        Failed,
        Read
    }

    public enum NotificationType
    {
        Email,
        SMS,
        Push,
        System
    }

    public class CreateNotificationRequest
    {
        [Required]
        [JsonPropertyName("user_id")]
        public Guid UserID { get; set; }

        [Required]
        [JsonPropertyName("message")]
        public string Message { get; set; } = string.Empty;

        [Required]
        [JsonPropertyName("type")]
        public string Type { get; set; } = string.Empty;
    }
}