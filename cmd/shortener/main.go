package main

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"io"
	"net/http"
)

var urls = make(map[string]string)

func getShortUrl(res http.ResponseWriter, req *http.Request) {
	if req.Method == http.MethodPost {
		body, err := io.ReadAll(req.Body)
		if err != nil {
			panic(err)
		}

		h := sha256.New()
		h.Write([]byte(body))
		hash := "/" + hex.EncodeToString(h.Sum(nil))

		result := fmt.Sprintf("http://localhost:8080" + hash[:8])

		res.Header().Set("content-type", "text/plain")
		res.WriteHeader(http.StatusCreated)

		res.Write([]byte(result))
		urls[hash[:8]] = string(body)
	} else {
		http.Error(res, "Only POST requests are allowed!", http.StatusMethodNotAllowed)
		return
	}
}

func redirectTo(res http.ResponseWriter, req *http.Request) {
	if req.Method == http.MethodGet {
		shortUrl := req.URL.String()

		result := urls[shortUrl]
		if result == "" {
			http.Error(res, "ShortURL not found", http.StatusNotFound)
			return
		}

		res.Header().Set("Location", result)
		res.WriteHeader(http.StatusTemporaryRedirect)
	} else {
		http.Error(res, "Only GET requests are allowed!", http.StatusMethodNotAllowed)
		return
	}
}

func main() {
	mux := http.NewServeMux()
	mux.HandleFunc(`/`, getShortUrl)
	mux.HandleFunc(`/{id}`, redirectTo)

	err := http.ListenAndServe(`:8080`, mux)
	if err != nil {
		panic(err)
	}
}
