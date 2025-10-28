# Build stage
FROM golang:1.25-alpine AS builder

WORKDIR /app

# Copiar arquivos de dependências
COPY go.mod go.sum ./

# Download de dependências
RUN go mod download

# Copiar código fonte
COPY . .

# Build da aplicação
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o loterias-api-golang ./cmd/server/main.go

# Runtime stage
FROM alpine:latest

RUN apk update && \
    apk add tzdata && \
    cp /usr/share/zoneinfo/America/Sao_Paulo /etc/localtime && \
    echo "America/Sao_Paulo" > /etc/timezone && \
    apk del tzdata

WORKDIR /root/

# Copiar binário do build stage
COPY --from=builder /app/loterias-api-golang .
COPY .env .env
COPY ./docs ./docs

# Expor porta
EXPOSE 9050

# Comando para executar
CMD ["./loterias-api-golang"]
