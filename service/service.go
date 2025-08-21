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

func buildResponse(url string, api model.CepApi, result model.RequestResult, elapsed time.Duration) model.CepResponseChannel {
	return model.CepResponseChannel{
		Body:    result.Body,
		Url:     url,
		Err:     result.Err,
		CepApi:  api,
		Elapsed: elapsed,
	}
}

func executeRequest(ctx context.Context, url string) model.RequestResult {
	req, err := http.NewRequestWithContext(ctx, "GET", url, nil)
	if err != nil {
		return model.RequestResult{Err: err}
	}

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return model.RequestResult{Err: err}
	}
	defer resp.Body.Close()

	// Verifica se houve timeout
	if ctx.Err() != nil {
		return model.RequestResult{Err: ctx.Err()}
	}

	body, err := io.ReadAll(resp.Body)
	return model.RequestResult{Body: body, Err: err}
}

func trySendResponse(ctx context.Context, channel chan<- model.CepResponseChannel, response model.CepResponseChannel) {
	select {
	case channel <- response:
		if response.Err == nil {
			log.Printf("chamada %s finalizada com sucesso em %v!\n", response.Url, response.Elapsed)
		}
	case <-ctx.Done():
		log.Printf("chamada %s foi cancelada durante envio\n", response.Url)
	}
}

func getModel[T model.ViacepResponse | model.BrasilApiResponse](bodyBytes []byte) T {
	var model T
	json.Unmarshal(bodyBytes, &model)
	return model
}

func displayViaCepResult(body []byte) {
	viacep := getModel[model.ViacepResponse](body)
	fields := map[string]string{
		"CEP":        viacep.Cep,
		"Logradouro": viacep.Logradouro,
		"Bairro":     viacep.Bairro,
		"Cidade":     viacep.Localidade,
		"UF":         viacep.Uf,
		"Estado":     viacep.Estado,
	}
	printFields(fields)
}

func displayBrasilApiResult(body []byte) {
	brasilApi := getModel[model.BrasilApiResponse](body)
	fields := map[string]string{
		"CEP":    brasilApi.Cep,
		"Rua":    brasilApi.Street,
		"Bairro": brasilApi.Neighborhood,
		"Cidade": brasilApi.City,
		"Estado": brasilApi.State,
	}
	printFields(fields)
}

func printFields(fields map[string]string) {
	for label, value := range fields {
		if value != "" {
			log.Printf("%s: %s\n", label, value)
		}
	}
}

func GetCep(ctx context.Context, channel chan<- model.CepResponseChannel, api model.CepApi, cep string) {
	start := time.Now()
	url := getApiUrl(api, cep)

	result := executeRequest(ctx, url)
	elapsed := time.Since(start)
	response := buildResponse(url, api, result, elapsed)
	trySendResponse(ctx, channel, response)
}

func DisplayResult(value model.CepResponseChannel) {
	log.Println()
	log.Printf("dados do endereço %v\n", value.CepApi)
	log.Println()

	switch value.CepApi {
	case model.ViaCep:
		displayViaCepResult(value.Body)
	case model.BrasilApi:
		displayBrasilApiResult(value.Body)
	}
}
