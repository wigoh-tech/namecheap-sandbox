package model

import (
	"time"
)

type DomainPurchase struct {
	ID        uint       `gorm:"primaryKey" json:"id"`
	Name      string     `gorm:"not null;unique" json:"name"`
	Purchased bool       `json:"purchased"`
	Revoked   bool       `json:"revoked"`
	RevokedAt *time.Time `json:"revokedAt"`
	Customer  string     `json:"customer"`

	WholesalePriceUSD float64 `json:"wholesalePriceUsd"`
	MarkupPercent     float64 `json:"markupPercent"`

	RetailPriceINR float64   `json:"retailPriceInr"`
	Price          float64   `json:"price"`
	Tax            float64   `json:"tax"`
	Total          float64   `json:"total"`
	CreatedAt      time.Time `json:"createdAt"`

	DNSRecord DNSRecord `gorm:"foreignKey:DomainPurchaseID" json:"dnsRecord"`
}

type DNSRecord struct {
	ID               uint   `gorm:"primaryKey" json:"id"`
	DomainPurchaseID uint   `gorm:"not null" json:"-"`
	ARecord          string `json:"aRecord"`
	CName            string `json:"cName"`
}
