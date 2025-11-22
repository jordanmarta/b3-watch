package monitoramento

import (
	"strings"
	"time"

	"github.com/jordanmarta/b3-watch/internal/api"
	"github.com/jordanmarta/b3-watch/internal/decision"
	"github.com/jordanmarta/b3-watch/internal/models"
)

type Resultado struct {
	Codigo           string
	PrecoAtual       float64
	PrecoMedio       float64
	AbaixoPM         bool
	AbaixoPercentual bool
	PercentualUsado  float64
	ApiOk            bool
}

func ExecutarMonitoramento(
	client *api.IBOVClient,
	carteira []models.Ticker,
	percentual float64,
) []Resultado {

	resultados := make([]Resultado, 0, len(carteira))

	for _, ativo := range carteira {

		// ---- 1) Primeira tentativa ----
		precoAtual, err := client.GetPrecoAtual(ativo.Codigo)

		if err != nil {

			msg := strings.ToLower(err.Error())

			// ticker inexistente → ignora de vez
			if strings.Contains(msg, "dados do ticker") ||
				strings.Contains(msg, "404") ||
				strings.Contains(msg, "not found") {
				continue
			}

			// ---- 2) Retry único ----
			time.Sleep(300 * time.Millisecond)
			precoAtual, err = client.GetPrecoAtual(ativo.Codigo)

			// se falhar de novo → ignora
			if err != nil {
				continue
			}
		}

		// ---- 3) Construir o resultado final ----
		r := Resultado{
			Codigo:           ativo.Codigo,
			PrecoAtual:       precoAtual,
			PrecoMedio:       ativo.PrecoMedio,
			PercentualUsado:  percentual,
			ApiOk:            true,
			AbaixoPM:         decision.PrecoAtualAbaixoDoPrecoMedio(precoAtual, ativo.PrecoMedio),
			AbaixoPercentual: decision.PrecoAtualAbaixoDoPercentual(precoAtual, ativo.PrecoMedio, percentual),
		}

		resultados = append(resultados, r)
	}

	return resultados
}
