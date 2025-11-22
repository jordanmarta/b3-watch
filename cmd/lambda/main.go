package main

import (
	"context"
	"fmt"
	"os"
	"strconv"

	"github.com/aws/aws-lambda-go/lambda"

	"github.com/jordanmarta/b3-watch/internal/alertas"
	"github.com/jordanmarta/b3-watch/internal/api"
	"github.com/jordanmarta/b3-watch/internal/carteira"
	"github.com/jordanmarta/b3-watch/internal/monitoramento"
)

func handler(ctx context.Context) error {
	token := os.Getenv("IBOV_TOKEN")
	if token == "" {
		return fmt.Errorf("IBOV_TOKEN não configurado")
	}

	from := os.Getenv("ALERT_FROM")
	to := os.Getenv("ALERT_TO")
	if from == "" || to == "" {
		return fmt.Errorf("ALERT_FROM / ALERT_TO não configurados")
	}

	percentStr := os.Getenv("ALERT_PERCENT")
	if percentStr == "" {
		percentStr = "5"
	}
	percent, _ := strconv.ParseFloat(percentStr, 64)

	cart, err := carteira.LoadCarteiraCSV("carteira.csv")
	if err != nil {
		return fmt.Errorf("erro ao carregar carteira: %w", err)
	}

	client := api.NewIBOVClient(token)

	resultados := monitoramento.ExecutarMonitoramento(client, cart, percent)

	alertasGerados := alertas.FiltrarAlertas(resultados)
	if len(alertasGerados) == 0 {
		return nil
	}

	html := alertas.MontarEmailHTML(alertasGerados)

	cfg := alertas.EmailConfig{
		From:    from,
		To:      to,
		Assunto: "B3-Watch – Oportunidades Encontradas",
	}

	return alertas.EnviarEmailSES(cfg, html)
}

func main() {
	lambda.Start(handler)
}
