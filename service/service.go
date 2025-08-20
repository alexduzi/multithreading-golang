package service

import (
	"context"
	"encoding/json"
	"errors"
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
	url := getApiUrl(api, cep)

	ctx, cancel := context.WithTimeout(ctx, time.Second)
	defer cancel()

	req, err := http.NewRequestWithContext(ctx, "GET", url, nil)
	if err != nil {
		channel <- model.CepResponseChannel{
			Body: nil,
			Url:  url,
			Err:  err,
		}
		close(channel)
		return
	}

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		channel <- model.CepResponseChannel{
			Body: nil,
			Url:  url,
			Err:  err,
		}
		close(channel)
		return
	}
	defer resp.Body.Close()

	select {
	case <-ctx.Done():
		if errors.Is(ctx.Err(), context.DeadlineExceeded) {
			log.Printf("tempo de chamada excedida %+v\n", ctx.Err())
		}
		channel <- model.CepResponseChannel{
			Body: nil,
			Url:  url,
			Err:  err,
		}
		close(channel)
		return
	default:
		log.Printf("chamada %s finalizada!\n", url)
	}

	bytes, _ := io.ReadAll(resp.Body)

	channel <- model.CepResponseChannel{
		Body: bytes,
		Url:  url,
		Err:  nil,
	}
	close(channel)
}

func GetModel[T model.ViacepResponse | model.BrasilApiResponse](bodyBytes []byte) T {
	var model T
	json.Unmarshal(bodyBytes, &model)
	return model
}
