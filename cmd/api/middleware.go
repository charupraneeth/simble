package main

import (
	"context"
	"log"
	"net/http"
)

type contextKey string

const userKey contextKey = "user"

func (app *App) requireAuth(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		token, err := r.Cookie("session")
		if err != nil {
			http.Error(w, "Invalid session cookie", http.StatusBadRequest)
			return
		}

		query := `
			SELECT u.id, u.username, u.email, s.expires_at
			FROM sessions s
			JOIN users u ON s.user_id = u.id
			WHERE s.token  = $1 
			AND s.expires_at > CURRENT_TIMESTAMP;
		`

		var user User
		err = app.DB.QueryRow(r.Context(), query, token.Value).Scan(&user.ID, &user.Username, &user.Email, &user.ExpiresAt)
		if err != nil {
			log.Println("Failed to get user with given session token")
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}

		ctx := context.WithValue(r.Context(), userKey, user)
		next.ServeHTTP(w, r.WithContext(ctx))
	}
}
