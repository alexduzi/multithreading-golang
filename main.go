package main

import (
	"context"
	"log"
	"os"
	"time"

	"github.com/alexduzi/multithreading/model"
	"github.com/alexduzi/multithreading/service"
)

func main() {
	cep := os.Args[1]

	ch := make(chan model.CepResponseChannel, 2)

	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	go service.GetCep(ctx, ch, model.ViaCep, cep)

	go service.GetCep(ctx, ch, model.BrasilApi, cep)

	select {
	case value := <-ch:
		if value.Err != nil {
			log.Printf("erro na requisição: %v\n", value.Err)
			return
		}

		log.Printf("resposta mais rápida veio da: %s com tempo de %v\n", value.CepApi, value.Elapsed)
		log.Printf("url: %s\n", value.Url)
		service.DisplayResult(value)
	case <-ctx.Done():
		log.Println("timeout: nenhuma API respondeu em 1 segundo")
	}
}
