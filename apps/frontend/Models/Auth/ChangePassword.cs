using System.ComponentModel.DataAnnotations;
namespace frontend.Models;
public class ChangePasswordModel
{
    [Required]
    public string CurrentPassword { get; set; } = "";

    [Required]
    [MinLength(6)]
    public string NewPassword { get; set; } = "";

    [Required]
    public string ConfirmPassword { get; set; } = "";
}