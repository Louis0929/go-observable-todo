package main

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"math/rand"
	"sync"
	"time"
)

// ==========================================
// 1. Struct & JSON (API 數據定義)
// ==========================================
// [任務] 定義 Transaction 結構體
// [提示] 記得加 JSON Tags，例如 `json:"tx_hash"`
type Transaction struct {
	TxHash      string `json:"tx_hash"`
	FromAddress string `json:"from_address"`
	Amount      int64  `json:"amount"`
	IsPending   bool   `json:"is_pending"`
}

func practiceStructs() {
	fmt.Println("\n--- 練習 1: Structs 與 JSON ---")

	// [任務] 初始化一個 Transaction 變量
	// [提示] 使用 := 來聲明並賦值
	// ??? tx := ...
	tx := Transaction{
		TxHash: "0xabc",
		Amount: 1000,
	}

	// [任務] 將 struct 轉成 JSON
	// [提示] 使用 json.MarshalIndent
	// jsonData, err := ???
	jsonData, err := json.MarshalIndent(tx, "", "")

	if err != nil {
		fmt.Println(err)
		return
	}

	fmt.Printf("Go Struct: %+v\n", tx)
	fmt.Printf("API JSON: %s\n", string(jsonData))
}

// ==========================================
// 2. Error Handling (錯誤處理)
// ==========================================
func broadcast(amount int64) (string, error) {
	// [任務] 檢查金額，如果 <= 0 就報錯
	// [提示] 使用 errors.New("...")
	if amount <= 0 {
		return "", errors.New("Invalid amount") // ???
	}
	return "0xSuccessHash", nil
}

func practiceErrorHandling() {
	fmt.Println("\n--- 練習 2: Error Handling ---")

	txAmount := int64(60000) // 故意設為負數測試錯誤

	// [任務] 調用 broadcast 函數，並接收返回值
	// hash, err := ???
	hash, err := broadcast(txAmount)

	// [任務] 判斷是否有錯誤
	// [提示] Go 的經典寫法 if err ...

	if err != nil {
		fmt.Println("廣播失敗:", err)
		return
	}

	fmt.Println("廣播成功! Hash:", hash)
}

// ==========================================
// 3. Goroutines & Channels (併發 - 重中之重!)
// ==========================================
func practiceConcurrency() {
	fmt.Println("\n--- 練習 3: Goroutines ---")

	nodes := []string{"Node A", "Node B", "Node C"}

	// [任務] 創建一個 string 類型的 Channel
	results := make(chan string)

	var wg sync.WaitGroup

	for _, node := range nodes {
		wg.Add(1)

		// [任務] 啟動一個 Goroutine (匿名函數)
		// [提示] 關鍵字是 go
		go func(n string) {
			defer wg.Done()
			time.Sleep(time.Duration(rand.Intn(500)) * time.Millisecond) // 模擬延遲

			// [任務] 將結果發送到 Channel
			// [提示] 使用 <- 箭頭
			results <- fmt.Sprintf("[%s] Done", n)
		}(node)
	}

	// [任務] 另開一個 Goroutine 等待所有任務完成並關閉 Channel
	go func() {
		wg.Wait()
		// [任務] 關閉 Channel
		close(results)
	}()

	// [任務] 從 Channel 讀取數據直到關閉
	for msg := range results {
		fmt.Println("收到:", msg)
	}

}

// ==========================================
// 4. Context (超時控制)
// ==========================================
func practiceContext() {
	fmt.Println("\n--- 練習 4: Context ---")

	// [任務] 創建一個 1 秒後超時的 Context
	// [提示] context.WithTimeout
	ctx, cancel := context.WithTimeout(context.Background(), 2 *time.Second)
	defer cancel()

	fmt.Println("請求中 (限時 1 秒)...")

	select {
	case <-time.After(5 * time.Millisecond):
		fmt.Println("請求完成 (太慢了)")
		// [任務] 捕捉 Context 超時信號
		// [提示] ctx.Done()
	case <- ctx.Done():
		fmt.Println("請求被取消:", ctx.Err())
	}
}

// ==========================================
// 5. Defer (資源清理)
// ==========================================
func practiceDefer() {
	fmt.Println("\n--- 練習 5: Defer ---")

	fmt.Println("1. 打開 DB 連接")

	// [任務] 確保 "3. 關閉 DB 連接" 是最後才執行的
	// [提示] 關鍵字 defer
	defer fmt.Println("3. 關閉 DB 連接")

	fmt.Println("2. 執行查詢...")
}

// ==========================================
// 主菜單 (不用改)
// ==========================================
func main() {
	for {
		fmt.Println("\n=== Go 練習模式 (填空題) ===")
		fmt.Println("1. Struct & JSON")
		fmt.Println("2. Error Handling")
		fmt.Println("3. Goroutines")
		fmt.Println("4. Context")
		fmt.Println("5. Defer")
		fmt.Println("0. Exit")
		fmt.Print("輸入編號: ")

		var choice int
		fmt.Scanln(&choice)

		switch choice {
		case 1:
			practiceStructs()
		case 2:
			practiceErrorHandling()
		case 3:
			practiceConcurrency()
		case 4:
			practiceContext()
		case 5:
			practiceDefer()
		case 0:
			return
		}
	}
}
