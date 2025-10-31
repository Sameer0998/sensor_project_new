package main

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"net"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	_ "github.com/go-sql-driver/mysql"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	grpcPkg "google.golang.org/grpc"

	"sensor_project/microservice-b/internal/config"
	grpcDelivery "sensor_project/microservice-b/internal/delivery/grpc"
	httpDelivery "sensor_project/microservice-b/internal/delivery/http"
	"sensor_project/microservice-b/internal/domain"
	"sensor_project/microservice-b/internal/repository/mysql"
	"sensor_project/microservice-b/internal/usecase"
)

func main() {
	// Load configuration
	cfg := config.LoadConfig()

	// Set up logger
	logger := log.New(os.Stdout, "[MICROSERVICE-B]", log.LstdFlags)
	logger.Println("Starting Microservice B...")

	// Connect to database
	logger.Println("Connecting to database...")
	db, err := connectToDatabase(cfg)
	if err != nil {
		logger.Fatalf("Failed to connect to database: %v", err)
	}
	defer db.Close()

	logger.Println("Database connection established")

	// Initialize repository
	sensorRepo := mysql.NewMySQLSensorRepository(db)

	// Initialize use case
	sensorUseCase := usecase.NewSensorDataUseCase(sensorRepo)

	// Start gRPC server in a goroutine
	go startGRPCServer(cfg, sensorUseCase, logger)

	// Start HTTP server
	startHTTPServer(cfg, sensorUseCase, logger)
}

func connectToDatabase(cfg *config.Config) (*sql.DB, error) {
	// Create a connection pool
	db, err := sql.Open("mysql", cfg.GetDSN())
	if err != nil {
		return nil, err
	}

	// Set connection pool parameters
	db.SetMaxOpenConns(cfg.MaxConnections)
	db.SetMaxIdleConns(cfg.MaxConnections / 2)
	db.SetConnMaxLifetime(time.Hour)

	// Verify connection
	err = db.Ping()
	if err != nil {
		return nil, err
	}

	return db, nil
}

func startGRPCServer(cfg *config.Config, sensorUseCase domain.SensorDataUseCase, logger *log.Logger) {
	lis, err := net.Listen("tcp", fmt.Sprintf(":%s", cfg.GRPCPort))
	if err != nil {
		logger.Fatalf("Failed to listen for gRPC: %v", err)
	}

	grpcServer := grpcPkg.NewServer()
	// Use the delivery layer's gRPC server adapter which implements the generated interface
	sensorServer := grpcDelivery.NewSensorServer(sensorUseCase)
	sensorServer.RegisterServer(grpcServer)

	logger.Printf("gRPC server starting on port %s", cfg.GRPCPort)
	if err := grpcServer.Serve(lis); err != nil {
		logger.Fatalf("Failed to serve gRPC: %v", err)
	}
}

func startHTTPServer(cfg *config.Config, sensorUseCase domain.SensorDataUseCase, logger *log.Logger) {
	// create Echo instance and wire up routes using the HTTP delivery layer
	// Create Echo instance
	e := echo.New()

	// Middleware
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())
	e.Use(middleware.CORS())

	// Initialize HTTP handler and setup routes
	httpHandler := httpDelivery.NewHandler(sensorUseCase)
	httpHandler.SetupRoutes(e)

	// Start server in a goroutine
	go func() {
		address := fmt.Sprintf(":%s", cfg.ServerPort)
		logger.Printf("HTTP server starting on port %s", cfg.ServerPort)
		if err := e.Start(address); err != nil && err != http.ErrServerClosed {
			logger.Fatalf("Failed to start HTTP server: %v", err)
		}
	}()

	// Wait for interrupt signal to gracefully shut down the server
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	// Graceful shutdown
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	if err := e.Shutdown(ctx); err != nil {
		logger.Fatalf("Server shutdown failed: %v", err)
	}
	logger.Println("Server gracefully stopped")
}
