package main

import (
	"fmt"
	"os"
	"strconv"

	"github.com/jordanmarta/b3-watch/internal/alertas"
	"github.com/jordanmarta/b3-watch/internal/api"
	"github.com/jordanmarta/b3-watch/internal/carteira"
	"github.com/jordanmarta/b3-watch/internal/monitoramento"
)

func main() {
	token := os.Getenv("IBOV_TOKEN")
	if token == "" {
		fmt.Println("Defina IBOV_TOKEN antes de rodar: export IBOV_TOKEN=...")
		return
	}

	cart, err := carteira.LoadCarteiraCSV("carteira.csv")
	if err != nil {
		fmt.Println("Erro ao carregar CSV:", err)
		return
	}

	fmt.Println("Carteira carregada:")
	for _, t := range cart {
		fmt.Printf(" - %s | Quantidade: %d | PM: %.2f\n",
			t.Codigo, t.Quantidade, t.PrecoMedio)
	}

	percentStr := os.Getenv("ALERT_PERCENT")
	if percentStr == "" {
		percentStr = "5"
	}
	percent, _ := strconv.ParseFloat(percentStr, 64)

	client := api.NewIBOVClient(token)

	resultados := monitoramento.ExecutarMonitoramento(client, cart, percent)

	fmt.Println("\nResultados do Monitoramento:")
	for _, r := range resultados {
		fmt.Printf(" [%s] Atual: %.2f | PM: %.2f | <PM? %v | <%.2f%%? %v\n",
			r.Codigo,
			r.PrecoAtual,
			r.PrecoMedio,
			r.AbaixoPM,
			percent,
			r.AbaixoPercentual,
		)
	}

	alertasGerados := alertas.FiltrarAlertas(resultados)
	fmt.Printf("\nAlertas encontrados: %d\n", len(alertasGerados))

	if len(alertasGerados) == 0 {
		fmt.Println("Nenhuma oportunidade.")
		return
	}

	html := alertas.MontarEmailHTML(alertasGerados)

	fmt.Println("\nPreview do HTML gerado:\n")
	fmt.Println(html)
}
