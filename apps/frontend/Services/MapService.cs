using LeafletForBlazor;

public class MapService
{
    public string OpenCycleMapAPIKey { get; set; } = "YOUR_API_KEY";

    public List<RealTimeMap.BasemapLayer> GetDefaultBasemapLayers() => new()
    {
        new RealTimeMap.BasemapLayer
        {
            url = "http://{s}.tile.openstreetmap.org/{z}/{x}/{y}.png",
            attribution = "©Open Street Map",
            title = "Open Street Map",
            detectRetina = true
        },
        new RealTimeMap.BasemapLayer
        {
            url = $"https://tile.thunderforest.com/cycle/{{z}}/{{x}}/{{y}}.png?apikey={OpenCycleMapAPIKey}",
            attribution = "©Open Cycle Map",
            title = "Open Cycle Map"
        }
    };
}
