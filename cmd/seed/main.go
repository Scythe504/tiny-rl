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
	"strings"
	"time"

	"github.com/scythe504/tiny-rl/internal"
	"github.com/scythe504/tiny-rl/internal/database"
	"github.com/scythe504/tiny-rl/internal/geodatabase"
)

const BATCH_SIZE = 500 // Adjust based on your needs (500-1000 works well)

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
		"direct",
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

	startTime := time.Now()
	log.Println("==== SEEDING STARTED ====")
	log.Printf("📅 Date Range: %s to %s (1 year)\n", time.Now().AddDate(-1, 0, 0).Format("2006-01-02"), time.Now().Format("2006-01-02"))
	log.Printf("📦 Batch Size: %d records per insert\n", BATCH_SIZE)
	
	ctx := context.Background()
	tx, err := conn.BeginTx(ctx, &sql.TxOptions{Isolation: sql.LevelSerializable})
	if err != nil {
		log.Fatal("❌ Transaction Start Failed:", err)
	}
	log.Println("✅ Transaction Started")
	defer tx.Rollback()

	// Insert link_map
	link_map := database.LinkMap{
		ShortCode: "demo-code",
		Url:       "https://github.com/scythe504/tiny-rl",
	}

	stmt := `INSERT INTO link_map (short_code, url) VALUES ($1, $2) 
         ON CONFLICT (short_code) DO NOTHING`

	_, err = tx.Exec(stmt, link_map.ShortCode, link_map.Url)
	if err != nil {
		log.Fatal("❌ Failed to insert link_map:", err)
	}
	log.Println("✅ Link Map Inserted/Updated")

	// Initialize GeoDb
	geodb := geodatabase.New()
	log.Println("✅ GeoDb Connection Opened")

	// Time Range
	endDate := time.Now()
	startDate := endDate.AddDate(-1, 0, 0)

	// Batch processing
	var batch []database.Clicks
	totalRecords := 0
	batchCount := 0
	dayCount := 0

	log.Println("\n==== GENERATING & INSERTING DATA ====")
	
	// Time Range Loop
	for d := startDate; !d.After(endDate); d = d.AddDate(0, 0, 1) {
		dayCount++
		randNum := RandomInt(30, 50)

		// Generate clicks for this day
		for range int(randNum) {
			ipAddr := RandomChoice(IPs)
			geoIpCountry, err := geodb.GetCountryByIP(net.ParseIP(ipAddr))
			if err != nil {
				log.Printf("⚠️  GeoIP Parsing Error for %s: %v\n", ipAddr, err)
			}

			countryName := "India"
			countryIsoCode := "IN"

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

			batch = append(batch, click)

			// Execute batch insert when batch size is reached
			if len(batch) >= BATCH_SIZE {
				rowsInserted, err := executeBatchInsert(tx, batch)
				if err != nil {
					log.Fatal("❌ Batch Insert Failed:", err)
				}
				totalRecords += rowsInserted
				batchCount++
				log.Printf("📊 Batch #%d: Inserted %d records | Total: %d | Days Processed: %d\n", 
					batchCount, rowsInserted, totalRecords, dayCount)
				batch = batch[:0] // Clear batch
			}
		}

		// Log progress every 30 days
		if dayCount%30 == 0 {
			elapsed := time.Since(startTime)
			log.Printf("🎯 Milestone: %d days completed in %s | Records: %d\n", 
				dayCount, elapsed.Round(time.Second), totalRecords)
		}
	}

	// Insert remaining records
	if len(batch) > 0 {
		rowsInserted, err := executeBatchInsert(tx, batch)
		if err != nil {
			log.Fatal("❌ Final Batch Insert Failed:", err)
		}
		totalRecords += rowsInserted
		batchCount++
		log.Printf("📊 Final Batch #%d: Inserted %d records | Total: %d\n", 
			batchCount, rowsInserted, totalRecords)
	}

	log.Println("\n==== COMMITTING TRANSACTION ====")
	if err = tx.Commit(); err != nil {
		_ = tx.Rollback()
		log.Fatal("❌ Commit Failed:", err)
	}

	endTime := time.Since(startTime)
	log.Println("✅ Transaction Committed Successfully")
	log.Println("\n==== SEEDING SUMMARY ====")
	log.Printf("📈 Total Records Inserted: %d\n", totalRecords)
	log.Printf("📦 Total Batches: %d\n", batchCount)
	log.Printf("📅 Days Processed: %d\n", dayCount)
	log.Printf("⏱️  Total Time: %s\n", endTime.Round(time.Second))
	log.Printf("⚡ Average Speed: %.0f records/second\n", float64(totalRecords)/endTime.Seconds())
	
	conn.Close()
	geodb.Close()
	log.Println("==== SEEDING COMPLETED ====")
}

// executeBatchInsert performs a batch insert using a single multi-value INSERT statement
func executeBatchInsert(tx *sql.Tx, clicks []database.Clicks) (int, error) {
	if len(clicks) == 0 {
		return 0, nil
	}

	// Build the query with placeholders
	valueStrings := make([]string, 0, len(clicks))
	valueArgs := make([]interface{}, 0, len(clicks)*8)

	for i, click := range clicks {
		valueStrings = append(valueStrings, fmt.Sprintf("($%d, $%d, $%d, $%d, $%d, $%d, $%d, $%d)",
			i*8+1, i*8+2, i*8+3, i*8+4, i*8+5, i*8+6, i*8+7, i*8+8))
		
		valueArgs = append(valueArgs,
			click.ShortCode,
			click.Browser,
			click.ClickedAt,
			click.UserAgent,
			click.IpAddr,
			click.Referrer,
			click.Country,
			click.CountryISOCode,
		)
	}

	stmt := fmt.Sprintf(`INSERT INTO clicks (
		short_code,
		browser, 
		clicked_at, 
		user_agent, 
		ip_addr, 
		referrer, 
		country, 
		country_iso_code
	) VALUES %s`, strings.Join(valueStrings, ","))

	result, err := tx.Exec(stmt, valueArgs...)
	if err != nil {
		return 0, err
	}

	rowsAffected, _ := result.RowsAffected()
	return int(rowsAffected), nil
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