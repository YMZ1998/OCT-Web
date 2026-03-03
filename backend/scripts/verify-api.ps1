param(
  [string]$BaseUrl = "http://localhost:8080",
  [string]$Username = "u$([int](Get-Date -UFormat %s))",
  [string]$Password = "P@ssw0rd123",
  [string]$Email = "test_$([int](Get-Date -UFormat %s))@example.com"
)

$ErrorActionPreference = 'Stop'

function Assert-BizSuccess($resp, $step) {
  if ($null -eq $resp) { throw "$step: empty response" }
  if ($resp.code -ne 0) {
    throw "$step failed: code=$($resp.code), msg=$($resp.msg)"
  }
}

Write-Host "[1/3] Register -> $BaseUrl/api/v1/register"
$registerBody = @{
  username = $Username
  password = $Password
  email    = $Email
  gender   = "other"
  age      = 20
} | ConvertTo-Json

$registerResp = Invoke-RestMethod -Method Post `
  -Uri "$BaseUrl/api/v1/register" `
  -ContentType "application/json" `
  -Body $registerBody
Assert-BizSuccess $registerResp "register"

Write-Host "[2/3] Login -> $BaseUrl/api/v1/login"
$loginBody = @{
  username = $Username
  password = $Password
} | ConvertTo-Json

$loginResp = Invoke-RestMethod -Method Post `
  -Uri "$BaseUrl/api/v1/login" `
  -ContentType "application/json" `
  -Body $loginBody
Assert-BizSuccess $loginResp "login"

$token = $loginResp.data.token
$id = $loginResp.data.id
if ([string]::IsNullOrWhiteSpace($token) -or [string]::IsNullOrWhiteSpace($id)) {
  throw "login response missing token or id"
}

Write-Host "[3/3] Get user (JWT) -> $BaseUrl/api/v1/user/$id"
$getResp = Invoke-RestMethod -Method Get `
  -Uri "$BaseUrl/api/v1/user/$id" `
  -Headers @{ Authorization = "Bearer $token" }
Assert-BizSuccess $getResp "get user"

Write-Host "SUCCESS: backend API is reachable and auth flow works."
Write-Host "username=$Username id=$id"
