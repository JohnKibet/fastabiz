namespace frontend.Models.Driver
{
    public class Driver
    {
        public Guid ID { get; set; }
        public string FullName { get; set; } = string.Empty;
        public string Email { get; set; } = string.Empty;
        public string VehicleInfo { get; set; } = string.Empty;
        public string CurrentLocation { get; set; } = string.Empty;
        public bool Available { get; set; }
        public DateTime CreatedAt { get; set; }
    }

    public class CreateDriverRequest
    {
        public string FullName { get; set; } = string.Empty;
        public string Email { get; set; } = string.Empty;
        public string VehicleInfo { get; set; } = string.Empty;
        public string CurrentLocation { get; set; } = string.Empty;
    }
}