package utils

import (
	"crypto/rand"
	"crypto/sha256"
	"encoding/base64"
	"fmt"
	"log"
	"net"
	"net/http"
	"strings"
	"time"
)

func GetRealIP(r *http.Request) string {
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

func GetDailyVisitorID(host, ua, salt string) string {
	date := time.Now().Format("2006-01-02")

	data := fmt.Sprintf("%s|%s|%s|%s", host, ua, salt, date)

	hash := sha256.Sum256([]byte(data))

	return fmt.Sprintf("%x", hash)[:16]
}

func GenerateRandomToken(n int) (string, error) {
	b := make([]byte, n)
	_, err := rand.Read(b)

	if err != nil {
		return "", err
	}

	return base64.URLEncoding.EncodeToString(b), nil
}
