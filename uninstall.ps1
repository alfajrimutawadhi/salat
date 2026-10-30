$version = "v1.0.4"

$arch = if ($env:PROCESSOR_ARCHITEW6432) { $env:PROCESSOR_ARCHITEW6432 } else { $env:PROCESSOR_ARCHITECTURE }
if ($arch -eq "AMD64") {
    $arch = "amd64"
} elseif ($arch -eq "ARM64") {
    $arch = "arm64"
} else {
    $arch = $env:PROCESSOR_ARCHITECTURE
}

Write-Host "==== Uninstalling salat $version for windows-$arch ===="

$destDir = Join-Path $env:USERPROFILE ".salat"
if (Test-Path $destDir) {
    Remove-Item $destDir -Recurse -Force
} else {
    Write-Host "~\.salat not found, nothing to remove"
}

# Remove ~/.salat from the user PATH
$currentUserPath = [Environment]::GetEnvironmentVariable('Path', 'User')
if ($currentUserPath) {
    $pathParts = $currentUserPath -split ';' | Where-Object { $_ -ne '' -and $_ -ne $destDir }
    $newPath = $pathParts -join ';'
    if ($newPath -ne $currentUserPath) {
        [Environment]::SetEnvironmentVariable('Path', $newPath, 'User')
        Write-Host "Note: removed $destDir from your user PATH - open a new terminal for it to take effect"
    }
}

Write-Host "==== Salat uninstalled successfully ===="
