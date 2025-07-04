package model

type DomainPurchaseRequest struct {
	Domain     string `json:"domain"`
	FirstName  string `json:"firstName"`
	LastName   string `json:"lastName"`
	Email      string `json:"email"`
	Address    string `json:"address"`
	City       string `json:"city"`
	Phone      string `json:"phone"`
	PostalCode string `json:"postalCode"`
	Country    string `json:"country"` // e.g., "IN"
}
