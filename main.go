package main

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/valyala/fasthttp"
)

// fetchKLineData : 呼叫 API 回傳昨日收盤資訊 arrary
func fetchKLineData() ([]interface{}, error) {
	// 取 k 線 api url
	baseURL := "https://max-api.maicoin.com/api/v2/k"
	now := time.Now()
	// 當地時間的午夜
	yesterdayMidnight := time.Date(now.Year(), now.Month(), now.Day()-1, 0, 0, 0, 0, time.Local)
	/*
		由起始時間取usdttwd的收盤價，週期為1440分鐘(1天)，意即昨日最後的收盤價(也等於今日開盤價)
	*/
	market := "usdttwd"                       // 台幣匯率
	limit := 1                                // 資料筆數
	period := 1440                            // 週期分鐘數
	timeStamp := yesterdayMidnight.Unix() - 1 // 起始時間時間戳
	url := fmt.Sprintf("%s?market=%s&limit=%d&period=%d&timestamp=%d", baseURL, market, limit, period, timeStamp)

	client := &fasthttp.Client{}
	statusCode, body, err := client.Get(nil, url)
	if err != nil {
		return nil, fmt.Errorf("取讀資料錯誤: %w", err)
	}

	if statusCode != fasthttp.StatusOK {
		return nil, fmt.Errorf("狀態回應錯誤: %d", statusCode)
	}

	var dataArray [][]interface{}
	if err := json.Unmarshal(body, &dataArray); err != nil {
		return nil, fmt.Errorf("JSON 轉換錯誤: %w", err)
	}

	return dataArray[0], nil
}

func main() {
	// 呼叫 API
	data, err := fetchKLineData()
	if err != nil {
		fmt.Println("Error:", err)
		return
	}

	/*
		回應資料結構 : res -> array of [timestamp, open, high, low, close, volume]
		取得收盤價 res[4]
	*/
	closePrice, ok := data[4].(float64)
	if !ok {
		fmt.Println("Error: 浮點數轉換錯誤")
		return
	}

	fmt.Printf("前日收盤價: %f\n", closePrice)
}
