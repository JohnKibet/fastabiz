namespace frontend.Models
{
    public class ErrorResponse
    {
        public string Error { get; set; } = string.Empty;
        public string Detail { get; set; } = string.Empty;

        // handle backend validation messages
        public Dictionary<string, string[]>? Errors { get; set; }
    }

    public class ServiceResult<T>
    {
        public bool Success { get; set; }
        public T? Data { get; set; }
        public string? ErrorMessage { get; set; }
        public bool FromCache { get; set; }

        public static ServiceResult<T> Ok(T data, bool fromCache = false) => new() { Success = true, Data = data, FromCache = fromCache };
        public static ServiceResult<T> Fail(string message) => new() { Success = false, ErrorMessage = message };
    };

    public class ServiceResult2<T>
    {
        public bool Success { get; set; }
        public T? Data { get; set; }

        // Full structured backend error
        public ErrorResponse? Error { get; set; }

        public bool FromCache { get; set; }

        // Helpers
        public bool HasError => !Success && Error != null;
        public bool HasValidationErrors => Error?.Errors != null && Error.Errors.Count > 0;

        public static ServiceResult2<T> Ok(T data, bool fromCache = false)
            => new() { Success = true, Data = data, FromCache = fromCache };

        public static ServiceResult2<T> Fail(ErrorResponse error)
            => new() { Success = false, Error = error };

        public static ServiceResult2<T> Fail(string message)
            => new()
            {
                Success = false,
                Error = new ErrorResponse
                {
                    Error = message,
                    Detail = message
                }
            };
    }
}

