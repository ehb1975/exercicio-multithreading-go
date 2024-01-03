package main

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"
)

type BrasilAPI struct {
	Cep          string `json:"cep"`
	State        string `json:"state"`
	City         string `json:"city"`
	Neighborhood string `json:"neighborhood"`
	Street       string `json:"street"`
	Service      string `json:"service"`
}

type ViaCep struct {
	Cep         string `json:"cep"`
	Logradouro  string `json:"logradouro"`
	Complemento string `json:"complemento"`
	Bairro      string `json:"bairro"`
	Localidade  string `json:"localidade"`
	Uf          string `json:"uf"`
	Ibge        string `json:"ibge"`
	Gia         string `json:"gia"`
	Ddd         string `string:"ddd"`
	Siafi       string `string:"siafi"`
}

type Resultado struct {
	endereco string
	api      string
}

const (
	brasilapi = "https://brasilapi.com.br/api/cep/v1/"
	viacep    = "http://viacep.com.br/ws/"
)

func main() {
	c1 := make(chan Resultado)
	c2 := make(chan Resultado)
	cep := "01153000"

	go func() {
		resp1, error := BuscaBrasilApi(cep)
		if error != nil {
			panic(error)
		}
		endereco := montaEnderecoBrasilAPI(resp1)
		resultado := Resultado{endereco: endereco, api: "BrasilAPI"}
		c1 <- resultado
	}()

	go func() {

		resp2, error := BuscaViaCep(cep)
		if error != nil {
			panic(error)
		}
		endereco := montaEnderecoViaCep(resp2)
		resultado := Resultado{endereco: endereco, api: "BrasilAPI"}
		c2 <- resultado
	}()

	select {
	case msg := <-c1:
		fmt.Printf("Endereco: %s - Fonte: %s", msg.endereco, msg.api)
	case msg := <-c2:
		fmt.Printf("Endereco: %s - Fonte: %s", msg.endereco, msg.api)
	case <-time.After(time.Second):
		println("timeout")
	}

}

func montaEnderecoBrasilAPI(resp *BrasilAPI) string {
	return resp.Street + ", " + resp.Neighborhood + ", " + resp.City + " - " + resp.State + ", " + resp.Cep
}

func montaEnderecoViaCep(resp *ViaCep) string {
	return resp.Logradouro + ", " + resp.Bairro + ", " + resp.Localidade + " - " + resp.Uf + ", " + resp.Cep
}

func BuscaBrasilApi(cep string) (*BrasilAPI, error) {
	uri := brasilapi + cep
	resp, error := http.Get(uri)
	if error != nil {
		return nil, http.ErrAbortHandler
	}
	defer resp.Body.Close()
	body, error := io.ReadAll(resp.Body)
	if error != nil {
		return nil, error
	}
	var c BrasilAPI
	error = json.Unmarshal(body, &c)
	if error != nil {
		return nil, error
	}
	return &c, nil
}

func BuscaViaCep(cep string) (*ViaCep, error) {
	uri := viacep + cep + "/json/"
	resp, error := http.Get(uri)
	if error != nil {
		return nil, http.ErrAbortHandler
	}
	defer resp.Body.Close()
	body, error := io.ReadAll(resp.Body)
	if error != nil {
		return nil, error
	}
	var c ViaCep
	error = json.Unmarshal(body, &c)
	if error != nil {
		return nil, error
	}
	return &c, nil
}
