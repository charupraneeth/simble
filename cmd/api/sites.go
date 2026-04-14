package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strings"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
)

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

// ── Sites CRUD ────────────────────────────────────────────────────────────────

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
