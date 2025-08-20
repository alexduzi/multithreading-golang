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
	fmt.Println(os.Args)

	cep := os.Args[1]

	fmt.Printf("consultando cep utilizando via cep: %s\n", cep)
	start := time.Now()
	modelViaCep, err := service.GetCep[model.ViacepResponse](context.Background(), model.ViaCep, cep)
	if err != nil {
		panic(err)
	}
	fmt.Printf("finalizado consulta cep utilizando via cep: %v\n", time.Until(start))
	fmt.Printf("retorno: %v\n", modelViaCep)
}
