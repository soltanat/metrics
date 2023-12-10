package main

import (
	"flag"
)

var flagAddr string
var flagReportInterval int
var flagPollInterval int

func parseFlags() {
	flag.StringVar(&flagAddr, "a", "localhost:8080", "address and port metrics http server")
	flag.IntVar(&flagReportInterval, "r", 10, "send metrics report interval")
	flag.IntVar(&flagPollInterval, "p", 2, "poll metrics interval")
	flag.Parse()
}
