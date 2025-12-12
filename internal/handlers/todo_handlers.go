package handlers

import (
	"context"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/your-username/go-observable-todo/internal/models"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/codes"
	"go.opentelemetry.io/otel/trace"
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

	// 3. 模擬非同步處理 (Asynchronous Processing)
	// 場景：我們想快速回應用戶 202 Accepted，然後在背景寫入 DB

	// FIX: 創建一個全新的 Background Context，避免被 HTTP Request Cancel
	// 但是！我們提取原始 Context 中的 Trace Span，作為 Link 加入新 Context
	// 這樣我們既能保持異步執行，又能讓 Jaeger 串連起這兩個 Span
	span := trace.SpanFromContext(c.Request.Context())
	traceContext := trace.ContextWithSpanContext(context.Background(), span.SpanContext())

	go func(ctx context.Context, t models.Todo) {
		// 重新啟動一個新的 Span，並指明它 "Links" 到原本的 Request
		tracer := otel.Tracer("todo-handler")
		ctx, span := tracer.Start(ctx, "background_job", trace.WithLinks(trace.LinkFromContext(c.Request.Context())))
		defer span.End()

		// 增加 Trace Attributes，讓 Jaeger 更好看
		span.SetAttributes(
			attribute.String("todo.title", t.Title),
			attribute.String("job.type", "async_creation"),
		)

		h.Logger.Info("Background job started", zap.String("title", t.Title))

		// 嘗試寫入 DB
		// 現在 ctx 是獨立的 context.Background() + Trace Info，不會被 Cancel
		if err := h.DB.WithContext(ctx).Create(&t).Error; err != nil {
			h.Logger.Error("Background job failed", zap.Error(err))
			span.RecordError(err)
			span.SetStatus(codes.Error, err.Error())
		} else {
			h.Logger.Info("Background job finished", zap.Uint("id", t.ID))
		}
	}(traceContext, todo)

	// 4. 快速回應
	c.JSON(http.StatusAccepted, gin.H{
		"status":  "queued",
		"message": "Todo creation is being processed in background",
	})
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
