package main

import (
	"flag"
)

var flagAddr string

func parseFlags() {
	flag.StringVar(&flagAddr, "a", "localhost:8080", "address and port metrics http server")
	flag.Parse()
}
