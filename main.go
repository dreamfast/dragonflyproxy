package main

import (
	"flag"
	"fmt"
	"io"
	"log"
	"net"
	"net/http"
	"net/url"
	"strings"
)

func main() {
	// Define command-line flags
	target := flag.String("target", "https://ironman.dragonflybsd.org", "The target URL to forward requests to")
	port := flag.Int("port", 8899, "The port to run the proxy server on")
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

		// Create a new request to forward
		proxyReq, err := http.NewRequest(r.Method, targetURL.String()+r.URL.Path, r.Body)
		if err != nil {
			http.Error(w, "Failed to create proxy request", http.StatusInternalServerError)
			log.Printf("Failed to create proxy request: %v", err)
			return
		}

		// Copy headers from original request
		for header, values := range r.Header {
			for _, value := range values {
				proxyReq.Header.Add(header, value)
			}
		}

		// Set X-Forwarded-For header
		if clientIP, _, err := net.SplitHostPort(r.RemoteAddr); err == nil {
			if prior, ok := proxyReq.Header["X-Forwarded-For"]; ok {
				clientIP = strings.Join(prior, ", ") + ", " + clientIP
			}
			proxyReq.Header.Set("X-Forwarded-For", clientIP)
		}

		// Send the request and get the response
		resp, err := http.DefaultClient.Do(proxyReq)
		if err != nil {
			http.Error(w, "Server Error", http.StatusInternalServerError)
			log.Printf("Failed to forward request: %v", err)
			return
		}
		defer resp.Body.Close()

		// Copy response headers
		for k, v := range resp.Header {
			w.Header()[k] = v
		}

		// Set the status code and copy the body
		w.WriteHeader(resp.StatusCode)
		io.Copy(w, resp.Body)
	})

	addr := fmt.Sprintf(":%d", *port)
	server := &http.Server{
		Addr:    addr,
		Handler: proxy,
	}

	log.Printf("Starting proxy server on %s", addr)
	if err := server.ListenAndServe(); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
