package main

import (
	"context"
	"log"
	"os"

	"github.com/aws/aws-lambda-go/lambda"
	"github.com/jordanmarta/b3-watch/internal/api"
)

func handler(ctx context.Context) error {
	log.Println("Monitoramento iniciado...")

	client := api.NewIBOVClient("a57aa95b57af20642e9d2f691d4edd7da572479d")

	preco, err := client.GetPrecoAtual("MXRF11")
	if err != nil {
		log.Println("Erro ao buscar preço:", err)
	} else {
		log.Println("Preço atual MXRF11:", preco)
	}

	return nil
}

func main() {
	// Detecta se estamos rodando local ou na Lambda
	if os.Getenv("AWS_LAMBDA_RUNTIME_API") != "" {
		lambda.Start(handler)
	} else {
		handler(context.Background())
	}
}
