package main

import (
	"log"
	"net/http"

	"status-board/internal/config"
	"status-board/internal/handlers"
)

func main() {
	cfg := config.Load()
	if missing := cfg.Missing(); len(missing) > 0 {
		log.Printf("WARNING: incomplete configuration, /settings will fail: missing %v", missing)
	}

	h := handlers.New(cfg)

	mux := http.NewServeMux()
	mux.HandleFunc("GET /healthz", handlers.Healthz)
	mux.HandleFunc("GET /hello-world", handlers.HelloWorld)
	mux.HandleFunc("GET /settings", h.Settings)
	mux.HandleFunc("GET /checks", h.Checks)

	log.Printf("Starting server on port %s...", cfg.Port)
	if err := http.ListenAndServe(":"+cfg.Port, mux); err != nil {
		log.Fatalf("Server failed to start: %v", err)
	}
}
