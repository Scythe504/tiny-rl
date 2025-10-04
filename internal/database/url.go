package database

import (
	"log"
	"time"

	"github.com/jackc/pgx/v5"
)

type LinkMap struct {
	ShortCode string     `db:"short_code"`
	Url       string     `db:"url"`
	CreatedAt *time.Time `db:"created_at"`
	UpdatedAt *time.Time `db:"updated_at"`
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

	link := LinkMap{
		ShortCode: "",
		Url:       "",
		CreatedAt: &time.Time{},
		UpdatedAt: &time.Time{},
	}

	err := row.Scan(&link.ShortCode, &link.Url, &link.CreatedAt, &link.UpdatedAt)

	if err != nil && err != pgx.ErrNoRows {
		log.Println("[GetLink] error occured while copying data: ", err)
		return nil, err
	}

	return &link, nil
}
