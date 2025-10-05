package geodatabase

import (
	"log"
	"net"

	"github.com/oschwald/geoip2-golang"
)

func (s *service) GetCountryByIP(ipAddr net.IP) (*geoip2.Country, error) {
	record, err := s.db.Country(ipAddr)

	if err != nil {
		log.Println("[GetCountryByIP] Failed to get country for ip", ipAddr, err)
		return nil, err
	}

	return record, nil
}
