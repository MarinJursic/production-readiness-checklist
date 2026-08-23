package riskyserver

import (
	"net/http"
	"time"
)

func serve() error {
	server := http.Server{
		Addr:              ":8080",
		ReadHeaderTimeout: 5 * time.Second,
		ReadTimeout:       30 * time.Second,
		WriteTimeout:      30 * time.Second,
		IdleTimeout:       60 * time.Second,
	}
	_ = server
	return server.ListenAndServe()
}
