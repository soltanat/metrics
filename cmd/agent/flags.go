package main

import (
	"flag"
	"time"
)

var flagAddr string
var flagReportInterval time.Duration
var flagPollInterval time.Duration

func parseFlags() {
	flag.StringVar(&flagAddr, "a", "http://localhost:8080", "address and port metrics http server")
	flag.DurationVar(&flagReportInterval, "r", time.Second*10, "send metrics report interval")
	flag.DurationVar(&flagPollInterval, "p", time.Second*2, "poll metrics interval")
	flag.Parse()
}
