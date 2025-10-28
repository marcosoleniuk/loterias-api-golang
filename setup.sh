#!/bin/bash

# ============================================
# Loterias API - Script de Inicialização (Linux/Mac)
# ============================================

echo "========================================"
echo "   Loterias API Golang - Inicialização"
echo "========================================"
echo ""

# Verificar se Go está instalado
echo "Verificando instalação do Go..."
if ! command -v go &> /dev/null; then
    echo "✗ Go não encontrado! Por favor, instale Go 1.25+"
    echo "  Download: https://go.dev/dl/"
    exit 1
fi

GO_VERSION=$(go version)
echo "✓ Go encontrado: $GO_VERSION"

# Verificar se arquivo .env existe
echo ""
echo "Verificando arquivo .env..."
if [ ! -f ".env" ]; then
    echo "! Arquivo .env não encontrado"
    if [ -f ".env.example" ]; then
        echo "  Criando .env a partir de .env.example..."
        cp .env.example .env
        echo "✓ Arquivo .env criado"
        echo "  Por favor, revise as configurações em .env"
    else
        echo "✗ Arquivo .env.example não encontrado"
        exit 1
    fi
else
    echo "✓ Arquivo .env encontrado"
fi

# Verificar MongoDB
echo ""
echo "Verificando MongoDB..."
MONGO_RUNNING=false

# Tentar conectar ao MongoDB
if command -v mongosh &> /dev/null; then
    if mongosh --eval "db.version()" --quiet &> /dev/null; then
        MONGO_RUNNING=true
        echo "✓ MongoDB está rodando"
    fi
elif pgrep mongod &> /dev/null; then
    MONGO_RUNNING=true
    echo "✓ MongoDB está rodando (processo detectado)"
fi

if [ "$MONGO_RUNNING" = false ]; then
    echo "! MongoDB não está rodando"
    echo "  Opções:"
    echo "  1. Iniciar serviço: sudo systemctl start mongodb"
    echo "  2. Docker: docker run -d -p 27017:27017 mongo:7.0"
    echo ""
    
    read -p "Deseja continuar mesmo assim? (s/N) " response
    if [ "$response" != "s" ] && [ "$response" != "S" ]; then
        echo "Execução cancelada"
        exit 0
    fi
fi

# Instalar dependências
echo ""
echo "Instalando dependências..."
go mod download
if [ $? -eq 0 ]; then
    echo "✓ Dependências instaladas"
else
    echo "✗ Erro ao instalar dependências"
    exit 1
fi

# Limpar módulos
echo ""
echo "Limpando módulos..."
go mod tidy
echo "✓ Módulos limpos"

# Verificar se swag está instalado (para Swagger)
echo ""
echo "Verificando Swagger..."
if ! command -v swag &> /dev/null; then
    echo "! Swag não encontrado"
    echo "  Instalando swag..."
    go install github.com/swaggo/swag/cmd/swag@latest
    if [ $? -eq 0 ]; then
        echo "✓ Swag instalado"
    else
        echo "! Aviso: Não foi possível instalar swag"
        echo "  Swagger pode não funcionar corretamente"
    fi
else
    echo "✓ Swag encontrado"
fi

# Gerar documentação Swagger
echo ""
echo "Gerando documentação Swagger..."
if command -v swag &> /dev/null; then
    swag init -g cmd/server/main.go -o docs 2> /dev/null
    if [ $? -eq 0 ]; then
        echo "✓ Documentação Swagger gerada"
    fi
else
    echo "! Aviso: Não foi possível gerar documentação Swagger"
fi

# Exibir informações
echo ""
echo "========================================"
echo "   Configuração Concluída!"
echo "========================================"
echo ""
echo "Para iniciar a aplicação:"
echo "  go run cmd/server/main.go"
echo ""
echo "Ou compilar e executar:"
echo "  go build -o loterias-api cmd/server/main.go"
echo "  ./loterias-api"
echo ""
echo "Endpoints:"
echo "  API: http://localhost:9050/api"
echo "  Swagger: http://localhost:9050/swagger/index.html"
echo ""

# Perguntar se deseja iniciar
read -p "Deseja iniciar a aplicação agora? (S/n) " start
if [ "$start" != "n" ] && [ "$start" != "N" ]; then
    echo ""
    echo "Iniciando aplicação..."
    echo "Pressione Ctrl+C para parar"
    echo ""
    go run cmd/server/main.go
fi
