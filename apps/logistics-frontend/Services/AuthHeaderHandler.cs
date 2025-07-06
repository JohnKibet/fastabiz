using System.Net.Http;
using System.Net.Http.Headers;
using System.Threading;
using System.Threading.Tasks;

namespace logistics_frontend.Services.AuthHeaderHandler;

public class AuthHeaderHandler : DelegatingHandler
{
    private readonly UserSessionService _session;
    private readonly ILogger<AuthHeaderHandler> _logger;
    public AuthHeaderHandler(UserSessionService session, ILogger<AuthHeaderHandler> logger)
    {
        _session = session;
        _logger = logger;
    }

    protected override async Task<HttpResponseMessage> SendAsync(HttpRequestMessage request, CancellationToken cancellationToken)
    {
        var token = await _session.GetTokenAsync();

        _logger.LogInformation("AuthHeaderHandler: Sending request to {Url} with token: {Token}", request.RequestUri, string.IsNullOrEmpty(token) ? "<null>" : token);

        if (!string.IsNullOrWhiteSpace(token))
        {
            request.Headers.Authorization = new AuthenticationHeaderValue("Bearer", token);
        }

        return await base.SendAsync(request, cancellationToken);
    }
}

