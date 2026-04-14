package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"simble/internal/utils"
)

func getGithubDetails(token string) (GitHubUser, error) {
	client := &http.Client{}
	url := "https://api.github.com/user"

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return GitHubUser{}, fmt.Errorf("error creating request to github user: %w", err)
	}

	req.Header.Set("Authorization", "token "+token)

	resp, err := client.Do(req)
	if err != nil {
		return GitHubUser{}, fmt.Errorf("error requesting github user: %w", err)
	}

	defer resp.Body.Close()
	var githubUser GitHubUser
	if err := json.NewDecoder(resp.Body).Decode(&githubUser); err != nil {
		return GitHubUser{}, fmt.Errorf("error decoding github user response: %w", err)
	}

	return githubUser, nil
}

func (app *App) handleGitHubLogin(w http.ResponseWriter, r *http.Request) {
	state, err := utils.GenerateRandomToken(32)
	if err != nil {
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	http.SetCookie(w, &http.Cookie{
		Name:     "oauth_state",
		Value:    state,
		HttpOnly: true,
		MaxAge:   300,
		SameSite: http.SameSiteLaxMode,
		Path:     "/",
	})

	url := app.OAuthConfig.AuthCodeURL(state)
	http.Redirect(w, r, url, http.StatusTemporaryRedirect)
}

func (app *App) handleGitHubCallback(w http.ResponseWriter, r *http.Request) {
	stateCookie, err := r.Cookie("oauth_state")

	if err != nil || stateCookie.Value != r.URL.Query().Get("state") {
		http.Error(w, "Invalid state", http.StatusBadRequest)
		return
	}

	code := r.URL.Query().Get("code")

	token, err := app.OAuthConfig.Exchange(r.Context(), code)
	if err != nil {
		http.Error(w, "Failed to exchange token", http.StatusInternalServerError)
		return
	}

	githubUser, err := getGithubDetails(token.AccessToken)
	if err != nil {
		http.Error(w, "Failed to fetch github user", http.StatusInternalServerError)
		return
	}

	userIDQuery := `
		INSERT INTO users (github_id, email, username)
		VALUES($1, $2, $3)
		on conflict(github_id) 
		do update set email = EXCLUDED.email
		returning id
	`

	var userID int64
	err = app.DB.QueryRow(r.Context(), userIDQuery, githubUser.ID, githubUser.Email, githubUser.Login).Scan(&userID)
	if err != nil {
		log.Printf("Failed to insert user into DB: %v", err)
		http.Error(w, "Failed to register github user", http.StatusInternalServerError)
		return
	}

	sessionToken, err := utils.GenerateRandomToken(32)
	if err != nil {
		http.Error(w, "Failed to generate session", http.StatusInternalServerError)
		return
	}

	expiresAt := time.Now().Add(30 * 24 * time.Hour) // 30 days from now

	sessionQuery := `
		INSERT INTO sessions (token, user_id , expires_at)
		VALUES ($1, $2, $3)
	`

	_, err = app.DB.Exec(r.Context(), sessionQuery, sessionToken, userID, expiresAt)
	if err != nil {
		http.Error(w, "Failed to create session", http.StatusInternalServerError)
		return
	}

	http.SetCookie(w, &http.Cookie{
		Name:     "session",
		Value:    sessionToken,
		HttpOnly: true,
		Expires:  expiresAt,
		Secure:   os.Getenv("ENV") == "production",
		SameSite: http.SameSiteLaxMode,
		Path:     "/",
	})

	frontendURL := os.Getenv("FRONTEND_URL")
	if frontendURL == "" {
		frontendURL = "/"
	}

	http.Redirect(w, r, frontendURL, http.StatusTemporaryRedirect)
}

func (app *App) handleGetMe(w http.ResponseWriter, r *http.Request) {
	user := r.Context().Value(userKey).(User)

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(user)
}

func (app *App) handleLogout(w http.ResponseWriter, r *http.Request) {
	token, err := r.Cookie("session")
	if err != nil {
		http.Error(w, "Failed to delete user session", http.StatusBadRequest)
		return
	}

	query := `
		DELETE from sessions
		where token = $1
	`

	_, err = app.DB.Exec(r.Context(), query, token.Value)
	if err != nil {
		http.Error(w, "Failed to remove user session", http.StatusInternalServerError)
		return
	}

	// Delete the cookie on the client side by expiring it immediately
	http.SetCookie(w, &http.Cookie{
		Name:     "session",
		Value:    "",
		HttpOnly: true,
		MaxAge:   -1,
		Expires:  time.Now().Add(-100 * time.Hour),
		Path:     "/",
	})

	w.WriteHeader(http.StatusOK)
}
