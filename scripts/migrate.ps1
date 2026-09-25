# Loads .env and runs the goose CLI for the given command (up, down, status, create).
# Usage:
#   scripts/migrate.ps1 up
#   scripts/migrate.ps1 down
#   scripts/migrate.ps1 status
#   scripts/migrate.ps1 create <name>
[CmdletBinding()]
param(
  [Parameter(Mandatory = $true, Position = 0)]
  [string] $Command,

  [Parameter(Position = 1)]
  [string] $Name
)

$errorActionPreference = "Stop"

$migrationDir = "internal/database/migrations"
$envFile = Join-Path (Get-Location) ".env"

# Carga .env sin pisar variables ya definidas en el entorno (paridad con godotenv)
if (Test-Path $envFile) {
  Get-Content $envFile | ForEach-Object {
    if ($_ -match '^\s*([A-Za-z_][A-Za-z0-9_]*)=(.*)$') {
      $key = $Matches[1]
      $value = $Matches[2].Trim()
      if (-not (Test-Path "Env:$key")) {
        Set-Item -Path "Env:$key" -Value $value
      }
    }
  }
}

if (-not $env:DB_NAME) {
  Write-Error "DB_NAME no definido. Verifica .env"
  exit 1
}

$env:GOOSE_DRIVER = "postgres"
$env:GOOSE_DBSTRING = "host=$env:DB_HOST port=$env:DB_PORT user=$env:DB_USER password=$env:DB_PASSWORD dbname=$env:DB_NAME sslmode=$env:DB_SSL_MODE"

switch ($Command) {
  "create" {
    if ([string]::IsNullOrWhiteSpace($Name)) {
      Write-Error "Uso: make migrate-create name=<descripcion>"
      exit 1
    }
    $gooseArgs = @("-s", "create", $Name, "sql")
  }
  default {
    $gooseArgs = @($Command)
  }
}

& go tool goose -dir $migrationDir @gooseArgs
exit $LASTEXITCODE