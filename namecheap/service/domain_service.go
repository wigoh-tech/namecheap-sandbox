package service

import (
	"fmt"
	"namecheap-microservice/database"
	"namecheap-microservice/model"
	"time"
)

func RevokeDomain(domainName string) error {
	var domain model.DomainPurchase
	if err := database.DB.First(&domain, "name = ?", domainName).Error; err != nil {
		return fmt.Errorf("domain not found")
	}

	now := time.Now()
	domain.Revoked = true
	domain.RevokedAt = &now

	if err := database.DB.Save(&domain).Error; err != nil {
		return fmt.Errorf("failed to revoke domain")
	}
	return nil
}

func UnrevokeOldDomains() {
	threshold := time.Now().Add(-1 * time.Minute)

	var domains []model.DomainPurchase
	err := database.DB.Where("revoked = ? AND revoked_at < ?", true, threshold).Find(&domains).Error
	if err != nil {
		fmt.Println("Failed to query old revoked domains:", err)
		return
	}

	for _, d := range domains {
		d.Revoked = false
		d.RevokedAt = nil
		if err := database.DB.Save(&d).Error; err != nil {
			fmt.Println("Failed to unrevoke domain:", d.Name)
		} else {
			fmt.Println("✅ Domain made available again:", d.Name)
		}
	}
}
