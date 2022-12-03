package main

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/route53"
)

func main() {
	clientProfile := "default"
	domainName := "jbecomputersolutions.com"
	cfg, err := config.LoadDefaultConfig(context.TODO(),
		config.WithSharedConfigProfile(clientProfile),
	)

	if err != nil {
		fmt.Errorf("error loading Default config for profile: %s. Message:\n%w", clientProfile, err)
		return
	}

	route53Client := route53.NewFromConfig(cfg)

	hostedZones, err := route53Client.ListHostedZonesByName(context.TODO(), &route53.ListHostedZonesByNameInput{
		DNSName: &domainName,
	})
	if err != nil {
		fmt.Errorf("failed to list hosted zones. Error msg:\n%w", err)
		return
	}
	marshalledJSON, err := json.Marshal(hostedZones)
	if err != nil {
		fmt.Errorf("failed to marshal response as JSON %s", err)
		return
	}
	fmt.Printf("Hosted Zone info:\n%s", marshalledJSON)
}
