$AppName   = "simple-monitor"
$OutputDir = "build"
$OSES      = @("windows", "linux", "darwin")
$ARCHS     = @("amd64", "arm64", "386")

# Define supported combinations to avoid unsupported GOOS/GOARCH pairs
$SupportedCombinations = @{
    "windows" = @("amd64", "arm64", "386")
    "linux"   = @("amd64", "arm64", "386")
    "darwin"  = @("amd64", "arm64")  # darwin/386 is not supported
}

if (Test-Path $OutputDir) {
    Remove-Item $OutputDir -Recurse -Force
}
New-Item -ItemType Directory -Path $OutputDir | Out-Null

foreach ($os in $OSES) {
    $supportedArchs = $SupportedCombinations[$os]
    foreach ($arch in $supportedArchs) {
        Write-Host "Building for $os/$arch ..."
        $ext = ""
        if ($os -eq "windows") { $ext = ".exe" }

        $env:GOOS = $os
        $env:GOARCH = $arch
        go build -o "$OutputDir\$AppName-$os-$arch$ext" main.go
    }
}

Write-Host "✅ All builds are in the '$OutputDir' folder."
