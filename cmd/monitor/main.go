package main

import (
	"fmt"

	"github.com/jordanmarta/b3-watch/internal/api"
	"github.com/jordanmarta/b3-watch/internal/carteira"
	"github.com/jordanmarta/b3-watch/internal/monitoramento"
)

// func handler(ctx context.Context) error {
// 	log.Println("Monitoramento iniciado...")

// 	client := api.NewIBOVClient("a57aa95b57af20642e9d2f691d4edd7da572479d")

// 	preco, err := client.GetPrecoAtual("MXRF11")
// 	if err != nil {
// 		log.Println("Erro ao buscar preço:", err)
// 	} else {
// 		log.Println("Preço atual MXRF11:", preco)
// 	}

// 	return nil
// }

func main() {

	// --- 1. Carrega CSV ---
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

	// --- 2. Cria client da API ---
	client := api.NewIBOVClient("a57aa95b57af20642e9d2f691d4edd7da572479d")

	// --- 3. Executa monitoramento ---
	resultados := monitoramento.ExecutarMonitoramento(client, cart, 5.0)

	// --- 4. Exibe resultados ---
	fmt.Println("\nResultados do Monitoramento:")
	for _, r := range resultados {
		fmt.Printf(" [%s] Atual: %.2f | PM: %.2f | <PM? %v | <5%%? %v\n",
			r.Codigo,
			r.PrecoAtual,
			r.PrecoMedio,
			r.AbaixoPM,
			r.AbaixoPercentual,
		)
	}
}

// func main() {
// 	carteira, err := carteira.LoadCarteiraCSV("carteira.csv")
// 	if err != nil {
// 		log.Println("Erro ao carregar carteira:", err)
// 	} else {
// 		log.Println("Carteira carregada:")
// 		for _, c := range carteira {
// 			log.Printf(" - %s | Quantidade: %d | Preço Médio: %.2f",
// 				c.Codigo, c.Quantidade, c.PrecoMedio)
// 		}
// 	}

// 	// // Detecta se estamos rodando local ou na Lambda
// 	// if os.Getenv("AWS_LAMBDA_RUNTIME_API") != "" {
// 	// 	lambda.Start(handler)
// 	// } else {
// 	// 	handler(context.Background())
// 	// }
// }
