package internal

import (
	"crypto/rand"
	"math/big"
	"net"
	"net/http"
	"net/url"
	"strings"
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
