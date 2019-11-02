//conversion of edd.py into golang

package main

import (
	"errors"
	"flag"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/route53"
)

var (
	Trace                             *log.Logger
	Info                              *log.Logger
	Warning                           *log.Logger
	Error                             *log.Logger
	domain, sessionProfile, awsRegion string
	traceHandle                       io.Writer
	infoHandle                        io.Writer = os.Stdout
	warningHandle                     io.Writer = os.Stderr
	errorHandle                       io.Writer = os.Stderr
	sess                              *session.Session
	maxItems                          = "100"
	moreIPv4                          = "https://ipv4.moreip.jbecomputersolutions.com"
	moreIPv6                          = "https://ipv6.moreip.jbecomputersolutions.com"
)

func init() {

	Trace = log.New(traceHandle,
		"TRACE: ",
		log.Ldate|log.Ltime|log.Lshortfile)

	Info = log.New(infoHandle,
		"INFO: ",
		log.Ldate|log.Ltime|log.Lshortfile)

	Warning = log.New(warningHandle,
		"WARNING: ",
		log.Ldate|log.Ltime|log.Lshortfile)

	Error = log.New(errorHandle,
		"ERROR: ",
		log.Ldate|log.Ltime|log.Lshortfile)

	flag.StringVar(&domain, "d", "example.com", "enter your fully qualified domain name here. Default: example.com")
	flag.StringVar(&domain, "domain", "example.com", "enter your fully qualified domain name here. Default: example.com")
	flag.StringVar(&sessionProfile, "profile", "default", "enter the profile you wish to use to connect. Default: default")
	flag.StringVar(&sessionProfile, "p", "default", "enter the profile you wish to use to connect. Default: default")
	flag.StringVar(&awsRegion, "region", "us-east-1", "Enter region you wish to connect with. Default: us-east-1")
}

//handleIPv4 handles cases where A record in r53. On match, do nothing. On mismatch, update to new val.
func handleIPv4(rrecord *route53.ResourceRecordSet, zoneID string) (update bool, err error) {
	var (
		ipv4String string
	)
	svc := route53.New(sess)

	for i := 0; i < 3; i++ {
		currentIPv4, err := http.Get(moreIPv4)
		defer currentIPv4.Body.Close()
		if err != nil && i < 3 {
			time.Sleep(time.Duration(i) * time.Second)
			err = nil
			continue
		} else if err != nil && i >= 3 {
			Error.Println("Error getting ipv4 address.")
			return false, err
		} else {
			body, _ := ioutil.ReadAll(currentIPv4.Body)
			ipv4String = string(body)
			break
		}
	}
	if *rrecord.ResourceRecords[0].Value == strings.Trim(ipv4String, "\n") {
		Info.Println("IPv4 Records match. Doing nothing.")
		return false, nil
	}

	Info.Println("IPv4 Records do not match. Updating.")
	Info.Println("Current IPv4 address: ", ipv4String)
	Info.Println("Current IPv4 record: ", *rrecord.ResourceRecords[0].Value)
	*rrecord.ResourceRecords[0].Value = ipv4String
	changeAction := "UPSERT"
	changeBatch := route53.ChangeBatch{Changes: []*route53.Change{&route53.Change{Action: &changeAction, ResourceRecordSet: rrecord}}}
	input := route53.ChangeResourceRecordSetsInput{
		HostedZoneId: &zoneID,
		ChangeBatch:  &changeBatch,
	}
	output, err := svc.ChangeResourceRecordSets(&input)
	if err != nil {
		return false, err
	}
	Info.Println(output)
	return true, nil
}

//handleIPv6 handles cases where AAAA record in r53. On match, do nothing. On mismatch, update to new val.
func handleIPv6(rrecord *route53.ResourceRecordSet, zoneID string) (update bool, err error) {
	var (
		ipv6String string
	)
	svc := route53.New(sess)
	for i := 0; i < 3; i++ {
		currentIPv6, err := http.Get(moreIPv6)
		defer currentIPv6.Body.Close()
		if err != nil && i < 3 {
			time.Sleep(time.Duration(i) * time.Second)
			err = nil
			continue
		} else if err != nil && i >= 3 {
			Error.Println("Error getting ipv6 address.")
			return false, err
		} else {
			body, _ := ioutil.ReadAll(currentIPv6.Body)
			ipv6String = string(body)
			break
		}
	}
	if *rrecord.ResourceRecords[0].Value == strings.Trim(ipv6String, "\n") {
		Info.Println("IPv6 Records match. Doing nothing.")
		return false, nil
	}

	Info.Println("IPv6 Records do not match. Updating.")
	Info.Println("Current IPv6 address: ", ipv6String)
	Info.Println("Current IPv6 record: ", *rrecord.ResourceRecords[0].Value)
	*rrecord.ResourceRecords[0].Value = ipv6String
	changeAction := "UPSERT"
	changeBatch := route53.ChangeBatch{Changes: []*route53.Change{&route53.Change{Action: &changeAction, ResourceRecordSet: rrecord}}}
	input := route53.ChangeResourceRecordSetsInput{
		HostedZoneId: &zoneID,
		ChangeBatch:  &changeBatch,
	}
	output, err := svc.ChangeResourceRecordSets(&input)
	if err != nil {
		return false, err
	}
	Info.Println(output)
	return true, nil
}

//checkHostDomainNameExists checks whether or not the host has a domain name set
func checkHostDomainNameExists(domainPtr *string) (err error) {
	hostname, err := os.Hostname()
	if err != nil {
		return err
	} else {
		Info.Println("Using system hostname: ", hostname)
		*domainPtr = hostname
	}
	return err
}

//checkDomainExists checks that the domain exists in Route53
func checkDomainExists(baseDomain string) (zoneID string, err error) {
	svc := route53.New(sess)

	input := route53.ListHostedZonesByNameInput{DNSName: &baseDomain}
	output, err := svc.ListHostedZonesByName(&input)

	if err != nil {
		return "", err
	}
	var zones string
	for _, hostedZone := range output.HostedZones {
		if *hostedZone.Name == baseDomain {
			Info.Println("Output:\n", hostedZone)
			zoneID = strings.Split(*hostedZone.Id, "/")[2]
			return zoneID, err
		}
		zones = zones + "\nName: " + *hostedZone.Name + "\n" + "Zone ID: " + *hostedZone.Id
	}

	hostedZonesError := "Domain not found in R53. HostedZones:\n" + zones
	err = errors.New(hostedZonesError)
	return zoneID, err
}

//checkZoneRecords finds resource record sets and returns the set of them that match the domain
func checkZoneRecords(zID string) (zonerecords []route53.ResourceRecordSet, err error) {
	svc := route53.New(sess)

	input := route53.ListResourceRecordSetsInput{HostedZoneId: &zID, MaxItems: &maxItems}
	output, err := svc.ListResourceRecordSets(&input)

	if err != nil {
		return zonerecords, err
	}
	for _, record := range output.ResourceRecordSets {
		if *record.Name == domain {
			zonerecords = append(zonerecords, *record)
		}
	}

	return zonerecords, err
}

func main() {
	flag.Parse()
	var (
		err        error
		baseDomain string
	)
	traceHandle = ioutil.Discard

	sess, err = session.NewSession(&aws.Config{
		Region:      aws.String(awsRegion),
		Credentials: credentials.NewSharedCredentials("", sessionProfile),
	})

	_, err = sess.Config.Credentials.Get()
	if err != nil {
		Error.Println("error retrieving credentials. Profile name: ", sessionProfile)
		Error.Println("Error msg: ", err)
		os.Exit(1)
	}

	Info.Println("Check if domain is default")
	if domain == "example.com" {
		Info.Println("Checking if there is a locally defined FQDN.")
		err = checkHostDomainNameExists(&domain)
		if err != nil {
			Error.Println("error checking hostname in OS. Error msg:\n", err)
		}
	}

	if domain[len(domain)-1:] != "." {
		domain += "."
	}
	domainLen := len(strings.Split(domain, "."))
	if domainLen > 2 {
		baseDomain = strings.Join(strings.Split(domain, ".")[domainLen-3:], ".")
	} else {
		baseDomain = domain
	}

	Info.Println("Checking if domain is in Route53")
	zoneID, err := checkDomainExists(baseDomain)
	if err != nil {
		Error.Println("Domain: ", domain)
		Error.Println(err)
		os.Exit(1)
	}
	Info.Println("ZoneID: ", zoneID)
	Info.Println("Checking to see what records exist for the domain.")
	records, err := checkZoneRecords(zoneID)
	if err != nil {
		Error.Println("Records: ", records)
		Error.Println(err)
		os.Exit(1)
	}
	Info.Println("Checking if domain needs to be updated")
	var updated = false
	for _, record := range records {
		if *record.Type == "A" {
			updated, err = handleIPv4(&record, zoneID)
			if err != nil {
				Error.Println(err)
				os.Exit(1)
			}
		}
		if *record.Type == "AAAA" {
			updated, err = handleIPv6(&record, zoneID)
			if err != nil {
				Error.Println(err)
				os.Exit(1)
			}
		}
	}
	records, err = checkZoneRecords(zoneID)
	if err != nil {
		Error.Println("Records: ", records)
		Error.Println(err)
		os.Exit(1)
	}
	if updated {
		Info.Println(records)
	}
	return
}
