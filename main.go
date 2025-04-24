package main

import (
	"flag"
	"fmt"

	"four-download/internal/data"
	"four-download/internal/logging"
	"four-download/internal/setup"
)

type FlagValues struct {
	inputURL string
	waitTime float64
}

var flagValues FlagValues

func main() {
	var cfg setup.Config = setup.ConfigSetup()
	logging.Enable(cfg.LogPath)

	flag.StringVar(&flagValues.inputURL, "u", "", "URL of the thread to download from")
	flag.Float64Var(&flagValues.waitTime, "t", 1, "Number of seconds to wait in-between requests")
	flag.Parse()

	boardName, threadNo := setup.PingURL(flagValues.inputURL)

	logging.InfoLogger.Printf("Starting on board: '%s', thread: '%s' with '%f' second interval.\n", boardName, threadNo, flagValues.waitTime)
	fmt.Printf("Starting on board: '%s', thread: '%s' with '%f' second interval.\n", boardName, threadNo, flagValues.waitTime)

	data.GetData(flagValues.waitTime, boardName, threadNo, flagValues.inputURL, cfg)
}
