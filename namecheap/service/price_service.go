package service

import (
	"encoding/xml"
	"fmt"

	"github.com/namecheap/go-namecheap-sdk/v2/namecheap"
)

func GetWholesalePrice(client *namecheap.Client, tld string) (float64, error) {
	params := map[string]string{
		"Command":         "namecheap.users.getPricing",
		"ProductType":     "DOMAIN",
		"ProductName":     tld,
		"ProductCategory": "REGISTER",
		"ActionName":      "REGISTER",
	}

	var resp struct {
		XMLName         xml.Name `xml:"ApiResponse"`
		Status          string   `xml:"Status,attr"`
		CommandResponse struct {
			Pricing struct {
				Product struct {
					Name  string `xml:"Name,attr"`
					Price struct {
						Price    float64 `xml:"Price,attr"`
						Currency string  `xml:"Currency,attr"`
					} `xml:"Price"`
				} `xml:"Product"`
			} `xml:"UserGetPricingResult>ProductType>ProductCategory>Product"`
		} `xml:"CommandResponse"`
		Errors struct {
			Error string `xml:"Error"`
		} `xml:"Errors"`
	}

	_, err := client.DoXML(params, &resp)
	if err != nil {
		return 0, err
	}
	if resp.Status == "ERROR" {
		return 0, fmt.Errorf("pricing error: %s", resp.Errors.Error)
	}

	return resp.CommandResponse.Pricing.Product.Price.Price, nil
}
