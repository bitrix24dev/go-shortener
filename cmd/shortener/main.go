package main

import (
	"fmt"
	"github.com/bitrix24dev/go-shortener/cmd/config"
	"github.com/gorilla/mux"
	"io"
	"log"
	"math/rand"
	"net/http"
	"net/url"
	"sync"
)

type SafeURLMap struct {
	mu sync.Mutex
	v  map[string]string
}

var urlMap = SafeURLMap{v: make(map[string]string)}

const charset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"

func handleURLShortening(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		log.Println("Method not allowed in handleURLShortening method")
		return
	}

	// Чтение данных из тела запроса
	body, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "Unable to read request body", http.StatusInternalServerError)
		return
	}
	defer r.Body.Close()

	shortenedURL := *config.ShortenerBasePath + "/" + generateRandomString(8)

	urlMap.mu.Lock()
	urlMap.v[shortenedURL] = string(body)
	urlMap.mu.Unlock()

	w.Header().Set("Content-Type", "text/plain")
	w.WriteHeader(http.StatusCreated)
	_, err = fmt.Fprint(w, shortenedURL)
	if err != nil {
		log.Println("Error:", err)
		return
	}
}

func handleRedirect(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		log.Println("Method not allowed in handleRedirect method")
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
		log.Println("Error: parsed URL is not found")
		return
	}

	// Получение оригинального URL из карты
	urlMap.mu.Lock()
	originalURL, ok := urlMap.v[*config.ShortenerBasePath+parsedURL.Path]
	if !ok {
		http.NotFound(w, r)
		log.Println("Error: original URL is not found")
		return
	}
	urlMap.mu.Unlock()

	w.Header().Set("Location", originalURL)
	w.WriteHeader(http.StatusTemporaryRedirect)
}

func generateRandomString(length uint) string {

	result := make([]byte, length)
	for i := range result {
		result[i] = charset[rand.Intn(len(charset))]
	}
	return string(result)

}

func main() {

	config.InitConfig()
	r := mux.NewRouter()

	r.HandleFunc("/", handleURLShortening).Methods(http.MethodPost)
	r.HandleFunc("/{variable}", handleRedirect).Methods(http.MethodGet)

	var ServerAddrPathStr string
	if config.ServerAddrPath != nil {
		ServerAddrPathStr = *config.ServerAddrPath
	} else {
		log.Println("Error: ServerAddrPathStr not found")
		return
	}

	fmt.Println("Server is running on " + ServerAddrPathStr)
	err := http.ListenAndServe(ServerAddrPathStr, r)
	if err != nil {
		log.Println("Error:", err)
		return
	}
}
