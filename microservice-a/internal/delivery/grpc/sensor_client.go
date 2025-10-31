package grpc

import (
	"context"
	"log"
	"time"

	pb "sensor_project/proto/sensor_project/proto/sensor"

	"sensor_project/microservice-a/internal/domain"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

// SensorClient implements the SensorSender interface using gRPC
type SensorClient struct {
	client pb.SensorServiceClient
	conn   *grpc.ClientConn
}

// NewSensorClient creates a new gRPC client for sending sensor data
func NewSensorClient(serverAddr string) (*SensorClient, error) {
	// Set up a connection to the server
	conn, err := grpc.Dial(serverAddr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return nil, err
	}

	client := pb.NewSensorServiceClient(conn)
	return &SensorClient{
		client: client,
		conn:   conn,
	}, nil
}

// SendSensorData sends a single sensor reading to the server
func (c *SensorClient) SendSensorData(data *domain.SensorData) error {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	// Convert domain model to protobuf message
	pbData := &pb.SensorData{
		SensorValue: float32(data.SensorValue),
		SensorType:  data.SensorType,
		Id1:         data.ID1,
		Id2:         int32(data.ID2),
		Timestamp:   time.Now().UnixMilli(),
	}

	// Send the data
	resp, err := c.client.SendSensorData(ctx, pbData)
	if err != nil {
		log.Printf("Failed to send sensor data: %v", err)
		return err
	}

	log.Println("Response from [MICROSERVICE - B] : ", resp)
	if !resp.Success {
		log.Printf("Server rejected sensor data: %s", resp.Message)
	}

	return nil
}

// Close closes the gRPC connection
func (c *SensorClient) Close() error {
	return c.conn.Close()
}
