package main

import (
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"log"
	"os"
	"parcel-service/internal/core/services/parcel_service"
	"parcel-service/internal/core/services/rabbitmq_service"
	"parcel-service/internal/handlers"
	"parcel-service/internal/repositories"
	"parcel-service/pkg/rabbitmq"

	"github.com/gin-gonic/gin"
)

const defaultPort = ":1234"
const defaultRmqConn = "amqp://user:password@localhost:5672/"
const defaultDbConn = "postgresql://user:password@localhost:5432/parcel"

func main() {
	dbConn := GetEnvOrDefault("DATABASE", defaultDbConn)

	db, err := gorm.Open(postgres.Open(dbConn))

	if err != nil {
		panic(err)
	}

	parcelRepository, err := repositories.NewCockroachDB(db)

	if err != nil {
		panic(err)
	}

	rmqConn := GetEnvOrDefault("RABBITMQ", defaultRmqConn)

	rmqServer, err := rabbitmq.NewRabbitMQ(rmqConn)

	if err != nil {
		panic(err)
	}

	rmqPublisher := rabbitmq_service.New(rmqServer)

	parcelService := parcel_service.New(parcelRepository, rmqPublisher)

	rmqSubscriber := handlers.NewRabbitMQ(rmqServer, parcelService)

	router := gin.New()

	riderHandler := handlers.NewRest(parcelService, router)
	riderHandler.SetupEndpoints()

	port := GetEnvOrDefault("PORT", defaultPort)

	go rmqSubscriber.Listen("parcelQueue")
	log.Fatal(router.Run(port))
}

func GetEnvOrDefault(environmentKey, defaultValue string) string {
	returnValue := os.Getenv(environmentKey)
	if returnValue == "" {
		returnValue = defaultValue
	}
	return returnValue
}
