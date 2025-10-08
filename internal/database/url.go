package database

import (
	"log"
	"time"

	"github.com/jackc/pgx/v5"
)

type LinkMap struct {
	ShortCode string     `db:"short_code" json:"short_code"`
	Url       string     `db:"url" json:"url"`
	CreatedAt time.Time `db:"created_at" json:"created_at,omitempty"`
	UpdatedAt time.Time `db:"updated_at" json:"updated_at,omitempty"`
}

func (s *service) InsertShortenedLink(link LinkMap) error {
	stmt := `INSERT INTO link_map (short_code, url) VALUES ($1, $2)`

	_, err := s.db.Exec(stmt, link.ShortCode, link.Url)

	if err != nil {
		log.Println("[InsertShortenedLink] Insert statment error: ", err)
		return err
	}

	return nil
}

func (s *service) GetLink(short_code string) (*LinkMap, error) {
	stmt := `SELECT 
	 short_code,
	 url, 
	 created_at, 
	 updated_at 
	 FROM link_map 
	 WHERE short_code = $1`

	row := s.db.QueryRow(stmt, short_code)

	var link LinkMap

	err := row.Scan(&link.ShortCode, &link.Url, &link.CreatedAt, &link.UpdatedAt)

	if err != nil && err != pgx.ErrNoRows {
		log.Println("[GetLink] error occured while copying data: ", err)
		return nil, err
	}

	return &link, nil
}

func (s *service) UpdateShortenedLink(shortCode string, destUrl string) error {
	stmt := `UPDATE link_map SET url=$1 WHERE short_code=$2`

	_, err := s.db.Exec(stmt, destUrl, shortCode)

	if err != nil {
		log.Println("[UpdateShortenedLink] Update statment error: ", err)
		return err
	}

	return nil
}
