package model

type DomainPurchaseRequest struct {
	Domain     string  `json:"domain"`
	EPPCode    string  `json:"eppCode"`
	FirstName  string  `json:"firstName"`
	LastName   string  `json:"lastName"`
	Email      string  `json:"email"`
	Address    string  `json:"address"`
	City       string  `json:"city"`
	Phone      string  `json:"phone"`
	PostalCode string  `json:"postalCode"`
	Country    string  `json:"country"`
	Price      float64 `json:"price"`
	Tax        float64 `json:"tax"`
	Total      float64 `json:"total"`

	// NEW: flexible list of DNS records
	DNSRecords []DNSInputRecord `json:"dnsRecords"`
}
