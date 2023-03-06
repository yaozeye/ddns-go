package dns

import (
	"context"
	"fmt"

	"github.com/Azure/azure-sdk-for-go/profiles/latest/dns/mgmt/dns"
	"github.com/Azure/go-autorest/autorest/azure/auth"
)

type AzureDNSProvider struct {
	subscriptionID string
	resourceGroup  string
	dnsZoneName    string
}

func NewAzureDNSProvider(subscriptionID, resourceGroup, dnsZoneName string) (*AzureDNSProvider, error) {
	return &AzureDNSProvider{
		subscriptionID: subscriptionID,
		resourceGroup:  resourceGroup,
		dnsZoneName:    dnsZoneName,
	}, nil
}

func (p *AzureDNSProvider) UpdateRecord(domain, recordType, recordName, recordValue string, ttl int) error {
	// Set up Azure credentials and DNS client
	authorizer, err := auth.NewAuthorizerFromEnvironment()
	if err != nil {
		return err
	}
	dnsClient := dns.NewRecordSetsClient(p.subscriptionID)
	dnsClient.Authorizer = authorizer

	// Create the record set object
	var recordSet *dns.RecordSet
	switch recordType {
	case "A":
		recordSet = &dns.RecordSet{
			Name: &recordName,
			Type: &recordType,
			TTL:  &int64(ttl),
			ARecords: &[]dns.ARecord{
				{Ipv4Address: &recordValue},
			},
		}
	case "AAAA":
		recordSet = &dns.RecordSet{
			Name: &recordName,
			Type: &recordType,
			TTL:  &int64(ttl),
			AAAARecords: &[]dns.AAAARecord{
				{Ipv6Address: &recordValue},
			},
		}
	default:
		return fmt.Errorf("unsupported record type '%s'", recordType)
	}

	// Update the DNS record
	_, err = dnsClient.CreateOrUpdate(context.Background(), p.resourceGroup, p.dnsZoneName, domain, dns.RecordType(recordType), *recordSet, "", "")
	if err != nil {
		return err
	}

	return nil
}
