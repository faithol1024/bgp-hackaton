package main

import (
	"net/http"
	"time"

	"github.com/faithol1024/bgp-hackathon/internal/config"
	"github.com/tokopedia/tdk/go/grace"
)

func startServer(handler http.Handler, cfg *config.Config) error {
	srv := http.Server{
		ReadTimeout:  10 * time.Second, // TODO: read it from config
		WriteTimeout: 10 * time.Second, // TODO: read it from config
		Handler:      handler,
	}

	// ServeHTTP will do these things:
	// - listen on the given port, use socketmaster if exists
	// - setup signal handler which compatible with upstart & ctrl-c
	// - call graceful shutdown when the signal come.
	return grace.ServeHTTP(&srv, cfg.Server.HTTP.Address, time.Second*2)
}
