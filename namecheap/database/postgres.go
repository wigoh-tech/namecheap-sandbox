package database

import (
	"fmt"
	"log"
	"os"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"

	"namecheap-microservice/model"
)

var DB *gorm.DB

func ConnectDB() {
	dsn := fmt.Sprintf(
		"host=localhost user=%s password=%s dbname=%s port=5432 sslmode=disable",
		os.Getenv("DB_USER"),
		os.Getenv("DB_PASSWORD"),
		os.Getenv("DB_NAME"),
	)

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatal("❌ Failed to connect to the database:", err)
	}

	DB = db
	err = db.AutoMigrate(&model.DomainPurchase{}, &model.DNSRecord{})
	if err != nil {
		log.Fatal("❌ AutoMigrate failed:", err)
	}

	fmt.Println("✅ Connected to the PostgreSQL database!")
}

func SaveDomainWithDNS(domain model.DomainPurchase, dns model.DNSRecord) error {
	tx := DB.Begin()

	if err := tx.Create(&domain).Error; err != nil {
		tx.Rollback()
		return err
	}

	dns.DomainPurchaseID = domain.ID
	if err := tx.Create(&dns).Error; err != nil {
		tx.Rollback()
		return err
	}

	return tx.Commit().Error
}

// UpdateDNSInDB updates A record and CNAME in the dns_records table for a domain
func UpdateDNSInDB(domainName, aRecord, cName string) error {
	var domain model.DomainPurchase

	// Step 1: Find the domain by name
	if err := DB.Where("name = ?", domainName).First(&domain).Error; err != nil {
		return fmt.Errorf("domain not found: %w", err)
	}

	// Step 2: Update DNS record where foreign key matches
	var dns model.DNSRecord
	if err := DB.Where("domain_purchase_id = ?", domain.ID).First(&dns).Error; err != nil {
		return fmt.Errorf("DNS record not found: %w", err)
	}

	dns.ARecord = aRecord
	dns.CName = cName

	// Save updates
	if err := DB.Save(&dns).Error; err != nil {
		return fmt.Errorf("failed to update DNS: %w", err)
	}

	return nil
}
