<#
.SYNOPSIS
  Build, package, and publish the payload-capture-wasm custom policy to
  Anypoint Exchange so it can be attached to APIs in API Manager.

.DESCRIPTION
  Three steps:
    1. cargo build --release --target wasm32-wasi
    2. package the .wasm + policy-manifest.yaml as a zip
    3. anypoint-cli-v4 exchange asset upload --type custom-policy

  Prerequisites:
    - Rust toolchain with wasm32-wasi target installed
        (rustup target add wasm32-wasi)
    - anypoint-cli-v4 installed and configured with credentials
    - Anypoint Business Group ID known (replace placeholder in
      policy-manifest.yaml before deploy)
    - The Service Bus upstream cluster declared in Flex Gateway config

.EXAMPLE
  # First-time setup
  rustup target add wasm32-wasi
  npm install -g anypoint-cli-v4

  # Run
  .\deploy-policy.ps1 -BusinessGroupId 12345678-1234-1234-1234-123456789012
#>
[CmdletBinding()]
param(
  [Parameter(Mandatory)]
  [string] $BusinessGroupId,

  [string] $Version = "0.1.0",
  [string] $AssetId = "payload-capture-wasm",
  [switch] $SkipBuild,
  [switch] $SkipUpload
)

$ErrorActionPreference = 'Stop'
$scriptDir = Split-Path -Parent $MyInvocation.MyCommand.Path
Set-Location $scriptDir

# -- 1. Build --
if (-not $SkipBuild) {
  Write-Host "1/3 cargo build --release --target wasm32-wasi" -ForegroundColor Cyan
  cargo build --release --target wasm32-wasi
  if ($LASTEXITCODE -ne 0) {
    throw "cargo build failed"
  }

  $wasmSrc = "target/wasm32-wasi/release/payload_capture_wasm.wasm"
  if (-not (Test-Path $wasmSrc)) {
    throw "Expected $wasmSrc not found after build"
  }
  $wasmDst = "payload-capture-wasm.wasm"
  Copy-Item $wasmSrc $wasmDst -Force
  $sizeKb = [math]::Round((Get-Item $wasmDst).Length / 1KB, 1)
  Write-Host "  -> built $wasmDst ($sizeKb KB)"
} else {
  Write-Host "1/3 build skipped (-SkipBuild)" -ForegroundColor Yellow
}

# -- 2. Package --
Write-Host "2/3 packaging .wasm + manifest" -ForegroundColor Cyan

# Replace business group placeholder in manifest with the real ID
$manifestSrc = "policy-manifest.yaml"
$manifestStaged = "policy-manifest-staged.yaml"
(Get-Content $manifestSrc -Raw) -replace "<YOUR_BUSINESS_GROUP_ID>", $BusinessGroupId `
  | Set-Content -Path $manifestStaged -Encoding UTF8

# Package as a zip (Anypoint accepts the custom-policy asset as a zip)
$zipPath = "$AssetId-$Version.zip"
if (Test-Path $zipPath) { Remove-Item $zipPath }
Compress-Archive -Path "payload-capture-wasm.wasm", $manifestStaged -DestinationPath $zipPath
Write-Host "  -> packaged $zipPath"

# -- 3. Upload to Exchange --
if (-not $SkipUpload) {
  Write-Host "3/3 anypoint-cli-v4 exchange asset upload" -ForegroundColor Cyan
  $existing = anypoint-cli-v4 exchange asset list --search $AssetId 2>$null
  if ($LASTEXITCODE -eq 0 -and $existing -match $AssetId) {
    Write-Host "  -> asset exists; uploading new version $Version"
  } else {
    Write-Host "  -> uploading new asset $AssetId v$Version"
  }

  anypoint-cli-v4 exchange asset upload `
    --name "Payload Capture WASM" `
    --type custom-policy `
    --version $Version `
    --groupId $BusinessGroupId `
    --classifier "policy-binding" `
    --properties.mainFile "policy-manifest-staged.yaml" `
    $zipPath

  if ($LASTEXITCODE -ne 0) {
    throw "anypoint-cli upload failed"
  }
  Write-Host ""
  Write-Host "Published $AssetId v$Version to Exchange (business group $BusinessGroupId)" -ForegroundColor Green
  Write-Host "Attach via API Manager -> API instance -> Policies -> Add Policy -> Custom"
} else {
  Write-Host "3/3 upload skipped (-SkipUpload)" -ForegroundColor Yellow
}

# Cleanup staged manifest
Remove-Item $manifestStaged -ErrorAction SilentlyContinue
