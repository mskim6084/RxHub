package middleware

import (
	"fmt"
	"net/http"
	"sync"
	"time"
)

var (
	rateLimit = 5
	rateWindow = 10 * time.Second
	mu sync.Mutex
	lastRequest = make(map[string]time.Time)
)

func RateLimiter(next http.Handler) http.Handler {
	return http.HandlerFunc(func(response http.ResponseWriter, request *http.Request){
		clientIP := request.RemoteAddr

		mu.Lock()
		defer mu.Unlock()

		fmt.Printf("Rate Limiter check!")

		if lastTime, exists := lastRequest[clientIP]; exists {
			
			if time.Since(lastTime) < rateWindow {
				response.WriteHeader(http.StatusTooManyRequests)
				fmt.Printf("Rate limit exceeded. Please try again in a bit")
				return
			}
		}

		lastRequest[clientIP] = time.Now()

		next.ServeHTTP(response, request)
	})
}