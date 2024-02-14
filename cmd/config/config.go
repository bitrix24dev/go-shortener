package config

import (
	"flag"
	"os"
)

var (
	ServerAddrPath    = flag.String("a", "localhost:8080", "HTTP server address")
	ShortenerBasePath = flag.String("b", "http://localhost:8080", "Base URL for shortened URLs")
)

func InitConfig() {
	flag.Parse()

	// Проверяем переменные окружения для адреса сервера
	ServerAddrPathEnv := os.Getenv("SERVER_ADDRESS")
	if ServerAddrPathEnv != "" {
		*ServerAddrPath = ServerAddrPathEnv
	}

	// Проверяем переменные окружения для базового URL
	ShortenerBasePathEnv := os.Getenv("BASE_URL")
	if ShortenerBasePathEnv != "" {
		*ShortenerBasePath = ShortenerBasePathEnv
	}
}
