package plugins

import (
	"cmd/edd/pkg/utils"
	"context"
	"encoding/json"
	"fmt"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/route53"
	"github.com/aws/aws-sdk-go-v2/service/route53/types"
	"log"
	"strings"
)

//Route53Upsert updates (or inserts) a resource record (if one does not exist) for a given domain name input and record type
func Route53Upsert(vars utils.ArgVars, localIP, recordtype string)(err error){
	vars.DomainID, err = getHostedZoneID(vars.Domain, vars.DNSMaxRecordsReturned)
	fmt.Printf("domainID: %s\n", vars.DomainID)
	//retrieve resource records and iterate through the list
	resourceRecord, err := getResourceRecordSets(vars.Domain, recordtype, vars)
	recordsToChange := []types.Change{}
	for _, record := range resourceRecord {
		if len(record.ResourceRecords) == 1 {
			record.ResourceRecords[0].Value = &localIP
		} else {
		log.Fatal("unexpected number of records to update")
		}
		changeRecord := types.Change{
			Action: types.ChangeActionUpsert,
			ResourceRecordSet: &record,
		}
		recordsToChange = append(recordsToChange, changeRecord)
	}
	err = updateRecordSet(&vars.DomainID, recordsToChange)
	if err != nil {
		log.Fatal(err)
	}

	return nil
}

func getHostedZoneID(domain string, maxItems int32)(domainID string, err error){
	hostedZoneInputs := route53.ListHostedZonesInput{
		MaxItems: &maxItems,
	}
	ctx := context.TODO()
	cfg, err := config.LoadDefaultConfig(ctx)
	if err != nil {
		log.Fatal(err)
	}
	r53client := route53.NewFromConfig(cfg)
	hostedZonesOutput, err := r53client.ListHostedZones(ctx, &hostedZoneInputs)
	if err != nil {
		log.Fatal(err)
	}

	tempZoneName := strings.Split(domain, ".")
	for {
		for _, zone := range hostedZonesOutput.HostedZones {
			if *zone.Name == (strings.Join(tempZoneName, ".") + ".") {
				tempDomainID := *zone.Id
				domainID = strings.TrimPrefix(tempDomainID, "/hostedzone/")
			}
		}
		fmt.Printf("tempZoneName before removing subdomain: %s", tempZoneName)
		tempZoneName = tempZoneName[1:len(tempZoneName)]
		fmt.Printf(" tempZoneName after removing subdomain: %s\n", tempZoneName)
		if domainID != "" {
			break
		}
	}

	fmt.Printf("hostedZonesOutput:%s\n", domainID)
	return domainID, nil
}

func getResourceRecordSets(domain, recordType string, vars utils.ArgVars)(RRSet []types.ResourceRecordSet, err error){
	rrInputs := route53.ListResourceRecordSetsInput{
		HostedZoneId:          &vars.DomainID,
		MaxItems:              &vars.DNSMaxRecordsReturned,
		StartRecordIdentifier: nil,
		StartRecordName: 	   &domain,
		StartRecordType:       types.RRType(recordType),

	}

	ctx := context.TODO()
	cfg, err := config.LoadDefaultConfig(ctx)
	if err != nil {
		log.Fatal(err)
	}
	r53client := route53.NewFromConfig(cfg)

	for {
		RRSetInitial, err := r53client.ListResourceRecordSets(ctx, &rrInputs)
		if err != nil {
			log.Fatal(err)
		}
		for _, resourceRecord := range RRSetInitial.ResourceRecordSets{
			if string(resourceRecord.Type) == recordType {
				fmt.Printf("resourceRecord Name: %s\ndomain: %s\n", *resourceRecord.Name, domain)
				if strings.Compare(*resourceRecord.Name, domain + ".") == 0{
					fmt.Printf("resourceRecord Name: %s\ndomain: %s\n", *resourceRecord.Name, domain)
					RRSet = append(RRSet, resourceRecord)
				}
			}
		}
		if RRSetInitial.NextRecordIdentifier == nil {
			break
		}
	}

	resourceRecordJSON, err := json.MarshalIndent(RRSet, "", "  ")
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("%s\n", resourceRecordJSON)

	return RRSet, nil
}

func updateRecordSet(zoneID *string, recordChanges []types.Change)(err error){
	changeBatchInput := types.ChangeBatch{
		Changes: recordChanges,
		Comment: nil,
	}
	updateRecordSetInput := route53.ChangeResourceRecordSetsInput{
		ChangeBatch:  &changeBatchInput,
		HostedZoneId: zoneID,
	}
	ctx := context.TODO()
	cfg, err := config.LoadDefaultConfig(ctx)
	if err != nil {
		log.Fatal(err)
	}
	r53client := route53.NewFromConfig(cfg)
	changeRRSetOutput, err := r53client.ChangeResourceRecordSets(ctx, &updateRecordSetInput)
	if err != nil {
		log.Fatal(err)
	}
	changeRRSetOutputJSON, err := json.MarshalIndent(changeRRSetOutput, "", "  ")
	if err != nil {
		fmt.Printf("unable to marshal JSON output from changeRRSetOutput")
	}
	fmt.Printf("ChangeRRSetOutput:\n%s\n", changeRRSetOutputJSON)
	return nil
}