package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/your-username/go-observable-todo/internal/models"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

// TodoHandler 是一個結構體，用來持有 Handler 所需的依賴 (Dependencies)
// 這就是 "Dependency Injection" 的容器
type TodoHandler struct {
	DB     *gorm.DB
	Logger *zap.Logger
}

// NewTodoHandler 是一個構造函數 (Constructor)
// 用來創建一個 TodoHandler 實例
func NewTodoHandler(db *gorm.DB, logger *zap.Logger) *TodoHandler {
	return &TodoHandler{
		DB:     db,
		Logger: logger,
	}
}

// CreateTodoRequest 定義了創建 Todo 時前端需要傳來的數據格式 (DTO)
// 我們把 DTO 和 Model 分開，這樣更安全
type CreateTodoRequest struct {
	Title string `json:"title" binding:"required,min=3"` // 必須且至少3個字
}

// Create 現在變成了 TodoHandler 的方法 (Method)
// 我們可以透過 h.DB 和 h.Logger 來訪問依賴
func (h *TodoHandler) Create(c *gin.Context) {
	// 1. 綁定數據到 DTO
	var req CreateTodoRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		h.Logger.Warn("Invalid create todo request", zap.Error(err))
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	h.Logger.Info("Creating new todo", zap.String("title", req.Title))

	// 2. 轉換 DTO -> Model
	todo := models.Todo{
		Title:  req.Title,
		Status: "pending", // 預設值
	}

	// 3. 存入數據庫 (使用 Context!)
	// h.DB 就是我們注入進來的 GORM 實例
	if err := h.DB.WithContext(c.Request.Context()).Create(&todo).Error; err != nil {
		h.Logger.Error("Failed to save todo to DB", zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create todo"})
		return
	}

	h.Logger.Info("Todo created successfully", zap.Uint("id", todo.ID))
	c.JSON(http.StatusCreated, todo)
}

// GetList 也是 TodoHandler 的方法
func (h *TodoHandler) GetList(c *gin.Context) {
	var todos []models.Todo

	// 使用 Context 查詢
	if err := h.DB.WithContext(c.Request.Context()).Find(&todos).Error; err != nil {
		h.Logger.Error("Failed to fetch todos", zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch todos"})
		return
	}

	h.Logger.Info("Todos fetched", zap.Int("count", len(todos)))
	c.JSON(http.StatusOK, todos)
}
