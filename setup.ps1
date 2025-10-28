# ==============================================
# Loterias API - Script de Inicialização
# ==============================================

Write-Host "========================================" -ForegroundColor Cyan
Write-Host "   Loterias API - Inicialização" -ForegroundColor Cyan
Write-Host "========================================" -ForegroundColor Cyan
Write-Host ""

# Verificar se Go está instalado
Write-Host "Verificando instalação do Go..." -ForegroundColor Yellow
try {
    $goVersion = go version
    Write-Host "✓ Go encontrado: $goVersion" -ForegroundColor Green
}
catch {
    Write-Host "✗ Go não encontrado! Por favor, instale Go 1.21+" -ForegroundColor Red
    Write-Host "  Download: https://go.dev/dl/" -ForegroundColor Yellow
    exit 1
}

# Verificar se arquivo .env existe
Write-Host ""
Write-Host "Verificando arquivo .env..." -ForegroundColor Yellow
if (-not (Test-Path ".env")) {
    Write-Host "! Arquivo .env não encontrado" -ForegroundColor Yellow
    if (Test-Path ".env.example") {
        Write-Host "  Criando .env a partir de .env.example..." -ForegroundColor Yellow
        Copy-Item ".env.example" ".env"
        Write-Host "✓ Arquivo .env criado" -ForegroundColor Green
        Write-Host "  Por favor, revise as configurações em .env" -ForegroundColor Cyan
    }
    else {
        Write-Host "✗ Arquivo .env.example não encontrado" -ForegroundColor Red
        exit 1
    }
}
else {
    Write-Host "✓ Arquivo .env encontrado" -ForegroundColor Green
}

# Verificar MongoDB
Write-Host ""
Write-Host "Verificando MongoDB..." -ForegroundColor Yellow
$mongoRunning = $false

# Tentar conectar ao MongoDB
try {
    $mongoTest = mongosh --eval "db.version()" --quiet 2>$null
    if ($LASTEXITCODE -eq 0) {
        $mongoRunning = $true
        Write-Host "✓ MongoDB está rodando" -ForegroundColor Green
    }
}
catch {
    # MongoDB CLI não encontrado, tentar verificar processo
    $mongoProcess = Get-Process mongod -ErrorAction SilentlyContinue
    if ($mongoProcess) {
        $mongoRunning = $true
        Write-Host "✓ MongoDB está rodando (processo detectado)" -ForegroundColor Green
    }
}

if (-not $mongoRunning) {
    Write-Host "! MongoDB não está rodando" -ForegroundColor Yellow
    Write-Host "  Opções:" -ForegroundColor Cyan
    Write-Host "  1. Iniciar serviço: net start MongoDB" -ForegroundColor White
    Write-Host "  2. Docker: docker run -d -p 27017:27017 mongo:7.0" -ForegroundColor White
    Write-Host ""
    
    $response = Read-Host "Deseja continuar mesmo assim? (s/N)"
    if ($response -ne "s" -and $response -ne "S") {
        Write-Host "Execução cancelada" -ForegroundColor Yellow
        exit 0
    }
}

# Instalar dependências
Write-Host ""
Write-Host "Instalando dependências..." -ForegroundColor Yellow
go mod download
if ($LASTEXITCODE -eq 0) {
    Write-Host "✓ Dependências instaladas" -ForegroundColor Green
}
else {
    Write-Host "✗ Erro ao instalar dependências" -ForegroundColor Red
    exit 1
}

# Limpar módulos
Write-Host ""
Write-Host "Limpando módulos..." -ForegroundColor Yellow
go mod tidy
Write-Host "✓ Módulos limpos" -ForegroundColor Green

# Verificar se swag está instalado (para Swagger)
Write-Host ""
Write-Host "Verificando Swagger..." -ForegroundColor Yellow
try {
    $swagVersion = swag --version 2>$null
    Write-Host "✓ Swag encontrado" -ForegroundColor Green
}
catch {
    Write-Host "! Swag não encontrado" -ForegroundColor Yellow
    Write-Host "  Instalando swag..." -ForegroundColor Yellow
    go install github.com/swaggo/swag/cmd/swag@latest
    if ($LASTEXITCODE -eq 0) {
        Write-Host "✓ Swag instalado" -ForegroundColor Green
    }
    else {
        Write-Host "! Aviso: Não foi possível instalar swag" -ForegroundColor Yellow
        Write-Host "  Swagger pode não funcionar corretamente" -ForegroundColor Yellow
    }
}

# Gerar documentação Swagger
Write-Host ""
Write-Host "Gerando documentação Swagger..." -ForegroundColor Yellow
try {
    swag init -g cmd/server/main.go -o docs 2>$null
    if ($LASTEXITCODE -eq 0) {
        Write-Host "✓ Documentação Swagger gerada" -ForegroundColor Green
    }
}
catch {
    Write-Host "! Aviso: Não foi possível gerar documentação Swagger" -ForegroundColor Yellow
}

# Exibir informações
Write-Host ""
Write-Host "========================================" -ForegroundColor Cyan
Write-Host "   Configuração Concluída!" -ForegroundColor Green
Write-Host "========================================" -ForegroundColor Cyan
Write-Host ""
Write-Host "Para iniciar a aplicação:" -ForegroundColor Yellow
Write-Host "  go run cmd/server/main.go" -ForegroundColor White
Write-Host ""
Write-Host "Ou compilar e executar:" -ForegroundColor Yellow
Write-Host "  go build -o loterias-api.exe cmd/server/main.go" -ForegroundColor White
Write-Host "  .\loterias-api.exe" -ForegroundColor White
Write-Host ""
Write-Host "Endpoints:" -ForegroundColor Yellow
Write-Host "  API: http://localhost:9050/api" -ForegroundColor White
Write-Host "  Swagger: http://localhost:9050/swagger/index.html" -ForegroundColor White
Write-Host ""

# Perguntar se deseja iniciar
$start = Read-Host "Deseja iniciar a aplicação agora? (S/n)"
if ($start -ne "n" -and $start -ne "N") {
    Write-Host ""
    Write-Host "Iniciando aplicação..." -ForegroundColor Green
    Write-Host "Pressione Ctrl+C para parar" -ForegroundColor Yellow
    Write-Host ""
    go run cmd/server/main.go
}
