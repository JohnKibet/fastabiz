using frontend.Models;
using frontend.Models.Storefront;

namespace frontend.Services;

public class FeedbackService
{
    private readonly ApiService _api;

    public FeedbackService(ApiService api)
    {
        _api = api;
    }

    public Task<ApiResult<ApiMessageResponse>> CreateFeedback(CreateFeedbackRequest req)
        => _api.PostAsync<CreateFeedbackRequest, ApiMessageResponse>("feedbacks/create", req);

    public Task<ApiResult<Feedback>> GetFeedbackById(Guid feedbackId)
        => _api.GetAsync<Feedback>($"feedbacks/{feedbackId}");

    public Task<ApiResult<List<Feedback>>> GetAllFeedbacks()
        => _api.GetAsync<List<Feedback>>("feedbacks/all_feedbacks");

    // public Task<ApiResult<List<Feedback>>> GetFeedbacksByOrderId(Guid orderId)
    //     => _api.GetAsync<List<Feedback>>($"feedbacks/order_id/{orderId}");
}
