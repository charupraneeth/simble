package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net"
	"net/http"
	"net/netip"
	"os"

	"github.com/medama-io/go-useragent"
	"github.com/oschwald/geoip2-golang/v2"
)

type ResponseData struct {
	Device  string `json:"device"`
	Browser string `json:"browser"`
	Os      string `json:"os"`
	Referer string `json:"referer"`
	Host    string `json:"host"`
	Country string `json:"country"`
	City    string `json:"city"`
}

type App struct {
	GeoDB    *geoip2.Reader
	UAParser *useragent.Parser
}

func (app *App) handleRequest(w http.ResponseWriter, r *http.Request) {

	uaString := r.Header.Get("User-Agent")
	agent := app.UAParser.Parse(uaString)

	browser := agent.Browser()
	device := agent.Device()
	os := agent.OS()
	referer := r.Referer()

	host, _, err := net.SplitHostPort(r.RemoteAddr)

	if err != nil {
		fmt.Println("error splitting host/port: ", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)

		return
	}

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

	country := ""
	city := ""
	if record.HasData() {
		country = record.Country.Names.English
		city = record.City.Names.English
	} else {
		fmt.Println("No data found for this IP")
	}

	response := ResponseData{
		Device:  device.String(),
		Browser: browser.String(),
		Os:      os.String(),
		Referer: referer,
		Host:    host,
		Country: country,
		City:    city,
	}

	fmt.Println(response)

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(response); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func main() {
	mux := http.NewServeMux()

	db, err := geoip2.Open("./GeoLite2-City/GeoLite2-City.mmdb")
	if err != nil {
		log.Fatal(err)
	}

	defer db.Close()

	ua := useragent.NewParser()

	app := &App{GeoDB: db, UAParser: ua}

	mux.HandleFunc("/", app.handleRequest)

	port := os.Getenv("PORT")

	if port == "" {
		port = "8080"
	}

	log.Fatal(http.ListenAndServe(":"+port, mux))
}
