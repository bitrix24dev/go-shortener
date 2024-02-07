package main

import (
	"fmt"
	"github.com/gorilla/mux"
	_ "github.com/gorilla/mux"
	"net/http"
)

func handleURLShortening(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	shortenedURL := "http://localhost:8080/EwHXdJfB"

	w.Header().Set("Content-Type", "text/plain")
	w.WriteHeader(http.StatusCreated)
	fmt.Fprintf(w, shortenedURL)
}

func handleRedirect(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	originalURL := "https://practicum.yandex.ru/"

	w.Header().Set("Location", originalURL)
	w.WriteHeader(http.StatusTemporaryRedirect)
}

func main() {
	r := mux.NewRouter()

	r.HandleFunc("/", handleURLShortening)
	r.HandleFunc("/{variable}", handleRedirect)

	err := http.ListenAndServe(`:8080`, r)
	if err != nil {
		panic(err)
	}
}
