package main

import (
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/medama-io/go-useragent"
	"github.com/oschwald/geoip2-golang/v2"
	"golang.org/x/oauth2"
)

type App struct {
	GeoDB       *geoip2.Reader
	UAParser    *useragent.Parser
	DB          *pgxpool.Pool
	OAuthConfig *oauth2.Config
	DemoSiteID  string
}
