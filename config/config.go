package config

import (
	"flag"
	"os"
)

type Configuration struct {
	ServerAddr string
	ResultAddr string
}

var Config Configuration
var serverAddr string
var resultAddr string

func SetConfig() {
	flag.StringVar(&serverAddr, "a", "localhost:8080", "address and port to run server")
	flag.StringVar(&resultAddr, "b", "http://localhost:8080", "base address of the shortened URL")

	flag.Parse()

	if envRunAddr := os.Getenv("SERVER_ADDRESS"); envRunAddr != "" {
		serverAddr = envRunAddr
	}

	if envResAddr := os.Getenv("BASE_URL"); envResAddr != "" {
		resultAddr = envResAddr
	}

	Config = Configuration{ServerAddr: serverAddr, ResultAddr: resultAddr}
}
