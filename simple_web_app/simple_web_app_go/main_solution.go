package main

import (
	"log"
	"net/http"

	"simple-web-app/internal/handlers"
	"simple-web-app/internal/middleware"
)

func main() {
	mux := http.NewServeMux()

	mux.HandleFunc("GET /hello-world", handlers.HelloWorld)
	// Using the solution implementation for testing
	mux.HandleFunc("GET /repo-list/{org_name}", handlers.GetRepoListSolution)
	mux.HandleFunc("GET /protected", middleware.BasicAuth(handlers.Protected))

	port := "8080"
	log.Printf("Starting server on port %s...", port)
	log.Printf("Available endpoints:")
	log.Printf("  GET /hello-world")
	log.Printf("  GET /repo-list/{org_name}[?repo_filter=filter]")
	log.Printf("  GET /protected (requires AUTH_USERNAME and AUTH_PASSWORD env vars)")
	if err := http.ListenAndServe(":"+port, mux); err != nil {
		log.Fatalf("Server failed to start: %v", err)
	}
}
