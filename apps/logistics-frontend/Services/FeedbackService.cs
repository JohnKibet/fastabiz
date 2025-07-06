using System.Net.Http.Json;
using logistics_frontend.Models.Feedback;

public class FeedbackService
{
    private readonly HttpClient _http;
    public FeedbackService(IHttpClientFactory httpClientFactory)
    {
        _http = httpClientFactory.CreateClient("AuthenticatedApi");
    }

    public async Task CreateFeedback(CreateFeedbackRequest feedback)
    {
        var response = await _http.PostAsJsonAsync("feedbacks/create", feedback);
        response.EnsureSuccessStatusCode();
    }

    public async Task<Feedback> GetFeedbackById(Guid feedbackId)
    {
        var feedback = await _http.GetFromJsonAsync<Feedback>($"feedbacks/{feedbackId}");
        return feedback ?? throw new Exception("Payment not found");
    }

    // public async Task<List<Feedback>> GetFeedbacksByOrderId(Guid orderId)
    // {
    //     var feedbacks = await _http.GetFromJsonAsync<List<Feedback>>($"feedbacks/order_id/{orderId}");
    //     return feedbacks ?? new List<Feedback>();
    // }

    public async Task<List<Feedback>> GetAllFeedbacks()
    {
        var feedbacks = await _http.GetFromJsonAsync<List<Feedback>>("feedbacks/all_feedbacks");
        return feedbacks ?? new List<Feedback>();
    }
}