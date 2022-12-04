package utils

import (
	"io"
	"log"
	"net"
	"net/http"
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

func LookUpIPClientIP(ipv4Endpoint, ipv6Endpoint string) (IPAddrs []string, err []error) {
	ipv4, errV4 := http.Get(ipv4Endpoint)
	if errV4 != nil {
		Error.Printf("Unable to retrieve IPv4 Addr. Error msg: %v. Setting to nil.\n", errV4)
		IPAddrs = append(IPAddrs, "nil")
		err = append(err, errV4)
	} else {
		defer ipv4.Body.Close()
		ipv4Bytes, errV4Bytes := io.ReadAll(ipv4.Body)
		if errV4Bytes != nil {
			Error.Printf("Unable to extract IPv4 String from response. Error msg: %s.\n", errV4Bytes)
			err = append(err, errV4Bytes)
		} else {
			IPAddrs = append(
				IPAddrs,
				string(ipv4Bytes[:]),
			)
		}
	}
	ipv6, errV6 := http.Get(ipv6Endpoint)
	if errV6 != nil {
		Error.Printf("Unable to retrieve IPv6 Addr. Error msg: %v. Setting to nil.\n", errV6)
		err = append(err, errV6)
	} else {
		defer ipv4.Body.Close()
		ipv6Bytes, errV6Bytes := io.ReadAll(ipv6.Body)
		if errV6Bytes != nil {
			Error.Printf("Unable to extract IPv4 String from Response. Error msg: %s.\n", errV6Bytes)
			err = append(err, errV6Bytes)
		} else {
			IPAddrs = append(
				IPAddrs,
				string(ipv6Bytes[:]),
			)
		}
	}
	return
}

// CheckDomainExists checks the DNS registrar for the zone
func CheckDomainExists(baseDomain, plugin string) (zoneID string, err error) {

	return "", nil
}
