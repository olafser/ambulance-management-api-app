param (
    $command
)

if (-not $command) {
    $command = "start"
}

$ProjectRoot = "${PSScriptRoot}/.."

$env:AMBULANCE_API_ENVIRONMENT = "Development"
$env:AMBULANCE_API_PORT = "8080"

$GeneratedPath = "${ProjectRoot}/internal/model"
$HandlerPath = "${ProjectRoot}/internal/handler"

switch ($command) {

    "start" {
        go run ${ProjectRoot}/cmd/ambulance-api-service
    }

    "openapi" {

        # Run OpenAPI generator inside Docker
        docker run --rm -ti -v ${ProjectRoot}:/local openapitools/openapi-generator-cli generate -c /local/scripts/generator-cfg.yaml

        Write-Host "OpenAPI generation completed"

        # Remove generated README file if it exists
        $readme = Join-Path $GeneratedPath "README.md"
        if (Test-Path $readme) {
            Remove-Item $readme -Force
            Write-Host "Removed README.md"
        }

        # Remove .openapi-generator directory from project root
        $openapiGenDir = Join-Path $ProjectRoot ".openapi-generator"
        if (Test-Path $openapiGenDir) {
            Remove-Item $openapiGenDir -Recurse -Force
            Write-Host "Removed .openapi-generator directory"
        }

        # Ensure handler directory exists
        if (-not (Test-Path $HandlerPath)) {
            New-Item -ItemType Directory -Path $HandlerPath | Out-Null
        }

        # Move all generated API interface files (api_*.go) to handler directory
        Get-ChildItem -Path $GeneratedPath -Filter "api_*.go" | ForEach-Object {
            Move-Item $_.FullName -Destination $HandlerPath -Force
            Write-Host "Moved $($_.Name) to handler directory"
        }

        # Move routers.go to handler directory
        $routersFile = Join-Path $GeneratedPath "routers.go"
        if (Test-Path $routersFile) {
            Move-Item $routersFile -Destination $HandlerPath -Force
            Write-Host "Moved routers.go to handler directory"
        }

        Write-Host "Post-processing completed"
    }

    default {
        throw "Unknown command: $command"
    }
}
