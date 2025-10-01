package internal

import (
	"crypto/rand"
	"math/big"
	"net"
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