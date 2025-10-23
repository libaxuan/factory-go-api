$body = @{
    model = "claude-sonnet-4-5-20250929"
    messages = @(
        @{
            role = "user"
            content = "Say hello in Chinese"
        }
    )
    stream = $false
} | ConvertTo-Json

$headers = @{
    "Content-Type" = "application/json"
    "Authorization" = "Bearer 0000"
}

$response = Invoke-RestMethod -Uri "http://localhost:8003/v1/chat/completions" -Method Post -Headers $headers -Body $body

Write-Host $response.choices[0].message.content
