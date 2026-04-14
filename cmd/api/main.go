package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/medama-io/go-useragent"
	"github.com/oschwald/geoip2-golang/v2"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/github"
)

func main() {
	mux := http.NewServeMux()
	ctx := context.Background()

	geoDB, err := geoip2.Open("./GeoLite2-City/GeoLite2-City.mmdb")
	if err != nil {
		log.Fatal(err)
	}
	defer geoDB.Close()

	dbString := os.Getenv("DATABASE_URL")

	if os.Getenv("AUTO_MIGRATE_DB") == "true" {
		log.Println("Running database migrations...")
		m, err := migrate.New("file://migrations", dbString)

		if err != nil && err != migrate.ErrNoChange {
			log.Fatalf("Failed to initialize migrations: %v", err)
		}
		if err := m.Up(); err != nil && err != migrate.ErrNoChange {
			log.Fatalf("Failed to run migrations: %v", err)
		}
		log.Println("Database Migrations applied successfully!")
	}

	db, err := pgxpool.New(ctx, dbString)
	if err != nil {
		log.Fatal(err)
	}

	ua := useragent.NewParser()

	oauthConfig := &oauth2.Config{
		ClientID:     os.Getenv("GITHUB_CLIENT_ID"),
		ClientSecret: os.Getenv("GITHUB_CLIENT_SECRET"),
		Scopes:       []string{"user:email"},
		Endpoint:     github.Endpoint,
		RedirectURL:  os.Getenv("GITHUB_REDIRECT_URL"),
	}

	app := &App{GeoDB: geoDB, UAParser: ua, DB: db, OAuthConfig: oauthConfig, DemoSiteID: os.Getenv("DEMO_SITE_ID")}

	mux.HandleFunc("GET /script.js", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "./public/script.js")
	})
	mux.HandleFunc("GET /auth/github", app.handleGitHubLogin)
	mux.HandleFunc("GET /auth/github/callback", app.handleGitHubCallback)
	mux.HandleFunc("/api/event", app.handleRequest)
	mux.HandleFunc("GET /api/me", app.requireAuth(app.handleGetMe))
	mux.HandleFunc("DELETE /api/session", app.requireAuth(app.handleLogout))

	mux.HandleFunc("GET /api/sites", app.requireAuth(app.handleGetSites))
	mux.HandleFunc("POST /api/sites", app.requireAuth(app.handleCreateSite))
	mux.HandleFunc("GET /api/sites/{id}/stats", app.requireAuth(app.handleGetSiteStats))
	mux.HandleFunc("GET /api/sites/{id}/traffic", app.requireAuth(app.handleGetSiteTraffic))
	mux.HandleFunc("GET /api/sites/{id}/pages", app.requireAuth(app.handleGetSitePages))
	mux.HandleFunc("GET /api/sites/{id}/countries", app.requireAuth(app.handleGetSiteCountries))
	mux.HandleFunc("GET /api/sites/{id}/referrers", app.requireAuth(app.handleGetSiteReferrers))

	// Public demo endpoints
	mux.HandleFunc("GET /api/demo/stats", app.handleDemoStats)
	mux.HandleFunc("GET /api/demo/traffic", app.handleDemoTraffic)
	mux.HandleFunc("GET /api/demo/pages", app.handleDemoPages)
	mux.HandleFunc("GET /api/demo/countries", app.handleDemoCountries)
	mux.HandleFunc("GET /api/demo/referrers", app.handleDemoReferrers)

	distFS := http.Dir("./public/dist")
	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		path := r.URL.Path
		if f, err := distFS.Open(path); err == nil {
			f.Close()
			http.FileServer(distFS).ServeHTTP(w, r)
			return
		}
		// Fallback to index.html for all unmatched routes and let vue router handle it
		http.ServeFile(w, r, "./public/dist/index.html")
	})

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	server := &http.Server{
		Addr:         ":" + port,
		Handler:      mux,
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 10 * time.Second,
		IdleTimeout:  60 * time.Second,
	}

	log.Fatal(server.ListenAndServe())
}
