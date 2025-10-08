package internal

import (
	"crypto/rand"
	"crypto/sha256"
	"encoding/hex"
	"math/big"
	"net"
	"net/http"
	"net/url"
	"strings"
	"time"
)

func ShortCode() string {
	charset := "abcdefghijklmnopqrstuvwxyz0123456789ABCDEFGHIJKLMNOPQRSTUVWXYZ_-"

	var shortCode = ""
	charsetLen := int64(len(charset))
	// 6 length key (64^6)
	for range 6 {
		rd, _ := rand.Int(rand.Reader, big.NewInt(charsetLen))
		shortCode += string(charset[int(rd.Int64())])
	}

	return shortCode
}

func ValidURL(rawUrl string) bool {
	uri, err := url.Parse(rawUrl)
	if err != nil {
		return false
	}

	if uri.Scheme != "http" && uri.Scheme != "https" {
		return false
	}

	port := uri.Port()

	if port != "" &&
		!(uri.Scheme != "http" && port != "80") &&
		!(uri.Scheme != "https" && port != "443") {
		return false
	}

	if uri.Hostname() == "" {
		return false
	}

	if net.ParseIP(uri.Hostname()) != nil {
		return false
	}

	if strings.HasSuffix(uri.Hostname(), ".local") ||
		uri.Hostname() == "localhost" {
		return false
	}

	return true
}

func GetClientIP(r *http.Request) string {
	// Check common proxy headers (in order of trust)
	headers := []string{
		"X-Forwarded-For",
		"X-Real-IP",
		"CF-Connecting-IP", // Cloudflare
	}

	for _, h := range headers {
		if ip := r.Header.Get(h); ip != "" {
			// Sometimes multiple IPs in X-Forwarded-For: "client, proxy1, proxy2"
			parts := strings.Split(ip, ",")
			return strings.TrimSpace(parts[0])
		}
	}

	// Fallback to remote address
	ip, _, err := net.SplitHostPort(r.RemoteAddr)
	if err != nil {
		return r.RemoteAddr
	}

	// Handle local testing (loopback)
	if ip == "127.0.0.1" || ip == "::1" {
		// Option 1: fallback to local network IP
		localIP := getLocalIP()
		return localIP
	}

	return ip
}
func getLocalIP() string {
	addrs, err := net.InterfaceAddrs()
	if err != nil {
		return "127.0.0.1"
	}

	for _, addr := range addrs {
		if ipnet, ok := addr.(*net.IPNet); ok && !ipnet.IP.IsLoopback() && ipnet.IP.To4() != nil {
			return ipnet.IP.String()
		}
	}

	return "127.0.0.1"
}

// HashIPWithDate hashes an IP address together with a secret salt and a date.
// This anonymizes the IP while still allowing per-day uniqueness.
func HashIPWithDate(ip string, salt string, t time.Time) string {
	if ip == "" || salt == "" {
		return ""
	}

	date := t.Format("2006-01-02") 

	// Combine IP + date + salt and hash
	data := ip + date + salt
	hash := sha256.Sum256([]byte(data))

	return hex.EncodeToString(hash[:])
}
