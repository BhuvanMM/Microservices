package main

import (
	"io"
	"log"
	"net/http"
	"strings"
)

// This function forwards the incoming request to the target service.
func proxyRequest(w http.ResponseWriter, r *http.Request, targetBase string) {
	targetURL := "http://" + targetBase + r.URL.Path

	req, err := http.NewRequest(r.Method, targetURL, r.Body)
	if err != nil {
		http.Error(w, "Failed to create request", http.StatusInternalServerError)
		return
	}

	req.Header = r.Header

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		http.Error(w, "Failed to forward request", http.StatusBadGateway)
		return
	}
	defer resp.Body.Close()

	w.WriteHeader(resp.StatusCode)
	io.Copy(w, resp.Body)
}

func routeHandler(w http.ResponseWriter, r *http.Request) {
	path := r.URL.Path

	switch {
	case strings.HasPrefix(path, "/user/"):
		proxyRequest(w, r, "localhost:8001")

	case strings.HasPrefix(path, "/company/"):
		proxyRequest(w, r, "localhost:8002")

	case strings.HasPrefix(path, "/job/"):
		proxyRequest(w, r, "localhost:8003")

	case strings.HasPrefix(path, "/apply/"):
		proxyRequest(w, r, "localhost:8004")

	default:
		http.Error(w, "Route not found", http.StatusNotFound)
	}
}

func main() {
	log.Println("API Gateway running at http://localhost:8000")
	http.HandleFunc("/", routeHandler)
	log.Fatal(http.ListenAndServe(":8000", nil))
}
