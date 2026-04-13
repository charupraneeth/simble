package main

import (
	"context"
	"crypto/rand"
	"crypto/sha256"
	"database/sql"
	"encoding/base64"
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
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/medama-io/go-useragent"
	"github.com/oschwald/geoip2-golang/v2"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/github"
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

type GitHubUser struct {
	ID    int64  `json:"id"`
	Login string `json:"login"` // the @username
	Email string `json:"email"` // might be empty string or null
}

type User struct {
	ID        int64     `json:"-"` // Hidden from frontend JSON responses
	Username  string    `json:"username"`
	Email     *string   `json:"email"` // this can be null if the user had made it private
	ExpiresAt time.Time `json:"expires_at"`
}

type Site struct {
	ID        int64     `json:"id" db:"id"`
	Domain    string    `json:"domain" db:"domain"`
	CreatedAt time.Time `json:"created_at" db:"created_at"`
	Visitors  int64     `json:"visitors" db:"visitors"`
}

type CreateSitePayload struct {
	Domain string `json:"domain"`
}

type SiteStats struct {
	UniqueVisitors int64 `json:"unique_visitors"`
	Pageviews      int64 `json:"pageviews"`
}

type TrafficPoint struct {
	Hour     time.Time `json:"hour"`
	Visitors int64     `json:"visitors"`
}

type TopPage struct {
	Path           string `json:"path"`
	Views          int64  `json:"views"`
	UniqueVisitors int64  `json:"unique_visitors"`
}

type TopCountry struct {
	CountryCode    string `json:"country_code"`
	Views          int64  `json:"views"`
	UniqueVisitors int64  `json:"unique_visitors"`
}

type TopReferrer struct {
	Referrer       string `json:"referrer"`
	Views          int64  `json:"views"`
	UniqueVisitors int64  `json:"unique_visitors"`
}

type App struct {
	GeoDB       *geoip2.Reader
	UAParser    *useragent.Parser
	DB          *pgxpool.Pool
	OAuthConfig *oauth2.Config
	DemoSiteID  string
}

type contextKey string

const userKey contextKey = "user"

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
		log.Println("error splitting host/port: ", err)
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

func generateRandomToken(n int) (string, error) {
	b := make([]byte, n)
	_, err := rand.Read(b)

	if err != nil {
		return "", err
	}

	return base64.URLEncoding.EncodeToString(b), nil
}

func getGithubDetails(token string) (GitHubUser, error) {
	client := &http.Client{}
	url := "https://api.github.com/user"

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return GitHubUser{}, fmt.Errorf("Error creating request to github user: %w", err)
	}

	req.Header.Set("Authorization", "token "+token)

	resp, err := client.Do(req)
	if err != nil {
		return GitHubUser{}, fmt.Errorf("Error requesting github user: %w", err)
	}

	defer resp.Body.Close()
	var githubUser GitHubUser
	if err := json.NewDecoder(resp.Body).Decode(&githubUser); err != nil {
		return GitHubUser{}, fmt.Errorf("Error decoding github user response: %w", err)
	}

	return githubUser, nil
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

	isLocalOrigin := originURL.Hostname() == "localhost" || originURL.Hostname() == "127.0.0.1"

	if originURL.Hostname() != payload.Domain && !isLocalOrigin {
		http.Error(w, "Invalid domain", http.StatusForbidden)
		return
	}

	parsedBodyURL, err := url.Parse(payload.URL)
	if err != nil {
		http.Error(w, "Invalid url", http.StatusForbidden)
		return
	}

	isLocalPage := parsedBodyURL.Hostname() == "localhost" || parsedBodyURL.Hostname() == "127.0.0.1"

	if strings.TrimPrefix(parsedBodyURL.Hostname(), "www.") != payload.Domain && !isLocalPage {
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
		log.Println("error getting ip from ip string", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)

		return
	}

	// Drop local development events so they don't pollute production data
	if isLocalOrigin && os.Getenv("ENV") != "development" {
		w.WriteHeader(http.StatusAccepted)
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
		log.Println("No data found for this IP")
	}

	salt := os.Getenv("VISITOR_ID_SALT")
	if salt == "" {
		salt = "fallback-temp-salt" // for local dev
	}

	visitorID := getDailyVisitorID(host, uaString, salt)

	insertQuery := `INSERT INTO analytics(
	 					site_id, visitor_id, path, browser_name, device_type, os_name, country_code, city_name, referrer
					)
					VALUES($1, $2, $3, $4, $5, $6, $7, $8, $9);`

	event := Event{
		Device:      agent.Device().String(),
		Browser:     agent.Browser().String(),
		Os:          agent.OS().String(),
		Referer:     payload.Referer,
		VisitorID:   visitorID,
		CountryCode: countryCode,
		City:        city,
		Name:        payload.Name,
		Path:        path,
		SiteID:      siteID,
	}

	_, err = app.DB.Exec(r.Context(), insertQuery, event.SiteID, event.VisitorID, event.Path, event.Browser, event.Device, event.Os, event.CountryCode, event.City, event.Referer)

	if err != nil {
		log.Printf("CRITICAL: Failed to insert analytics event for site %d: %v", siteID, err)
		w.WriteHeader(http.StatusAccepted)
		return
	}

	log.Printf("Event inserted: %v\n", event)

	w.WriteHeader(http.StatusAccepted)
}

func (app *App) handleGitHubLogin(w http.ResponseWriter, r *http.Request) {
	state, err := generateRandomToken(32)
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

	sessionToken, err := generateRandomToken(32)
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
		http.Error(w, "Failed to created session", http.StatusInternalServerError)
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

func (app *App) handleGetSites(w http.ResponseWriter, r *http.Request) {
	userID := r.Context().Value(userKey).(User).ID

	query := `
		SELECT
			s.id,
			s.domain,
			s.created_at,
			COALESCE(live.visitors, 0) AS visitors
		FROM sites s
		LEFT JOIN (
			SELECT site_id, COUNT(DISTINCT visitor_id) AS visitors
			FROM analytics
			WHERE created_at > NOW() - INTERVAL '5 minutes'
			GROUP BY site_id
		) live ON live.site_id = s.id
		WHERE s.user_id = $1
		ORDER BY s.created_at DESC
	`

	rows, err := app.DB.Query(r.Context(), query, userID)
	if err != nil {
		http.Error(w, "Failed to query sites", http.StatusInternalServerError)
		return
	}

	sites, err := pgx.CollectRows(rows, pgx.RowToStructByName[Site])
	if err != nil {
		log.Printf("pgx CollectRows error: %v", err)
		http.Error(w, fmt.Sprintf("Failed to map sites: %v", err), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(sites)
}

// parseDateRange reads the ?range= query param and returns a Postgres interval
// string and a truncation unit for GROUP BY queries.
// Supported values: "24h" (default), "7d", "30d".
func parseDateRange(r *http.Request) (interval string, trunc string) {
	switch r.URL.Query().Get("range") {
	case "7d":
		return "7 days", "day"
	case "30d":
		return "30 days", "day"
	default:
		return "24 hours", "hour"
	}
}

func (app *App) serveStats(w http.ResponseWriter, r *http.Request, siteID string) {
	interval, _ := parseDateRange(r)
	var stats SiteStats
	err := app.DB.QueryRow(r.Context(), fmt.Sprintf(`
		SELECT
			COUNT(DISTINCT visitor_id) AS unique_visitors,
			COUNT(*) AS pageviews
		FROM analytics
		WHERE site_id = $1
		  AND created_at > NOW() - INTERVAL '%s'
	`, interval), siteID).Scan(&stats.UniqueVisitors, &stats.Pageviews)
	if err != nil {
		log.Printf("serveStats error: %v", err)
		http.Error(w, "Failed to get stats", http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(stats)
}

func (app *App) serveTraffic(w http.ResponseWriter, r *http.Request, siteID string) {
	interval, trunc := parseDateRange(r)
	rows, err := app.DB.Query(r.Context(), fmt.Sprintf(`
		SELECT
			DATE_TRUNC('%s', created_at) AS hour,
			COUNT(DISTINCT visitor_id) AS visitors
		FROM analytics
		WHERE site_id = $1
		  AND created_at > NOW() - INTERVAL '%s'
		GROUP BY hour
		ORDER BY hour ASC
	`, trunc, interval), siteID)
	if err != nil {
		log.Printf("serveTraffic error: %v", err)
		http.Error(w, "Failed to get traffic", http.StatusInternalServerError)
		return
	}
	defer rows.Close()
	var points []TrafficPoint
	for rows.Next() {
		var p TrafficPoint
		if err := rows.Scan(&p.Hour, &p.Visitors); err != nil {
			continue
		}
		points = append(points, p)
	}
	if points == nil {
		points = []TrafficPoint{}
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(points)
}

func (app *App) servePages(w http.ResponseWriter, r *http.Request, siteID string) {
	interval, _ := parseDateRange(r)
	rows, err := app.DB.Query(r.Context(), fmt.Sprintf(`
		SELECT
			path,
			COUNT(*) AS views,
			COUNT(DISTINCT visitor_id) AS unique_visitors
		FROM analytics
		WHERE site_id = $1
		  AND created_at > NOW() - INTERVAL '%s'
		GROUP BY path
		ORDER BY views DESC
		LIMIT 10
	`, interval), siteID)
	if err != nil {
		log.Printf("servePages error: %v", err)
		http.Error(w, "Failed to get pages", http.StatusInternalServerError)
		return
	}
	defer rows.Close()
	var pages []TopPage
	for rows.Next() {
		var p TopPage
		if err := rows.Scan(&p.Path, &p.Views, &p.UniqueVisitors); err != nil {
			continue
		}
		pages = append(pages, p)
	}
	if pages == nil {
		pages = []TopPage{}
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(pages)
}

func (app *App) serveCountries(w http.ResponseWriter, r *http.Request, siteID string) {
	interval, _ := parseDateRange(r)
	rows, err := app.DB.Query(r.Context(), fmt.Sprintf(`
		SELECT
			COALESCE(country_code, 'Unknown') AS country_code,
			COUNT(*) AS views,
			COUNT(DISTINCT visitor_id) AS unique_visitors
		FROM analytics
		WHERE site_id = $1
		  AND created_at > NOW() - INTERVAL '%s'
		GROUP BY country_code
		ORDER BY unique_visitors DESC
		LIMIT 10
	`, interval), siteID)
	if err != nil {
		log.Printf("serveCountries error: %v", err)
		http.Error(w, "Failed to get countries", http.StatusInternalServerError)
		return
	}
	defer rows.Close()
	var countries []TopCountry
	for rows.Next() {
		var c TopCountry
		var cc sql.NullString
		if err := rows.Scan(&cc, &c.Views, &c.UniqueVisitors); err != nil {
			continue
		}
		if cc.Valid && cc.String != "" {
			c.CountryCode = cc.String
		} else {
			c.CountryCode = "Unknown"
		}
		countries = append(countries, c)
	}
	if countries == nil {
		countries = []TopCountry{}
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(countries)
}

func (app *App) serveReferrers(w http.ResponseWriter, r *http.Request, siteID string) {
	interval, _ := parseDateRange(r)
	rows, err := app.DB.Query(r.Context(), fmt.Sprintf(`
		SELECT
			COALESCE(NULLIF(referrer, ''), 'Direct / None') AS referrer,
			COUNT(*) AS views,
			COUNT(DISTINCT visitor_id) AS unique_visitors
		FROM analytics
		WHERE site_id = $1
		  AND created_at > NOW() - INTERVAL '%s'
		GROUP BY 1
		ORDER BY unique_visitors DESC
		LIMIT 10
	`, interval), siteID)
	if err != nil {
		log.Printf("serveReferrers error: %v", err)
		http.Error(w, "Failed to get referrers", http.StatusInternalServerError)
		return
	}
	defer rows.Close()
	var referrers []TopReferrer
	for rows.Next() {
		var ref TopReferrer
		if err := rows.Scan(&ref.Referrer, &ref.Views, &ref.UniqueVisitors); err != nil {
			continue
		}
		if ref.Referrer != "Direct / None" {
			ref.Referrer = strings.TrimPrefix(ref.Referrer, "https://")
			ref.Referrer = strings.TrimPrefix(ref.Referrer, "http://")
			ref.Referrer = strings.TrimSuffix(ref.Referrer, "/")
		}
		referrers = append(referrers, ref)
	}
	if referrers == nil {
		referrers = []TopReferrer{}
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(referrers)
}

// ── Authenticated handlers (verify site ownership, then delegate) ─────────────

func (app *App) checkSiteOwnership(w http.ResponseWriter, r *http.Request, siteID string) bool {
	userID := r.Context().Value(userKey).(User).ID
	var ownerID int64
	err := app.DB.QueryRow(r.Context(), `SELECT user_id FROM sites WHERE id = $1`, siteID).Scan(&ownerID)
	if err != nil || ownerID != userID {
		http.Error(w, "Not found", http.StatusNotFound)
		return false
	}
	return true
}

func (app *App) handleGetSiteStats(w http.ResponseWriter, r *http.Request) {
	siteID := r.PathValue("id")
	if !app.checkSiteOwnership(w, r, siteID) {
		return
	}
	app.serveStats(w, r, siteID)
}

func (app *App) handleGetSiteTraffic(w http.ResponseWriter, r *http.Request) {
	siteID := r.PathValue("id")
	if !app.checkSiteOwnership(w, r, siteID) {
		return
	}
	app.serveTraffic(w, r, siteID)
}

func (app *App) handleGetSitePages(w http.ResponseWriter, r *http.Request) {
	siteID := r.PathValue("id")
	if !app.checkSiteOwnership(w, r, siteID) {
		return
	}
	app.servePages(w, r, siteID)
}

func (app *App) handleGetSiteCountries(w http.ResponseWriter, r *http.Request) {
	siteID := r.PathValue("id")
	if !app.checkSiteOwnership(w, r, siteID) {
		return
	}
	app.serveCountries(w, r, siteID)
}

func (app *App) handleGetSiteReferrers(w http.ResponseWriter, r *http.Request) {
	siteID := r.PathValue("id")
	if !app.checkSiteOwnership(w, r, siteID) {
		return
	}
	app.serveReferrers(w, r, siteID)
}

// ── Public demo handlers (check DEMO_SITE_ID is set, then delegate) ───────────

func (app *App) demoGuard(w http.ResponseWriter) bool {
	if app.DemoSiteID == "" {
		http.Error(w, "Demo not configured", http.StatusServiceUnavailable)
		return false
	}
	return true
}

func (app *App) handleDemoStats(w http.ResponseWriter, r *http.Request) {
	if !app.demoGuard(w) {
		return
	}
	app.serveStats(w, r, app.DemoSiteID)
}

func (app *App) handleDemoTraffic(w http.ResponseWriter, r *http.Request) {
	if !app.demoGuard(w) {
		return
	}
	app.serveTraffic(w, r, app.DemoSiteID)
}

func (app *App) handleDemoPages(w http.ResponseWriter, r *http.Request) {
	if !app.demoGuard(w) {
		return
	}
	app.servePages(w, r, app.DemoSiteID)
}

func (app *App) handleDemoCountries(w http.ResponseWriter, r *http.Request) {
	if !app.demoGuard(w) {
		return
	}
	app.serveCountries(w, r, app.DemoSiteID)
}

func (app *App) handleDemoReferrers(w http.ResponseWriter, r *http.Request) {
	if !app.demoGuard(w) {
		return
	}
	app.serveReferrers(w, r, app.DemoSiteID)
}

func (app *App) handleCreateSite(w http.ResponseWriter, r *http.Request) {
	var payload CreateSitePayload
	if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
		http.Error(w, "Invalid JSON", http.StatusBadRequest)
		return
	}

	// Block users from registering localhost addresses, stopping local traffic theft
	payload.Domain = strings.ToLower(strings.TrimSpace(payload.Domain))
	if payload.Domain == "localhost" || payload.Domain == "127.0.0.1" || payload.Domain == "0.0.0.0" {
		http.Error(w, "Invalid domain", http.StatusBadRequest)
		return
	}

	userID := r.Context().Value(userKey).(User).ID

	query := `
		INSERT INTO sites (domain, user_id)
		VALUES($1, $2)
		RETURNING id, domain, created_at
	`
	var site Site
	err := app.DB.QueryRow(r.Context(), query, payload.Domain, userID).Scan(&site.ID, &site.Domain, &site.CreatedAt)
	if err != nil {
		// Check for unique constraint violation (domain already exists)
		if pgErr, ok := err.(*pgconn.PgError); ok && pgErr.Code == "23505" {
			http.Error(w, "This domain is already registered", http.StatusConflict)
			return
		}
		log.Printf("handleCreateSite error: %v", err)
		http.Error(w, "Failed to create site", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(site)
}

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
