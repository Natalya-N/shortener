package main

import (
	"fmt"
	"io"
	"log"
	"math/rand"
	"net/http"
	"net/url"
	"strings"
	"time"
)

var (
	urlMap = make(map[string]string)
)

func generateRandomString() string {
	const charset = "abcdefghijklmnopqrstuvwxyz"
	var result string
	for i := 0; i < 6; i++ {
		randomIndex := rand.Intn(len(charset))
		result += string(charset[randomIndex])
	}
	return result
}

func handleRequests(res http.ResponseWriter, req *http.Request) {
	if req.Method == http.MethodPost {
		body, err := io.ReadAll(req.Body)
		if err != nil {
			res.WriteHeader(http.StatusBadRequest)
			return
		}

		if len(body) == 0 {
			res.WriteHeader(http.StatusBadRequest)
			return
		}

		originalURL := strings.TrimSpace(string(body))
		parsedURL, err := url.Parse(originalURL)
		if err != nil || parsedURL.Scheme == "" {
			res.WriteHeader(http.StatusBadRequest)
			return
		}

		id := generateRandomString()
		urlMap[id] = originalURL
		log.Printf("%s %s\n", id, originalURL)

		shortenedURL := fmt.Sprintf("http://localhost:8080/%s", id)
		res.Header().Set("Content-Type", "text/plain")
		res.WriteHeader(http.StatusCreated)
		fmt.Fprintf(res, "%s", shortenedURL)
	} else if req.Method == http.MethodGet {
		parts := strings.Split(req.URL.Path, "/")
		if len(parts) < 2 {
			http.Error(res, "url not found", http.StatusNotFound)
			return
		}

		idUrl := parts[1]
		originalURL, exists := urlMap[idUrl]
		if !exists {
			http.Error(res, "url not found", http.StatusNotFound)
			return
		}

		http.Redirect(res, req, originalURL, http.StatusTemporaryRedirect)
	} else {
		http.Error(res, "Invalid request method", http.StatusMethodNotAllowed)
	}
}

func main() {
	fmt.Println("Starting server...")
	rand.Seed(time.Now().UnixNano())
	mux := http.NewServeMux()
	mux.HandleFunc("/", handleRequests)

	err := http.ListenAndServe(":8081", mux)
	if err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
