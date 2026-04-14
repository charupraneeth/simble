package main

import "time"

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
	Email     *string   `json:"email"` // can be null if user made it private
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

type TopEvent struct {
	Name           string `json:"name"`
	Events         int64  `json:"events"`
	UniqueVisitors int64  `json:"unique_visitors"`
}
