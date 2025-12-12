# Go Context 深度解析與實戰筆記

本筆記總結了關於 Go `context` 的核心概念、與 `otelgin` Trace ID 的整合，以及在異步任務（Goroutine）中的正確使用模式。

## 1. Context 與 Channel 的本質區別

Goroutine 之間溝通的兩大支柱：

| 特性 | Channel (`chan`) | Context (`ctx`) |
| :--- | :--- | :--- |
| **主要職責** | **傳遞數據** (Data Flow) | **控制生命週期** (Life Cycle) |
| **比喻** | 工廠的輸送帶 (傳送零件) | 工廠的廣播系統 (喊停工、趕進度) |
| **核心功能** | 發送/接收數據、同步執行 | 取消信號 (Cancel)、超時控制 (Timeout)、攜帶範疇數據 (Values) |
| **方向性** | 點對點 (Point-to-Point) | 樹狀廣播 (One-to-All) |

**協作模式**：通常使用 `Channel` 傳遞任務，配合 `select { case <-ctx.Done(): ... }` 來監聽是否該停止任務。

---

## 2. Context 的繼承結構 (Immutability)

Context 是**不可變的 (Immutable)**。我們不能「修改」一個 Context，只能基於父 Context 「衍生 (Derive)」出一個新的子 Context。

*   **Root Context**: `context.Background()`
    *   全域的源頭，是一張「白紙」。
    *   永遠不會超時，無法被取消，沒有值。
    *   用於 `main` 函數、初始化、或作為新請求的根基。

*   **衍生流程**:
    ```text
    [ context.Background() ] (Root: 永生)
               |
               v
    [ WithValue(Root, "ID", 1) ] (Child 1: 帶有 ID)
               |
               v
    [ WithTimeout(Child 1, 2s) ] (Child 2: 帶有 ID + 限時 2秒)
    ```
    *   當 Child 2 超時死亡，Child 1 和 Root **不受影響**。
    *   每個 Job/Request 應該擁有自己獨立衍生的 Context，互不干擾。

---

## 3. HTTP Request Context 與 Trace ID

在 Web Server (如 Gin) 中，Context 的行為尤為重要。

### 3.1 Trace ID 的來源
*   **Trace ID 不是 Context 原生的**，但在 Observable 系統中，它被儲存在 Context 內。
*   使用 `otelgin` 中間件時：
    1.  Middleware 自動生成 Trace ID (或沿用上游的)。
    2.  將 Trace ID 封裝進 Span。
    3.  將 Span **注入 (Inject)** 到 `c.Request.Context()` 中。

### 3.2 HTTP Cancel 機制
*   `c.Request.Context()` 的生命週期 **綁定 TCP 連線**。
*   **觸發 Cancel 的時機**：
    1.  Handler 執行完畢，Response 回傳後 (正常結束)。
    2.  **User Client 取消** (關閉瀏覽器、按停止、斷網)。
*   **影響**：一旦 Cancel，所有綁定此 Context 的操作 (如 DB 查詢) 都會立即收到 `ctx.Done()` 信號並終止 (回傳 `context canceled` 錯誤)。

---

## 4. 異步任務 (Goroutine) 的 Context 陷阱與解法

在 HTTP Handler 中啟動背景任務 (`go func`) 時，**絕對不能直接使用 `c.Request.Context()`**。

### ❌ 錯誤寫法
```go
go func() {
    // 危險！一旦 HTTP 請求結束 (200 OK)，這個 Context 就會被 Cancel。
    // 導致背景的 DB 寫入失敗。
    db.WithContext(c.Request.Context()).Create(...) 
}()
```

### ✅ 正確寫法：換頭 (Context Detachment) + 繼承 Trace
我們需要一個「不會死」的 Context，但又要「繼承」Trace ID 以便追蹤。

```go
// 1. 取出 Span (包含 Trace ID)
span := trace.SpanFromContext(c.Request.Context())

// 2. 嫁接：用 Background (不死之身) 作為新基底，但注入原本的 Span
// 這樣生成的 traceContext 既獨立於 HTTP 連線，又共享同一個 Trace ID
traceContext := trace.ContextWithSpanContext(context.Background(), span.SpanContext())

go func(ctx context.Context) {
    // 3. 安全執行
    // 這裡的 ctx 是全新的，不會因為 User 關閉瀏覽器而被 Cancel
    db.WithContext(ctx).Create(...) 
}(traceContext)
```

### 總結
*   **同步操作** (直接回應 User 的)：直接用 `c.Request.Context()`，讓 User 能取消。
*   **異步操作** (背景任務)：必須切斷與 Request Context 的聯繫 (使用 `Background()`)，但透過 `trace.ContextWithSpanContext` 保持可觀測性 (Observability) 的連結。

