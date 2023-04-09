//conversion of edd.py into golang

package main

import (
	"cmd/edd/pkg/edd"
	"cmd/edd/pkg/utils"
	"encoding/json"
	"flag"
	"fmt"
	"net"
	"os"
	"strings"
	"sync"
)

var (
	domain, moreIPv4, moreIPv6,dnsProviderPlugin string
	maxItems                          int32     = 100
)

func init() {

	flag.StringVar(&domain, "d", "example.com", "enter your fully qualified domain name here. Default: example.com")
	flag.StringVar(&domain, "domain", "example.com", "enter your fully qualified domain name here. Default: example.com")
	flag.StringVar(&moreIPv4, "v4", "https://ipv4.moreip.jbecomputersolutions.com", "enter an IPv4 endpoint to use to check your IPv4 address. Default:" +
		"https://ipv4.moreip.jbecomputersolutions.com")
	flag.StringVar(&moreIPv6, "v6", "https://ipv6.moreip.jbecomputersolutions.com", "enter an IPv6 endpoint to use to check your IPv6 address. Default:" +
		"https://ipv6.moreip.jbecomputersolutions.com")
	flag.StringVar(&dnsProviderPlugin, "p", "route53", "enter the provider plugin you wish to use for the dns API endpoint." +
		" Default: route53")
}

func worker(inputArgs utils.ArgVars)(err error){
	var wg = sync.WaitGroup{}
	//Marshal inputs to json to log as info logging later.
	inputArgsJSON, err := json.MarshalIndent(inputArgs, "", "  ")
	if err != nil {
		fmt.Printf("unable to marshal input args to JSON. data")
		return err
	}

	//Check for existence of AAAA record and A records
	records, err := edd.DNSRecordExists(domain)
	if err != nil {
		fmt.Printf("unable to resolve domain and return list of records.")
		return err
	}

	for _, ip := range records {
		recordType := ""
		addr := net.ParseIP(ip)
		if addr.To4() != nil && addr != nil {
			recordType = "A"
		} else if strings.Contains(ip, ":") && addr != nil {
			recordType = "AAAA"
		} else {
			fmt.Printf("ip record is not of type A or AAAA")
			continue
		}
		wg.Add(1)
		go edd.CheckAndUpdateRecords(ip, recordType, inputArgs, &wg)
	}

	wg.Wait()

	//debug section remove later
	fmt.Printf(string(inputArgsJSON))
	fmt.Printf("records for domain (%s):\n %v\n", domain, records)
	//end debug section
	return nil
}

func main(){
	//parse input flags
	flag.Parse()

	inputVars := utils.ArgVars{
		Domain: domain,
		V4Endpoint: moreIPv4,
		V6Endpoint: moreIPv6,
		DNSProviderPlugin: dnsProviderPlugin,
		DNSMaxRecordsReturned: 100,
	}

	err := worker(inputVars)
	if err != nil {
		fmt.Printf("Error executing main worker function. Error returned %v", err)
	}

	os.Exit(0)
}