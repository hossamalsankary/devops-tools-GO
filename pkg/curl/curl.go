package curl

import (
	"encoding/json"
	"net/http"
	"net/url"
	"time"
)

type Request struct {
	URL     string  `json:"url"`
	Proxy   string  `json:"proxy,omitempty"`
	Timeout float64 `json:"timeout,omitempty"`
}

type Response struct {
	Reachable  bool   `json:"reachable"`
	StatusCode int    `json:"statusCode,omitempty"`
	Error      string `json:"error,omitempty"`
}

func HandleCurl(w http.ResponseWriter, r *http.Request) {
	var req Request
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		sendResponse(w, Response{Error: "Invalid request body"})
		return
	}

	client := &http.Client{Timeout: 30 * time.Second}

	if req.Timeout > 0 {
		client.Timeout = time.Duration(req.Timeout * float64(time.Second))
	}

	if req.Proxy != "" {
		proxyURL, err := url.Parse(req.Proxy)
		if err != nil {
			sendResponse(w, Response{Error: "Invalid proxy URL"})
			return
		}
		client.Transport = &http.Transport{Proxy: http.ProxyURL(proxyURL)}
	}

	resp, err := client.Get(req.URL)
	if err != nil {
		sendResponse(w, Response{Error: err.Error()})
		return
	}
	defer resp.Body.Close()

	sendResponse(w, Response{Reachable: true, StatusCode: resp.StatusCode})
}

func sendResponse(w http.ResponseWriter, resp Response) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(resp)
}