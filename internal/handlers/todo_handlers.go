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

// Job 定義背景任務的結構
type Job struct {
	Ctx  context.Context
	Todo models.Todo
}

// TodoHandler 是一個結構體，用來持有 Handler 所需的依賴 (Dependencies)
// 這就是 "Dependency Injection" 的容器
type TodoHandler struct {
	DB       *gorm.DB
	Logger   *zap.Logger
	JobQueue chan Job // Worker Pool 隊列
}

// NewTodoHandler 是一個構造函數 (Constructor)
// 用來創建一個 TodoHandler 實例
func NewTodoHandler(db *gorm.DB, logger *zap.Logger) *TodoHandler {
	h := &TodoHandler{
		DB:       db,
		Logger:   logger,
		JobQueue: make(chan Job, 100), // Buffer=100 的任務隊列
	}

	// 啟動 3 個 Worker
	for i := 0; i < 3; i++ {
		go h.worker(i)
	}

	return h
}

// worker 是背景工作者，負責消費 JobQueue
func (h *TodoHandler) worker(id int) {
	h.Logger.Info("Worker started", zap.Int("worker_id", id))

	for job := range h.JobQueue {
		h.processJob(id, job)
	}
}

// processJob 處理單個任務邏輯
func (h *TodoHandler) processJob(workerID int, job Job) {
	// 重新啟動一個新的 Span (因為我們在新的 goroutine 中)
	tracer := otel.Tracer("todo-handler")
	// 使用 job.Ctx (它已經包含了 Trace ID)
	ctx, span := tracer.Start(job.Ctx, "background_job")
	defer span.End()

	span.SetAttributes(
		attribute.String("todo.title", job.Todo.Title),
		attribute.String("job.type", "async_creation"),
		attribute.Int("worker.id", workerID),
	)

	h.Logger.Info("Background job processing",
		zap.Int("worker_id", workerID),
		zap.String("title", job.Todo.Title),
	)

	// 嘗試寫入 DB
	if err := h.DB.WithContext(ctx).Create(&job.Todo).Error; err != nil {
		h.Logger.Error("Background job failed", zap.Error(err))
		span.RecordError(err)
		span.SetStatus(codes.Error, err.Error())
	} else {
		h.Logger.Info("Background job finished",
			zap.Int("worker_id", workerID),
			zap.Uint("id", job.Todo.ID),
		)
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

	// 封裝任務
	job := Job{
		Ctx:  traceContext,
		Todo: todo,
	}

	// 4. 嘗試丟進隊列 (Backpressure 保護)
	select {
	case h.JobQueue <- job:
		h.Logger.Info("Job queued successfully", zap.String("title", todo.Title))
		c.JSON(http.StatusAccepted, gin.H{
			"status":  "queued",
			"message": "Todo creation is being processed in background",
		})
	default:
		// 隊列滿了！這就是 Backpressure (背壓) 保護
		h.Logger.Warn("Job queue full, rejecting request", zap.String("title", todo.Title))
		c.JSON(http.StatusServiceUnavailable, gin.H{
			"error": "Server is busy, please try again later",
		})
	}
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
