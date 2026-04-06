package main

import (
	"context"
	crand "crypto/rand"
	"encoding/hex"
	"fmt"
	"log"
	"math"
	"math/rand/v2"
	"os"
	"strconv"
	"time"

	"github.com/jackc/pgx/v5"
)

var paths = []string{"/", "/", "/", "/about", "/pricing", "/blog", "/blog/hello-world", "/contact"}
var browsers = []string{"Chrome", "Chrome", "Safari", "Safari", "Firefox", "Edge", "Other"}
var devices = []string{"Desktop", "Desktop", "Mobile", "Mobile", "Tablet"}
var osNames = []string{"Windows", "Mac OS", "Mac OS", "iOS", "Android", "Android", "Linux"}
var countries = []string{"US", "US", "US", "GB", "GB", "CA", "IN", "IN", "DE", "FR", "BR", "JP"}
var referrers = []string{"https://google.com", "https://google.com", "https://twitter.com", "https://reddit.com", "https://news.ycombinator.com", "", "", "", "https://github.com"}

func randomVisitorID() string {
	b := make([]byte, 8)
	crand.Read(b)
	return hex.EncodeToString(b)
}

func main() {
	siteID, err := strconv.ParseInt(os.Args[1], 10, 64)
	if err != nil {
		log.Fatalf("Invalid site_id: %v", err)
	}

	dbURL := os.Getenv("DATABASE_URL")
	if dbURL == "" {
		dbURL = "postgres://postgres:password@localhost:5432/simble?sslmode=disable"
		fmt.Println("Warning: DATABASE_URL not set, defaulting to local postgres URL.")
	}

	ctx := context.Background()
	conn, err := pgx.Connect(ctx, dbURL)
	if err != nil {
		log.Fatalf("Unable to connect to database: %v\n", err)
	}
	defer conn.Close(ctx)

	// verify site exists
	var domain string
	err = conn.QueryRow(ctx, "SELECT domain FROM sites WHERE id = $1", siteID).Scan(&domain)
	if err != nil {
		log.Fatalf("Site ID %d not found in database.", siteID)
	}
	fmt.Printf("Seeding data for site %s (ID: %d)\n", domain, siteID)

	// Clean existing data for this site
	_, err = conn.Exec(ctx, "DELETE FROM analytics WHERE site_id = $1", siteID)
	if err != nil {
		log.Fatalf("Failed to clear old data: %v", err)
	}
	fmt.Println("Cleared existing data for site.")

	now := time.Now()
	daysToGenerate := 30
	eventsToGenerate := 15000 // Total fake events

	fmt.Printf("Generating %d events over %d days...\n", eventsToGenerate, daysToGenerate)

	var inserted int
	batchSize := 1000

	batch := &pgx.Batch{}

	// We'll generate a realistic-ish wave pattern
	// Visitor ids will be reused some of the time to simulate return visitors
	var visitorIDs []string
	for i := 0; i < 2000; i++ {
		visitorIDs = append(visitorIDs, randomVisitorID())
	}

	for i := 0; i < eventsToGenerate; i++ {
		// Time logic: mostly recent, with less traffic in the past, plus some daily wave structure (sine wave)
		daysAgo := float64(daysToGenerate) * rand.Float64()
		timeOffsetHours := daysAgo * 24.0

		// Adds a daily peak at noon
		wave := math.Sin((timeOffsetHours/24.0 - 0.5) * math.Pi * 2) // -1 to 1
		// bias traffic towards the "wave peak" by adjusting probability (quick and dirty hack)
		if rand.Float64() > (wave+1.0)/3.0 {
			continue // skip this random point, try again
		}

		eventTime := now.Add(-time.Duration(timeOffsetHours*float64(time.Hour)))

		// 80% new visitors, 20% returning
		vid := randomVisitorID()
		if rand.Float64() < 0.2 {
			vid = visitorIDs[rand.IntN(len(visitorIDs))]
		}

		path := paths[rand.IntN(len(paths))]
		browser := browsers[rand.IntN(len(browsers))]
		device := devices[rand.IntN(len(devices))]
		osName := osNames[rand.IntN(len(osNames))]
		country := countries[rand.IntN(len(countries))]
		referrer := referrers[rand.IntN(len(referrers))]

		batch.Queue(
			"INSERT INTO analytics (site_id, visitor_id, path, browser_name, device_type, os_name, country_code, city_name, referrer, created_at) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10)",
			siteID, vid, path, browser, device, osName, country, "City", referrer, eventTime,
		)

		inserted++

		if inserted%batchSize == 0 {
			br := conn.SendBatch(ctx, batch)
			_, err = br.Exec()
			if err != nil {
				log.Fatalf("Batch insert failed: %v", err)
			}
			br.Close()
			batch = &pgx.Batch{}
			fmt.Printf("Inserted %d events...\n", inserted)
		}
	}

	// flush remaining
	if batch.Len() > 0 {
		br := conn.SendBatch(ctx, batch)
		_, err = br.Exec()
		if err != nil {
			log.Fatalf("Batch insert failed: %v", err)
		}
		br.Close()
	}

	fmt.Printf("Successfully inserted %d events.\n", inserted)
}
