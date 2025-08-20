package main

import (
	"context"
	"fmt"
	"os"
	"time"

	"github.com/alexduzi/multithreading/model"
	"github.com/alexduzi/multithreading/service"
)

// https://brasilapi.com.br/api/cep/v1/01153000 + cep

// http://viacep.com.br/ws/" + cep + "/json/

func main() {
	cep := os.Args[1]

	ch := make(chan model.CepResponseChannel, 1)

	go func(cep string) {
		service.GetCep(context.Background(), ch, model.ViaCep, cep)
	}(cep)

	go func(cep string) {
		service.GetCep(context.Background(), ch, model.BrasilApi, cep)
	}(cep)

	for {
		time.Sleep(time.Second)
		select {
		case value := <-ch:
			fmt.Printf("%+v", value)
		default:
			println("finalizado")
		}
	}
}
