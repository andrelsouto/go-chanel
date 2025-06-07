package main

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"os"
	"time"
)

const viaCepApi = "https://viacep.com.br/ws/%s/json/"
const brasilApi = "https://brasilapi.com.br/api/cep/v1/%s"

type Request struct {
	url string
	cep string
	ch  chan<- string
}

func main() {

	cep := os.Args[1:][0]

	chViaCep := make(chan string)
	chBrasilApi := make(chan string)

	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	viaCepRequest := Request{
		url: viaCepApi,
		cep: cep,
		ch:  chViaCep,
	}
	brasilApiRequest := Request{
		url: brasilApi,
		cep: cep,
		ch:  chBrasilApi,
	}

	go getCep(ctx, viaCepRequest)
	go getCep(ctx, brasilApiRequest)

	select {
	case response := <-chViaCep:
		fmt.Println("Response from ViaCEP: ", response)
	case response := <-chBrasilApi:
		fmt.Println("Response from BrasilApi: ", response)
	case <-ctx.Done():
		fmt.Println("Timeout: No response received within 1 second")
	}

}

func getCep(ctx context.Context, request Request) error {
	url := fmt.Sprintf(request.url, request.cep)
	req, err := http.NewRequestWithContext(ctx, "GET", url, nil)
	if err != nil {
		return err
	}
	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return err
	}
	defer res.Body.Close()
	jsonByte, err := io.ReadAll(res.Body)
	if err != nil {
		return err
	}
	request.ch <- string(jsonByte)
	return nil
}
