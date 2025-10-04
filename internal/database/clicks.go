package database

import (
	"log"
	"time"
)

type Clicks struct {
	Id        string    `db:"id" json:"id"`
	ShortCode string    `db:"short_code" json:"short_code"`
	ClickedAt time.Time `db:"clicked_at" json:"clicked_at"`
	UserAgent string    `db:"user_agent" json:"user_agent"`
	IpAddr    string    `db:"ip_addr" json:"ip_addr"`
	Referrer  string    `db:"referrer" json:"referrer"`
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
