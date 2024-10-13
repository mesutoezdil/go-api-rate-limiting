package main

import (
	"fmt"
	"net/http"
	"sync"
	"time"
)

var (
	mu          sync.Mutex
	requests    = 0
	maxRequests = 5        // Max rate limit
	resetTime   = 10 * time.Second
)

func rateLimiter(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		mu.Lock()
		defer mu.Unlock()

		if requests >= maxRequests {
			http.Error(w, "Rate limit exceeded. Try again later.", http.StatusTooManyRequests)
			return
		}

		requests++
		next.ServeHTTP(w, r)
	})
}

func resetRequests() {
	for {
		time.Sleep(resetTime)
		mu.Lock()
		requests = 0
		mu.Unlock()
	}
}

func helloHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Hello, World!")
}

func main() {
	go resetRequests()

	mux := http.NewServeMux()
	mux.HandleFunc("/", helloHandler)

	limitedMux := rateLimiter(mux)

	http.ListenAndServe(":8080", limitedMux)
}
