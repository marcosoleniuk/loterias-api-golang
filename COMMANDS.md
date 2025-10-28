# 📋 Comandos Úteis - Loterias API Golang

## 🚀 Desenvolvimento

```bash
# Executar em modo desenvolvimento
go run cmd/server/main.go

# Executar com hot reload (requer air)
go install github.com/cosmtrek/air@latest
air

# Formatar código
go fmt ./...

# Verificar código
go vet ./...

# Analisar código com linter
golangci-lint run
```

## 🧪 Testes

```bash
# Executar todos os testes
go test ./...

# Testes verbosos
go test -v ./...

# Testes com cobertura
go test -cover ./...

# Gerar relatório de cobertura HTML
go test -coverprofile=coverage.out ./...
go tool cover -html=coverage.out

# Testar pacote específico
go test ./internal/service/...

# Testes em paralelo
go test -parallel 4 ./...

# Benchmark
go test -bench=. ./...
```

## 🔨 Build

```bash
# Build simples
go build -o loterias-api-golang cmd/server/main.go

# Build otimizado (menor tamanho)
go build -ldflags="-s -w" -o loterias-api-golang cmd/server/main.go

# Build para diferentes sistemas
# Linux
GOOS=linux GOARCH=amd64 go build -o loterias-api-golang-linux cmd/server/main.go

# Windows
GOOS=windows GOARCH=amd64 go build -o loterias-api-golang.exe cmd/server/main.go

# macOS
GOOS=darwin GOARCH=amd64 go build -o loterias-api-golang-mac cmd/server/main.go

# ARM (Raspberry Pi, etc)
GOOS=linux GOARCH=arm64 go build -o loterias-api-golang-arm cmd/server/main.go
```

## 📦 Dependências

```bash
# Baixar dependências
go mod download

# Adicionar dependência
go get github.com/package/name

# Remover dependências não utilizadas
go mod tidy

# Verificar dependências
go list -m all

# Atualizar dependência
go get -u github.com/package/name

# Atualizar todas as dependências
go get -u ./...

# Verificar vulnerabilidades
go list -json -m all | nancy sleuth
```

## 🐳 Docker

```bash
# Build da imagem
docker build -t loterias-api-golang:latest .

# Build com tag de versão
docker build -t loterias-api-golang:1.0.0 .

# Executar container
docker run -p 9050:9050 --env-file .env loterias-api-golang

# Executar em background
docker run -d -p 9050:9050 --env-file .env --name loterias loterias-api-golang

# Ver logs
docker logs -f loterias

# Parar container
docker stop loterias

# Remover container
docker rm loterias

# Docker Compose
docker-compose up -d          # Iniciar
docker-compose logs -f api    # Ver logs
docker-compose ps             # Listar containers
docker-compose down           # Parar e remover
docker-compose restart api    # Reiniciar API
```

## 🗄️ MongoDB

```bash
# Conectar ao MongoDB
mongosh

# Ou especificar URI
mongosh "mongodb://localhost:27017/loterias"

# Comandos úteis dentro do mongosh:
use loterias                  # Usar database
show collections              # Listar coleções
db.resultados.count()         # Contar documentos

# Ver últimos resultados
db.resultados.find().sort({"_id.concurso": -1}).limit(5)

# Buscar por loteria
db.resultados.find({"_id.loteria": "megasena"})

# Contar por loteria
db.resultados.count({"_id.loteria": "megasena"})

# Apagar todos os resultados (cuidado!)
db.resultados.deleteMany({})

# Backup
mongodump --db=loterias --out=backup/

# Restore
mongorestore --db=loterias backup/loterias/
```

## 📊 Swagger

```bash
# Gerar documentação
swag init -g cmd/server/main.go -o docs

# Instalar swag se não tiver
go install github.com/swaggo/swag/cmd/swag@latest

# Verificar versão
swag --version

# Formatar comentários
swag fmt
```

## 🔍 Debug e Profiling

```bash
# Executar com race detector
go run -race cmd/server/main.go

# Build com símbolos de debug
go build -gcflags="all=-N -l" -o loterias-api-golang cmd/server/main.go

# CPU profiling
go test -cpuprofile cpu.prof ./...
go tool pprof cpu.prof

# Memory profiling
go test -memprofile mem.prof ./...
go tool pprof mem.prof

# Ver estatísticas de memória em tempo real
curl http://localhost:9050/debug/pprof/heap > heap.prof
go tool pprof heap.prof
```

## 🌐 Testando API

```bash
# Listar loterias
curl http://localhost:9050/api

# Último resultado Mega-Sena
curl http://localhost:9050/api/megasena/latest

# Resultado específico
curl http://localhost:9050/api/megasena/2650

# Todos os resultados da Quina
curl http://localhost:9050/api/quina

# Com pretty print (jq)
curl http://localhost:9050/api/megasena/latest | jq

# Benchmark com Apache Bench
ab -n 1000 -c 10 http://localhost:9050/api/megasena/latest

# Teste de carga com wrk
wrk -t12 -c400 -d30s http://localhost:9050/api
```

## 📝 Logs

```bash
# Ver logs em tempo real (Linux/Mac)
tail -f build-errors.log

# Buscar erros nos logs
grep "error" build-errors.log

# Contar requisições
grep "GET" build-errors.log | wc -l
```

## 🔧 Manutenção

```bash
# Verificar versão do Go
go version

# Limpar cache
go clean -cache
go clean -modcache

# Verificar espaço usado
du -sh ~/go/pkg/mod

# Verificar variáveis de ambiente
go env

# Definir GOPATH
export GOPATH=$HOME/go

# Adicionar ao PATH
export PATH=$PATH:$GOPATH/bin
```

## 🚀 Deploy

```bash
# Build para produção
CGO_ENABLED=0 go build -ldflags="-s -w" -o loterias-api-golang cmd/server/main.go

# Criar tarball
tar -czf loterias-api-golang-v1.0.0.tar.gz loterias-api-golang .env.example

# Deploy via SCP
scp loterias-api-golang user@server:/opt/loterias-api-golang/

# Systemd service (Linux)
sudo systemctl start loterias-api-golang
sudo systemctl enable loterias-api-golang
sudo systemctl status loterias-api-golang

# Ver logs do systemd
journalctl -u loterias-api-golang -f
```

## 🎯 Produtividade

```bash
# Usar Make para comandos comuns
make run          # Executar
make build        # Compilar
make test         # Testar
make clean        # Limpar
make swagger      # Gerar Swagger
make docker-build # Build Docker
make docker-run   # Rodar Docker Compose

# Aliases úteis (adicionar ao .bashrc ou .zshrc)
alias gor='go run cmd/server/main.go'
alias gob='go build -o loterias-api-golang cmd/server/main.go'
alias got='go test ./...'
alias gof='go fmt ./...'
```

## 🔐 Segurança

```bash
# Verificar vulnerabilidades
go list -json -m all | nancy sleuth

# Instalar nancy
go install github.com/sonatype-nexus-community/nancy@latest

# Scan com govulncheck
govulncheck ./...

# Instalar govulncheck
go install golang.org/x/vuln/cmd/govulncheck@latest
```

## 📈 Performance

```bash
# Benchmark de endpoints
go-wrk -c 100 -d 30 http://localhost:9050/api/megasena/latest

# Instalar go-wrk
go install github.com/tsliwowicz/go-wrk@latest

# Verificar goroutines
curl http://localhost:9050/debug/pprof/goroutine

# Memory stats
curl http://localhost:9050/debug/pprof/heap
```

---

**Dica**: Salve este arquivo como referência rápida! 📚
