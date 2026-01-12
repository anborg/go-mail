$ErrorActionPreference = "Stop"

# Define directories
$coverDir = "coverage_data"
$binaryName = "go-mail_cover.exe"
$tempConfig = "config_temp.yml"

Write-Host "Cleaning up previous coverage data..." -ForegroundColor Cyan
if (Test-Path $coverDir) { Remove-Item -Recurse -Force $coverDir }
New-Item -ItemType Directory -Path $coverDir | Out-Null
if (Test-Path $binaryName) { Remove-Item $binaryName }
if (Test-Path "coverage_binary.out") { Remove-Item "coverage_binary.out" }

Write-Host "Building binary with coverage instrumentation..." -ForegroundColor Cyan
go build -cover -o $binaryName ./cmd/go-mail

try {
    Write-Host "Preparing test environment..." -ForegroundColor Cyan
    # Prepare a test run
    if (-not (Test-Path "input")) { New-Item -ItemType Directory -Path "input" | Out-Null }
    if (-not (Test-Path "testfiles/cayinput566000.csv")) {
        Write-Host "Warning: testfiles/cayinput566000.csv not found. Creating a dummy one." -ForegroundColor Yellow
        if (-not (Test-Path "testfiles")) { New-Item -ItemType Directory -Path "testfiles" | Out-Null }
        "Header`nEmail,s@e.com,ID,Name,Addr,Contact,$100,2024-01-01,Acc,Invoices" | Out-File "testfiles/cayinput566000.csv"
    }
    Copy-Item "testfiles/cayinput566000.csv" "input/AP566000.csv" -Force

    # Ensure directories exist for config
    if (-not (Test-Path "done")) { New-Item -ItemType Directory -Path "done" | Out-Null }
    if (-not (Test-Path "error")) { New-Item -ItemType Directory -Path "error" | Out-Null }

    # Create a temporary config if config.yml is missing
    if (-not (Test-Path "config.yml")) {
        Write-Host "Creating temporary config.yml..." -ForegroundColor Cyan
        $configContent = @"
app:
  lumberJackLogging:
    filename: "logs/eft.log"
    maxSize: 10
    maxBackups: 3
    maxAge: 28
    compress: true
fileProcessor:
  globPath: "AP566*.csv"
  olderThanSeconds: -1
  inputDir: "input/"
  doneDir: "done/"
  errorDir: "error/"
mailServer:
  host: "localhost"
  port: 1025
"@
        $configContent | Out-File "config.yml" -Encoding utf8
    }

    # Set GOCOVERDIR and run
    $env:GOCOVERDIR = $coverDir
    Write-Host "Running binary..." -ForegroundColor Cyan
    & "./$binaryName"

    Write-Host "`nReporting coverage percentages..." -ForegroundColor Cyan
    # Use quotes to ensure parameters stay together
    go tool covdata percent "-i=$coverDir"

    Write-Host "`nGenerating text report (coverage_binary.out)..." -ForegroundColor Cyan
    go tool covdata textfmt "-i=$coverDir" "-o=coverage_binary.out"
    
    if (Test-Path "coverage_binary.out") {
        Write-Host "`nTop 10 functions by coverage:" -ForegroundColor Cyan
        go tool cover -func coverage_binary.out | Select-Object -First 10
    }
}
finally {
    $env:GOCOVERDIR = $null
    if (Test-Path $binaryName) { Remove-Item $binaryName }
    # We leave config.yml if it existed, otherwise we should probably clean it up if we created it.
    # But for a test script, having it might be better. 
}
