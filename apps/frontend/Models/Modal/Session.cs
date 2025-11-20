namespace frontend.Models
{
    public class SessionModel
    {
        public Guid ID { get; set; }
        public string Device { get; set; } = string.Empty;
        public string Browser { get; set; } = string.Empty;
        public string IP { get; set; } = string.Empty;

        public DateTime LastActive { get; set; }
        public bool IsCurrent { get; set; }
    }
}
