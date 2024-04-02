package main

import (
	"context"
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	"time"
)

type Currency struct {
	Bid string `json:"Bid"`
}

func main() {
	ctx, cancel := context.WithTimeout(context.Background(), 300*time.Millisecond)
	defer cancel()

	req, err := http.NewRequest(http.MethodGet, "http://localhost:8080/cotacao", nil)
	if err != nil {
		log.Println("Cannot form request:", err)
	}
	req = req.WithContext(ctx)

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		log.Println("Cannot get response(http.DefaultClient.Do):", err)
		return
	}
	defer resp.Body.Close()

	var curr Currency
	err = json.NewDecoder(resp.Body).Decode(&curr)
	if err != nil {
		log.Println("Cannot decode response:", err)
	}

	err = ioutil.WriteFile("cotacao.txt", []byte("Dólar: "+curr.Bid), 0644)
	if err != nil {
		log.Println("Cannot write to file:", err)
	}
}
