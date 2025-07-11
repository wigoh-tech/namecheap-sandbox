package model

type DomainPurchaseRequest struct {
	Domain     string  `json:"domain"`
	FirstName  string  `json:"firstName"`
	LastName   string  `json:"lastName"`
	Email      string  `json:"email"`
	Address    string  `json:"address"`
	City       string  `json:"city"`
	Phone      string  `json:"phone"`
	PostalCode string  `json:"postalCode"`
	Country    string  `json:"country"` // e.g., "IN"
	Price      float64 `json:"price"`   // Base price of the domain
	Tax        float64 `json:"tax"`     // Tax amount
	Total      float64 `json:"total"`   // Total amount after tax

	// DNS records
	ARecord string `json:"aRecord"` // A record for the domain
	CName   string `json:"cName"`   // CNAME record for the domain
}
