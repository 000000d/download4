package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"strings"
	"sync"
	"time"
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

var (
	inputURL    string
	workerCount int
	boardName   string
	threadNo    string
	count       int = 0
	totalCount  int
)

func SetupLogging() {
	const logDirName string = "logs"
	var logFileName string = time.Now().Format("2006-01-02_15.04.05")

	err := os.MkdirAll(logDirName, 0755)
	if err != nil {
		log.Fatalln("Error creating log directory:", err)
	}

	logFile, err := os.Create(fmt.Sprintf("./%s/%s.log", logDirName, logFileName))
	if err != nil {
		log.Fatalln("Error creating log file:", err)
	}
	log.SetOutput(logFile)
	log.SetFlags(log.Lmicroseconds | log.Lshortfile)
}

func DegenCheck(boardName string) {
	var boards []string = []string{"s", "hc", "hm", "h", "e", "u", "d", "y", "t", "hr", "gif", "aco", "r"}
	for _, board := range boards {
		if boardName == board {
			fmt.Println("Stop being a degenerate.")
			os.Exit(0)
		}
	}
}

func main() {
	SetupLogging()

	flag.StringVar(&inputURL, "u", "", "URL of the thread to download from.")
	flag.IntVar(&workerCount, "t", 1, "Number of CPU workers to use when concurrently downloading.")
	flag.Parse()

	if inputURL == "" {
		log.Fatalln("Input an URL to download from.")
	}

	request, err := http.Get(inputURL)
	if err != nil || request.StatusCode != http.StatusOK {
		log.Fatalln("Invalid input URL or IP banned. Request status:", request.Status)
	}

	var splitInput []string = strings.Split(inputURL, "/")
	boardName = splitInput[3]
	// DegenCheck(boardName)
	threadNo = splitInput[5]

	log.Printf("Starting on board: '%s', thread: '%s' with '%d' workers.", boardName, threadNo, workerCount)
	fmt.Printf("Starting on board: '%s', thread: '%s' with '%d' workers.\n", boardName, threadNo, workerCount)

	getFiles(workerCount)
}

func getFiles(workerCount int) {
	const downloadsDir string = "downloads"
	var threadEndpoint string = "https://a.4cdn.org/" + boardName + "/thread/" + threadNo + ".json"
	var workers chan struct{} = make(chan struct{}, workerCount)
	var wg sync.WaitGroup

	request, err := http.Get(threadEndpoint)
	if err != nil {
		log.Fatalln("Failed to query URL. Request status:", request.Status)
	}

	requestBody, err := io.ReadAll(request.Body)
	if err != nil {
		log.Fatalln("Failed to read queried body. Error:", err)
	}

	var threadData Thread
	err = json.Unmarshal(requestBody, &threadData)
	if err != nil {
		log.Fatalln("Failed during JSON unmarshal. Error:", err)
	}

	var threadName string = threadData.Posts[0].SemanticURL
	totalCount = threadData.Posts[0].Images

	var storeDownloads string = fmt.Sprintf("./%s/%s/%s", downloadsDir, boardName, threadName)
	_, err = os.Stat(storeDownloads)
	if os.IsNotExist(err) {
		err = os.MkdirAll(storeDownloads, 0755)
		if err != nil {
			log.Fatalln("Error creating downloads directory. Error:", err)
		}
	}
	log.Printf("Downloads path: %s\n", storeDownloads)
	fmt.Printf("Downloads path: %s\n", storeDownloads)

	err = os.WriteFile(fmt.Sprintf("%s/link.txt", storeDownloads), []byte(inputURL), 0644)
	if err != nil {
		log.Println("Error creating text file to write input URL into. Input URL for refrence:", inputURL)
	}

	for i, post := range threadData.Posts {
		if post.Ext == "" {
			log.Printf("Skipping: [%d] is not a file.", i)
			continue
		}

		var fileEndpoint string = fmt.Sprintf("https://i.4cdn.org/%s/%d%s", boardName, post.Tim, post.Ext)
		var filePath string = fmt.Sprintf("%s/%d%s", storeDownloads, post.Tim, post.Ext)

		_, err := os.Stat(filePath)
		if !os.IsNotExist(err) {
			log.Printf("Skipping: [%d] already downloaded.", i)
			continue
		}

		wg.Add(1)
		go Download(fileEndpoint, filePath, &wg, workers)

		log.Printf("[%d] Saved file %s%s", i, post.Filename, post.Ext)
	}
	wg.Wait()
	log.Println("Finished.")
	fmt.Println("\nFinished.")
}

func Download(fileEndpoint string, filePath string, wg *sync.WaitGroup, workers chan struct{}) {
	workers <- struct{}{}
	defer func() {
		<-workers
		wg.Done()
		fmt.Printf("\rDownloaded: %d / %d", count, totalCount)
		count++
	}()

	out, err := os.Create(filePath)
	if err != nil {
		log.Print("Error with creating file.")
	}
	defer out.Close()

	fileRequest, err := http.Get(fileEndpoint)
	if err != nil || fileRequest.StatusCode != http.StatusOK {
		log.Printf("Error connecting to server. Status: %s", fileRequest.Status)
	}
	defer fileRequest.Body.Close()

	_, err = io.Copy(out, fileRequest.Body)
	if err != nil {
		log.Printf("Error saving file. Direct link to content: %s", fileEndpoint)
	}
}
