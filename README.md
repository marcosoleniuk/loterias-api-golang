# 🎰 Loterias API Caixa Econômica - Golang

<div align="center">

![Go Version](https://img.shields.io/badge/Go-1.25+-00ADD8?style=for-the-badge&logo=go)
![MongoDB](https://img.shields.io/badge/MongoDB-4.4+-47A248?style=for-the-badge&logo=mongodb&logoColor=white)
![License](https://img.shields.io/badge/License-Apache2-green?style=for-the-badge)
![Docker](https://img.shields.io/badge/Docker-Ready-2496ED?style=for-the-badge&logo=docker&logoColor=white)

**API REST para consulta de resultados das loterias da Caixa Econômica Federal**

🌐 **[API em Produção](https://api-loterias.moleniuk.com/)** | 📚 **[Documentação Swagger](https://api-loterias.moleniuk.com/swagger/index.html)**

[Documentação](#-documentação) •
[Instalação](#-instalação) •
[Uso](#-uso) •
[API Endpoints](#-endpoints) •
[Contribuir](#-contribuindo)

</div>

---

## ⚡ Quick Start

Não quer instalar nada? Use a API diretamente em produção:

```bash
# Listar todas as loterias
curl https://api-loterias.moleniuk.com/api

# Resultado mais recente da Mega-Sena
curl https://api-loterias.moleniuk.com/api/megasena/latest

# Documentação interativa
https://api-loterias.moleniuk.com/swagger/index.html
```

---

## 📖 Sobre o Projeto

A **Loterias API Golang** é uma API REST desenvolvida em Go que fornece acesso programático aos resultados históricos e atuais das principais loterias da Caixa Econômica Federal. O projeto inclui:

- ✅ Consulta de resultados de **10 loterias diferentes**
- ✅ Atualização automática via **scheduler** (cron jobs)
- ✅ Persistência de dados com **MongoDB**
- ✅ Documentação interativa com **Swagger**
- ✅ Suporte a **Docker** e **Docker Compose**
- ✅ API RESTful com padrões modernos
- ✅ CORS configurado para integração frontend

### 🎲 Loterias Disponíveis

| Loteria      | ID               | Descrição                        |
| ------------ | ---------------- | -------------------------------- |
| Mega-Sena    | `megasena`       | A loteria mais popular do Brasil |
| Lotofácil    | `lotofacil`      | Facilita com 25 números          |
| Quina        | `quina`          | Sorteios diários                 |
| Lotomania    | `lotomania`      | 50 números para escolher         |
| Timemania    | `timemania`      | A loteria dos times de futebol   |
| Dupla Sena   | `duplasena`      | Dois sorteios em um              |
| Federal      | `federal`        | Prêmios fixos garantidos         |
| Dia de Sorte | `diadesorte`     | Escolha seu mês da sorte         |
| Super Sete   | `supersete`      | 7 colunas de números             |
| +Milionária  | `maismilionaria` | Prêmios milionários              |

---

## 🚀 Tecnologias

Este projeto foi desenvolvido com as seguintes tecnologias:

- **[Go 1.25+](https://go.dev/)** - Linguagem de programação
- **[Gin Web Framework](https://gin-gonic.com/)** - Framework HTTP web
- **[MongoDB](https://www.mongodb.com/)** - Banco de dados NoSQL
- **[Swagger/OpenAPI](https://swagger.io/)** - Documentação da API
- **[Docker](https://www.docker.com/)** - Containerização
- **[Cron](https://github.com/robfig/cron)** - Agendamento de tarefas

### 📦 Principais Dependências

```go
github.com/gin-gonic/gin           // Framework web
go.mongodb.org/mongo-driver        // Driver MongoDB
github.com/swaggo/gin-swagger      // Documentação Swagger
github.com/robfig/cron/v3          // Scheduler
github.com/joho/godotenv           // Gerenciamento de variáveis de ambiente
```

---

## 📋 Pré-requisitos

Antes de começar, você precisará ter instalado:

- **Go 1.25 ou superior** - [Download](https://go.dev/dl/)
- **MongoDB 4.4 ou superior** - [Download](https://www.mongodb.com/try/download/community)
- **Git** - [Download](https://git-scm.com/)
- **Docker** (opcional) - [Download](https://www.docker.com/)

---

## 🔧 Instalação

### Método 1: Instalação Local

#### 1. Clone o repositório

```powershell
git clone https://github.com/marcosoleniuk/loterias-api-golang.git
cd loterias-api-golang
```

#### 2. Configure as variáveis de ambiente

Crie um arquivo `.env` na raiz do projeto:

```env
# Configuração do Servidor
PORT=9050
GIN_MODE=debug

# Configuração do MongoDB
MONGODB_URI=mongodb://localhost:27017/loterias

# Configuração do Scheduler (formato cron)
CRON_SCHEDULE=0 22 * * *
```

#### 3. Instale as dependências

```powershell
go mod download
go mod tidy
```

#### 4. Inicie o MongoDB

```powershell
# Windows - Serviço local
net start MongoDB

# Ou via Docker
docker run -d -p 27017:27017 --name loterias-mongodb mongo:latest
```

#### 5. Execute a aplicação

```powershell
# Executar diretamente
go run cmd/server/main.go

# Ou compilar e executar
go build -o loterias-api.exe cmd/server/main.go
./loterias-api.exe
```

### Método 2: Docker Compose (Recomendado)

#### 1. Clone o repositório

```powershell
git clone https://github.com/marcosoleniuk/loterias-api-golang.git
cd loterias-api-golang
```

#### 2. Configure o `.env` (opcional)

O Docker Compose já vem com configurações padrão.

#### 3. Execute com Docker Compose

```powershell
# Iniciar serviços
docker-compose up -d

# Visualizar logs
docker-compose logs -f

# Parar serviços
docker-compose down
```

---

## 🎯 Uso

### 🌐 API em Produção

Você pode usar a API diretamente em produção sem precisar instalar nada:

**URL Base:** `https://api-loterias.moleniuk.com/api`

**Swagger Docs:** `https://api-loterias.moleniuk.com/swagger/index.html`

### 🏠 Uso Local

#### Verificar Status da API

Acesse no navegador:

```
http://localhost:9050/
```

#### Documentação Interativa (Swagger)

```
http://localhost:9050/swagger/index.html
```

### Exemplos de Requisições

#### 🌐 Usando a API em Produção

##### Listar todas as loterias disponíveis

```bash
curl https://api-loterias.moleniuk.com/api
```

##### Buscar resultado mais recente da Mega-Sena

```bash
curl https://api-loterias.moleniuk.com/api/megasena/latest
```

##### Buscar resultado de concurso específico

```bash
curl https://api-loterias.moleniuk.com/api/megasena/2500
```

##### Buscar todos os resultados da Lotofácil

```bash
curl https://api-loterias.moleniuk.com/api/lotofacil
```

#### 🏠 Usando a API Local

##### Listar todas as loterias disponíveis

```bash
curl http://localhost:9050/api
```

##### Buscar resultado mais recente da Mega-Sena

```bash
curl https://api-loterias.moleniuk.com/api/megasena/latest
```

### Formato de Resposta

**Exemplo de resposta (resultado mais recente):**

```json
{
  "loteria": "megasena",
  "concurso": 2932,
  "data": "25/10/2025",
  "local": "ESPAÇO DA SORTE em SÃO PAULO, SP",
  "dezenasOrdemSorteio": ["40", "04", "53", "25", "36", "13"],
  "dezenas": ["04", "13", "25", "36", "40", "53"],
  "premiacoes": [
    {
      "descricao": "6 acertos",
      "faixa": 1,
      "numeroDeGanhadores": 1,
      "valor": 96166949.14
    },
    {
      "descricao": "5 acertos",
      "faixa": 2,
      "numeroDeGanhadores": 144,
      "valor": 27798.3
    },
    {
      "descricao": "4 acertos",
      "faixa": 3,
      "numeroDeGanhadores": 6869,
      "valor": 960.58
    }
  ],
  "municipiosUFGanhadores": [
    {
      "ganhadores": 1,
      "municipio": "OURINHOS",
      "posicao": 1,
      "uf": "SP"
    }
  ],
  "acumulou": false,
  "proximoConcurso": 2933,
  "dataProximoConcurso": "28/10/2025",
  "valorArrecadado": 100453338,
  "valorAcumuladoConcurso_0_5": 16774002.79,
  "valorAcumuladoConcursoEspecial": 133392230.93,
  "valorEstimadoProximoConcurso": 3500000
}
```

---

## 📚 Endpoints

### Base URL

**Produção:** `https://api-loterias.moleniuk.com/api`

**Local:** `http://localhost:9050/api`

### Rotas Disponíveis

| Método | Endpoint                    | Descrição                                   |
| ------ | --------------------------- | ------------------------------------------- |
| `GET`  | `/api`                      | Lista todas as loterias disponíveis         |
| `GET`  | `/api/{loteria}`            | Retorna todos os resultados de uma loteria  |
| `GET`  | `/api/{loteria}/latest`     | Retorna o resultado mais recente            |
| `GET`  | `/api/{loteria}/{concurso}` | Retorna resultado de um concurso específico |

### Parâmetros

- `{loteria}`: ID da loteria (ex: `megasena`, `lotofacil`)
- `{concurso}`: Número do concurso (ex: `2650`)

### Respostas

#### Sucesso (200)

```json
{
  "loteria": "string",
  "concurso": "number",
  "data": "date",
  "dezenas": ["string"],
  "premiacoes": [
    {
      "acertos": "number",
      "vencedores": "number",
      "valorPremio": "number"
    }
  ]
}
```

#### Erro (404)

```json
{
  "error": "Resource Not Found",
  "message": "Invalid lottery ID. Available lotteries: megasena, lotofacil, ..."
}
```

---

## 🏗️ Arquitetura do Projeto

```
loterias-api-golang/
├── cmd/
│   └── server/
│       └── main.go                 # Entry point da aplicação
├── internal/
│   ├── config/
│   │   └── cors.go                 # Configuração CORS
│   ├── controller/
│   │   ├── api_controller.go       # Handlers HTTP
│   │   └── root_controller.go      # Rota raiz
│   ├── model/
│   │   ├── loteria.go              # Modelo de loterias
│   │   ├── resultado.go            # Modelo de resultados
│   │   └── exceptions.go           # Tratamento de erros
│   ├── repository/
│   │   └── resultado_repository.go # Acesso ao MongoDB
│   ├── scheduler/
│   │   └── scheduled_consumer.go   # Cron jobs
│   └── service/
│       ├── consumer.go             # Consumo da API Caixa
│       ├── resultado_service.go    # Lógica de negócio
│       └── loterias_update.go      # Atualização de dados
├── docs/
│   ├── docs.go                     # Documentação Swagger
│   ├── swagger.json
│   └── swagger.yaml
├── docker-compose.yml              # Orquestração Docker
├── Dockerfile                      # Imagem Docker
├── Makefile                        # Comandos make
├── go.mod                          # Dependências Go
└── .env                            # Variáveis de ambiente
```

### Fluxo de Dados

```
┌─────────────┐
│   Cliente   │
└──────┬──────┘
       │ HTTP Request
       ▼
┌─────────────────┐
│  Gin Router     │
│  (CORS/Logger)  │
└────────┬────────┘
         │
         ▼
┌─────────────────┐
│  Controller     │
│  (Handlers)     │
└────────┬────────┘
         │
         ▼
┌─────────────────┐
│    Service      │
│ (Business Logic)│
└────────┬────────┘
         │
         ▼
┌─────────────────┐
│   Repository    │
│   (Data Access) │
└────────┬────────┘
         │
         ▼
┌─────────────────┐
│    MongoDB      │
└─────────────────┘
```

### Scheduler (Atualização Automática)

O sistema possui um **scheduler** que executa automaticamente:

- **Horário padrão**: Diariamente às 22:00 (após os sorteios)
- **Configurável via**: Variável `CRON_SCHEDULE` no `.env`
- **Função**: Busca os resultados mais recentes e atualiza o banco

```go
// Formato Cron: minuto hora dia mês dia-da-semana
// Exemplo: "0 22 * * *" = Todos os dias às 22:00
```

---

## 🛠️ Comandos Úteis

### Usando Make

```powershell
# Ver todos os comandos disponíveis
make help

# Executar aplicação
make run

# Compilar binário
make build

# Executar testes
make test

# Testes com cobertura
make test-cover

# Gerar documentação Swagger
make swagger

# Build Docker
make docker-build

# Executar com Docker Compose
make docker-run
```

### Comandos Go Diretos

```powershell
# Desenvolvimento
go run cmd/server/main.go

# Build
go build -o loterias-api.exe cmd/server/main.go

# Testes
go test ./...
go test -v ./internal/service/...

# Formatação
go fmt ./...

# Linter
go vet ./...

# Atualizar dependências
go mod tidy
go mod download
```

### Gerar Documentação Swagger

```powershell
# Instalar swag
go install github.com/swaggo/swag/cmd/swag@latest

# Gerar documentação
swag init -g cmd/server/main.go -o docs

# A documentação estará disponível em /swagger/index.html
```

---

## 🧪 Testes

### Executar Testes

```powershell
# Todos os testes
go test ./...

# Testes verbosos
go test -v ./...

# Testes com cobertura
go test -cover ./...

# Gerar relatório HTML de cobertura
go test -coverprofile=coverage.out ./...
go tool cover -html=coverage.out
```

### Estrutura de Testes

```
internal/
└── service/
    ├── consumer.go
    └── consumer_test.go      # Testes do consumer
```

---

## 🐳 Docker

### Build da Imagem

```powershell
docker build -t loterias-api-golang .
```

### Executar Container

```powershell
docker run -p 9050:9050 --env-file .env loterias-api-golang
```

### Docker Compose

```powershell
# Iniciar serviços (API + MongoDB)
docker-compose up -d

# Ver logs
docker-compose logs -f api

# Parar serviços
docker-compose down

# Remover volumes
docker-compose down -v
```

---

## ⚙️ Configuração

### Variáveis de Ambiente

Crie um arquivo `.env` na raiz do projeto:

```env
# Porta do servidor
PORT=9050

# Modo do Gin (debug, release, test)
GIN_MODE=debug

# URI do MongoDB
MONGODB_URI=mongodb://localhost:27017/loterias

# Schedule do cron (formato: minuto hora dia mês dia-da-semana)
# Padrão: Todos os dias às 22:00
CRON_SCHEDULE=0 22 * * *
```

### Configuração de CORS

O CORS já está configurado no arquivo `internal/config/cors.go` para aceitar:

- ✅ Todas as origens (`*`)
- ✅ Métodos: GET, POST, PUT, DELETE, OPTIONS
- ✅ Headers comuns

Para restringir origens em produção, edite:

```go
// internal/config/cors.go
config.AllowOrigins = []string{"https://seu-dominio.com"}
```

---

## 📊 Monitoramento

### Logs

A aplicação usa o logger padrão do Gin que exibe:

- Método HTTP
- Path da requisição
- Status code
- Latência
- IP do cliente

Exemplo:

```
[GIN] 2024/10/27 - 14:30:45 | 200 |    2.456789ms |    192.168.1.10 | GET      "/api/megasena/latest"
```

### Health Check

```bash
curl http://localhost:9050/
```

---

## 🚀 Deploy

### Deploy em Produção

1. **Configure as variáveis de ambiente**

```env
GIN_MODE=release
MONGODB_URI=mongodb://seu-servidor:27017/loterias
PORT=9050
```

2. **Compile o binário otimizado**

```powershell
go build -ldflags="-s -w" -o loterias-api cmd/server/main.go
```

3. **Execute o binário**

```powershell
./loterias-api
```

### Deploy com Docker

```powershell
# Build da imagem
docker build -t loterias-api:v1.0 .

# Push para registry (opcional)
docker tag loterias-api:v1.0 seu-registry/loterias-api:v1.0
docker push seu-registry/loterias-api:v1.0

# Deploy
docker run -d \
  -p 9050:9050 \
  -e GIN_MODE=release \
  -e MONGODB_URI=mongodb://mongo:27017/loterias \
  --name loterias-api \
  loterias-api:v1.0
```

---

## 🤝 Contribuindo

Contribuições são sempre bem-vindas!

1. Faça um Fork do projeto
2. Crie uma branch para sua feature (`git checkout -b feature/MinhaFeature`)
3. Commit suas mudanças (`git commit -m 'Adiciona nova feature'`)
4. Push para a branch (`git push origin feature/MinhaFeature`)
5. Abra um Pull Request

### Padrões de Código

- Siga as convenções do Go ([Effective Go](https://go.dev/doc/effective_go))
- Use `go fmt` para formatar o código
- Execute `go vet` antes de commitar
- Escreva testes para novas funcionalidades
- Atualize a documentação Swagger quando necessário

---

## 📝 Licença

Este projeto está sob a licença Apache-2.0. Veja o arquivo [LICENSE](LICENSE) para mais detalhes.

---

## 👨‍💻 Autor

**Marcos Oleniuk**

- 📧 Email: marcos@moleniuk.com
- 💬 WhatsApp: [+55 44 9 98425-745](https://wa.me/554998425745)

---

## 🙏 Agradecimentos

- Caixa Econômica Federal pela API pública de loterias
- Comunidade Go pela excelente documentação

---

## 📚 Documentação Adicional

- [SETUP.md](SETUP.md) - Guia detalhado de instalação
- [COMMANDS.md](COMMANDS.md) - Lista completa 
- [Swagger Docs](https://api-loterias.moleniuk.com/swagger/index.html) - API interativa

---

<div align="center">

**⭐ Se este projeto foi útil, considere dar uma estrela!**

Made with ❤️ and Go

</div>
