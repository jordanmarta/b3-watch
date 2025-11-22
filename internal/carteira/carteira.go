package carteira

import (
	"encoding/csv"
	"os"
	"strconv"
	"strings"

	"github.com/jordanmarta/b3-watch/internal/models"
)

// normalize deixa comparações mais consistentes
func normalize(s string) string {
	return strings.ToLower(strings.TrimSpace(s))
}

func LoadCarteiraCSV(path string) ([]models.Ticker, error) {

	f, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	r := csv.NewReader(f)
	r.FieldsPerRecord = -1 // aceita número variável de colunas

	rows, err := r.ReadAll()
	if err != nil {
		return nil, err
	}

	if len(rows) < 2 {
		return []models.Ticker{}, nil
	}

	header := rows[0]

	colTipo := -1
	colStatus := -1
	colAplic := -1
	colQtd := -1
	colPreco := -1

	// Detecta índices das colunas automaticamente
	for i, h := range header {
		hNorm := normalize(h)

		switch hNorm {
		case "tipo":
			colTipo = i
		case "status":
			colStatus = i
		case "aplicação", "aplicacao":
			colAplic = i
		case "quantidade":
			colQtd = i
		case "preço", "preco":
			// evita capturar colunas tipo Preço.1
			if !strings.Contains(hNorm, ".") && colPreco == -1 {
				colPreco = i
			}
		}
	}

	if colTipo < 0 || colStatus < 0 || colAplic < 0 || colQtd < 0 || colPreco < 0 {
		return nil, err
	}

	// Estrutura interna para acumular posição
	type acumulado struct {
		totalCompraQtd int
		totalCompraVal float64
		totalVendaQtd  int
	}

	temp := make(map[string]*acumulado)

	for _, row := range rows[1:] {

		// Tipo → só ações e FIIs
		tipo := normalize(row[colTipo])
		if tipo != "ações" && tipo != "fii" {
			continue
		}

		// Aplicação → ticker
		ticker := strings.ToUpper(strings.TrimSpace(row[colAplic]))
		if ticker == "" {
			continue
		}

		// Status → Compra / Venda
		status := normalize(row[colStatus])
		if status != "compra" && status != "venda" {
			continue
		}

		// Quantidade
		qtdStr := strings.TrimSpace(row[colQtd])
		qtd, _ := strconv.Atoi(qtdStr)
		if qtd <= 0 {
			continue
		}

		// Preço
		precoStr := strings.ReplaceAll(row[colPreco], "R$", "")
		precoStr = strings.ReplaceAll(precoStr, ",", ".")
		precoStr = strings.TrimSpace(precoStr)

		preco, _ := strconv.ParseFloat(precoStr, 64)
		if preco <= 0 {
			continue
		}

		// Inicializa caso não exista
		if _, ok := temp[ticker]; !ok {
			temp[ticker] = &acumulado{}
		}

		acc := temp[ticker]

		if status == "compra" {
			acc.totalCompraQtd += qtd
			acc.totalCompraVal += preco * float64(qtd)
		} else if status == "venda" {
			acc.totalVendaQtd += qtd
		}
	}

	// Converter para slice final
	carteira := make([]models.Ticker, 0)

	for ticker, acc := range temp {

		qtdFinal := acc.totalCompraQtd - acc.totalVendaQtd
		if qtdFinal <= 0 {
			continue // ativo zerado → não faz parte da carteira atual
		}

		precoMedio := acc.totalCompraVal / float64(acc.totalCompraQtd)

		carteira = append(carteira, models.Ticker{
			Codigo:     ticker,
			Quantidade: qtdFinal,
			PrecoMedio: precoMedio,
		})
	}

	return carteira, nil
}
