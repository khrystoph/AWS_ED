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

func LookUpIPClientIP(endpoint string) (IPAddr string, err error) {
	endpointResp, err := http.Get(endpoint)
	if err != nil {
		Error.Printf("Unable to retrieve IP address from %s. Error msg: %s.\n", endpoint, err)
		return "", err
	}
	defer endpointResp.Body.Close()

	//read response body and extract the bytes array, then check for errors
	IPAddrBytes, err := io.ReadAll(endpointResp.Body)
	if err != nil {
		Error.Printf("Unable to read endpoint response %s", err)
		return "", err
	}

	//drop any trailing newline character
	if IPAddrBytes[len(IPAddrBytes)-1] == '\n' {
		IPAddrBytes = IPAddrBytes[:len(IPAddrBytes)-1]
	}

	//Convert Bytes returned from request into string
	IPAddr = string(IPAddrBytes)
	return IPAddr, nil
}

// CheckDomainExists checks the DNS registrar for the zone
func CheckDomainExists(baseDomain, plugin string) (zoneID string, err error) {

	return "", nil
}
