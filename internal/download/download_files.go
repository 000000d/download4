package download

import (
	"four-download/internal/logging"

	"fmt"
	"io"
	"net/http"
	"os"
	"sync"
)

var (
	count int = 0
)

func DownloadFile(fileEndpoint string, filePath string, totalCount int, wg *sync.WaitGroup, workers chan struct{}) {
	workers <- struct{}{}
	defer func() {
		<-workers
		wg.Done()
		fmt.Printf("\rDownloaded: %d / %d", count, totalCount)
		count++
	}()

	out, err := os.Create(filePath)
	if err != nil {
		logging.WarningLogger.Printf("Error with creating file.")
	}
	defer out.Close()

	fileRequest, err := http.Get(fileEndpoint)
	if err != nil || fileRequest.StatusCode != http.StatusOK {
		logging.WarningLogger.Printf("Error connecting to server. Status: %s", fileRequest.Status)
	}
	defer fileRequest.Body.Close()

	_, err = io.Copy(out, fileRequest.Body)
	if err != nil {
		logging.WarningLogger.Printf("Error saving file. Direct link to content: %s", fileEndpoint)
	}
}
