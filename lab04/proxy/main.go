package main

import (
	"bytes"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
	"os"
	"strings"
)

const (
	PORT     = ":8080"
	LOG_PATH = "proxy.log"
)

var cache = InitCache()

func main() {
	http.HandleFunc("/", handleRequest)
	log.Printf("Starting proxy server on %s...", PORT)
	if err := http.ListenAndServe(PORT, nil); err != nil {
		log.Fatalf("Failed to start proxy server: %v", err)
	}
}

func handleRequest(w http.ResponseWriter, r *http.Request) {
	if r.Method != "GET" && r.Method != "POST" {
		http.Error(w, "Bad HTTP method: only GET and POST available", http.StatusBadRequest)
		return
	}
	targetURL := strings.TrimPrefix(r.URL.Path, "/")
	if targetURL == "" {
		http.Error(w, "Target URL is required", http.StatusBadRequest)
		return
	}
	file, err := os.OpenFile(LOG_PATH, os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		log.Fatalf("[ERROR]: Failed to open log file: %v", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}
	defer file.Close()
	log.SetOutput(file)
	log.Println("[INFO]: Forward request to " + targetURL)

	if !strings.HasPrefix(targetURL, "http://") && !strings.HasPrefix(targetURL, "https://") {
		targetURL = "http://" + targetURL
	}

	log.Printf("[INFO]: Searching for %s in cache...\n", targetURL)
	if cache.Contains(targetURL) {
		log.Printf("[INFO]: Cache hit for %s.\n", targetURL)

		lm, _ := cache.GetRespMetadata(targetURL)
		resp, respStatusCode := condGET(lm, targetURL)
		if resp == nil {
			log.Printf("[INFO]: Cache for %s is valid.\n", targetURL)
			cache.ReadCachedResponse(targetURL, w)
			return
		} else {
			log.Printf("[INFO]: Cache for %s is expired, updating...\n", targetURL)
			w.WriteHeader(respStatusCode)
			io.Copy(w, bytes.NewReader(resp))
			return
		}
	}
	log.Printf("[INFO]: Cache miss for %s. Sending request to %s.\n", targetURL, targetURL)

	parsedURL, err := url.Parse(targetURL)
	if err != nil {
		http.Error(w, "Invalid target URL", http.StatusBadRequest)
		return
	}

	req := &http.Request{
		Method: r.Method,
		URL:    parsedURL,
		Header: r.Header,
		Body:   r.Body,
	}

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		http.Error(w, "Failed to forward request", http.StatusInternalServerError)
		return
	}
	log.Printf("[INFO]: Response from %s: %s", targetURL, resp.Status)
	defer resp.Body.Close()

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Fatalf("Failed to read response body: %v", err)
	}

	if resp.StatusCode >= 200 && resp.StatusCode < 400 {
		for key, values := range resp.Header {
			for _, value := range values {
				w.Header().Add(key, value)
			}
		}
		w.WriteHeader(resp.StatusCode)
		io.Copy(w, bytes.NewReader(respBody))
		cache.SaveResponse(targetURL, resp.Header.Get("Last-Modified"), resp.Header.Get("ETag"), bytes.NewReader(respBody))
		return
	}
	w.WriteHeader(resp.StatusCode)
	cache.SaveResponse(targetURL, resp.Header.Get("Last-Modified"), resp.Header.Get("ETag"), bytes.NewReader(respBody))
}

func condGET(lm, url string) ([]byte, int) {
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		log.Fatalf("Failed to create request: %v", err)
	}

	req.Header.Set("If-Modified-Since", lm)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		log.Fatalf("Failed to send request: %v", err)
	}
	defer resp.Body.Close()

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Fatalf("Failed to read response body: %v", err)
	}

	switch resp.StatusCode {
	case http.StatusNotModified:
		fmt.Println("Resource has not been modified. Cache is valid.")
		return nil, 0
	case http.StatusOK:
		fmt.Println("Resource has been modified. Updating cache...")
		cache.SaveResponse(url, resp.Header.Get("Last-Modified"), resp.Header.Get("ETag"), bytes.NewReader(respBody))
		return respBody, resp.StatusCode
	default:
		fmt.Printf("Unexpected status code: %d\n", resp.StatusCode)
		return nil, 0
	}
}
