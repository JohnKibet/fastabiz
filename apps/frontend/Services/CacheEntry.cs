namespace frontend.Services;

/// <summary>
/// Lightweight in-memory cache for a single value with a TTL.
/// Eliminates the repeated _cachedX / _lastFetchTime / _cacheDuration
/// pattern across every service.
///
/// USAGE:
///     private readonly CacheEntry&lt;List&lt;User&gt;&gt; _cache = new(TimeSpan.FromMinutes(5));
///
///     public async Task&lt;ApiResult&lt;List&lt;User&gt;&gt;&gt; GetAllCachedUsers(bool forceRefresh = false)
///     {
///         if (!forceRefresh && _cache.TryGet(out var cached))
///             return ApiResult&lt;List&lt;User&gt;&gt;.Ok(cached!, fromCache: true);
///
///         var result = await GetAllUsers();
///         if (result.Success && result.Data != null)
///             _cache.Set(result.Data);
///
///         return result;
///     }
///
///     public void InvalidateCache() => _cache.Invalidate();
/// </summary>
public sealed class CacheEntry<T>
{
    private T?       _value;
    private DateTime _setAt  = DateTime.MinValue;
    private readonly TimeSpan _ttl;

    public CacheEntry(TimeSpan ttl) => _ttl = ttl;

    /// <summary>Returns true and populates <paramref name="value"/> if the cache is valid.</summary>
    public bool TryGet(out T? value)
    {
        if (_value is not null && DateTime.UtcNow - _setAt < _ttl)
        {
            value = _value;
            return true;
        }
        value = default;
        return false;
    }

    /// <summary>Stores a value and resets the TTL clock.</summary>
    public void Set(T value)
    {
        _value = value;
        _setAt = DateTime.UtcNow;
    }

    /// <summary>Clears the cache so the next call hits the network.</summary>
    public void Invalidate()
    {
        _value = default;
        _setAt = DateTime.MinValue;
    }

    public bool IsValid => _value is not null && DateTime.UtcNow - _setAt < _ttl;
}
