package main

import (
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
	w.WriteHeader(http.StatusOK)

	response := map[string]string{"status": "ok"}
	json.NewEncoder(w).Encode(response)
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

	proxy := httputil.NewSingleHostReverseProxy(remote)

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		r.Host = remote.Host
		proxy.ServeHTTP(w, r)
	})

	http.HandleFunc("/health", helloWorldHandler)

	log.Printf("Reverse proxy running on :%s, forwarding to %s", proxyPort, targetURL)
	err = http.ListenAndServe(":"+proxyPort, nil)
	if err != nil {
		log.Fatal("Server failed:", err)
	}
}
