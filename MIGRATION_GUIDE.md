# Guia de Migração: Java/Spring Boot → Go/Gin

Este documento detalha a conversão do projeto Loterias API de Java/Spring Boot para Go/Gin.

## 📊 Visão Geral da Conversão

### Arquitetura

A arquitetura hexagonal foi mantida com as seguintes camadas:

| Java/Spring Boot            | Go/Gin                        |
| --------------------------- | ----------------------------- |
| `@RestController`           | `controller/` package         |
| `@Service`                  | `service/` package            |
| `@Repository` / Spring Data | `repository/` package         |
| `@Document` / Entity        | `model/` package              |
| `@Configuration`            | `config/` package             |
| `@Scheduled`                | `scheduler/` package com cron |

### Dependências

| Java                      | Go                                      |
| ------------------------- | --------------------------------------- |
| Spring Boot Starter Web   | github.com/gin-gonic/gin                |
| Spring Data MongoDB       | go.mongodb.org/mongo-driver             |
| Spring Boot Starter Cache | Cache em memória (pode adicionar Redis) |
| Springfox/Swagger         | github.com/swaggo/gin-swagger           |
| JSoup                     | net/http + encoding/json                |
| Jackson                   | encoding/json (nativo)                  |
| Spring Scheduling         | github.com/robfig/cron                  |

## 🔄 Mapeamento de Conceitos

### 1. Injeção de Dependências

**Java (Spring):**

```java
@Autowired
private ResultadoService resultadoService;
```

**Go:**

```go
// Dependências passadas via construtor
func NewApiController(resultadoService *service.ResultadoService) *ApiController {
    return &ApiController{
        resultadoService: resultadoService,
    }
}
```

### 2. Configuração

**Java (application.properties):**

```properties
server.port=9050
spring.data.mongodb.uri=mongodb://localhost/loterias
```

**Go (.env):**

```env
PORT=9050
MONGODB_URI=mongodb://localhost:27017/loterias
```

### 3. Controllers/Handlers

**Java:**

```java
@RestController
@RequestMapping("api")
public class ApiRestController {
    @GetMapping("{loteria}")
    public ResponseEntity<List<Resultado>> getResults(@PathVariable String loteria) {
        return ResponseEntity.ok(service.findByLoteria(loteria));
    }
}
```

**Go:**

```go
func (c *ApiController) GetResultsByLottery(ctx *gin.Context) {
    loteria := ctx.Param("loteria")
    resultados, err := c.resultadoService.FindByLoteria(loteria)
    if err != nil {
        ctx.JSON(http.StatusInternalServerError, ErrorResponse{...})
        return
    }
    ctx.JSON(http.StatusOK, resultados)
}
```

### 4. Models/Entities

**Java:**

```java
@Document("resultados")
public class Resultado {
    @Id
    private ResultadoId id;
    private String data;
    // getters e setters
}
```

**Go:**

```go
type Resultado struct {
    ID   ResultadoID `bson:"_id" json:"-"`
    Data string      `bson:"data" json:"data"`
}
```

### 5. Repository

**Java (Spring Data):**

```java
@Repository
public interface ResultadoRepository extends MongoRepository<Resultado, ResultadoId> {
    List<Resultado> findById_Loteria(String loteria);
}
```

**Go:**

```go
func (r *ResultadoRepository) FindByLoteria(loteria string) ([]model.Resultado, error) {
    filter := bson.M{"_id.loteria": loteria}
    cursor, err := r.collection.Find(ctx, filter, opts)
    // ...
}
```

### 6. Agendamento

**Java:**

```java
@Scheduled(cron = "0 0 * * * *")
public void checkForUpdates() {
    loteriasUpdate.checkForUpdates();
}
```

**Go:**

```go
c := cron.New()
c.AddFunc("0 * * * *", func() {
    loteriasUpdate.UpdateAll()
})
c.Start()
```

### 7. HTTP Client

**Java (JSoup):**

```java
Document doc = Jsoup.connect(url).get();
JSONObject json = new JSONObject(doc.select("body").text());
```

**Go:**

```go
resp, err := http.Get(url)
body, err := io.ReadAll(resp.Body)
var data CaixaResponse
json.Unmarshal(body, &data)
```

## 🎯 Principais Diferenças

### Tratamento de Erros

**Java** usa exceptions:

```java
try {
    resultado = service.findById(id);
} catch (Exception e) {
    throw new ResourceNotFoundException("Not found");
}
```

**Go** retorna erros explicitamente:

```go
resultado, err := service.FindById(id)
if err != nil {
    return nil, fmt.Errorf("not found: %w", err)
}
```

### Concorrência

**Java:**

```java
@Async
public void updateLoteria(String loteria) {
    // código assíncrono
}
```

**Go (mais simples e eficiente):**

```go
go func() {
    updateLoteria(loteria)
}()
```

### Null Safety

**Java:**

```java
if (objeto != null && objeto.getValor() != null) {
    // uso seguro
}
```

**Go (com ponteiros):**

```go
if objeto != nil && objeto.Valor != nil {
    // uso seguro
}
```

## 📦 Estrutura de Pastas

```
Java/Spring Boot:
src/main/java/com/gutotech/loteriasapi/
├── controller/
├── service/
├── repository/
├── model/
└── config/

Go:
.
├── cmd/server/          # Ponto de entrada
└── internal/            # Código privado
    ├── controller/
    ├── service/
    ├── repository/
    ├── model/
    ├── config/
    └── scheduler/
```

## 🚀 Performance

Go oferece vantagens significativas:

1. **Startup mais rápido**: Segundos vs minutos do Spring Boot
2. **Menor uso de memória**: ~20-50MB vs ~300-500MB do Spring Boot
3. **Concorrência nativa**: Goroutines vs Thread pools
4. **Compilação**: Binário único vs JAR com JVM

## ✅ Checklist de Migração

- [x] Configuração de ambiente (.env)
- [x] Modelos de dados
- [x] Repository layer com MongoDB
- [x] Service layer com lógica de negócio
- [x] Controllers/Handlers HTTP
- [x] Consumer da API da Caixa
- [x] Agendamento de tarefas
- [x] CORS middleware
- [x] Documentação Swagger
- [x] Docker e Docker Compose
- [x] README e documentação
- [ ] Testes unitários (a implementar)
- [ ] Cache Redis (opcional)
- [ ] Métricas e monitoring (opcional)

## 🧪 Testando a Migração

```bash
# 1. Instalar dependências
go mod download

# 2. Iniciar MongoDB
docker run -d -p 27017:27017 mongo:7.0

# 3. Executar aplicação
go run cmd/server/main.go

# 4. Testar endpoints
curl http://localhost:9050/api
curl http://localhost:9050/api/megasena/latest
```

## 📝 Próximos Passos

1. **Implementar testes**: Usar `testing` package do Go
2. **Adicionar métricas**: Prometheus/Grafana
3. **Implementar cache**: Redis para melhor performance
4. **CI/CD**: GitHub Actions ou similar
5. **Logging estruturado**: Usar zap ou logrus
6. **Configuração**: Viper para gerenciamento avançado

## 💡 Dicas

1. **Use go fmt**: Formatação automática
2. **Use golangci-lint**: Linter robusto
3. **Evite panic**: Prefira retornar erros
4. **Use context**: Para timeouts e cancelamento
5. **Estruture bem os pacotes**: Mantenha código privado em `internal/`

## 🔗 Recursos Úteis

- [Effective Go](https://go.dev/doc/effective_go)
- [Gin Documentation](https://gin-gonic.com/docs/)
- [MongoDB Go Driver](https://www.mongodb.com/docs/drivers/go/current/)
- [Go by Example](https://gobyexample.com/)

---

**Resultado**: Projeto totalmente funcional em Go, mantendo compatibilidade com a API original Java/Spring Boot! 🎉
