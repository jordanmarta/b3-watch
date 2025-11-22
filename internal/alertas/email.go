package alertas

import (
	"context"
	"fmt"
	"strings"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	ses "github.com/aws/aws-sdk-go-v2/service/ses"
	"github.com/aws/aws-sdk-go-v2/service/ses/types"
	"github.com/jordanmarta/b3-watch/internal/monitoramento"
)

// Estrutura que controla os detalhes do envio
type EmailConfig struct {
	From    string
	To      string
	Assunto string
}

// Filtra apenas os tickers onde as regras dispararam
func FiltrarAlertas(resultados []monitoramento.Resultado) []monitoramento.Resultado {
	alertas := []monitoramento.Resultado{}

	for _, r := range resultados {
		if r.AbaixoPM || r.AbaixoPercentual {
			alertas = append(alertas, r)
		}
	}

	return alertas
}

// Constrói o HTML do e-mail
func MontarEmailHTML(alertas []monitoramento.Resultado) string {

	var sb strings.Builder

	sb.WriteString(`
	<html>
	<body style="font-family: Arial, sans-serif; color: #333;">
		<h2 style="margin-bottom: 0;">📊 B3-Watch</h2>
		<p style="margin-top: 4px;">Relatório de Oportunidades</p>

		<p>Olá, Jordan! Os seguintes ativos estão abaixo dos seus critérios:</p>

		<table style="border-collapse: collapse; width: 100%; max-width: 600px;">
			<tr style="background: #f0f0f0;">
				<th style="padding: 8px; border: 1px solid #ddd;">Ticker</th>
				<th style="padding: 8px; border: 1px solid #ddd;">Atual</th>
				<th style="padding: 8px; border: 1px solid #ddd;">PM</th>
				<th style="padding: 8px; border: 1px solid #ddd;">Diferença</th>
				<th style="padding: 8px; border: 1px solid #ddd;">Critério</th>
			</tr>
	`)

	for _, r := range alertas {

		difPercent := ((r.PrecoAtual - r.PrecoMedio) / r.PrecoMedio) * 100

		criterios := []string{}
		if r.AbaixoPM {
			criterios = append(criterios, "Abaixo PM")
		}
		if r.AbaixoPercentual {
			criterios = append(criterios, fmt.Sprintf("Abaixo %.2f%%", r.PercentualUsado))
		}

		sb.WriteString(fmt.Sprintf(`
			<tr>
				<td style="padding: 8px; border: 1px solid #ddd;">%s</td>
				<td style="padding: 8px; border: 1px solid #ddd;">%.2f</td>
				<td style="padding: 8px; border: 1px solid #ddd;">%.2f</td>
				<td style="padding: 8px; border: 1px solid #ddd;">%.2f%%</td>
				<td style="padding: 8px; border: 1px solid #ddd;">%s</td>
			</tr>`,
			r.Codigo,
			r.PrecoAtual,
			r.PrecoMedio,
			difPercent,
			strings.Join(criterios, " / "),
		))
	}

	sb.WriteString(`
		</table>

		<p style="margin-top: 20px; font-size: 12px; color: #777;">
			Enviado automaticamente pelo B3-Watch via AWS Lambda + SES.
		</p>

	</body>
	</html>
	`)

	return sb.String()
}

// Envia o e-mail via Amazon SES
func EnviarEmailSES(cfg EmailConfig, htmlBody string) error {
	// Carrega config padrão
	awsCfg, err := config.LoadDefaultConfig(context.TODO())
	if err != nil {
		return err
	}

	client := ses.NewFromConfig(awsCfg)

	input := &ses.SendEmailInput{
		Destination: &types.Destination{
			ToAddresses: []string{cfg.To},
		},
		Message: &types.Message{
			Subject: &types.Content{
				Data: aws.String(cfg.Assunto),
			},
			Body: &types.Body{
				Html: &types.Content{
					Data: aws.String(htmlBody),
				},
			},
		},
		Source: aws.String(cfg.From),
	}

	_, err = client.SendEmail(context.TODO(), input)
	return err
}
