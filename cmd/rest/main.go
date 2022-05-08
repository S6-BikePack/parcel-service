package main

import (
	"fmt"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"log"
	"os"
	"parcel-service/config"
	"parcel-service/internal/core/services/parcel_service"
	"parcel-service/internal/core/services/rabbitmq_service"
	"parcel-service/internal/handlers"
	"parcel-service/internal/repositories"
	"parcel-service/pkg/rabbitmq"

	"github.com/gin-gonic/gin"
)

const defaultConfig = "./config/local.config"

func main() {
	cfgPath := GetEnvOrDefault("config", defaultConfig)
	cfg, err := config.UseConfig(cfgPath)

	if err != nil {
		panic(err)
	}

	dsn := fmt.Sprintf("host=%s port=%d user=%s "+
		"password=%s dbname=%s sslmode=disable",
		cfg.Database.Host, cfg.Database.Port, cfg.Database.User, cfg.Database.Password, cfg.Database.Database)
	db, err := gorm.Open(postgres.Open(dsn))

	if err != nil {
		panic(err)
	}

	if cfg.Database.Debug {
		db.Debug()
	}

	parcelRepository, err := repositories.NewCockroachDB(db)

	if err != nil {
		panic(err)
	}

	rmqServer, err := rabbitmq.NewRabbitMQ(cfg)

	if err != nil {
		panic(err)
	}

	rmqPublisher := rabbitmq_service.New(rmqServer, cfg)

	parcelService := parcel_service.New(parcelRepository, rmqPublisher)

	rmqSubscriber := handlers.NewRabbitMQ(rmqServer, parcelService, cfg)

	router := gin.New()

	riderHandler := handlers.NewRest(parcelService, router)
	riderHandler.SetupEndpoints()

	go rmqSubscriber.Listen()
	log.Fatal(router.Run(cfg.Server.Port))
}

func GetEnvOrDefault(environmentKey, defaultValue string) string {
	returnValue := os.Getenv(environmentKey)
	if returnValue == "" {
		returnValue = defaultValue
	}
	return returnValue
}
