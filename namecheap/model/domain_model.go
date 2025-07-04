package model

import "encoding/xml"

type NamecheapResponse struct {
	XMLName         xml.Name        `xml:"ApiResponse"`
	Status          string          `xml:"Status,attr"`
	Errors          []ErrorMessage  `xml:"Errors>Error"`
	CommandResponse CommandResponse `xml:"CommandResponse"`
}

type ErrorMessage struct {
	Number string `xml:"Number,attr"`
	Text   string `xml:",chardata"`
}

type CommandResponse struct {
	DomainCheckResult []DomainCheckResult `xml:"DomainCheckResult"`
}

type DomainCheckResult struct {
	Domain      string `xml:"Domain,attr"`
	Available   string `xml:"Available,attr"`
	ErrorNo     string `xml:"ErrorNo,attr,omitempty"`
	Description string `xml:"Description,attr,omitempty"`
}
