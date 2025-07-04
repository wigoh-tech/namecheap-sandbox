package service

import (
	"fmt"
	"namecheap-microservice/database"
	"namecheap-microservice/model"

	"github.com/namecheap/go-namecheap-sdk/v2/namecheap"
)

func RevokeDomain(domainName string) error {
	var domain model.DomainPurchase
	if err := database.DB.First(&domain, "name = ?", domainName).Error; err != nil {
		return fmt.Errorf("domain not found")
	}

	domain.Revoked = true
	if err := database.DB.Save(&domain).Error; err != nil {
		return fmt.Errorf("failed to revoke domain")
	}
	return nil
}

func MoveDomain(client *namecheap.Client, domainName string, newARecord string, newCNAME string) error {
	// Step 1: Update DNS on Namecheap
	if err := SetDNSRecords(client, domainName, newARecord); err != nil {
		return fmt.Errorf("dns update failed: %v", err)
	}

	// Step 2: Fetch the domain purchase record
	var domain model.DomainPurchase
	if err := database.DB.First(&domain, "name = ?", domainName).Error; err != nil {
		return fmt.Errorf("domain not found")
	}

	// Step 3: Fetch its linked DNS record
	var dns model.DNSRecord
	if err := database.DB.First(&dns, "domain_purchase_id = ?", domain.ID).Error; err != nil {
		return fmt.Errorf("dns record not found")
	}

	// Step 4: Update values
	dns.ARecord = newARecord
	dns.CName = newCNAME

	// Step 5: Save the updated DNS record
	if err := database.DB.Save(&dns).Error; err != nil {
		return fmt.Errorf("failed to update dns record")
	}

	return nil
}
