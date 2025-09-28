$AppName   = "simple-monitor"
$OutputDir = "build"
$OSES      = @("windows", "linux", "darwin")
$ARCHS     = @("amd64", "arm64", "386")

if (Test-Path $OutputDir) {
    Remove-Item $OutputDir -Recurse -Force
}
New-Item -ItemType Directory -Path $OutputDir | Out-Null

foreach ($os in $OSES) {
    foreach ($arch in $ARCHS) {
        Write-Host "Building for $os/$arch ..."
        $ext = ""
        if ($os -eq "windows") { $ext = ".exe" }

        $env:GOOS = $os
        $env:GOARCH = $arch
        go build -o "$OutputDir\$AppName-$os-$arch$ext" main.go
    }
}

Write-Host "✅ All builds are in the '$OutputDir' folder."
