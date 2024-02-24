package bybitConnect

import (
	"context"
	"encoding/json"
	"fmt"
	"strconv"
	"time"

	bybit "github.com/wuhewuhe/bybit.go.api"
)

type MarketKline struct {
	RetCode    int             `json:"retCode"`
	RetMsg     string          `json:"retMsg"`
	Result     json.RawMessage `json:"result"`
	RetExtInfo struct{}        `json:"retExtInfo"`
	Time       int64           `json:"time"`
}

type MarketResult struct {
	Category string     `json:"category"`
	List     [][]string `json:"list"`
	Symbol   string     `json:"symbol"`
}

func ParseKlineJSON(data []byte) (MarketKline, error) {
	var kline MarketKline
	err := json.Unmarshal(data, &kline)
	return kline, err
}

func ParseResultJSON(resultData []byte) (MarketResult, error) {
	var result MarketResult
	err := json.Unmarshal(resultData, &result)
	return result, err
}

// Khởi tạo client Bybit
func InitClient(apiKey, apiSecret string) *bybit.Client {
	client := bybit.NewBybitHttpClient(apiKey, apiSecret, bybit.WithBaseURL(bybit.MAINNET))
	return client
}

// Lấy dữ liệu Kline từ Bybit
func FetchMarketKline(client *bybit.Client, symbol, interval string, limit int) (*bybit.ServerResponse, error) {
	marketKline, err := client.NewMarketKlineService("kline", "linear", symbol, interval).Limit(limit).Do(context.Background())
	if err != nil {
		fmt.Errorf("Error FetchMarketKline: %+v\n", err)
		return nil, err
	}
	return marketKline, nil
}

// Phân tích dữ liệu Kline
func AnalyzeKlineData(marketKline *bybit.ServerResponse) (MarketResult, error) {
	var result MarketResult
	kline, err := ParseKlineJSON([]byte(bybit.PrettyPrint(marketKline)))
	if err != nil {
		return result, err
	}

	if kline.RetCode != 0 {
		return result, fmt.Errorf("error: %s", kline.RetMsg)
	}

	results, err := ParseResultJSON(kline.Result)
	if err != nil {
		return result, err
	}
	result = results
	return result, nil
}

// Tính toán biến động giá
func CalculatePriceChange(result MarketResult) (float64, float64, string, error) {
	newPrice, err := strconv.ParseFloat(result.List[0][4], 64) // giá gần nhất
	if err != nil {
		return 0, 0, "", err
	}

	oldPrice, err := strconv.ParseFloat(result.List[len(result.List)-1][4], 64) // giá cách đây 5*10 phút
	if err != nil {
		return 0, 0, "", err
	}

	priceChange := newPrice - oldPrice
	percentageChange := (priceChange / oldPrice) * 100

	return newPrice, oldPrice, fmt.Sprintf("%.2f%%", percentageChange), nil
}

// Chuyển đổi và in thời gian theo UTC+7
func PrintTimeInUTCPlus7(unixTimestampMs int64) {
	unixTimestamp := unixTimestampMs / 1000
	t := time.Unix(unixTimestamp, 0)
	location, err := time.LoadLocation("Asia/Bangkok")
	if err != nil {
		fmt.Println("Error loading location:", err)
		return
	}
	tInUTCPlus7 := t.In(location)
	formattedTime := tInUTCPlus7.Format("02/01/2006 15:04:05")
	fmt.Println("Thời gian theo UTC+7:", formattedTime)
}
