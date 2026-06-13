param(
    [string]$OutputDir = "",
    [ValidateSet("amd64", "arm64", "386")]
    [string]$Arch = "amd64"
)

$ErrorActionPreference = "Stop"

$RootDir = Split-Path -Parent $PSScriptRoot
if ([string]::IsNullOrWhiteSpace($OutputDir)) {
    $OutputDir = Join-Path $RootDir "bin"
}

$OutFile = Join-Path $OutputDir "cognitor.exe"

New-Item -ItemType Directory -Force -Path $OutputDir | Out-Null
Push-Location $RootDir
try {
    $env:GOOS = "windows"
    $env:GOARCH = $Arch
    go build -trimpath -ldflags="-s -w" -o $OutFile ./cmd/cognitor
}
finally {
    Remove-Item Env:\GOOS -ErrorAction SilentlyContinue
    Remove-Item Env:\GOARCH -ErrorAction SilentlyContinue
    Pop-Location
}

Write-Host "built $OutFile for windows/$Arch"
