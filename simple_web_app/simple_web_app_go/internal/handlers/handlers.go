package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"
)

type HelloWorldResponse struct {
	Message string `json:"message"`
}

type ProtectedResponse struct {
	Message string `json:"message"`
}

// HelloWorld handles the /hello-world endpoint
func HelloWorld(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(HelloWorldResponse{Message: "Hello World"})
}

// GetRepoList handles the /repo-list/{org_name}[?repo_filter=filter] endpoint
func GetRepoList(w http.ResponseWriter, r *http.Request) {
	//orgName := r.PathValue("org_name")
	//repoFilter := r.URL.Query().Get("repo_filter")
}

// Protected handles the /protected endpoint
func Protected(w http.ResponseWriter, r *http.Request) {
	username, _, _ := r.BasicAuth()
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(ProtectedResponse{
		Message: fmt.Sprintf("Login successful for %s", username),
	})
}
