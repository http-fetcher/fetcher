package main

import (
	"fetcher/fetcher"
	"github.com/go-chi/chi"
	"log"
	"net/http"
	"time"
)

func main() {
	addr := ":8080"
	maxBodySize := int64(1024 * 1024)
	client := http.Client{Timeout: 5 * time.Second}

	router := chi.NewRouter()
	crawler := fetcher.NewCrawler(&client)
	srv := fetcher.NewServer(maxBodySize, router, crawler)

	log.Printf("Listening on %s", addr)
	log.Fatal(http.ListenAndServe(addr, srv))
}
