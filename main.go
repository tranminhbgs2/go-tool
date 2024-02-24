package main

import (
	"fmt"
	"go-tool/bybitConnect" // Sửa lại đường dẫn phù hợp với package của bạn
	"os"
	"strconv"
	"sync"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	bybit "github.com/wuhewuhe/bybit.go.api"
)

var (
	latestData string
	dataMutex  sync.Mutex
	client     *bybit.Client
)

func fetchData() {
	ticker := time.NewTicker(1 * time.Minute)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			// Thực hiện việc cập nhật dữ liệu ở đây
			marketKline, err := bybitConnect.FetchMarketKline(client, "BTCUSDT", "5", 10)
			if err != nil {
				fmt.Printf("Error retrieving Kline data: %v\n", err)
				continue
			}

			result, err := bybitConnect.AnalyzeKlineData(marketKline)
			if err != nil {
				fmt.Printf("Error analyzing Kline data: %v\n", err)
				continue
			}

			newPrice, oldPrice, percentageChange, err := bybitConnect.CalculatePriceChange(result)
			if err != nil {
				fmt.Printf("Error calculating price change: %v\n", err)
				continue
			}

			dataMutex.Lock()
			latestData = fmt.Sprintf("NEW Price: %f, OLD Price: %f, Percentage Change: %s", newPrice, oldPrice, percentageChange)
			dataMutex.Unlock()

			// In thông báo sau mỗi lần cập nhật dữ liệu thành công
			unixTimestampMs, _ := strconv.ParseInt(result.List[0][0], 10, 64)
			bybitConnect.PrintTimeInUTCPlus7(unixTimestampMs)
			fmt.Printf("Data updated successfully: NEW Price: %f, OLD Price: %f, Percentage Change: %s\n", newPrice, oldPrice, percentageChange)
		}
	}
}

func setupRouter() *gin.Engine {
	r := gin.Default()
	r.GET("/data", func(c *gin.Context) {
		dataMutex.Lock()
		defer dataMutex.Unlock()
		c.JSON(200, gin.H{
			"latestData": latestData,
		})
	})
	return r
}

func main() {
	// Load .env file
	err := godotenv.Load()
	if err != nil {
		fmt.Printf("Error loading .env file")
	}
	// Khởi tạo Bybit client
	apiKey := os.Getenv("BYBIT_API_KEY")
	apiSecret := os.Getenv("BYBIT_API_SECRET")
	fmt.Println("apiSecret:", apiSecret)
	client = bybitConnect.InitClient(apiKey, apiSecret)

	// Bắt đầu goroutine để cập nhật dữ liệu định kỳ
	go fetchData()

	// Cấu hình và khởi động HTTP server với GIN
	r := setupRouter()
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080" // Port mặc định
	}
	r.Run(":" + port) // listen and serve on 0.0.0.0:8080 (default)
}
