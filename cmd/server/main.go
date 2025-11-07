package main

import (
	"context"
	_ "fmt"
	"log"
	"os"
	"time"

	"loterias-api-golang/internal/config"
	"loterias-api-golang/internal/controller"
	"loterias-api-golang/internal/repository"
	"loterias-api-golang/internal/scheduler"
	"loterias-api-golang/internal/service"

	_ "loterias-api-golang/docs"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

//	@title			Loterias API Golang
//	@version		1.0
//	@description	API para consulta de resultados de loterias da Caixa Econômica Federal
//	@termsOfService	http://swagger.io/terms/

//	@contact.name	Marcos  Oleniuk
//	@contact.url	https://wa.me/554498425745
//	@contact.email	marcos@moleniuk.comm

//	@license.name	MIT
//	@license.url	https://opensource.org/licenses/MIT

// @host		api-loterias.moleniuk.com
// @BasePath	/api
// @schemes	https http
func main() {
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found, using system environment variables")
	}

	mongoClient := connectMongoDB()
	defer func() {
		if err := mongoClient.Disconnect(context.Background()); err != nil {
			log.Fatal(err)
		}
	}()

	db := mongoClient.Database("loterias")
	resultadoRepo := repository.NewResultadoRepository(db)
	consumerService := service.NewConsumer()
	defer consumerService.CloseBrowser() // Garantir que browser seja fechado
	resultadoService := service.NewResultadoService(resultadoRepo)
	loteriasUpdate := service.NewLoteriasUpdate(consumerService, resultadoService)

	schedulerLoteria := scheduler.NewScheduledConsumer(loteriasUpdate)
	schedulerLoteria.Start()
	defer schedulerLoteria.Stop()

	router := setupRouter(resultadoService, loteriasUpdate)

	port := getEnv("PORT", "9050")
	log.Printf("Starting server on port %s", port)
	if err := router.Run(":" + port); err != nil {
		log.Fatal("Failed to start server:", err)
	}
}

func connectMongoDB() *mongo.Client {
	mongoURI := getEnv("MONGODB_URI", "mongodb://localhost:27017/loterias")
	log.Printf("Conectando ao MongoDB: %s", mongoURI)

	// Timeout aumentado para 30 segundos (conexão remota pode ser mais lenta)
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	clientOptions := options.Client().ApplyURI(mongoURI)
	// Configurações adicionais para melhor performance
	clientOptions.SetMaxPoolSize(100)
	clientOptions.SetMinPoolSize(10)
	clientOptions.SetMaxConnIdleTime(5 * time.Minute)

	client, err := mongo.Connect(ctx, clientOptions)
	if err != nil {
		log.Fatalf("❌ Falha ao conectar ao MongoDB: %v", err)
	}

	// Aumentar timeout do Ping também
	pingCtx, pingCancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer pingCancel()

	if err := client.Ping(pingCtx, nil); err != nil {
		log.Fatalf("❌ Falha ao fazer Ping no MongoDB: %v", err)
	}

	log.Println("✅ Conectado ao MongoDB com sucesso!")
	return client
}

func setupRouter(resultadoService *service.ResultadoService, loteriasUpdate *service.LoteriasUpdate) *gin.Engine {
	ginMode := getEnv("GIN_MODE", "debug")
	gin.SetMode(ginMode)

	router := gin.New()

	router.Use(gin.Logger())
	router.Use(gin.Recovery())
	router.Use(config.CORSMiddleware())

	rootController := controller.NewRootController()
	router.GET("/", rootController.Root)

	apiController := controller.NewApiController(resultadoService)
	api := router.Group("/api")
	{
		api.GET("", apiController.GetLotteries)
		api.GET("/:loteria", apiController.GetResultsByLottery)
		api.GET("/:loteria/:concurso", apiController.GetResultByID)
		api.GET("/:loteria/latest", apiController.GetLatestResult)
	}

	// Endpoint administrativo para forçar atualização
	admin := router.Group("/admin")
	{
		admin.POST("/update", func(c *gin.Context) {
			log.Println("Manual update triggered via /admin/update")
			go loteriasUpdate.UpdateAll()
			c.JSON(200, gin.H{
				"message": "Update triggered successfully",
				"status":  "processing",
			})
		})
		admin.POST("/update/:loteria", func(c *gin.Context) {
			loteria := c.Param("loteria")
			log.Printf("Manual update triggered for %s via /admin/update/%s", loteria, loteria)
			go func() {
				err := loteriasUpdate.UpdateOne(loteria)
				if err != nil {
					log.Printf("Error updating %s: %v", loteria, err)
				}
			}()
			c.JSON(200, gin.H{
				"message": "Update triggered for " + loteria,
				"status":  "processing",
			})
		})
		admin.GET("/status", func(c *gin.Context) {
			// Verificar status de bloqueio
			service.BlockMutex.Lock()
			blocked := time.Now().Before(service.BlockedUntil)
			var blockedUntilStr string
			if blocked {
				blockedUntilStr = service.BlockedUntil.Format("2006-01-02 15:04:05")
			}
			service.BlockMutex.Unlock()

			c.JSON(200, gin.H{
				"api_blocked": blocked,
				"blocked_until": blockedUntilStr,
				"current_time": time.Now().Format("2006-01-02 15:04:05"),
			})
		})
		admin.POST("/reset-block", func(c *gin.Context) {
			// Reset bloqueio (cuidado: não fazer sem necessidade)
			service.BlockMutex.Lock()
			service.BlockedUntil = time.Now()
			service.BlockMutex.Unlock()
			
			log.Println("API block status reset")
			c.JSON(200, gin.H{
				"message": "Block status reset",
				"status": "ok",
			})
		})
	}

	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	return router
}

func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}
