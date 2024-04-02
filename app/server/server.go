package main

import (
	"context"
	"database/sql"
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	"time"

	_ "modernc.org/sqlite"
)

type APIResponse struct {
	USD_BRL struct {
		Bid string `json:"bid"`
	} `json:"USDBRL"`
}

func main() {
	http.HandleFunc("/cotacao", handler)
	log.Fatal(http.ListenAndServe(":8080", nil))
}

func handler(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(r.Context(), 200*time.Millisecond)
	defer cancel()

	resultC := make(chan APIResponse)
	go func() {
		resp, err := httpClient().Get("https://economia.awesomeapi.com.br/json/last/USD-BRL")
		if err != nil {
			log.Println("Error while getting API data:", err)
			return
		}
		defer resp.Body.Close()
		body, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			log.Println("Error while reading API data:", err)
			return
		}
		result := APIResponse{}
		err = json.Unmarshal(body, &result)
		if err != nil {
			log.Println("Error while parsing API data:", err)
			return
		}
		resultC <- result
	}()

	select {
	case res := <-resultC:
		saveToDB(res)
		json.NewEncoder(w).Encode(res.USD_BRL)
	case <-ctx.Done():
		log.Println("Request took too long. Aborting...")
	}
}

func httpClient() *http.Client {
	return &http.Client{
		Timeout: time.Second * 2,
	}
}

func saveToDB(res APIResponse) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Millisecond)
	defer cancel()

	db, err := sql.Open("sqlite", "./cotacao.db")
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	statement, err := db.Prepare("CREATE TABLE IF NOT EXISTS cotacao (id INTEGER PRIMARY KEY, bid TEXT)")
	if err != nil {
		log.Printf("Error while creating table: %v", err)
	}
	statement.Exec()

	stmt, err := db.Prepare("INSERT INTO cotacao (bid) VALUES (?)")
	if err != nil {
		log.Printf("Error while preparing DB statement: %v", err)
	}
	_, err = stmt.ExecContext(ctx, res.USD_BRL.Bid)
	if err != nil {
		log.Printf("Error while writing to DB: %v", err)
	} else {
		log.Printf("Sucesso stmt.ExecContext(ctx, res.USD_BRL.Bid")
	}
}
