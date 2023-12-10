package main

import (
	"flag"
)

var flagAddr string

func parseFlags() {
	flag.StringVar(&flagAddr, "a", ":8080", "address and port metrics http server")
	flag.Parse()
}
