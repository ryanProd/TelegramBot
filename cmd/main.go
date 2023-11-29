package main

import (
	"log"
	"net/http"

	handler "github.com/ryanProd/TelegramBot"
)

func main() {
	mux := http.NewServeMux()

	mux.Handle("/", http.HandlerFunc(handler.Handler))

	log.Print("Listening...")

	http.ListenAndServe(":3000", mux)
}
