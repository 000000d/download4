package setup

import (
	"net/http"
	"strings"

	"four-download/internal/default_client"
	"four-download/internal/logging"
)

func PingURL(inputURL string) (string, string) {
	if inputURL == "" {
		logging.ErrorLogger.Fatalf("Input an URL to download from.\n")
	}

	request, err := default_client.HttpDefaultClientDo(http.MethodGet, inputURL)
	if err != nil || request.StatusCode != http.StatusOK {
		logging.ErrorLogger.Fatalf("Invalid input URL or IP banned. Request status: %s. ERROR: %v\n", request.Status, err)
	}
	defer request.Body.Close()

	return splitURL(inputURL)
}

func splitURL(inputURL string) (string, string) {
	var res []string = strings.Split(inputURL, "/")
	return res[3], res[5]
}
