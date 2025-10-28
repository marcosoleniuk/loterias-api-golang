# Instruções de Execução - Loterias API (Go)

## 🚀 Início Rápido

### 1. Pré-requisitos

Certifique-se de ter instalado:

- **Go 1.25+**: [Download aqui](https://go.dev/dl/)
- **MongoDB 4.4+**: [Download aqui](https://www.mongodb.com/try/download/community)
- **Git**: Para clonar o repositório

### 2. Configuração Inicial

```powershell
# Clone o repositório (se ainda não tiver)
git clone https://github.com/marcosoleniuk/loterias-api-golang.git
cd loterias-api-golang

# Crie o arquivo .env a partir do exemplo
copy .env.example .env

# Edite o .env conforme necessário (opcional)
notepad .env
```

### 3. Instalar Dependências

```powershell
# Download de todas as dependências
go mod download

# Verificar e limpar dependências
go mod tidy
```

### 4. Iniciar MongoDB

**Opção A - MongoDB Local:**

```powershell
# Iniciar serviço do MongoDB
net start MongoDB
```

**Opção B - MongoDB via Docker:**

```powershell
docker run -d -p 27017:27017 --name loterias-mongodb mongo:latest
```

### 5. Executar a Aplicação

**Modo Desenvolvimento:**

```powershell
# Executar diretamente com Go
go run cmd/server/main.go
```

**Compilar e Executar:**

```powershell
# Compilar
go build -o loterias-api-golang.exe cmd/server/main.go

# Executar binário
.\loterias-api-golang.exe
```

### 6. Verificar se está Funcionando

Abra seu navegador em:

- **API Root**: http://localhost:9050/
- **Swagger**: http://localhost:9050/swagger/index.html
- **Lista de Loterias**: http://localhost:9050/api

## 📝 Comandos Úteis

### Desenvolvimento

```powershell
# Executar com hot reload (usando air)
go install github.com/cosmtrek/air@latest
air

# Formatar código
go fmt ./...

# Verificar código
go vet ./...

# Executar testes
go test ./...

# Executar testes com cobertura
go test -cover ./...
```

### Build

```powershell
# Build para Windows
go build -o loterias-api-golang.exe cmd/server/main.go

# Build para Linux
$env:GOOS="linux"; $env:GOARCH="amd64"; go build -o loterias-api-golang-linux cmd/server/main.go

# Build otimizado (menor tamanho)
go build -ldflags="-s -w" -o loterias-api-golang.exe cmd/server/main.go
```

### Docker

```powershell
# Build da imagem
docker build -t loterias-api-golang .

# Executar com Docker Compose (MongoDB + API)
docker-compose up -d

# Ver logs
docker-compose logs -f

# Parar containers
docker-compose down
```

## 🧪 Testando a API

### Com curl (PowerShell)

```powershell
# Listar todas as loterias
curl http://localhost:9050/api

# Último resultado da Mega-Sena
curl http://localhost:9050/api/megasena/latest

# Resultado específico
curl http://localhost:9050/api/megasena/2650

# Todos os resultados da Quina
curl http://localhost:9050/api/quina
```

### Com Invoke-WebRequest (PowerShell nativo)

```powershell
# Listar loterias
Invoke-WebRequest -Uri "http://localhost:9050/api" | Select-Object -Expand Content

# Último resultado
$response = Invoke-WebRequest -Uri "http://localhost:9050/api/megasena/latest"
$response.Content | ConvertFrom-Json | ConvertTo-Json -Depth 10
```

### Com Postman

1. Importe a coleção Swagger: http://localhost:9050/swagger/doc.json
2. Ou crie requisições manualmente para os endpoints

## 🔧 Troubleshooting

### Erro: "go: command not found"

**Solução**: Instale o Go e adicione ao PATH

```powershell
# Verificar instalação
go version
```

### Erro: "cannot connect to MongoDB"

**Solução 1**: Verificar se MongoDB está rodando

```powershell
# Windows
Get-Service MongoDB

# Iniciar serviço
net start MongoDB
```

**Solução 2**: Verificar MONGODB_URI no .env

```env
MONGODB_URI=mongodb://localhost:27017/loterias
```

### Erro: "port 9050 already in use"

**Solução**: Alterar porta no .env

```env
PORT=8091
```

### Erro ao compilar: "package not found"

**Solução**: Reinstalar dependências

```powershell
go mod download
go mod tidy
```

### Swagger não aparece

**Solução**: Gerar documentação Swagger

```powershell
# Instalar swag
go install github.com/swaggo/swag/cmd/swag@latest

# Gerar docs
swag init -g cmd/server/main.go -o docs
```

## 📊 Monitoramento

### Ver logs em tempo real

```powershell
# A aplicação imprime logs no console
# Busque por mensagens como:
# - "Connected to MongoDB successfully"
# - "Starting server on port 9050"
# - "Running scheduled lottery update..."
```

### Verificar banco de dados

```powershell
# Conectar ao MongoDB
mongosh

# Usar database
use loterias

# Ver coleções
show collections

# Contar documentos
db.resultados.countDocuments()

# Ver último resultado da Mega-Sena
db.resultados.find({"_id.loteria": "megasena"}).sort({"_id.concurso": -1}).limit(1)
```

## 🎯 Próximos Passos

1. **Explore a API**: Teste todos os endpoints via Swagger
2. **Monitore atualizações**: Logs mostram quando novos resultados são salvos
3. **Customize**: Ajuste o cron schedule em `scheduler/scheduled_consumer.go`
4. **Adicione features**: Cache, métricas, autenticação, etc.

## 📚 Documentação Adicional

- **Swagger UI**: http://localhost:9050/swagger/index.html
- **README.md**: Visão geral do projeto

## 💬 Suporte

Em caso de dúvidas ou problemas:

1. Verifique os logs da aplicação
2. Consulte o README.md
3. Abra uma issue no GitHub

---

**Boa sorte e bom desenvolvimento! 🚀**
