package main

import (
	"encoding/json"
	"log"
	"net/http"
	"net/netip"
	"net/url"
	"os"
	"strings"
	"sync"

	"simble/internal/utils"

	"github.com/jackc/pgx/v5"
	"github.com/oschwald/geoip2-golang/v2"
)

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

	host := utils.GetRealIP(r)

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

	visitorID := utils.GetDailyVisitorID(host, uaString, salt)

	insertQuery := `INSERT INTO analytics(
 					site_id, visitor_id, path, browser_name, device_type, os_name, country_code, city_name, referrer, name
					)
					VALUES($1, $2, $3, $4, $5, $6, $7, $8, $9, $10);`

	if payload.Name == "" {
		payload.Name = "pageview"
	}
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

	_, err = app.DB.Exec(r.Context(), insertQuery, event.SiteID, event.VisitorID, event.Path, event.Browser, event.Device, event.Os, event.CountryCode, event.City, event.Referer, event.Name)

	if err != nil {
		log.Printf("CRITICAL: Failed to insert analytics event for site %d: %v", siteID, err)
		w.WriteHeader(http.StatusAccepted)
		return
	}

	log.Printf("Event inserted: %v\n", event)

	w.WriteHeader(http.StatusAccepted)
}
