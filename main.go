package main

import (
	"crypto/tls"
	"encoding/json"
	"log"
	"net/http"
	"net/http/httputil"
	"net/url"
	"os"

	"github.com/joho/godotenv"
)

func helloWorldHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{"status": "ok"})
}

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	targetURL := os.Getenv("TARGET_URL")
	proxyPort := os.Getenv("PROXY_PORT")

	if targetURL == "" || proxyPort == "" {
		log.Fatal("TARGET_URL or PROXY_PORT is missing from .env")
	}

	remote, err := url.Parse(targetURL)
	if err != nil {
		log.Fatal("Invalid TARGET_URL:", err)
	}

	// Set up the reverse proxy
	proxy := httputil.NewSingleHostReverseProxy(remote)

	// Override the transport to skip TLS verification
	proxy.Transport = &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	}

	// Modify the request to match the expected host
	director := proxy.Director
	proxy.Director = func(r *http.Request) {
		director(r)
		r.Host = remote.Host
	}

	// Routes
	http.HandleFunc("/health", helloWorldHandler)
	http.HandleFunc("/", proxy.ServeHTTP)

	log.Printf("Reverse proxy running on :%s, forwarding to %s", proxyPort, targetURL)
	err = http.ListenAndServe(":"+proxyPort, nil)
	if err != nil {
		log.Fatal("Server failed:", err)
	}
}
