.PHONY: run build test clean swagger docker-build docker-run help

# Variáveis
BINARY_NAME=loterias-api-golang
MAIN_PATH=./cmd/server/main.go

## help: Exibir este menu de ajuda
help:
	@echo "Comandos disponíveis:"
	@echo "  make run           - Executar a aplicação em modo desenvolvimento"
	@echo "  make build         - Compilar o binário"
	@echo "  make test          - Executar testes"
	@echo "  make test-cover    - Executar testes com cobertura"
	@echo "  make clean         - Limpar binários e arquivos temporários"
	@echo "  make swagger       - Gerar documentação Swagger"
	@echo "  make docker-build  - Build da imagem Docker"
	@echo "  make docker-run    - Executar com Docker Compose"
	@echo "  make docker-down   - Parar containers Docker"
	@echo "  make lint          - Executar linter"
	@echo "  make fmt           - Formatar código"

## run: Executar a aplicação
run:
	go run $(MAIN_PATH)

## build: Compilar o binário
build:
	@echo "Building..."
	go build -o $(BINARY_NAME) $(MAIN_PATH)
	@echo "Build complete!"

## build-linux: Compilar para Linux
build-linux:
	@echo "Building for Linux..."
	GOOS=linux GOARCH=amd64 go build -o $(BINARY_NAME)-linux $(MAIN_PATH)
	@echo "Build complete!"

## build-windows: Compilar para Windows
build-windows:
	@echo "Building for Windows..."
	GOOS=windows GOARCH=amd64 go build -o $(BINARY_NAME).exe $(MAIN_PATH)
	@echo "Build complete!"

## test: Executar testes
test:
	@echo "Running tests..."
	go test -v ./...

## test-cover: Executar testes com cobertura
test-cover:
	@echo "Running tests with coverage..."
	go test -cover ./...
	go test -coverprofile=coverage.out ./...
	go tool cover -html=coverage.out -o coverage.html
	@echo "Coverage report generated: coverage.html"

## clean: Limpar binários e arquivos temporários
clean:
	@echo "Cleaning..."
	@if exist $(BINARY_NAME) del $(BINARY_NAME)
	@if exist $(BINARY_NAME).exe del $(BINARY_NAME).exe
	@if exist $(BINARY_NAME)-linux del $(BINARY_NAME)-linux
	@if exist coverage.out del coverage.out
	@if exist coverage.html del coverage.html
	@echo "Clean complete!"

## swagger: Gerar documentação Swagger
swagger:
	@echo "Generating Swagger documentation..."
	swag init -g cmd/server/main.go -o docs
	@echo "Swagger documentation generated!"

## docker-build: Build da imagem Docker
docker-build:
	@echo "Building Docker image..."
	docker build -t $(BINARY_NAME):latest .
	@echo "Docker image built!"

## docker-run: Executar com Docker Compose
docker-run:
	@echo "Starting Docker Compose..."
	docker-compose up -d
	@echo "Docker Compose started!"

## docker-down: Parar containers Docker
docker-down:
	@echo "Stopping Docker Compose..."
	docker-compose down
	@echo "Docker Compose stopped!"

## lint: Executar linter
lint:
	@echo "Running linter..."
	golangci-lint run
	@echo "Linting complete!"

## fmt: Formatar código
fmt:
	@echo "Formatting code..."
	go fmt ./...
	@echo "Formatting complete!"

## deps: Baixar dependências
deps:
	@echo "Downloading dependencies..."
	go mod download
	@echo "Dependencies downloaded!"

## tidy: Limpar dependências não utilizadas
tidy:
	@echo "Tidying dependencies..."
	go mod tidy
	@echo "Dependencies tidied!"

## install-tools: Instalar ferramentas de desenvolvimento
install-tools:
	@echo "Installing development tools..."
	go install github.com/swaggo/swag/cmd/swag@latest
	go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest
	@echo "Tools installed!"
