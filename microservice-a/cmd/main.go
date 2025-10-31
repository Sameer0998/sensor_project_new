package main

import (
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"sensor_project/microservice-a/internal/config"
	grpcClient "sensor_project/microservice-a/internal/delivery/grpc"
	httpDelivery "sensor_project/microservice-a/internal/delivery/http"
	"sensor_project/microservice-a/internal/domain"
	"sensor_project/microservice-a/internal/usecase"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

func main() {
	// Load configuration
	cfg := config.LoadConfig()

	// Create sensor generator
	generator := usecase.NewSensorGenerator(cfg.SensorConfig)

	// Create gRPC client
	client, err := grpcClient.NewSensorClient(cfg.GRPCServerAddr)
	if err != nil {
		log.Fatalf("Failed to create gRPC client: %v", err)
	}
	defer client.Close()

	// Create HTTP server
	e := echo.New()
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())

	// Set up HTTP routes
	handler := httpDelivery.NewHandler(generator)
	handler.SetupRoutes(e)

	// Start HTTP server in a goroutine
	go func() {
		if err := e.Start(":" + cfg.ServerPort); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Failed to start server: %v", err)
		}
	}()

	// Start data generation in a goroutine
	stopChan := make(chan struct{})
	go generateAndSendData(generator, client, stopChan)

	// Wait for interrupt signal to gracefully shut down the server
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	// Stop data generation
	close(stopChan)

	log.Println("Shutting down server...")
}

// generateAndSendData continuously generates and sends sensor data
func generateAndSendData(generator domain.SensorGenerator, client *grpcClient.SensorClient, stopChan <-chan struct{}) {
	for {
		select {
		case <-stopChan:
			return
		default:
			// Generate sensor data
			data := generator.GenerateSensorData()

			// Send data to Microservice B
			if err := client.SendSensorData(data); err != nil {
				log.Printf("Error sending sensor data: %v", err)
			} else {
				log.Printf("Sent sensor data: %+v", data)
			}

			// Wait for the next generation cycle
			time.Sleep(generator.GetGenerationRate())
		}
	}
}
