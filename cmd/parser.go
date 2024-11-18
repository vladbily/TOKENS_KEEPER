package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"

	"gorm.io/gorm"
)

type CryptoResponse struct {
	Data map[string]struct {
		Quote map[string]struct {
			Price float64 `json:"price"`
		} `json:"quote"`
	} `json:"data"`
}

func parser(symbol string, DB_A *gorm.DB) float64 {
	var apiKeyRecord API
	DB_A.Where("Id = ?", 1).First(&apiKeyRecord)
	apiKey := apiKeyRecord.Api
	url := fmt.Sprintf("https://pro-api.coinmarketcap.com/v1/cryptocurrency/quotes/latest?symbol=%s", symbol)

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		fmt.Println("Ошибка при создании запроса:", err)
		return 0
	}

	req.Header.Add("Accepts", "application/json")
	req.Header.Add("X-CMC_PRO_API_KEY", apiKey)

	client := http.Client{}
	res, err := client.Do(req)
	if err != nil {
		fmt.Println("Ошибка при выполнении запроса:", err)
		return 0
	}
	defer res.Body.Close()

	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		fmt.Println("Ошибка при чтении ответа:", err)
		return 0
	}

	var cryptoResponse CryptoResponse
	if err := json.Unmarshal(body, &cryptoResponse); err != nil {
		fmt.Println("Ошибка при разборе JSON:", err)
		return 0
	}

	if cryptoData, exists := cryptoResponse.Data[symbol]; exists {
		price := cryptoData.Quote["USD"].Price
		return price
	} else {
		return 0
	}
}
