package config

import (
	"flag"
)

var FlagServerAddr string
var FlagResultAddr string

func ParseFlags() {
	flag.StringVar(&FlagServerAddr, "a", "localhost:8080", "address and port to run server")
	flag.StringVar(&FlagResultAddr, "b", "http://localhost:8080", "base address of the shortened URL")

	flag.Parse()
}
