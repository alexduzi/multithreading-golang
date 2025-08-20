package service

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"time"

	"github.com/alexduzi/multithreading/model"
)

func getApiUrl(api model.CepApi, cep string) string {
	url := fmt.Sprintf("https://brasilapi.com.br/api/cep/v1/%s", cep)
	if api == model.ViaCep {
		url = fmt.Sprintf("http://viacep.com.br/ws/%s/json/", cep)
	}
	return url
}

func GetCep(ctx context.Context, channel chan<- model.CepResponseChannel, api model.CepApi, cep string) {
	start := time.Now()
	url := getApiUrl(api, cep)

	req, err := http.NewRequestWithContext(ctx, "GET", url, nil)
	if err != nil {
		select {
		case channel <- model.CepResponseChannel{
			Body:    nil,
			Url:     url,
			Err:     err,
			CepApi:  api,
			Elapsed: time.Until(start),
		}:
		case <-ctx.Done():
		}
		return
	}

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		select {
		case channel <- model.CepResponseChannel{
			Body:    nil,
			Url:     url,
			Err:     err,
			CepApi:  api,
			Elapsed: time.Until(start),
		}:
		case <-ctx.Done():
		}
		return
	}
	defer resp.Body.Close()

	select {
	case <-ctx.Done():
		log.Printf("Timeout na chamada %s: %v\n", url, ctx.Err())
		select {
		case channel <- model.CepResponseChannel{
			Body:    nil,
			Url:     url,
			Err:     ctx.Err(),
			CepApi:  api,
			Elapsed: time.Until(start),
		}:
		default:
		}
		return
	default:
	}

	bytes, err := io.ReadAll(resp.Body)
	if err != nil {
		select {
		case channel <- model.CepResponseChannel{
			Body:    nil,
			Url:     url,
			Err:     err,
			CepApi:  api,
			Elapsed: time.Until(start),
		}:
		case <-ctx.Done():
		}
		return
	}

	log.Printf("Chamada %s finalizada com sucesso!\n", url)

	select {
	case channel <- model.CepResponseChannel{
		Body:    bytes,
		Url:     url,
		Err:     nil,
		CepApi:  api,
		Elapsed: time.Until(start),
	}:
	case <-ctx.Done():
		log.Printf("chamada %s foi cancelada ou teve timeout!\n", url)
	}
}

func GetModel[T model.ViacepResponse | model.BrasilApiResponse](bodyBytes []byte) T {
	var model T
	json.Unmarshal(bodyBytes, &model)
	return model
}

func DisplayResult(value model.CepResponseChannel) {
	log.Println()
	log.Printf("Dados do endereço %v\n", value.CepApi)
	log.Println()
	switch value.CepApi {
	case "VIACEP":
		viacep := GetModel[model.ViacepResponse](value.Body)
		log.Printf("CEP: %s\n", viacep.Cep)
		log.Printf("Logradouro: %s\n", viacep.Logradouro)
		log.Printf("Bairro: %s\n", viacep.Bairro)
		log.Printf("Cidade: %s\n", viacep.Localidade)
		log.Printf("UF: %s\n", viacep.Uf)
		log.Printf("Estado: %s\n", viacep.Estado)

	case "BRAZILAPI":
		brasilApi := GetModel[model.BrasilApiResponse](value.Body)
		log.Printf("CEP: %s\n", brasilApi.Cep)
		log.Printf("Rua: %s\n", brasilApi.Street)
		log.Printf("Bairro: %s\n", brasilApi.Neighborhood)
		log.Printf("Cidade: %s\n", brasilApi.City)
		log.Printf("Estado: %s\n", brasilApi.State)
	}
}
