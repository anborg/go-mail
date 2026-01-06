$source = "\\markham.ca\apps\eft-test\test_files"
$destination = "\\markham.ca\apps\eft-test\input"

Get-ChildItem -Path $source -Filter *.csv -File | ForEach-Object {
    $destFile = Join-Path $destination $_.Name

    if (-not (Test-Path $destFile)) {
        Copy-Item -Path $_.FullName -Destination $destination
    }
}
