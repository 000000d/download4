package data

import (
	"four-download/internal/download"
	"four-download/internal/logging"

	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"sync"
)

type Thread struct {
	Posts []struct {
		Filename    string `json:"filename"`
		Ext         string `json:"ext"`
		Tim         int64  `json:"tim"`
		SemanticURL string `json:"semantic_url"`
		Images      int    `json:"images"`
	} `json:"posts"`
}

func GetData(workerCount int, boardName, threadNo, inputURL string) {
	const downloadsDir string = "downloads"
	var threadEndpoint string = "https://a.4cdn.org/" + boardName + "/thread/" + threadNo + ".json"
	var workers chan struct{} = make(chan struct{}, workerCount)
	var wg sync.WaitGroup

	request, err := http.Get(threadEndpoint)
	if err != nil {
		logging.ErrorLogger.Fatalln("Failed to query URL. Request status:", request.Status)
	}

	requestBody, err := io.ReadAll(request.Body)
	if err != nil {
		logging.ErrorLogger.Fatalln("Failed to read queried body. Error:", err)
	}

	var threadData Thread
	err = json.Unmarshal(requestBody, &threadData)
	if err != nil {
		logging.ErrorLogger.Fatalln("Failed during JSON unmarshal. Error:", err)
	}

	var threadName string = threadData.Posts[0].SemanticURL
	var totalCount = threadData.Posts[0].Images

	var storeDownloads string = fmt.Sprintf("./%s/%s/%s", downloadsDir, boardName, threadName)
	_, err = os.Stat(storeDownloads)
	if os.IsNotExist(err) {
		err = os.MkdirAll(storeDownloads, 0755)
		if err != nil {
			logging.ErrorLogger.Fatalln("Error creating downloads directory. Error:", err)
		}
	}
	logging.InfoLogger.Printf("Downloads path: %s\n", storeDownloads)
	fmt.Printf("Downloads path: %s\n", storeDownloads)

	err = os.WriteFile(fmt.Sprintf("%s/link.txt", storeDownloads), []byte(inputURL), 0644)
	if err != nil {
		logging.WarningLogger.Println("Error creating text file to write input URL into. Input URL for refrence:", inputURL)
	}

	for i, post := range threadData.Posts {
		if post.Ext == "" {
			logging.InfoLogger.Printf("Skipping: [%d] is not a file.", i)
			continue
		}

		var fileEndpoint string = fmt.Sprintf("https://i.4cdn.org/%s/%d%s", boardName, post.Tim, post.Ext)
		var filePath string = fmt.Sprintf("%s/%d%s", storeDownloads, post.Tim, post.Ext)

		_, err := os.Stat(filePath)
		if !os.IsNotExist(err) {
			logging.InfoLogger.Printf("Skipping: [%d] already downloaded.", i)
			continue
		}

		wg.Add(1)
		go download.DownloadFile(fileEndpoint, filePath, totalCount, &wg, workers)

		logging.InfoLogger.Printf("[%d] Saved file %s%s", i, post.Filename, post.Ext)
	}
	wg.Wait()
	logging.InfoLogger.Printf("Finished.")
	fmt.Println("\nFinished.")
}
