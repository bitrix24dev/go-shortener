package config

import "flag"

var (
	ServerAddrPath    = flag.String("a", "localhost:8080", "HTTP server address")
	ShortenerBasePath = flag.String("b", "http://localhost:8080", "Base URL for shortened URLs")
)

func InitConfig() {
	flag.Parse()
}
