package main

import (
	"context"
	"crypto/rand"
	"database/sql"
	"fmt"
	"log"
	"math/big"
	"net"
	"os"
	"time"

	"github.com/scythe504/tiny-rl/internal"
	"github.com/scythe504/tiny-rl/internal/database"
	"github.com/scythe504/tiny-rl/internal/geodatabase"
)

var (
	db        = os.Getenv("DB_DATABASE")
	password  = os.Getenv("DB_PASSWORD")
	username  = os.Getenv("DB_USERNAME")
	port      = os.Getenv("DB_PORT")
	host      = os.Getenv("DB_HOST")
	schema    = os.Getenv("DB_SCHEMA")
	conn_str  = os.Getenv("POSTGRES_CONN_URL")
	HASH_SALT = os.Getenv("HASH_SALT")
	Browsers  = []string{
		"Chrome",
		"Firefox",
		"Safari",
		"Edge",
		"Opera",
		"Other",
	}

	UserAgents = []string{
		"Mozilla/5.0 (Windows NT 10.0; Win64; x64)",
		"Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7)",
		"Mozilla/5.0 (X11; Linux x86_64)",
		"Mozilla/5.0 (iPhone; CPU iPhone OS 16_0 like Mac OS X)",
		"Mozilla/5.0 (iPad; CPU OS 15_6 like Mac OS X)",
	}

	IPs = []string{
		"14.36.183.3",
		"87.114.231.235",
		"167.188.71.253",
		"97.98.109.35",
		"150.77.224.161",
		"222.38.154.82",
		"2.192.52.149",
		"94.82.109.5",
		"83.78.142.31",
		"18.140.192.63",
		"182.55.155.239",
		"152.173.194.120",
		"151.55.119.115",
		"94.136.22.162",
		"227.239.220.84",
		"118.77.77.150",
		"89.151.49.152",
		"223.225.3.224",
		"216.28.206.188",
		"192.227.165.142",
		"59.201.216.70",
		"50.91.187.146",
		"96.99.144.211",
		"223.101.109.26",
		"164.218.136.55",
		"100.197.186.62",
		"193.38.170.2",
		"170.202.114.30",
		"201.139.219.68",
		"134.204.215.203",
	}

	Referrers = []string{
		"https://google.com",
		"https://bing.com",
		"https://yahoo.com",
		"https://duckduckgo.com",
		"https://reddit.com",
		"https://twitter.com",
		"https://facebook.com",
		"https://linkedin.com",
		"https://stackoverflow.com",
		"direct", // for no referrer / direct traffic
	}
)

func main() {
	connStr := fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=disable&search_path=%s", username, password, host, port, db, schema)

	if conn_str != "" {
		connStr = conn_str
	}

	conn, err := sql.Open("pgx", connStr)
	if err != nil {
		log.Fatal(err)
	}

	ctx := context.Background()
	tx, err := conn.BeginTx(ctx, &sql.TxOptions{Isolation: sql.LevelSerializable})
	if err != nil {
		log.Fatal("TxStartFailed", err)
	}
	defer tx.Rollback()

	// Insert Database With Demo short_code
	link_map := database.LinkMap{
		ShortCode: "demo-code",
		Url:       "https://github.com/scythe504/tiny-rl",
	}

	stmt := `INSERT INTO link_map (short_code, url) VALUES ($1, $2) 
         ON CONFLICT (short_code) DO NOTHING` // prevents exit if already exists

	_, err = tx.Exec(stmt, link_map.ShortCode, link_map.Url)
	if err != nil {
		log.Println("[SeedDB] Failed to execute query", err)
		return
	}

	// SEED 6 months worth of data range (30-50 clicks in a day)
	// Insert statement
	geodb := geodatabase.New()
	singleClickStmt, err := tx.Prepare(`INSERT INTO clicks (
		short_code,
		browser, 
		clicked_at, 
		user_agent, 
		ip_addr, 
		referrer, 
		country, 
		country_iso_code
		) VALUES (
			$1, 
			$2,
			$3,
			$4,
			$5,
			$6, 
			$7,
			$8
		)`)
	if err != nil {
		log.Fatal("[ClicksStatementPreparation]", err)
	}
	defer singleClickStmt.Close()

	// Time Range
	endDate := time.Now()
	startDate := endDate.AddDate(-1, 0, 0)

	// Time Range Loop
	for d := startDate; !d.After(endDate); d = d.AddDate(0, 0, 1) {

		// Single Day Clicks Insertion (30-50)
		randNum := RandomInt(30, 50)
		for range int(randNum) {
			ipAddr := RandomChoice(IPs)

			geoIpCountry, err := geodb.GetCountryByIP(net.ParseIP(ipAddr))
			if err != nil {
				log.Println("[GoroutineLogClick] error parsing ipaddr", ipAddr, err)
				return
			}

			countryName := "India" // default fallback for local/dev
			countryIsoCode := "IN" // default fallback

			if geoIpCountry != nil && geoIpCountry.Country.IsoCode != "" {
				name := geoIpCountry.Country.Names["en"]
				if name != "" {
					countryName = name
				}
				countryIsoCode = geoIpCountry.Country.IsoCode
			}

			dayStart := time.Date(d.Year(), d.Month(), d.Day(), 0, 0, 0, 0, d.Location())

			secondsInDay := int64(24 * 60 * 60)
			secOffset := RandomInt(0, secondsInDay)
			timestamp := dayStart.Add(time.Duration(secOffset) * time.Second)
			hashedIp := internal.HashIPWithDate(ipAddr, HASH_SALT, timestamp)

			click := database.Clicks{
				ShortCode:      "demo-code",
				Browser:        RandomChoice(Browsers),
				UserAgent:      RandomChoice(UserAgents),
				IpAddr:         hashedIp,
				Referrer:       RandomChoice(Referrers),
				Country:        countryName,
				CountryISOCode: countryIsoCode,
				ClickedAt:      timestamp,
			}

			_, err = singleClickStmt.Exec(
				click.ShortCode,
				click.Browser,
				click.ClickedAt,
				click.UserAgent,
				click.IpAddr,
				click.Referrer,
				click.Country,
				click.CountryISOCode,
			)

			if err != nil {
				log.Fatal("Failed to insert data", err, d, randNum, timestamp, click)
			}
		}
	}

	if err = tx.Commit(); err != nil {
		_ = tx.Rollback()
		log.Fatal(err)
	}
	conn.Close()
	geodb.Close()
}

// RandomChoice picks a random string from a slice
func RandomChoice(choices []string) string {
	idx := RandomInt(0, int64(len(choices)))
	return choices[idx]
}

// RandomInt returns a cryptographically secure random integer in [min, max)
func RandomInt(min, max int64) int64 {
	nBig, err := rand.Int(rand.Reader, big.NewInt(max-min))
	if err != nil {
		panic(err)
	}
	return nBig.Int64() + min
}
