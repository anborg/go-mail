$source = "\\markham.ca\apps\eft-test\test_files"
$destination = "\\markham.ca\apps\eft-test\input"

Get-ChildItem -Path $source -Filter *.csv -File | ForEach-Object {
    $destFile = Join-Path $destination $_.Name

    if (-not (Test-Path $destFile)) {
        Copy-Item -Path $_.FullName -Destination $destination
    }
}
$BASE_PATH = "\\markham.ca\data\ITS\Common\devops\dist\eft"

# Executable path
$exePath = Join-Path $BASE_PATH "go-mail.exe"

# Config file path
$configFile = Join-Path $BASE_PATH "config.yml"

# Execute the command
& $exePath -configFile $configFile