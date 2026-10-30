$version = "v1.0.4"

# Prefer PROCESSOR_ARCHITEW6432: a 32-bit PowerShell on a 64-bit OS reports
# x86 in PROCESSOR_ARCHITECTURE but the real arch here
$arch = if ($env:PROCESSOR_ARCHITEW6432) { $env:PROCESSOR_ARCHITEW6432 } else { $env:PROCESSOR_ARCHITECTURE }
if ($arch -eq "AMD64") {
    $arch = "amd64"
} elseif ($arch -eq "ARM64") {
    $arch = "arm64"
} elseif ($arch -eq "x86") {
    Write-Host "Error: 32-bit Windows (x86) is not supported - salat releases amd64 and arm64 only"
    exit 1
} else {
    Write-Host "Unsupported architecture: $arch"
    exit 1
}

$assetName = "salat-$version-windows-$arch.zip"
$url = "https://github.com/alfajrimutawadhi/salat/releases/download/$version/$assetName"
$zipPath = Join-Path $env:TEMP $assetName

Write-Host "==== Downloading salat $version for windows-$arch ===="
$ProgressPreference = 'SilentlyContinue'
try {
    if (Get-Command curl.exe -ErrorAction SilentlyContinue) {
        & curl.exe -fSL -o $zipPath $url
        if ($LASTEXITCODE -ne 0) { exit 1 }
    } else {
        Invoke-WebRequest -Uri $url -OutFile $zipPath -UseBasicParsing
    }
} catch {
    Write-Host "Error: failed to download $url"
    Write-Host $_.Exception.Message
    exit 1
}

$destDir = Join-Path $env:USERPROFILE ".salat"
if (-not (Test-Path $destDir)) {
    New-Item -ItemType Directory -Path $destDir | Out-Null
}

Expand-Archive -Path $zipPath -DestinationPath $destDir -Force
Remove-Item $zipPath

$exePath = Join-Path $destDir "salat.exe"
if (-not (Test-Path $exePath)) {
    Write-Host "Error: salat.exe not found in downloaded archive"
    exit 1
}

# Add ~/.salat to the user PATH so 'salat' can be run from anywhere
$currentUserPath = [Environment]::GetEnvironmentVariable('Path', 'User')
$pathParts = @()
if ($currentUserPath) {
    $pathParts = $currentUserPath -split ';' | Where-Object { $_ -ne '' }
}
if (-not ($pathParts -contains $destDir)) {
    $newPath = ($( ($pathParts + $destDir) -join ';'))
    [Environment]::SetEnvironmentVariable('Path', $newPath, 'User')
    Write-Host "Note: added $destDir to your user PATH - open a new terminal to use 'salat'"
}

Write-Host "==== Salat installed successfully ===="
