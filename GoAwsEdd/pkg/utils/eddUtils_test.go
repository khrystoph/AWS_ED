package utils

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestResolveFQDNResolveFQDN(t *testing.T) {
	var testDomainName = "jbecomputersolutions.com"
	//assert := assert.New(t)
	got, err := ResolveFQDN(&testDomainName)
	assert.Equal(t, nil, err, "Check if err returned a value other than nil. %v", err)
	//assert.Equal(t, nil, got, "Both values should be nil as domain should exist")
	fmt.Printf("Domain %s response is: %v\n", testDomainName, got)
}

func TestLookUpIPClientIP(t *testing.T) {
	var (
		v4Endpoint = "https://ipv4.moreip.jbecomputersolutions.com"
		v6Endpoint = "https://ipv6.moreip.jbecomputersolutions.com"
	)
	got, err := LookUpIPClientIP(v4Endpoint, v6Endpoint)
	assert.Equal(t, 0, len(err), "Expecting 0 length for errors. Got %d errors.\n", len(err))
	for _, ip := range got {
		fmt.Printf("%s", ip)
	}
}
