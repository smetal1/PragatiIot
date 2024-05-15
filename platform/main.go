package main

import (
	"context"
	"log"
	"net/http"
	"os"

	"PragatiIot/platform/handlers"
	"PragatiIot/platform/mqtt"
	"PragatiIot/platform/rabbitmq"
	"PragatiIot/platform/repositories"
	"PragatiIot/platform/services"
	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/joho/godotenv"
)

type DeviceMessageHandler struct {
	deviceService *services.DeviceService
}

func (h *DeviceMessageHandler) HandleMessage(message []byte) error {
	// Example logic to handle a device message
	log.Printf("Handling message: %s", message)
	return nil
}

// @title Pragati IoT Platform API
// @description This is a generated API documentation for the Pragati IoT Platform.
// @version 1.0
// @host localhost:8080
// @BasePath /
// @schemes http
func main() {
	// Load .env file if it exists
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found")
	}

	dbURL := os.Getenv("DATABASE_URL")
	if dbURL == "" {
		log.Fatal("DATABASE_URL environment variable is not set")
	}
	pool, err := pgxpool.New(context.Background(), dbURL)
	if err != nil {
		log.Fatalf("Unable to connect to the database: %v\n", err)
	}
	defer pool.Close()

	userRepo := repositories.NewUserRepository(pool)
	homeRepo := repositories.NewHomeRepository(pool)
	deviceRepo := repositories.NewDeviceRepository(pool)

	userService := services.NewUserService(userRepo)
	roleRepo := repositories.NewRoleRepository(pool)
	roleService := services.NewRoleService(roleRepo)
	homeService := services.NewHomeService(homeRepo, roleService)
	deviceService := services.NewDeviceService(deviceRepo, homeService)

	userHandler := handlers.NewUserHandler(userService)
	homeHandler := handlers.NewHomeHandler(homeService)
	deviceHandler := handlers.NewDeviceHandler(deviceService, homeService)
	analyticsHandler := handlers.NewAnalyticsHandler(deviceService, homeService)

	rabbitMQURL := os.Getenv("RABBITMQ_URL")
	if rabbitMQURL == "" {
		log.Fatal("RABBITMQ_URL environment variable is not set")
	}
	producer, err := rabbitmq.NewProducer(rabbitMQURL, "device_data")
	if err != nil {
		log.Fatalf("Failed to initialize RabbitMQ producer: %v", err)
	}
	defer producer.Close()

	mqttFactory := mqtt.NewProtocolFactory(deviceService, producer)
	mqttClient := mqtt.NewMQTTClient(deviceService, producer, mqttFactory)

	deviceMessageHandler := &DeviceMessageHandler{deviceService: deviceService}
	consumer, err := rabbitmq.NewConsumer(rabbitMQURL, "device_data", deviceMessageHandler)
	if err != nil {
		log.Fatalf("Failed to initialize RabbitMQ consumer: %v", err)
	}
	consumer.StartConsuming()
	defer consumer.Close()

	router := gin.Default()
	handlers.SetupRoutes(router, userHandler, homeHandler, deviceHandler, analyticsHandler)

	// Adjust certificate paths as required
	//caCert := "platform/mosquitto/certs/ca.crt"
	//clientCert := "platform/mosquitto/certs/server.crt"
	//clientKey := "platform/mosquitto/certs/server.key"
	broker := "ssl://localhost:8883"
	clientID := "platform-client"

	go mqttClient.StartMQTT(broker, clientID, "", "", "", true)

	log.Println("Server started on :8080")
	log.Fatal(http.ListenAndServe(":8080", router))
}
