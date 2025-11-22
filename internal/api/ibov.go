package api

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"
)

type IBOVClient struct {
	Token      string
	HTTPClient *http.Client
	BaseURL    string
}

func NewIBOVClient(token string) *IBOVClient {
	return &IBOVClient{
		Token:   token,
		BaseURL: "http://www.ibovfinancials.com/api",
		HTTPClient: &http.Client{
			Timeout: 3 * time.Second,
		},
	}
}

type quoteResponse struct {
	Ticker string  `json:"ticker"`
	Price  float64 `json:"price"`
	Time   string  `json:"time"`
}

func (c *IBOVClient) GetPrecoAtual(ticker string) (float64, error) {
	// Monta a URL com query params
	url := fmt.Sprintf("%s/ibov/quotes/?symbol=%s&token=%s", c.BaseURL, ticker, c.Token)

	resp, err := c.HTTPClient.Get(url)
	if err != nil {
		return 0, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return 0, fmt.Errorf("erro da API: status %d", resp.StatusCode)
	}

	var raw map[string]interface{}
	if err := json.NewDecoder(resp.Body).Decode(&raw); err != nil {
		return 0, err
	}

	// O preço está em raw["data"][ticker]["last"]
	data, ok := raw["data"].(map[string]interface{})
	if !ok {
		return 0, fmt.Errorf("campo data não encontrado")
	}

	tickerData, ok := data[ticker].(map[string]interface{})
	if !ok {
		return 0, fmt.Errorf("dados do ticker não encontrados")
	}

	last, ok := tickerData["last"].(float64)
	if !ok {
		return 0, fmt.Errorf("campo last não encontrado")
	}

	return last, nil
}
