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
