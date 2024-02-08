package main

import (
	"fmt"
	"github.com/gorilla/mux"
	"io"
	"math/rand"
	"net/http"
	"net/url"
)

var urlMap = make(map[string]string)

const charset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"

func handleURLShortening(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// Чтение данных из тела запроса
	body, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "Unable to read request body", http.StatusInternalServerError)
		return
	}

	shortenedURL := "http://localhost:8080/" + generateRandomString(8)

	urlMap[shortenedURL] = string(body)

	w.Header().Set("Content-Type", "text/plain")
	w.WriteHeader(http.StatusCreated)
	_, err = fmt.Fprint(w, shortenedURL)
	if err != nil {
		return
	}

}

func handleRedirect(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	/*
		vars := mux.Vars(r)
		shortenedURL := vars["variable"]
	*/
	// Парсим URL
	parsedURL, err := url.Parse(r.URL.String())
	if err != nil {
		http.NotFound(w, r)
		return
	}

	// Получение оригинального URL из карты
	originalURL, ok := urlMap["http://localhost:8080"+parsedURL.Path]
	if !ok {
		http.NotFound(w, r)
		return
	}

	w.Header().Set("Location", originalURL)
	w.WriteHeader(http.StatusTemporaryRedirect)
}

func generateRandomString(length int) string {

	result := make([]byte, length)
	for i := range result {
		result[i] = charset[rand.Intn(len(charset))]
	}
	return string(result)

}

func main() {
	r := mux.NewRouter()

	r.HandleFunc("/", handleURLShortening).Methods(http.MethodPost)
	r.HandleFunc("/{variable}", handleRedirect)

	fmt.Println("Server is running on http://localhost:8080")
	err := http.ListenAndServe(":8080", r)
	if err != nil {
		return
	}
}
