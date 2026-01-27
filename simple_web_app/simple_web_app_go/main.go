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
	mux.HandleFunc("GET /repo-list/{org_name}", handlers.GetRepoList)
	mux.HandleFunc("GET /protected", middleware.BasicAuth(handlers.Protected))

	port := "8080"
	log.Printf("Starting server on port %s...", port)
	if err := http.ListenAndServe(":"+port, mux); err != nil {
		log.Fatalf("Server failed to start: %v", err)
	}
}
