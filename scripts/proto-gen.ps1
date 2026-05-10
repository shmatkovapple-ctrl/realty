$PROTO_SRC = "api/proto"
$PROTO_GEN = "api/gen"
$ROOT = (Get-Location).Path

if (-not (Test-Path $PROTO_GEN)) {
    New-Item -ItemType Directory -Path $PROTO_GEN | Out-Null
}

$protoFiles = Get-ChildItem -Path $PROTO_SRC -Recurse -Filter "*.proto"

foreach ($file in $protoFiles) {
    Write-Host "  -> $($file.FullName)"
    & protoc `
        --proto_path="$ROOT\$PROTO_SRC" `
        --go_out="$ROOT\$PROTO_GEN" `
        --go_opt=paths=source_relative `
        --go-grpc_out="$ROOT\$PROTO_GEN" `
        --go-grpc_opt=paths=source_relative `
        "$($file.FullName)"

    if ($LASTEXITCODE -ne 0) {
        Write-Host "OSHIBKA pri obrabotke $($file.FullName)"
        exit 1
    }
}

Write-Host "Gotovo. Faily v ./$PROTO_GEN"