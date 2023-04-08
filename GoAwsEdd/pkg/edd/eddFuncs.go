package edd

import (
	"cmd/edd/pkg/plugins"
	"fmt"
	"io"
	"net"
	"net/http"
	"strings"
	"sync"
)

type ArgVars struct {
	Domain string `json:"domain"`
	V4Endpoint string `json:"v4endpoint"`
	V6Endpoint string `json:"v6endpoint"`
	DNSProviderPlugin string `json:"dnsproviderplugin"`
}

func DNSRecordExists(domain string)(dnsRecordVal []string, err error){
	dnsRecordVal, err = net.LookupHost(domain)
	return dnsRecordVal, nil
}

func CheckAndUpdateRecords(ipAddress, recordType string, args ArgVars, group *sync.WaitGroup)(err error){
	defer group.Done()

	var(myIP, endpoint string
	myIPBytes []byte
	)

	switch recordType {
	case "A":
		endpoint = args.V4Endpoint
	case "AAAA":
		endpoint = args.V6Endpoint
	default:
		return fmt.Errorf("CheckAndUpdateRecords() did not match record tpe")
	}

	//fetch IP address for the record type we're requesting
	resp, err := http.Get(endpoint)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	//convert response to byte array by reading the io stream
	myIPBytes, err = io.ReadAll(resp.Body)
	if err != nil {
		return err
	}
	//convert the ip byte array back into string
	myIP = strings.TrimSuffix(string(myIPBytes), "\n")

	fmt.Printf("My IP: %s\nRecord IP: %s\n", myIP, ipAddress)
	if myIP == ipAddress {
		fmt.Printf("No need to update %s record for domain: %s. Records match.\n", myIP, args.Domain)
		return nil
	} else {
		fmt.Printf("IP Addresses do not match. Updating Record.\n")
		switch args.DNSProviderPlugin {
		case "route53":
			plugins.Route53Upsert(args, myIP, recordType)
		}
	}

	return nil
}