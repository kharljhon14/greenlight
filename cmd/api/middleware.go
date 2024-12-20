package main

import (
	"fmt"
	"net"
	"net/http"
	"sync"
	"time"

	"golang.org/x/time/rate"
)

func (app *application) recoverPanic(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Deferred function (which will always be run in the event of a panic)
		defer func() {

			if err := recover(); err != nil {
				// Set "Connection: close"  header on the response
				w.Header().Set("Connection", "close")

				app.serverErrorResponse(w, r, fmt.Errorf("%s", err))
			}
		}()

		next.ServeHTTP(w, r)
	})
}

func (app *application) rateLimit(next http.Handler) http.Handler {
	type client struct {
		limiter  *rate.Limiter
		lastSeen time.Time
	}

	var (
		mu      sync.Mutex
		clients = make(map[string]*client)
	)

	// Background goroutine which removes old entries from the clients map once every minute
	go func() {
		for {
			time.Sleep(time.Minute)

			mu.Lock()

			// Loop through all clients if they haven't been seen within the last three minutes
			for ip, client := range clients {
				if time.Since(client.lastSeen) > 3*time.Minute {
					delete(clients, ip)
				}
			}

			mu.Unlock()
		}
	}()

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Extract the clients IP address
		ip, _, err := net.SplitHostPort(r.RemoteAddr)
		if err != nil {
			app.serverErrorResponse(w, r, err)
			return
		}

		// Lock the mutex
		mu.Lock()

		// Check if the Ip address exists in the map.
		if _, found := clients[ip]; !found {
			clients[ip] = &client{limiter: rate.NewLimiter(2, 4)}
		}

		// IF the request isn't allowed, unlock the mutex and send 429
		if !clients[ip].limiter.Allow() {
			mu.Unlock()
			app.rateLimitExceededResponse(w, r)
			return
		}
		//that the mutex isn't unlocked until all the handlers downstream of this
		// middleware have also returned
		mu.Unlock()

		next.ServeHTTP(w, r)
	})

}
