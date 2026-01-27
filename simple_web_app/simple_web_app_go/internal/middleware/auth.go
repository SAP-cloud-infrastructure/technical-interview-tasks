package middleware

import (
	"crypto/subtle"
	"net/http"
	"os"
)

// BasicAuth middleware for HTTP Basic Authentication
func BasicAuth(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		username, password, ok := r.BasicAuth()
		if !ok {
			w.Header().Set("WWW-Authenticate", `Basic realm="Restricted"`)
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}

		correctUsername := os.Getenv("AUTH_USERNAME")
		correctPassword := os.Getenv("AUTH_PASSWORD")

		if correctUsername == "" || correctPassword == "" {
			http.Error(w, "Server configuration error", http.StatusInternalServerError)
			return
		}

		usernameMatch := subtle.ConstantTimeCompare([]byte(username), []byte(correctUsername)) == 1
		passwordMatch := subtle.ConstantTimeCompare([]byte(password), []byte(correctPassword)) == 1

		if !usernameMatch || !passwordMatch {
			w.Header().Set("WWW-Authenticate", `Basic realm="Restricted"`)
			http.Error(w, "Incorrect username or password", http.StatusUnauthorized)
			return
		}

		next(w, r)
	}
}
