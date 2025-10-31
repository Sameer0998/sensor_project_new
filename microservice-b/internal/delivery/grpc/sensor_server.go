package grpc

import (
	"context"
	"log"
	"time"

	"sensor_project/microservice-b/internal/domain"
	pb "sensor_project/proto/sensor_project/proto/sensor"

	"google.golang.org/grpc"
)

// SensorServer implements the SensorService gRPC server
type SensorServer struct {
	pb.UnimplementedSensorServiceServer
	sensorUseCase domain.SensorDataUseCase
}

// NewSensorServer creates a new gRPC sensor server
func NewSensorServer(sensorUseCase domain.SensorDataUseCase) *SensorServer {
	return &SensorServer{
		sensorUseCase: sensorUseCase,
	}
}

// RegisterServer registers the gRPC server to the provided gRPC server instance
func (s *SensorServer) RegisterServer(grpcServer *grpc.Server) {
	pb.RegisterSensorServiceServer(grpcServer, s)
}

// SendSensorData handles incoming sensor data from Microservice A
func (s *SensorServer) SendSensorData(ctx context.Context, req *pb.SensorData) (*pb.SensorResponse, error) {
	log.Printf("Received sensor data: %+v", req)

	// Convert protobuf message to domain model
	sensorData := &domain.SensorData{
		SensorValue: float64(req.SensorValue),
		SensorType:  req.SensorType,
		ID1:         req.Id1,
		ID2:         int(req.Id2),
		CreatedAt:   time.Now().UnixMilli(),
	}

	// Store the sensor data
	err := s.sensorUseCase.Store(sensorData)
	if err != nil {
		log.Printf("Error storing sensor data: %v", err)
		return &pb.SensorResponse{
			Success: false,
			Message: "Failed to store sensor data",
		}, err
	}

	return &pb.SensorResponse{
		Success: true,
		Message: "Sensor data stored successfully",
	}, nil
}

// StreamSensorData handles streaming sensor data from Microservice A
func (s *SensorServer) StreamSensorData(stream pb.SensorService_StreamSensorDataServer) error {
	for {
		req, err := stream.Recv()
		if err != nil {
			return err
		}

		// Convert protobuf message to domain model
		sensorData := &domain.SensorData{
			SensorValue: float64(req.SensorValue),
			SensorType:  req.SensorType,
			ID1:         req.Id1,
			ID2:         int(req.Id2),
			CreatedAt:   time.Now().UnixMilli(),
		}

		// Store the sensor data
		err = s.sensorUseCase.Store(sensorData)
		if err != nil {
			log.Printf("Error storing sensor data: %v", err)
			continue
		}
	}
}
