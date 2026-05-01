package services

import (
	"fmt"
	"github.com/go-resty/resty/v2"
)

type ExchangeService struct {
	client *resty.Client
}

func NewExchangeService() *ExchangeService {
	return &ExchangeService{
		client: resty.New(),
	}
}

// GetUSDExchangeRate — пример функции, которая идет во внешний мир за данными
func (s *ExchangeService) GetUSDExchangeRate() (float64, error) {
	// Используем публичное API для примера (например, из КБ РК или аналоги)
	// Для теста используем эндпоинт, который возвращает простые данные
	resp, err := s.client.R().
		SetResult(map[string]interface{}{}).
		Get("https://api.exchangerate-api.com/v4/latest/KZT")

	if err != nil {
		return 0, err
	}

	if resp.IsError() {
		return 0, fmt.Errorf("api error: %s", resp.Status())
	}

	data := resp.Result().(*map[string]interface{})
	rates := (*data)["rates"].(map[string]interface{})

	// Берем курс доллара к тенге (условно)
	usdRate := rates["USD"].(float64)
	return usdRate, nil
}
