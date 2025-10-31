# Sensor Data Processing System

## Services Overview

- **microservice-a**: Sensor data generator and gRPC client. Sends simulated sensor readings to microservice-b and exposes a REST API for configuration.
- **microservice-b**: Data storage and API service. Receives sensor data via gRPC, stores it in MySQL, and provides REST endpoints for querying, filtering, and deleting sensor data.
- **mysql**: Database service for persistent storage.

## Key Functions

### microservice-a
- `NewSensorClient(serverAddr string)`: Connects to microservice-b via gRPC.
- `SendSensorData(data *domain.SensorData)`: Sends a sensor reading to microservice-b.
- `/config/frequency`: REST endpoint to update sensor generation frequency.

### microservice-b
- `startGRPCServer(cfg, sensorUseCase, logger)`: Starts the gRPC server for receiving sensor data.
- `startHTTPServer(cfg, sensorUseCase, authUseCase, logger)`: Starts the REST API server.
- `/api/sensor-data`: REST endpoints for CRUD operations and filtering sensor data.
- `/health`: Health check endpoint.

## How to Start Services

### 1. Start MySQL (if not using Docker Compose)
```bash
mysql -u root -p -h 127.0.0.1 < schema.sql
```

### 2. Start microservice-b
```powershell
cd microservice-b
update your config file based on parameters
go run ./cmd
```

### 3. Start microservice-a
```powershell
cd microservice-a
update your config file based on parameters
go run ./cmd
```

## How to Test

### Health Check
```bash
curl --noproxy localhost 'http://localhost:8080/health'
```

### Get Sensor Data by ID
```bash
curl --noproxy localhost 'http://localhost:8080/api/sensor-data/123'
```

### Get Sensor Data with Filters
```bash
curl --noproxy localhost 'http://localhost:8080/api/sensor-data?sensor_type=temperature&page=1&page_size=10'
```

### Delete Sensor Data by ID
```bash
curl --noproxy localhost -X DELETE 'http://localhost:8080/api/sensor-data/123'
```

### Update Sensor Generation Frequency (microservice-a)
```bash
curl -x "" -X POST http://localhost:8090/config/frequency -H "Content-Type: application/json" -d '{"interval_ms": 1000}'
```

## Development Notes

> **Note:** Apologies for not having atomic commits; this is being pushed as a single commit for now. Test cases will be added and pushed in future updates.

## Future Improvements
- Add comprehensive test suite
- Add Docker Compose for easy deployment
- Improve logging and error handling
- Add API documentation (Swagger/OpenAPI)
- Add metrics and monitoring

## License
MIT