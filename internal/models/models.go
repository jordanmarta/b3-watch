package models

// Ticker representa um ativo da carteira.
type Ticker struct {
	Codigo     string  `json:"codigo"`
	PrecoMedio float64 `json:"preco_medio"`
	Quantidade int     `json:"quantidade"`
}

// PrecoAtual representa o retorno da API de preço.
type PrecoAtual struct {
	Codigo string  `json:"codigo"`
	Preco  float64 `json:"preco"`
}

// ResultadoMonitoramento é usado para avisos.
type ResultadoMonitoramento struct {
	Codigo          string  `json:"codigo"`
	PrecoAtual      float64 `json:"preco_atual"`
	PrecoReferencia float64 `json:"preco_referencia"`
	Diferenca       float64 `json:"diferenca"`
	DeveAlertar     bool    `json:"deve_alertar"`
}

// SugestaoCompra é o retorno da função sob demanda.
type SugestaoCompra struct {
	Codigo     string  `json:"codigo"`
	Preco      float64 `json:"preco"`
	Quantidade int     `json:"quantidade"`
	Score      float64 `json:"score"`
}
