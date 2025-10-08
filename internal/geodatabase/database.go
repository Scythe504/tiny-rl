package geodatabase

import (
	"log"
	"net"

	"github.com/oschwald/geoip2-golang"
)

type Service interface {
	GetCountryByIP(ipAddr net.IP) (*geoip2.Country, error)
	Close() error
}

type service struct {
	db *geoip2.Reader
}

var (
	geodbPath  = "./data/GeoLite2-Country.mmdb"
	dbInstance *service
)

func New() Service {
	if dbInstance != nil {
		return dbInstance
	}
	geodbInstance, err := geoip2.Open(geodbPath)
	if err != nil {
		log.Fatal(err)
	}

	dbInstance = &service{
		db: geodbInstance,
	}

	return dbInstance
}

func (s *service) Close() error {
	log.Println("Disconnected from geodb: ", dbInstance.db)
	
	if err := s.db.Close(); err != nil {
		log.Fatal(err)
	}
	return nil
}
