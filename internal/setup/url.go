package setup

import (
	"four-download/internal/degenerate"
	"four-download/internal/logging"

	"net/http"
	"strings"
)

func PingURL(inputURL string, d bool) (string, string) {
	if inputURL == "" {
		logging.ErrorLogger.Fatalln("Input an URL to download from.")
	}

	request, err := http.Get(inputURL)
	if err != nil || request.StatusCode != http.StatusOK {
		logging.ErrorLogger.Fatalln("Invalid input URL or IP banned. Request status:", request.Status)
	}

	return splitURL(inputURL, d)
}

func splitURL(inputURL string, d bool) (string, string) {
	var res []string = strings.Split(inputURL, "/")

	if !d {
		degenerate.DegenCheck(res[3])
	}

	return res[3], res[5]
}
