package main

// TODO: Import the necessary packages
// You'll need:
// - "log" (standard library for basic logging)
// - "github.com/gin-gonic/gin" (web framework)
// - "gorm.io/driver/sqlite" (SQLite driver)
// - "gorm.io/gorm" (ORM library)
// - Your internal packages: handlers and models
// - "github.com/prometheus/client_golang/prometheus/promhttp" (Prometheus metrics)
// - "go.opentelemetry.io/otel" (OpenTelemetry core)
// - "go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracehttp" (OTLP HTTP exporter)
// - "go.opentelemetry.io/contrib/instrumentation/github.com/gin-gonic/gin/otelgin" (Gin instrumentation)
// - "go.opentelemetry.io/otel/sdk/trace" (Tracing SDK)
import (
	"context"
	"log"

	"github.com/gin-gonic/gin"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/your-username/go-observable-todo/internal/handlers"
	"github.com/your-username/go-observable-todo/internal/models"
	"go.opentelemetry.io/contrib/instrumentation/github.com/gin-gonic/gin/otelgin"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracehttp"
	"go.opentelemetry.io/otel/sdk/resource"
	"go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.27.0"
	"go.uber.org/zap"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

// initTracer initializes OpenTelemetry tracing and sets the service name.
func initTracer() {
	// Create a new OTLP HTTP exporter.
	exporter, err := otlptracehttp.New(context.Background(),
		otlptracehttp.WithEndpoint("localhost:4318"),
		otlptracehttp.WithInsecure(),
	)
	if err != nil {
		log.Fatalf("failed to create OTLP exporter: %v", err)
	}

	// Create a new resource only with the service name.
	// This avoids any schema conflicts from resource.Default().
	res := resource.NewWithAttributes(
		semconv.SchemaURL,
		semconv.ServiceNameKey.String("go-observable-todo"),
	)

	// Create a new TracerProvider with the exporter and our custom resource.
	tp := trace.NewTracerProvider(
		trace.WithBatcher(exporter),
		trace.WithResource(res),
	)

	// Set the global TracerProvider.
	otel.SetTracerProvider(tp)
}

func main() {
	// TODO: Step 0 - Initialize OpenTelemetry Tracer
	// Call initTracer() to set up tracing
	initTracer()

	logger, err := zap.NewProduction()
	if err != nil {
		log.Fatal("Failed to initialize logger:", zap.Error(err))
	}

	defer logger.Sync()
	// TODO: Step 1 - Initialize SQLite database connection
	// Use gorm.Open() with sqlite.Open("todos.db")
	// Store the connection in a variable called 'db'
	// Handle errors - if connection fails, use log.Fatal()

	db, err := gorm.Open(sqlite.Open("todos"), &gorm.Config{})
	if err != nil {
		log.Fatal("failed to connect to database", zap.Error(err))
	}

	// TODO: Step 2 - Auto-migrate the database schema
	// Use db.AutoMigrate() with a pointer to models.Todo struct
	// This will create the 'todos' table if it doesn't exist
	// Handle errors with log.Fatal()
	db.AutoMigrate(&models.Todo{})

	// TODO: Step 3 - Initialize Gin router with OpenTelemetry middleware
	// Use gin.Default() to create a router with default middleware
	// Store it in a variable called 'router'
	// Add OpenTelemetry middleware: router.Use(otelgin.Middleware("go-observable-todo"))
	router := gin.Default()
	router.Use(otelgin.Middleware("go-observable-todo"))

	// TODO: Step 4 - Define your API routes
	// Create a POST /todos route that calls handlers.CreateTodo(db, logger)
	// Create a GET /todos route that calls handlers.GetTodos(db, logger)
	router.POST("/todos", handlers.CreateTodo(db, logger))
	router.GET("/todos", handlers.GetTodos(db, logger))

	// TODO: Step 4.5 - Add Prometheus metrics endpoint
	// Add a GET /metrics route that exposes Prometheus metrics
	// Use: router.GET("/metrics", gin.WrapH(promhttp.Handler()))
	// This will expose Go runtime metrics and custom metrics.
	router.GET("/metrics", gin.WrapH(promhttp.Handler()))

	// TODO: Step 5 - Start the HTTP server
	// Log a message saying the server is starting
	// Use router.Run(":8080") to start on port 8080
	// Handle the error with log.Fatal() if it fails to start
	logger.Info("Starting server on :8080...")
	if err := router.Run(":8080"); err != nil {
		log.Fatal("failed to start server:", err)
	}

}
