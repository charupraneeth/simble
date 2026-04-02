package main

import (
	"context"
	"crypto/sha256"
	"encoding/json"
	"fmt"
	"log"
	"net"
	"net/http"
	"net/netip"
	"net/url"
	"os"
	"strings"
	"sync"
	"time"

	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/medama-io/go-useragent"
	"github.com/oschwald/geoip2-golang/v2"
)

type Event struct {
	Device      string `json:"device"`
	Browser     string `json:"browser"`
	Os          string `json:"os"`
	Referer     string `json:"referer"`
	VisitorID   string `json:"visitor_id"`
	CountryCode string `json:"country_code"`
	City        string `json:"city"`
	Name        string `json:"name"`
	Path        string `json:"path"`
	SiteID      int64  `json:"site_id"`
}

type RequestPayload struct {
	Domain  string `json:"domain"`
	Name    string `json:"name"`
	Referer string `json:"referer"`
	URL     string `json:"url"`
}

type App struct {
	GeoDB    *geoip2.Reader
	UAParser *useragent.Parser
	DB       *pgxpool.Pool
}

func getRealIP(r *http.Request) string {
	// Check for X-Forwared-For header
	if xff := r.Header.Get("X-Forwarded-For"); xff != "" {
		return strings.TrimSpace(strings.Split(xff, ",")[0])
	}

	// Check for X-Real-IP header
	if xri := r.Header.Get("X-Real-IP"); xri != "" {
		return xri
	}

	host, _, err := net.SplitHostPort(r.RemoteAddr)

	if err != nil {
		fmt.Println("error splitting host/port: ", err)
		return ""
	}

	return host

}

func getDailyVisitorID(host, ua, salt string) string {
	date := time.Now().Format("2006-01-02")

	data := fmt.Sprintf("%s|%s|%s|%s", host, ua, salt, date)

	hash := sha256.Sum256([]byte(data))

	return fmt.Sprintf("%x", hash)[:16]
}

func (app *App) handleRequest(w http.ResponseWriter, r *http.Request) {

	// Set CORS header
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type")

	if r.Method == http.MethodOptions {
		w.WriteHeader(http.StatusOK)
		return
	}

	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	reqOrigin := r.Header.Get("Origin")

	originURL, err := url.Parse(reqOrigin)
	if err != nil {
		http.Error(w, "Bad request", http.StatusBadRequest)
		return
	}

	r.Body = http.MaxBytesReader(w, r.Body, 1024) // 1KB
	var payload RequestPayload
	if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
		http.Error(w, "Invalid JSON", http.StatusBadRequest)
		return
	}

	if originURL.Hostname() != payload.Domain {
		http.Error(w, "Invalid domain", http.StatusForbidden)
		return
	}

	parsedBodyURL, err := url.Parse(payload.URL)
	if err != nil {
		http.Error(w, "Invalid url", http.StatusForbidden)
		return
	}

	if strings.TrimPrefix(parsedBodyURL.Hostname(), "www.") != payload.Domain {
		http.Error(w, "URL host does not match payload domain", http.StatusForbidden)
		return
	}

	path := parsedBodyURL.Path
	if path == "" {
		path = "/"
	}

	uaString := r.Header.Get("User-Agent")
	agent := app.UAParser.Parse(uaString)

	host := getRealIP(r)

	ip, err := netip.ParseAddr(host)
	if err != nil {
		fmt.Println("error getting ip from ip string", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)

		return
	}

	var siteID int64
	var siteErr error
	var record *geoip2.City
	var geoErr error

	var wg sync.WaitGroup
	wg.Add(2)

	go func() {
		defer wg.Done()
		siteQuery := `SELECT id from sites where domain = $1`
		siteErr = app.DB.QueryRow(r.Context(), siteQuery, payload.Domain).Scan(&siteID)
	}()

	go func() {
		defer wg.Done()
		record, geoErr = app.GeoDB.City(ip)
	}()

	wg.Wait()

	if siteErr != nil {
		if siteErr == pgx.ErrNoRows {
			http.Error(w, "Site not registered", http.StatusNotFound)
			return
		}

		log.Printf("Error querying site %v\n", siteErr)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	if geoErr != nil {
		log.Printf("GeoIP lookup failed: %v", geoErr)
		// dont return we can still record the event without the geo data
	}

	countryCode := ""
	city := ""
	if record.HasData() {
		countryCode = record.Country.ISOCode
		city = record.City.Names.English
	} else {
		fmt.Println("No data found for this IP")
	}

	salt := os.Getenv("VISITOR_ID_SALT")
	if salt == "" {
		salt = "fallback-temp-salt" // for local dev
	}

	visitorID := getDailyVisitorID(host, uaString, salt)

	insertQuery := `INSERT INTO analytics(
	 					site_id, visitor_id, path, browser_name, device_type, os_name, country_code, city_name
					)
					VALUES($1, $2, $3, $4, $5, $6, $7, $8);`

	event := Event{
		Device:      agent.Device().String(),
		Browser:     agent.Browser().String(),
		Os:          agent.OS().String(),
		Referer:     r.Referer(),
		VisitorID:   visitorID,
		CountryCode: countryCode,
		City:        city,
		Name:        payload.Name,
		Path:        path,
		SiteID:      siteID,
	}

	_, err = app.DB.Exec(r.Context(), insertQuery, event.SiteID, event.VisitorID, event.Path, event.Browser, event.Device, event.Os, event.CountryCode, event.City)

	if err != nil {
		log.Printf("CRITICAL: Failed to insert analytics event for site %d: %v", siteID, err)
		w.WriteHeader(http.StatusAccepted)
		return
	}

	log.Printf("Event inserted: %v\n", event)

	w.WriteHeader(http.StatusAccepted)
}

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

	app := &App{GeoDB: geoDB, UAParser: ua, DB: db}

	mux.HandleFunc("/script.js", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "./public/script.js")
	})
	mux.HandleFunc("/api/event", app.handleRequest)

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
