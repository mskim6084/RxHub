package middleware

import (
	"fmt"
	"net/http"
	"sync"
	"time"
)

var (
	rateLimit   = 5
	rateWindow  = 10 * time.Second
	mu          sync.Mutex
	lastRequest = make(map[string]time.Time)
	MAX_TOKENS  = 5
	tokensUsed  = 0
)

func TimeBasedRateLimiter(next http.Handler) http.Handler {
	return http.HandlerFunc(func(response http.ResponseWriter, request *http.Request) {
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

func TokenBucketRateLimiter(next http.Handler) http.Handler {
	return http.HandlerFunc(func(response http.ResponseWriter, request *http.Request) {
		if tokensUsed == MAX_TOKENS {
			refillTokenBucket()
			fmt.Printf("Max tokens used reject request until it is filled up\n")
			response.WriteHeader(http.StatusTooManyRequests)
			return
		}

		tokensUsed += 1
		tokenLeft := MAX_TOKENS - tokensUsed
		fmt.Printf("Tokens left %d\n", tokenLeft)

		next.ServeHTTP(response, request)
	})
}

func refillTokenBucket() {
	go func() {
		timer := time.NewTimer(2 * time.Second)
		<-timer.C
		fmt.Printf("Timer is done now!\n Filling up more tokens\n")
		tokensUsed = 0
		tokensLeft := MAX_TOKENS - tokensUsed

		fmt.Printf("Tokens left now are %d\n", tokensLeft)
	}()
}
