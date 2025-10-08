package handlers

// TODO: Import the necessary packages:
// - "net/http" (for HTTP status codes)
// - "github.com/gin-gonic/gin"
// - "gorm.io/gorm"
// - "github.com/your-username/go-observable-todo/internal/models"
import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/your-username/go-observable-todo/internal/models"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

// TODO: Modify CreateTodo function signature to accept logger
// Change from: func CreateTodo(db *gorm.DB) gin.HandlerFunc
// Change to:   func CreateTodo(db *gorm.DB, logger *zap.Logger) gin.HandlerFunc
// TODO: Add structured logging in CreateTodo:
// - Log when starting to create todo: logger.Info("Creating new todo", zap.String("title", todo.Title))
// - Log when creation fails: logger.Error("Failed to create todo", zap.Error(err))
// - Log when creation succeeds: logger.Info("Successfully created todo", zap.Int("id", todo.ID))

// This function should:
//  1. Return a gin.HandlerFunc (which is: func(c *gin.Context))
//  2. Inside the returned function:
//     a. Declare a variable to hold the todo (var todo models.Todo)
//     b. Use c.ShouldBindJSON(&todo) to parse the request body
//     - If error: return JSON with http.StatusBadRequest and error message
//     c. Use db.Create(&todo) to save to database
//     - If error: return JSON with http.StatusInternalServerError
//     d. If successful: return JSON with http.StatusCreated and the created todo
//
// Hint: Return JSON using: c.JSON(statusCode, data)
// Hint: For error messages use: gin.H{"error": "message"}
func CreateTodo(db *gorm.DB, Logger *zap.Logger) gin.HandlerFunc {
	return func(c *gin.Context) {
		Logger.Info("Starting create new todo")
		var todo models.Todo
		if err := c.ShouldBindJSON(&todo); err != nil {
			Logger.Error("Failed to parse JSON", zap.Error(err))
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		if err := db.Create(&todo).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create todo"})
		}

		c.JSON(http.StatusCreated, todo)
		Logger.Info("Successfully parsed todo from JSON", zap.String("title", todo.Title))
	}
}

// TODO: Modify GetTodos function signature to accept logger
// Change from: func GetTodos(db *gorm.DB) gin.HandlerFunc
// Change to:   func GetTodos(db *gorm.DB, logger *zap.Logger) gin.HandlerFunc
// TODO: Add structured logging in GetTodos:
// - Log when starting to fetch todos: logger.Info("Fetching all todos")
// - Log when fetch fails: logger.Error("Failed to fetch todos", zap.Error(err))
// - Log when fetch succeeds: logger.Info("Successfully fetched todos", zap.Int("count", len(todos)))
// This function should:
//  1. Return a gin.HandlerFunc
//  2. Inside the returned function:
//     a. Declare a slice to hold todos (var todos []models.Todo)
//     b. Use db.Find(&todos) to fetch all todos
//     - If error: return JSON with http.StatusInternalServerError
//     c. If successful: return JSON with http.StatusOK and the todos slice
func GetTodos(db *gorm.DB, Logger *zap.Logger) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		var todos []models.Todo
		if err := db.Find(&todos).Error; err != nil {
			Logger.Error("Failed to find database", zap.Error(err))
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Cabbit fetch the db"})
			return
		}
		Logger.Info("Successfully fetched todos", zap.Int("count", len(todos)))
		ctx.JSON(http.StatusOK, todos)

	}
}
