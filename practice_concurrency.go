package main

import (
	"fmt"
	"math/rand"
	"sync"
	"time"
)

// ==========================================
// 模擬場景：高性能訂單處理系統
// ==========================================

// Order 代表一張訂單
type Order struct {
	ID     int
	Amount int
}

// 1. [任務] 實現 Producer (生成訂單)
// 參數：
// - count: 要生成的訂單數量
// - out: 這是「唯寫」Channel (chan<-)，只能往裡面塞訂單
func generateOrders(count int, out chan<- Order) {
	fmt.Println("🛒 [系統] 開始接收訂單...")
	
	// [提示] 使用 for 迴圈生成 count 個訂單
	// 每個訂單的 Amount 可以用 rand.Intn(100) + 1
	// 生成後發送到 out channel
	// 最後記得 close(out) 告訴工人沒單了
	
	// 你的代碼...
	for  i:=1 ; i<=count;i++ {
		order := Order {
			ID :i,
			Amount : rand.Intn(100)+1,
		}
		out <-  order
		fmt.Printf("New order # %d coming \n", i, order.Amount)
	}
	close(out)
	fmt.Println("the system is not taking any order")
}

// 2. [任務] 實現 Worker (處理訂單的工人)
// 參數：
// - id: 工人編號
// - in: 這是「唯讀」Channel (<-chan)，只能從裡面拿訂單
// - wg: 用來告訴老闆我做完了
func orderProcessor(id int, in <-chan Order, wg *sync.WaitGroup) {
	defer wg.Done() // 確保工人下班會打卡

	// [提示] 使用 for range 循環不斷從 in channel 拿訂單
	// 模擬處理時間：time.Sleep(time.Millisecond * time.Duration(rand.Intn(100)))
	// 打印日誌：fmt.Printf("👷 工人 %d 處理訂單 #%d (金額: $%d)\n", id, order.ID, order.Amount)
	
	// 你的代碼...
	for order := range in {
		time.Sleep(time.Duration(rand.Intn(500)) * time.Millisecond) 
		fmt.Printf("👷 工人 %d 處理訂單 #%d (金額: $%d)\n", id, order.ID, order.Amount)
	}
}

func main() {
	rand.Seed(time.Now().UnixNano())

	// 設定參數
	totalOrders := 20  // 總共有 20 張訂單
	workerCount := 3   // 只有 3 個工人 (Worker Pool)

	// [任務] 創建一個帶有緩衝區 (Buffer) 的 Channel
	// 這樣即使工人來不及處理，也可以先暫存 10 張訂單
	// orders := ... 
	orders := make(chan Order, 10) // 修改這裡

	// 用來等待所有工人完成
	var wg sync.WaitGroup

	// 1. 啟動訂單生成器 (Producer)
	// 使用 go 關鍵字啟動 generateOrders
	// 你的代碼...
	go generateOrders(totalOrders, orders)

	// 2. 啟動工人池 (Worker Pool)
	// 使用 for 循環啟動 workerCount 個 orderProcessor
	// 記得 wg.Add(1)
	// 你的代碼...
	for	i:=1 ; i<= workerCount; i++{
		wg.Add(1)
		go orderProcessor(i, orders, &wg)
	}

	// 3. 等待所有工人下班
	fmt.Println("⏳ [系統] 等待所有訂單處理完成...")
	wg.Wait()
	fmt.Println("✅ [系統] 所有訂單處理完畢，下班！")
}

