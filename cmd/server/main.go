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

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	clientOptions := options.Client().ApplyURI(mongoURI)
	client, err := mongo.Connect(ctx, clientOptions)
	if err != nil {
		log.Fatal("Failed to connect to MongoDB:", err)
	}

	if err := client.Ping(ctx, nil); err != nil {
		log.Fatal("Failed to ping MongoDB:", err)
	}

	log.Println("Connected to MongoDB successfully")
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
			go loteriasUpdate.UpdateOne(loteria)
			c.JSON(200, gin.H{
				"message": "Update triggered for " + loteria,
				"status":  "processing",
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
