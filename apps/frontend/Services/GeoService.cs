using System.Net.Http.Json;

public class GeoService
{
    private readonly HttpClient _http;

    public GeoService(HttpClient http)
    {
        _http = http;
    }

    public record ReverseResult(Address address);
    public record Address(string road, string hamlet, string suburb, string city, string county, string state, string country);

    public async Task<string?> GetAddressAsync(double lat, double lng)
    {
        var url = $"https://nominatim.openstreetmap.org/reverse?lat={lat}&lon={lng}&format=json";

        var result = await _http.GetFromJsonAsync<ReverseResult>(url);

        if (result?.address == null) return null;

        // Formatting example (pick what exists)
        return $"{result.address.road}, {result.address.suburb}, {result.address.city}, {result.address.country}";
    }
}
