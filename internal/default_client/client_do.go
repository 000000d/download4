package default_client

import "net/http"

func HttpDefaultClientDo(method, url string) (*http.Response, error) {
	req, err := http.NewRequest(method, url, nil)
	if err != nil {
		return nil, err
	}
	var commonUA string = "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/132.0.0.0 Safari/537.36"
	req.Header.Add("User-Agent", commonUA)
	req.Header.Add("Accept", "text/html,application/xhtml+xml,application/xml;q=0.9,image/avif,image/png,image/jpeg,image/svg+xml,*/*;q=0.8")
	req.Header.Add("Accept-Language", "en-US,en;q=0.5,*;q=0.5")
	req.Header.Add("Accept-Encoding", "*")
	req.Header.Add("Dnt", "1")
	req.Header.Add("Sec-Gpc", "1")
	req.Header.Add("Connection", "keep-alive")
	req.Header.Add("Upgrade-Insecure-Requests", "1")
	req.Header.Add("Sec-Fetch-Dest", "document")
	req.Header.Add("Sec-Fetch-Mode", "navigate")
	req.Header.Add("Sec-Fetch-Site", "cross-origin")
	req.Header.Add("Priority", "u=0, i")
	req.Header.Add("Pragma", "no-cache")
	req.Header.Add("Cache-Control", "no-cache")
	req.Header.Add("Te", "trailers")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}

	return resp, nil
}
