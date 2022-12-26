package utils

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/route53"
	"strings"
)

// Route53Client returns an AWS client object that can interact with the Route53 APIs
func Route53Client() (r53Client *route53.Client, err error) {
	cfg, err := config.LoadDefaultConfig(context.TODO())
	if err != nil {
		Error.Fatal(err)
		return nil, err
	}
	r53Client = route53.NewFromConfig(cfg)
	return
}

func CheckRoute53DomainExists(baseDomain string) (zoneID string, err error) {
	svcClient, err := Route53Client()

	if err != nil {
		fmt.Printf("Unable to create client from default config. Error msg: %s\n", err)
		return "", err
	}

	hostedZonesResponseOutput, err := svcClient.ListHostedZonesByName(context.Background(),
		&route53.ListHostedZonesByNameInput{
			DNSName: &baseDomain,
		},
	)
	if err != nil {
		return "", err
	}
	hostedZonesJSON, err := json.Marshal(hostedZonesResponseOutput)
	if err != nil {
		fmt.Printf("Unable to marshal hostedzonesoutput.")
	} else {
		fmt.Printf("hostedZonesOutput:\n%s", hostedZonesJSON)
	}
	for _, zone := range hostedZonesResponseOutput.HostedZones {
		hostedZoneString := *zone.Name
		hostedZoneString = hostedZoneString[:len(hostedZoneString)-1]
		if baseDomain == hostedZoneString {
			zoneID = *zone.Id
			zoneID = strings.Split(zoneID, "/")[2]
			fmt.Printf("matched baseDomain: %s to hostedZoneString: %s\n", baseDomain,
				hostedZoneString)
		}
	}
	fmt.Printf("zoneID: %s\n", zoneID)
	return zoneID, nil
}

// RetrieveResourceRecordSet retrieves the resource records from a given zoneID and returns IP addresses of the RR set
func RetrieveResourceRecordSet(zoneID, baseDomain, rrType string) (rrID string, err error) {
	svcClient, err := Route53Client()
	if err != nil {
		fmt.Printf("unable to initialize route53 client. Error msg:%s\n", err)
		return "", err
	}
	rrSetsInput := &route53.ListResourceRecordSetsInput{
		HostedZoneId: &zoneID,
	}
	resourceRecords, err := svcClient.ListResourceRecordSets(context.Background(), rrSetsInput)
	if err != nil {
		fmt.Printf("unable to retrieve resource records. Error msg: %s\n", err)
		return "", err
	}
	rrJSON, err := json.MarshalIndent(resourceRecords, "", "  ")
	if err != nil {
		fmt.Printf("unable to marshal resource records. Error msg: %s", err)
		return "", err
	}
	fmt.Printf("Resource records:\n%s", rrJSON)
	return
}
