package utils

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"net"
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
	//check IPv4 Validity
	gotv4, err := LookUpIPClientIP(v4Endpoint)
	assert.Equal(t, nil, err, "Expecting nil error. Err msg: %s.\n", err)
	fmt.Printf("gotv4val: %v\n", gotv4)
	isNotNilIP := net.ParseIP(gotv4)
	assert.NotEqual(t, nil, isNotNilIP, "Expecting non-nil value, got: %v", isNotNilIP)
	fmt.Printf("result of isNotNilIP: %v\n", isNotNilIP)

	//Check IPv6 validity
	gotv6, err := LookUpIPClientIP(v6Endpoint)
	assert.Equal(t, nil, err, "Expecting 0 errors. Err msg: %s.\n", err)
	isNotNilIPv6 := net.ParseIP(gotv6)
	assert.NotEqual(t, nil, isNotNilIPv6, "Expecting non-nil value, got: %v", isNotNilIPv6)
	fmt.Printf("gotv6val: %v\n", gotv6)
	fmt.Printf("result of isNotNilIPv6: %v\n", isNotNilIPv6)
}
