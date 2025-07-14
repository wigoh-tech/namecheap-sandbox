package model

type DNSInputRecord struct {
	Host  string `json:"host"`
	Type  string `json:"type"`
	Value string `json:"value"`
	TTL   int    `json:"ttl"`

	// Optional for specific record types
	MXPref int    `json:"mxPref,omitempty"`
	Flag   int    `json:"flag,omitempty"`
	Tag    string `json:"tag,omitempty"`
}
