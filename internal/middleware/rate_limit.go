package middleware

import (
	"net"
	"net/http"
	"sync"
	"time"

	erp "bookloop.net/pkg/error_responses"
	"golang.org/x/time/rate"
)

func (mw *MiddlewareManager) RateLimit(next http.Handler) http.Handler {
	type client struct {
		limiter  *rate.Limiter
		lastSeen time.Time
	}

	var (
		mu      sync.Mutex
		clients = make(map[string]*client)
	)

	go func() {
		for {
			time.Sleep(time.Minute)

			mu.Lock()

			for ip, client := range clients {
				if time.Since(client.lastSeen) > 3*time.Minute {
					delete(clients, ip)
				}

			}

			mu.Unlock()
		}
	}()

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ip, _, err := net.SplitHostPort(r.RemoteAddr)
		if err != nil {
			erp.ServerErrorResponse(w, r, mw.logger, err)
			return
		}

		mu.Lock()

		if _, found := clients[ip]; !found {
			clients[ip] = &client{limiter: rate.NewLimiter(rate.Limit(mw.cfg.Limiter.RPS), mw.cfg.Limiter.Burst)}
		}

		if !clients[ip].limiter.Allow() {
			mu.Unlock()
			erp.RateLimitExceededResponse(w, r, mw.logger)
			return
		}

		mu.Unlock()

		next.ServeHTTP(w, r)
	})
}
