package main

import (
	"four-download/internal/data"
	"four-download/internal/logging"
	"four-download/internal/setup"

	"flag"
	"fmt"
)

var (
	inputURL    string
	workerCount int
	d           bool
	boardName   string
	threadNo    string
)

func main() {
	var cfg setup.Config = setup.ConfigSetup()
	logging.Enable(cfg.LogPath)

	flag.StringVar(&inputURL, "u", "", "URL of the thread to download from.")
	flag.IntVar(&workerCount, "t", 1, "Number of CPU workers to use when concurrently downloading.")
	flag.BoolVar(&d, "d", false, "")
	flag.Parse()

	boardName, threadNo = setup.PingURL(inputURL, d)

	logging.InfoLogger.Printf("Starting on board: '%s', thread: '%s' with '%d' workers.", boardName, threadNo, workerCount)
	fmt.Printf("Starting on board: '%s', thread: '%s' with '%d' workers.\n", boardName, threadNo, workerCount)

	data.GetData(workerCount, boardName, threadNo, inputURL, cfg)
}
