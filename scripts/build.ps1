$services = @(
    "realty/services/user-service/cmd",
    "realty/services/listing-service/cmd",
    "realty/services/deal-service/cmd",
    "realty/services/search-service/cmd",
    "realty/services/notification-service/cmd",
    "realty/services/api-gateway/cmd"
)

$outDir = "bin"
if (-not (Test-Path $outDir)) {
    New-Item -ItemType Directory -Path $outDir | Out-Null
}

foreach ($svc in $services) {
    $name = ($svc -split "/")[2]
    Write-Host "  -> $name"
    go build -o "$outDir/$name.exe" $svc
    if ($LASTEXITCODE -ne 0) {
        Write-Host "OSHIBKA: $name"
        exit 1
    }
}

Write-Host "Vse servisy sobrany v ./bin/"