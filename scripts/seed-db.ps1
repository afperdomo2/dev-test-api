#Requires -Version 5.0
# Ejecuta el seed de preguntas via docker exec (no requiere psql en el host)
$ErrorActionPreference = "Stop"

# Cargar .env si existe
if (Test-Path ".env") {
  Get-Content ".env" | ForEach-Object {
    if ($_ -match "^\s*#" -or $_ -match "^\s*$") { return }
    $kv = $_ -split "=",2
    if ($kv.Count -eq 2) {
      $k = $kv[0].Trim()
      $v = $kv[1].Trim().Trim('"').Trim("'")
      Set-Item -Path "env:$k" -Value $v
    }
  }
}

$DB_CONTAINER = if ($env:DB_CONTAINER) { $env:DB_CONTAINER } else { "dev-postgres" }
$DB_USER = if ($env:DB_USER) { $env:DB_USER } else { "devuser" }
$DB_NAME = if ($env:DB_NAME) { $env:DB_NAME } else { "dev_test_api" }
$SEED_FILE = if ($args.Count -gt 0) { $args[0] } else { "db/seeds/seed_questions.sql" }

if (-not (Test-Path $SEED_FILE)) {
  Write-Host "No se encontro el archivo seed: $SEED_FILE" -ForegroundColor Red
  exit 1
}

$running = docker ps --format "{{.Names}}"
if ($running -notcontains $DB_CONTAINER) {
  Write-Host "Contenedor no encontrado o no esta corriendo: $DB_CONTAINER" -ForegroundColor Red
  Write-Host "Contenedores activos:"
  docker ps --format "  - {{.Names}} ({{.Image}})"
  exit 1
}

Write-Host "Ejecutando seed: $SEED_FILE -> $DB_CONTAINER ($DB_USER/$DB_NAME)..." -ForegroundColor Cyan
Get-Content $SEED_FILE -Raw | docker exec -i $DB_CONTAINER psql -U $DB_USER -d $DB_NAME -v ON_ERROR_STOP=1 -q
if ($LASTEXITCODE -ne 0) {
  Write-Host "Error al ejecutar el seed" -ForegroundColor Red
  exit 1
}
Write-Host "Seed ejecutado correctamente" -ForegroundColor Green
Write-Host "Verificacion rapida:" -ForegroundColor Cyan
docker exec -i $DB_CONTAINER psql -U $DB_USER -d $DB_NAME -c "SELECT 'questions' AS tabla, count(*) FROM questions UNION ALL SELECT 'question_options', count(*) FROM question_options UNION ALL SELECT 'question_topics', count(*) FROM question_topics UNION ALL SELECT 'code_challenges', count(*) FROM code_challenges ORDER BY tabla;"
