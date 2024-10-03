package main

import (
	"flag"
	"io"
	"log"
	"net/http"
	"net/url"
)

func main() {
	// Define a command-line flag for the target URL
	target := flag.String("target", "https://ironman.dragonflybsd.org", "The target URL to forward requests to")
	flag.Parse()

	targetURL, err := url.Parse(*target)
	if err != nil {
		log.Fatalf("Failed to parse target URL: %v", err)
	}

	proxy := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Set CORS headers
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")

		// Handle preflight requests
		if r.Method == http.MethodOptions {
			w.WriteHeader(http.StatusOK)
			return
		}

		r.URL.Scheme = targetURL.Scheme
		r.URL.Host = targetURL.Host
		r.Host = targetURL.Host

		resp, err := http.DefaultTransport.RoundTrip(r)
		if err != nil {
			http.Error(w, "Server Error", http.StatusInternalServerError)
			log.Printf("Failed to forward request: %v", err)
			return
		}
		defer resp.Body.Close()

		for k, v := range resp.Header {
			w.Header()[k] = v
		}
		w.WriteHeader(resp.StatusCode)
		io.Copy(w, resp.Body)
	})

	log.Println("Starting proxy server on :8899")
	if err := http.ListenAndServe(":8899", proxy); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
