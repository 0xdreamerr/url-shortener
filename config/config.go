package config

import (
	"flag"
)

type Configuration struct {
	ServerAddr string
	ResultAddr string
}

var Config Configuration

func SetConfig() {
	serverAddr := flag.String("a", "localhost:8080", "address and port to run server")
	resultAddr := flag.String("b", "http://localhost:8080", "base address of the shortened URL")

	flag.Parse()

	Config = Configuration{ServerAddr: *serverAddr, ResultAddr: *resultAddr}
}
