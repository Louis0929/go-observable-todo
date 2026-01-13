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

### 3. Run Observability Stack (Jaeger, Prometheus, Grafana)
We use Docker Compose to spin up all observability tools at once.
```bash
docker-compose up -d
```
This will start:
- **Jaeger** (Tracing): http://localhost:16686
- **Prometheus** (Metrics): http://localhost:9090
- **Grafana** (Dashboards): http://localhost:3000 (User: `admin`, Password: `admin`)

## Available API Endpoints

> **Note**: All API endpoints are prefixed with `/api/v1`.

### Create a To-Do (Async)

- **Endpoint**: `POST /api/v1/todos`
- **Description**: Submits a request to create a new to-do item. This is processed asynchronously.
- **Body** (JSON):
  ```json
  {
      "title": "Learn Observability"
  }
  ```
- **Response** (202 Accepted):
  ```json
  {
      "status": "queued",
      "message": "Todo creation is being processed in background"
  }
  ```

### Get All To-Dos

- **Endpoint**: `GET /api/v1/todos`
- **Description**: Retrieves a list of all to-do items.

### Update a To-Do

- **Endpoint**: `PUT /api/v1/todos/:id`
- **Description**: Updates an existing to-do item.
- **Body** (JSON):
  ```json
  {
      "title": "Updated Title",
      "status": "completed"
  }
  ```

### Delete a To-Do

- **Endpoint**: `DELETE /api/v1/todos/:id`
- **Description**: Soft deletes a to-do item.

## Observability Endpoints

### Prometheus Metrics

- **Endpoint**: `GET /metrics`
- **URL**: [http://localhost:8080/metrics](http://localhost:8080/metrics)
- **Description**: Exposes application metrics in a format that a Prometheus server can scrape. This includes default Go process metrics and will include custom application metrics.

### Jaeger Tracing UI

- **URL**: [http://localhost:16686](http://localhost:16686)
- **Description**: After running the Jaeger Docker container, you can access the UI at this address to view traces. Select the `go-observable-todo` service to see the request traces as they flow through the application.
