package database

import (
	"log"
	"time"

	"github.com/jackc/pgx/v5"
)

type Clicks struct {
	Id        string    `db:"id" json:"id"`
	ShortCode string    `db:"short_code" json:"short_code"`
	ClickedAt time.Time `db:"clicked_at" json:"clicked_at"`
	UserAgent string    `db:"user_agent" json:"user_agent"`
	IpAddr    string    `db:"ip_addr" json:"ip_addr"`
	Referrer  string    `db:"referrer" json:"referrer"`
}

type ClicksOnDay struct {
	Day        time.Time `db:"day" json:"day"`
	ClickCount int       `db:"click_count" json:"click_count"`
}

func (s *service) LogClick(click Clicks) error {
	stmt := `INSERT INTO clicks (
			short_code, 
			ip_addr, 
			user_agent, 
			referrer, 
			clicked_at
		) VALUES (
			$1,
			$2, 
			$3, 
			$4, 
			$5
		)`

	_, err := s.db.Exec(stmt, click.ShortCode, click.IpAddr, click.UserAgent, click.Referrer, click.ClickedAt)
	if err != nil {
		log.Println("[LogClick] Error occured when Executing statement: ", err)
		return err
	}

	return nil
}

func (s *service) GetClicksOverTime(shortCode string) ([]ClicksOnDay, error) {
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

	var ClicksOnDays []ClicksOnDay = make([]ClicksOnDay, 0)

	for rows.Next() {
		var ClicksOnDay ClicksOnDay
		if err := rows.Scan(&ClicksOnDay.Day, &ClicksOnDay.ClickCount); err != nil {
			log.Println("[GetClicksOverTime] error occured while scanning to variable", err)
			return nil, err
		}

		ClicksOnDays = append(ClicksOnDays, ClicksOnDay)
	}

	return ClicksOnDays, nil
}
