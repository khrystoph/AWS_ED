package utils

type ArgVars struct {
	Domain string `json:"domain"`
	V4Endpoint string `json:"v4endpoint"`
	V6Endpoint string `json:"v6endpoint"`
	DNSProviderPlugin string `json:"dnsproviderplugin"`
	DNSMaxRecordsReturned int32 `json:"dnsmaxrecordsreturned"`
	DomainID string `json:"domainid"`
}