package main

import (
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
	"time"

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
	reqOrigin := r.Header.Get("Origin")

	originURL, err := url.Parse(reqOrigin)
	if err != nil {
		http.Error(w, "Bad request", http.StatusBadRequest)
		return
	}

	var payload RequestPayload
	if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
		http.Error(w, "Invalid JSON", http.StatusBadRequest)
		return
	}

	if originURL.Host != payload.Domain {
		http.Error(w, "Invalid domain", http.StatusForbidden)
		return
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

	record, err := app.GeoDB.City(ip)
	if err != nil {
		fmt.Println("error getting city from ip", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)

		return
	}

	country_code := ""
	city := ""
	if record.HasData() {
		country_code = record.Country.ISOCode
		city = record.City.Names.English
	} else {
		fmt.Println("No data found for this IP")
	}

	salt := os.Getenv("VISITOR_ID_SALT")
	if salt == "" {
		salt = "fallback-temp-salt" // for local dev
	}

	visitorID := getDailyVisitorID(host, uaString, salt)

	response := Event{
		Device:      agent.Device().String(),
		Browser:     agent.Browser().String(),
		Os:          agent.OS().String(),
		Referer:     r.Referer(),
		VisitorID:   visitorID,
		CountryCode: country_code,
		City:        city,
		Name:        payload.Name,
	}

	log.Printf("Response: %v\n", response)

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(response); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func main() {
	mux := http.NewServeMux()

	geoDB, err := geoip2.Open("./GeoLite2-City/GeoLite2-City.mmdb")
	if err != nil {
		log.Fatal(err)
	}

	defer geoDB.Close()

	ua := useragent.NewParser()

	app := &App{GeoDB: geoDB, UAParser: ua}

	mux.HandleFunc("POST /api/event", app.handleRequest)

	port := os.Getenv("PORT")

	if port == "" {
		port = "8080"
	}

	log.Fatal(http.ListenAndServe(":"+port, mux))
}
