package server

import (
	"log"
	"net/http"
	"time"
)

const ReadTimeout = 5 * time.Second

// New constructs built-in standard *http.Server.
func New(h http.Handler, addr string, l *log.Logger) *http.Server {
	srv := &http.Server{
		Handler:     h,
		ErrorLog:    l,
		Addr:        addr,
		ReadTimeout: ReadTimeout,
	}

	return srv
}
