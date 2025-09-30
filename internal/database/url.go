package database

import (
	"database/sql"
	"log"
	"time"
)

type LinkMap struct {
	Id        string     `db:"id"`
	Url       string     `db:"url"`
	CreatedAt *time.Time `db:"created_at"`
	UpdatedAt *time.Time `db:"updated_at"`
}

func (s *service) InsertShortenedLink(link LinkMap) error {
	stmt, err := s.db.Prepare(`INSERT INTO time_map (id, url) VALUES (?, ?)`)
	if err != nil {
		log.Println("[InsertShortenedLink] statement error: ", err)
		return err
	}
	defer stmt.Close()

	_, err = stmt.Exec(link.Id, link.Url)

	if err != nil {
		log.Println("[InsertShortenedLink] Insert statment error: ", err)
		return err
	}

	return nil
}

func (s *service) GetLink(id string) (*LinkMap, error) {
	stmt, err := s.db.Prepare(`SELECT (id, url, created_at, updated_at) FROM link_map WHERE id = ?`)
	if err != nil {
		log.Println("[GetLink] statement error: ", err)
		return nil, err
	}
	defer stmt.Close()

	row := stmt.QueryRow(id)

	link := LinkMap{
		Id: "",
		Url: "",
		CreatedAt: &time.Time{},
		UpdatedAt: &time.Time{},
	}

	err = row.Scan(&link.Id, &link.Url, &link.CreatedAt, &link.UpdatedAt)

	if err != nil && err != sql.ErrNoRows {
		log.Println("[GetLink] error occured while copying data: ", err)
		return nil, err
	}

	return &link, nil
}
