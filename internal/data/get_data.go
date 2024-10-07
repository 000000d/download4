package data

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"

	"four-download/internal/default_client"
	"four-download/internal/download"
	"four-download/internal/logging"
	"four-download/internal/setup"
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

func GetData(waitTime float64, boardName, threadNo, inputURL string, cfg setup.Config) {
	var threadEndpoint string = "https://a.4cdn.org/" + boardName + "/thread/" + threadNo + ".json"

	request, err := default_client.HttpDefaultClientDo(http.MethodGet, threadEndpoint)
	if err != nil {
		logging.ErrorLogger.Fatalln("Failed to query URL. Request status:", request.Status)
	}
	defer request.Body.Close()

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
	var totalCount int = threadData.Posts[0].Images

	var storeDownloads string = fmt.Sprintf("%s/%s/%s", cfg.DownloadPath, boardName, threadName)
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
			logging.InfoLogger.Printf("Skipping: [%d] is not a file.\n", i)
			continue
		}

		var fileEndpoint string = fmt.Sprintf("https://i.4cdn.org/%s/%d%s", boardName, post.Tim, post.Ext)
		var filePath string = fmt.Sprintf("%s/%d%s", storeDownloads, post.Tim, post.Ext)

		_, err := os.Stat(filePath)
		if !os.IsNotExist(err) {
			logging.InfoLogger.Printf("Skipping: [%d] already downloaded.\n", i)
			continue
		}

		download.DownloadFile(fileEndpoint, filePath, totalCount, waitTime)

		logging.InfoLogger.Printf("[%d] Saved file %s%s\n", i, post.Filename, post.Ext)
	}
	logging.InfoLogger.Printf("Finished.\n")
	fmt.Println("\nFinished.")
}
