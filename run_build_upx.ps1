# Build and compress
./run_build.ps1

Write-Host "Compressing with UPX..."
upx --best go-mail.exe
