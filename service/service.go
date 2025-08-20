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

func GetCep[T model.ViacepResponse | model.BrasilApiResponse](ctx context.Context, api model.CepApi, cep string) (T, error) {
	url := getApiUrl(api, cep)

	ctx, cancel := context.WithTimeoutCause(ctx, time.Second, fmt.Errorf("timeout na chamada %s", url))
	defer cancel()

	req, err := http.NewRequestWithContext(ctx, "GET", url, nil)
	if err != nil {
		var model T
		return model, err
	}

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		var model T
		return model, err
	}
	defer resp.Body.Close()

	select {
	case <-time.After(time.Second):
		log.Println("request processada com sucesso")
	case <-ctx.Done():
		log.Println(ctx.Err())
		var model T
		return model, err
	}

	bytes, _ := io.ReadAll(resp.Body)

	model := getModel[T](bytes)

	return model, nil
}

func getApiUrl(api model.CepApi, cep string) string {
	url := fmt.Sprintf("https://brasilapi.com.br/api/cep/v1/%s", cep)
	if api == model.ViaCep {
		url = fmt.Sprintf("http://viacep.com.br/ws/%s/json/", cep)
	}
	return url
}

func getModel[T model.ViacepResponse | model.BrasilApiResponse](bodyBytes []byte) T {
	var model T
	json.Unmarshal(bodyBytes, &model)
	return model
}
