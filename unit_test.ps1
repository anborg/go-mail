$ErrorActionPreference = "Stop"

Write-Host "Running Go Tests with Coverage..." -ForegroundColor Cyan
# Cleanup previous runs
if (Test-Path coverage) { Remove-Item coverage }
if (Test-Path coverage.out) { Remove-Item coverage.out }

go test -coverprofile coverage.out ./...

if ($LASTEXITCODE -eq 0) {
    if (-not (Test-Path coverage.out) -and (Test-Path coverage)) {
        Move-Item coverage coverage.out
    }

    if (Test-Path coverage.out) {
        Write-Host "`nGenerating Coverage Report..." -ForegroundColor Cyan
        cmd /c "go tool cover -func=coverage.out"
        # coverage.out is kept for inspection as requested
    } else {
        Write-Host "`nWarning: coverage.out not found." -ForegroundColor Yellow
    }
} else {
    Write-Host "`nTests Failed." -ForegroundColor Red
    exit $LASTEXITCODE
}
