package model

import (
	"time"
)

type DomainPurchase struct {
	ID        uint   `gorm:"primaryKey"`
	Name      string `gorm:"not null;unique"`
	Purchased bool
	Revoked   bool
	Customer  string
	Price     float64
	Tax       float64
	Total     float64
	CreatedAt time.Time

	DNSRecord DNSRecord `gorm:"foreignKey:DomainPurchaseID"`
}
type DNSRecord struct {
	ID               uint `gorm:"primaryKey"`
	DomainPurchaseID uint `gorm:"not null"` // Foreign key to DomainPurchase
	ARecord          string
	CName            string
}
