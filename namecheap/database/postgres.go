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
func UpdateDNSInDB(domainName string, newRecords []model.DNSRecord) error {
	var domain model.DomainPurchase

	// Find domain
	if err := DB.Where("name = ?", domainName).First(&domain).Error; err != nil {
		return fmt.Errorf("domain not found: %w", err)
	}

	// Delete existing DNS records for this domain
	if err := DB.Where("domain_purchase_id = ?", domain.ID).Delete(&model.DNSRecord{}).Error; err != nil {
		return fmt.Errorf("failed to clear old DNS records: %w", err)
	}

	// Assign DomainPurchaseID and insert new records
	for i := range newRecords {
		newRecords[i].DomainPurchaseID = domain.ID
	}

	if err := DB.Create(&newRecords).Error; err != nil {
		return fmt.Errorf("failed to save new DNS records: %w", err)
	}

	return nil
}
func GetDNSInputRecords(domain string) ([]model.DNSInputRecord, error) {
	var domainObj model.DomainPurchase
	if err := DB.Where("name = ?", domain).First(&domainObj).Error; err != nil {
		return nil, err
	}

	var records []model.DNSRecord
	if err := DB.Where("domain_purchase_id = ?", domainObj.ID).Find(&records).Error; err != nil {
		return nil, err
	}

	var input []model.DNSInputRecord
	for _, r := range records {
		input = append(input, model.DNSInputRecord{
			Type:   r.Type,
			Host:   r.Host,
			Value:  r.Value,
			TTL:    r.TTL,
			MXPref: r.MXPref,
			Flag:   r.Flag,
			Tag:    r.Tag,
		})
	}
	return input, nil
}
