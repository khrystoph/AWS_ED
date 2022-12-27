package utils

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/route53"
	"github.com/aws/aws-sdk-go-v2/service/route53/types"
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
func RetrieveResourceRecordSet(zoneID, baseDomain, rrType string) (resourceRecord types.ResourceRecordSet, err error) {
	svcClient, err := Route53Client()
	if err != nil {
		fmt.Printf("unable to initialize route53 client. Error msg:%s\n", err)
		return types.ResourceRecordSet{}, err
	}
	rrSetsInput := &route53.ListResourceRecordSetsInput{
		HostedZoneId: &zoneID,
	}
	resourceRecords, err := svcClient.ListResourceRecordSets(context.Background(), rrSetsInput)
	if err != nil {
		fmt.Printf("unable to retrieve resource records. Error msg: %s\n", err)
		return types.ResourceRecordSet{}, err
	}

	for _, records := range resourceRecords.ResourceRecordSets {
		if *records.Name == baseDomain && string(records.Type) == rrType {
			resourceRecord = records
			break
		}
	}

	return resourceRecord, nil
}

// UpdateResourceRecordSet takes the resourceRecordSet ID and current IP address to update the mismatched record
// in route53 with the correct IP for the resource record set type and domain/subdomain.
func UpdateResourceRecordSet(resourceRecordSets []types.ResourceRecordSet, hostedZoneID *string) (requestID string, err error) {
	var (
		rrSetChanges            []types.Change
		rrSetChangeBatch        = &types.ChangeBatch{Changes: rrSetChanges}
		resourceRecordSetInputs = &route53.ChangeResourceRecordSetsInput{
			ChangeBatch:  rrSetChangeBatch,
			HostedZoneId: hostedZoneID}
	)
	svcClient, err := Route53Client()
	if err != nil {
		fmt.Printf("error setting up route53 client to update resource record sets. Error msg: %s\n", err)
		return "", err
	}
	if !(len(resourceRecordSets) >= 1) {
		err = errors.New("no resource record sets provided as input to function. bailing out")
		return "", err
	}
	for index := range resourceRecordSets {
		change := types.Change{
			Action:            "UPSERT",
			ResourceRecordSet: &resourceRecordSets[index],
		}
		rrSetChanges = append(rrSetChanges, change)
	}

	resourceRecordSetInputs.ChangeBatch.Changes = rrSetChanges

	//print intput as JSON as a debug message to ensure the inputs are going in correctly
	rrSetInputsJSON, err := json.MarshalIndent(resourceRecordSetInputs, "", "  ")
	if err != nil {
		fmt.Printf("Warning: unable to marshal inputs as JSON. Error msg: %s\n", err)
	} else {
		fmt.Printf("Resource Record Set Inputs: \n%s\n", rrSetInputsJSON)
	}

	resourceRecordSetUpdateOutput, err := svcClient.ChangeResourceRecordSets(context.Background(), resourceRecordSetInputs)
	if err != nil {
		fmt.Printf("error updating resource record sets. Error msg: %s", err)
		return "", err
	}

	resourceRecordSetUpdateOutputJSON, err := json.MarshalIndent(resourceRecordSetUpdateOutput, "", "  ")
	if err != nil {
		fmt.Printf("unable to response output JSON. Error msg: %s\n", err)
	} else {
		fmt.Printf("resource record sets updated. Following output returned: %s\n", resourceRecordSetUpdateOutputJSON)
	}

	requestID = *resourceRecordSetUpdateOutput.ChangeInfo.Id
	return requestID, nil
}
