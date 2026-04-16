param (
    $command
)

if (-not $command) {
    $command = "start"
}

$ProjectRoot = "${PSScriptRoot}/.."

$GeneratedPath = "${ProjectRoot}/internal/model"
$HandlerPath = "${ProjectRoot}/internal/handler"

# local environment variables for development purposes ONLY
$env:AMBULANCE_MANAGEMENT_API_PORT = "8080"
$env:AMBULANCE_MANAGEMENT_API_ENVIRONMENT = "Development"
$env:AMBULANCE_MANAGEMENT_API_MONGODB_USERNAME="root"
$env:AMBULANCE_MANAGEMENT_API_MONGODB_PASSWORD="neUhaDnes"

function mongo {
    docker compose --file ${ProjectRoot}/deployments/docker-compose/compose.yaml $args
}

switch ($command) {

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

        # Fix package name from model/api to handler after moving files
        Get-ChildItem -Path $HandlerPath -Filter "*.go" | ForEach-Object {

            $filePath = $_.FullName
            $content = Get-Content $filePath

            # Replace package declaration
            $updatedContent = $content -replace '^package\s+\w+', 'package handler'

            Set-Content -Path $filePath -Value $updatedContent

            Write-Host "Fixed package name in $($_.Name)"
        }

        Write-Host "Post-processing completed"
    }

    "start" {
        try {
            mongo up --detach
            go run ${ProjectRoot}/cmd
        } finally {
            mongo down
        }
    }

    "mongo" {
        mongo up
    }

    "docker" {
         docker build -t myrres/ambulance-management-webapi:local-build -f ${ProjectRoot}/build/docker/Dockerfile .
    }

    "test" {
        go test -v ./...
    }

    default {
        throw "Unknown command: $command"
    }
}
