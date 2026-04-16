package main

import (
	"context"
	"log"
	"os"
	"strings"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/olafser/ambulance-management-api-app/api"
	"github.com/olafser/ambulance-management-api-app/internal/config"
	"github.com/olafser/ambulance-management-api-app/internal/handler"
	"github.com/olafser/ambulance-management-api-app/internal/repository"
	"github.com/olafser/ambulance-management-api-app/internal/service"
)

func main() {
	log.Printf("Server started")
	port := os.Getenv("AMBULANCE_MANAGEMENT_API_PORT")
	if port == "" {
		port = "8080"
	}
	environment := os.Getenv("AMBULANCE_MANAGEMENT_API_ENVIRONMENT")
	if !strings.EqualFold(environment, "production") { // case insensitive comparison
		gin.SetMode(gin.DebugMode)
	}

	mongoConfig := config.LoadMongoConfig()
	mongoClient, err := config.NewMongoClient(context.Background(), mongoConfig)
	if err != nil {
		log.Fatalf("MongoDB init failed: %v", err)
	}
	defer func() {
		if disconnectErr := mongoClient.Disconnect(context.Background()); disconnectErr != nil {
			log.Printf("MongoDB disconnect failed: %v", disconnectErr)
		}
	}()

	vehicleRepository, err := repository.NewVehicleRepository(mongoClient.Database(mongoConfig.Database), mongoConfig)
	if err != nil {
		log.Fatalf("Vehicle repository init failed: %v", err)
	}
	vehicleService := service.NewVehicleService(vehicleRepository)
	vehicleHandler := handler.NewVehicleManagementAPI(vehicleService)

	engine := gin.New()
	engine.Use(gin.Recovery())
	engine.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"*"},
		AllowMethods:     []string{"GET", "POST", "PUT", "PATCH", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Authorization", "Content-Type", "Accept"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: false,
	}))
	handler.NewRouterWithGinEngine(engine, handler.ApiHandleFunctions{VehicleManagementAPI: vehicleHandler})
	// request routings
	engine.GET("/openapi", api.HandleOpenApi)
	engine.Run(":" + port)
}
