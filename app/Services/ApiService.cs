using System;
using System.Net.Http;
using System.Net.Http.Headers;
using System.Text;
using System.Threading.Tasks;
using Newtonsoft.Json;

namespace app;

public class ApiService
{
    private readonly HttpClient _client;

    public ApiService()
    {
        _client = new HttpClient { BaseAddress = new Uri(AppConfig.ApiBaseUrl) };
    }

    public async Task<(bool success, string token, string error)> LoginAsync(string email, string password)
    {
        try
        {
            var payload = new { email, password };
            var content = new StringContent(JsonConvert.SerializeObject(payload), Encoding.UTF8, "application/json");
            var response = await _client.PostAsync("/api/auth/login", content);
            
            if (response.IsSuccessStatusCode)
            {
                var json = await response.Content.ReadAsStringAsync();
                var result = JsonConvert.DeserializeObject<dynamic>(json);
                return (true, result.token.ToString(), null);
            }
            
            return (false, null, await response.Content.ReadAsStringAsync());
        }
        catch (Exception ex)
        {
            return (false, null, ex.Message);
        }
    }

    public async Task<(bool success, string error)> RegisterAsync(string email, string password)
    {
        try
        {
            var payload = new { email, password };
            var content = new StringContent(JsonConvert.SerializeObject(payload), Encoding.UTF8, "application/json");
            var response = await _client.PostAsync("/api/auth/register", content);
            return (response.IsSuccessStatusCode, response.IsSuccessStatusCode ? null : await response.Content.ReadAsStringAsync());
        }
        catch (Exception ex)
        {
            return (false, ex.Message);
        }
    }

    public async Task<dynamic> GetUserStatsAsync()
    {
        try
        {
            _client.DefaultRequestHeaders.Authorization = new AuthenticationHeaderValue("Bearer", AppConfig.JwtToken);
            var response = await _client.GetAsync("/api/user/stats");
            if (response.IsSuccessStatusCode)
            {
                var json = await response.Content.ReadAsStringAsync();
                return JsonConvert.DeserializeObject<dynamic>(json);
            }
            return null;
        }
        catch
        {
            return null;
        }
    }

    public async Task<dynamic> GetFilesAsync()
    {
        try
        {
            _client.DefaultRequestHeaders.Authorization = new AuthenticationHeaderValue("Bearer", AppConfig.JwtToken);
            var response = await _client.GetAsync("/api/files");
            if (response.IsSuccessStatusCode)
            {
                var json = await response.Content.ReadAsStringAsync();
                return JsonConvert.DeserializeObject<dynamic>(json);
            }
            return null;
        }
        catch
        {
            return null;
        }
    }

    public async Task<bool> SendHeartbeatAsync()
    {
        try
        {
            _client.DefaultRequestHeaders.Authorization = new AuthenticationHeaderValue("Bearer", AppConfig.JwtToken);
            var response = await _client.PostAsync("/api/peer/heartbeat", null);
            return response.IsSuccessStatusCode;
        }
        catch
        {
            return false;
        }
    }

    public async Task<(bool success, string error)> UploadChunkAsync(string chunkHash, byte[] chunkData)
    {
        try
        {
            _client.DefaultRequestHeaders.Authorization = new AuthenticationHeaderValue("Bearer", AppConfig.JwtToken);
            var content = new ByteArrayContent(chunkData);
            var response = await _client.PostAsync($"/api/chunks/upload/{chunkHash}", content);
            return (response.IsSuccessStatusCode, response.IsSuccessStatusCode ? null : await response.Content.ReadAsStringAsync());
        }
        catch (Exception ex)
        {
            return (false, ex.Message);
        }
    }

    public async Task<byte[]> DownloadChunkAsync(string chunkHash)
    {
        try
        {
            _client.DefaultRequestHeaders.Authorization = new AuthenticationHeaderValue("Bearer", AppConfig.JwtToken);
            var response = await _client.GetAsync($"/api/chunks/{chunkHash}");
            if (response.IsSuccessStatusCode)
            {
                return await response.Content.ReadAsByteArrayAsync();
            }
            return null;
        }
        catch
        {
            return null;
        }
    }
}