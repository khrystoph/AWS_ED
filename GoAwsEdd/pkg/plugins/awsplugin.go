package plugins

import(
	"cmd/edd/pkg/edd"
	"github.com/aws/aws-sdk-go-v2/service/route53"
)

//Route53Upsert updates (or inserts) a resource record (if one does not exist) for a given domain name input and record type
func Route53Upsert(vars edd.ArgVars, localIP, recordtype string)(err error){
	//retrieve resource records and iterate through the list
	resourceRecord, err := getResourceRecordSets(vars.Domain, recordtype)
	return nil
}

func getResourceRecordSets(domain, recordType string)(RRSet route53.ListResourceRecordSetsOutput, err error){
	
}