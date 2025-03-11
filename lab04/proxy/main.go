package main

import (
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

func main() {
	http.HandleFunc("/", handleRequest)
	log.Printf("Starting proxy server on %s...", PORT)
	if err := http.ListenAndServe(PORT, nil); err != nil {
		log.Fatalf("Failed to start proxy server: %v", err)
	}
}

func handleRequest(w http.ResponseWriter, r *http.Request) {
	if r.Method != "GET" {
		http.Error(w, "Failed to handle non-GET request", http.StatusBadRequest)
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

	for key, values := range resp.Header {
		for _, value := range values {
			w.Header().Add(key, value)
		}
	}
	w.WriteHeader(resp.StatusCode)
	io.Copy(w, resp.Body)
}
