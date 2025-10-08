# go-observable-todo

A complete Go backend project from scratch, demonstrating a To-Do List API equipped with a full suite of observability tools. This project serves as a learning sandbox and a portfolio piece for modern backend development practices.

## Tech Stack

- **Language**: Go
- **Web Framework**: Gin
- **ORM**: GORM with SQLite
- **Logging**: Zap (structured logging)
- **Metrics**: Prometheus
- **Tracing**: OpenTelemetry with Jaeger

## How to Run

### 1. Install Dependencies
This command will download all the necessary libraries defined in `go.mod`.
```bash
go mod tidy
```

### 2. Run the Application
This command starts the API server on `http://localhost:8080`.
```bash
go run ./cmd/api/main.go
```

### 3. Run Jaeger for Tracing
To view the distributed traces, you need to run a Jaeger instance. The simplest way is using Docker.
```bash
docker run -d --name jaeger \
  -e COLLECTOR_OTLP_ENABLED=true \
  -p 16686:16686 \
  -p 4318:4318 \
  jaegertracing/all-in-one:latest
```

## Available API Endpoints

### Create a To-Do

- **Endpoint**: `POST /todos`
- **Description**: Creates a new to-do item.
- **Body** (JSON):
  ```json
  {
      "title": "Learn Observability",
      "status": "pending"
  }
  ```

### Get All To-Dos

- **Endpoint**: `GET /todos`
- **Description**: Retrieves a list of all to-do items.

## Observability Endpoints

### Prometheus Metrics

- **Endpoint**: `GET /metrics`
- **URL**: [http://localhost:8080/metrics](http://localhost:8080/metrics)
- **Description**: Exposes application metrics in a format that a Prometheus server can scrape. This includes default Go process metrics and will include custom application metrics.

### Jaeger Tracing UI

- **URL**: [http://localhost:16686](http://localhost:16686)
- **Description**: After running the Jaeger Docker container, you can access the UI at this address to view traces. Select the `go-observable-todo` service to see the request traces as they flow through the application.
