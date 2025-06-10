using System.ComponentModel.DataAnnotations;

namespace logistics_frontend.Models.Notification
{
    public class Notification
    {
        public Guid ID { get; set; }
        public Guid UserID { get; set; }
        public string Message { get; set; } = string.Empty;
        public string Type { get; set; } = string.Empty;
        public DateTime SentAt { get; set; }
    }

    public class CreateNotificationRequest
    {
        [Required]
        public Guid UserID { get; set; }

        [Required]
        public string Message { get; set; } = string.Empty;

        [Required]
        public string Type { get; set; } = string.Empty;
    }
}