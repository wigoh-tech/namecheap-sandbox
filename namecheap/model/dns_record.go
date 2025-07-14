package model

type DNSRecord struct {
	ID               uint `gorm:"primaryKey" json:"id"`
	DomainPurchaseID uint `gorm:"not null" json:"-"`

	Host  string `json:"host"`  // @, www, mail, etc.
	Type  string `json:"type"`  // A, AAAA, MX, CNAME, etc.
	Value string `json:"value"` // e.g., IP or hostname
	TTL   int    `json:"ttl"`   // Time to live (e.g., 1800)

	// Optional fields
	MXPref int    `json:"mxPref,omitempty"` // for MX records
	Flag   int    `json:"flag,omitempty"`   // for CAA records
	Tag    string `json:"tag,omitempty"`    // for CAA records
}
