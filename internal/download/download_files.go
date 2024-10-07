package download

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"time"

	"four-download/internal/default_client"
	"four-download/internal/logging"
)

var count int = 0

func DownloadFile(fileEndpoint string, filePath string, totalCount int, waitTime float64) {
	out, err := os.Create(filePath)
	if err != nil {
		logging.WarningLogger.Printf("Error with creating file. ERROR: %v\n", err)
	}
	defer out.Close()

	fileRequest, err := default_client.HttpDefaultClientDo(http.MethodGet, fileEndpoint)
	if err != nil || fileRequest.StatusCode != http.StatusOK {
		logging.WarningLogger.Printf("Error connecting to server. Status: %s. ERROR: %v\n", fileRequest.Status, err)
	}
	defer fileRequest.Body.Close()

	_, err = io.Copy(out, fileRequest.Body)
	if err != nil {
		logging.WarningLogger.Printf("Error saving file. Direct link to content: %s. ERROR: %v\n", fileEndpoint, err)
	}

	fmt.Printf("\rDownloaded: %d / %d", count, totalCount)
	count++
	time.Sleep(time.Duration(waitTime) * time.Second)
}
