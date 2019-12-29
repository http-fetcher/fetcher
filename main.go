package main

import (
	"fetcher/fetcher"
	"log"
	"net/http"
)

func main() {
	addr := ":8080"
	srv := fetcher.NewServer(1024 * 1204)
	log.Printf("Listening on %s", addr)
	log.Fatal(http.ListenAndServe(addr, srv))
}
