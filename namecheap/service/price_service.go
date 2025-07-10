package service

import (
	"encoding/xml"
	"fmt"
	"namecheap-microservice/config"
	"strings"

	"github.com/namecheap/go-namecheap-sdk/v2/namecheap"
)

func GetWholesalePrice(client *namecheap.Client, domain string) (float64, error) {
	params := map[string]string{
		"Command":     "namecheap.domains.getPricing",
		"ApiUser":     config.ApiUser,
		"ApiKey":      config.ApiKey,
		"Username":    config.ApiUser,
		"ClientIp":    config.ClientIp,
		"ProductType": "DOMAIN",
	}

	var result struct {
		XMLName         xml.Name `xml:"ApiResponse"`
		Status          string   `xml:"Status,attr"`
		CommandResponse struct {
			ProductType []struct {
				Name  string `xml:"Name,attr"`
				Price struct {
					YourPrice float64 `xml:"YourPrice"`
				} `xml:"Price"`
			} `xml:"ProductType"`
		} `xml:"CommandResponse"`
		Errors struct {
			Error string `xml:"Error"`
		} `xml:"Errors"`
	}

	_, err := client.DoXML(params, &result)
	if err != nil {
		return 0, err
	}

	if result.Status == "ERROR" {
		return 0, fmt.Errorf("Namecheap Error: %s", result.Errors.Error)
	}

	// Extract TLD
	parts := strings.Split(domain, ".")
	tld := parts[len(parts)-1]

	for _, pt := range result.CommandResponse.ProductType {
		if strings.EqualFold(pt.Name, tld) {
			return pt.Price.YourPrice, nil
		}
	}

	return 0, fmt.Errorf("TLD %s not found in pricing", tld)
}
