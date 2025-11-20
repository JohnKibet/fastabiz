using System.ComponentModel.DataAnnotations;

namespace frontend.Models
{
    public class EditUserModel
    {
        [Required]
        public string FullName { get; set; } = string.Empty;

        [Required, EmailAddress]
        public string Email { get; set; } = string.Empty;

        [Phone]
        public string? Phone { get; set; }

        [Required]
        public string Slug { get; set; } = string.Empty;
    }
}
