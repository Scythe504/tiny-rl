package database

import (
	"log"
	"time"

	"github.com/jackc/pgx/v5"
)

type Clicks struct {
	Id        string    `db:"id" json:"id"`
	ShortCode string    `db:"short_code" json:"short_code"`
	Browser   string    `db:"browser" json:"browser"`
	ClickedAt time.Time `db:"clicked_at" json:"clicked_at"`
	UserAgent string    `db:"user_agent" json:"user_agent"`
	IpAddr    string    `db:"ip_addr" json:"ip_addr"`
	Referrer  string    `db:"referrer" json:"referrer"`
}

type ClicksPerDay struct {
	Day        time.Time `db:"day" json:"day"`
	ClickCount int       `db:"click_count" json:"click_count"`
}

type ClicksPerBrowser struct {
	Browser    string `db:"browser" json:"browser"`
	ClickCount int    `db:"click_count" json:"click_count"`
}

type TrafficFromReferrer struct {
	Referrer   string `db:"referrer" json:"referrer"`
	ClickCount int    `db:"click_count" json:"click_count"`
}

func (s *service) LogClick(click Clicks) error {
	stmt := `INSERT INTO clicks (
			short_code, 
			ip_addr, 
			user_agent,
			browser, 
			referrer, 
			clicked_at
		) VALUES (
			$1,
			$2, 
			$3, 
			$4, 
			$5,
			$6
		)`

	_, err := s.db.Exec(stmt, click.ShortCode, click.IpAddr, click.UserAgent, click.Browser, click.Referrer, click.ClickedAt)
	if err != nil {
		log.Println("[LogClick] Error occured when Executing statement: ", err)
		return err
	}

	return nil
}

func (s *service) GetClicksOverTime(shortCode string) ([]ClicksPerDay, error) {
	stmt := `SELECT 
		DATE_TRUNC('day', clicked_at) AS day, 
		COUNT(*) AS click_count 
		FROM clicks 
		WHERE short_code = $1
		GROUP BY day 
		ORDER BY day;`

	rows, err := s.db.Query(stmt, shortCode)
	if err != nil && err != pgx.ErrNoRows {
		log.Println("[GetClicksOverTime] error occured while querying", rows)
		return nil, err
	}
	defer rows.Close()

	var clicksPerDays []ClicksPerDay = make([]ClicksPerDay, 0)

	for rows.Next() {
		var clicksPerDay ClicksPerDay
		if err := rows.Scan(&clicksPerDay.Day, &clicksPerDay.ClickCount); err != nil {
			log.Println("[GetClicksOverTime] error occured while scanning to variable", err)
			return nil, err
		}

		clicksPerDays = append(clicksPerDays, clicksPerDay)
	}

	return clicksPerDays, nil
}

func (s *service) GetBrowserStats(shortCode string) ([]ClicksPerBrowser, error) {
	stmt := `SELECT browser, COUNT(*) AS click_count
					 FROM clicks
					 WHERE short_code=$1
					 GROUP BY browser
					 ORDER BY click_count DESC;`
	rows, err := s.db.Query(stmt, shortCode)
	if err != nil && err != pgx.ErrNoRows {
		log.Println("[GetBrowserStats] error occured while querying", rows)
		return nil, err
	}
	defer rows.Close()

	var clicksPerBrowsers []ClicksPerBrowser = make([]ClicksPerBrowser, 0)

	for rows.Next() {
		var clicksPerBrowser ClicksPerBrowser
		if err := rows.Scan(&clicksPerBrowser.Browser, &clicksPerBrowser.ClickCount); err != nil {
			log.Println("[GetBrowserStats] error occured while scanning to variable", err)
			return nil, err
		}

		clicksPerBrowsers = append(clicksPerBrowsers, clicksPerBrowser)
	}

	return clicksPerBrowsers, nil
}

func (s *service) GetReferrerStats(shortCode string) ([]TrafficFromReferrer, error) {
	stmt := `SELECT referrer, COUNT(*) AS click_count
						FROM clicks
						WHERE short_code=$1
						GROUP BY referrer
						ORDER BY click_count DESC;`
	rows, err := s.db.Query(stmt, shortCode)
	if err != nil && err != pgx.ErrNoRows {
		log.Println("[GetReferrerStats] error occured while querying", rows)
		return nil, err
	}
	defer rows.Close()

	var trafficFromReferrers []TrafficFromReferrer = make([]TrafficFromReferrer, 0)

	for rows.Next() {
		var trafficFromReferrer TrafficFromReferrer
		if err := rows.Scan(&trafficFromReferrer.Referrer, &trafficFromReferrer.ClickCount); err != nil {
			log.Println("[GetReferrerStats] error occured while scanning to variable", err)
			return nil, err
		}

		trafficFromReferrers = append(trafficFromReferrers, trafficFromReferrer)
	}

	return trafficFromReferrers, nil
}
