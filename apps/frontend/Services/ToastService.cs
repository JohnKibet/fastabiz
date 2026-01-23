public class ToastService
{
    public event Action<ToastMessage>? OnShow;

    public void ShowToast(
        string message,
        string errDetail,
        ToastLevel level = ToastLevel.Info,
        int durationMs = 4000,
        bool allowHtml = false,
        string? actionLabel = null,
        Action? onAction = null)
    {
        var toast = new ToastMessage
        {
            Message = message,
            ErrDetail = errDetail,
            Level = level,
            DurationMs = durationMs,
            AllowHtml = allowHtml,
            ActionLabel = actionLabel,
            OnAction = onAction
        };

        OnShow?.Invoke(toast);
    }
    public enum ToastLevel
    {
        Info,
        Success,
        Warning,
        Error
    }

    public class ToastMessage
    {
        public string Message { get; set; } = string.Empty;
        public string ErrDetail { get; set; } = string.Empty;
        public ToastLevel Level { get; set; }
        public int DurationMs { get; set; }
        public bool AllowHtml { get; set; }
        public string? ActionLabel { get; set; }
        public Action? OnAction { get; set; }
    }
}