package utils

import (
	"log"
	"net"
)

var (
	Trace   *log.Logger
	Info    *log.Logger
	Warning *log.Logger
	Error   *log.Logger
)

// ResolveFQDN resolves the input domain name to any IP addresses
func ResolveFQDN(domain *string) (ipList []net.IP, err error) {
	ipList, err = net.LookupIP(*domain)
	if err != nil {
		Error.Printf("Unable to look up domain: %s. Err Code: %v\n", *domain, err)
		return nil, err
	}
	return
}

func LookUpIPClientIP() () {

}

// CheckDomainExists checks the DNS registrar for the zone
func CheckDomainExists(baseDomain, plugin string) (zoneID string, err error) {

	return "", nil
}
