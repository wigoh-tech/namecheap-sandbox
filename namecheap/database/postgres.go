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
