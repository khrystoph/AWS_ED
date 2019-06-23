//conversion of edd.py into golang

package main

import (
	"errors"
	"flag"
	"io"
	"io/ioutil"
	"log"
	"os"
	"strings"

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

func handleIPv4() {

}

func handleIPv6() {

}

func checkHostDomainNameExists(domainPtr *string) (err error) {
	hostname, err := os.Hostname()
	if err != nil {
		return err
	} else {
		Info.Println("Using system hostname: ", hostname)
		domainPtr = &hostname
	}
	return err
}

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

	Info.Println(zonerecords)

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

	Info.Println("baseDomain: ", baseDomain)
	Info.Println("Check if domain is default")
	if baseDomain == "example.com" {
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
	println(records)
	Info.Println("Checking if domain needs to be updated")
}
