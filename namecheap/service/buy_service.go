package service

import (
	"encoding/xml"
	"fmt"

	"namecheap-microservice/config"
	"namecheap-microservice/database"
	"namecheap-microservice/model"

	"github.com/namecheap/go-namecheap-sdk/v2/namecheap"
)

// }
func BuyDomain(client *namecheap.Client, req model.DomainPurchaseRequest) (bool, string, error) {
	params := map[string]string{
		"Command":                 "namecheap.domains.create",
		"ApiUser":                 config.ApiUser,
		"ApiKey":                  config.ApiKey,
		"Username":                config.ApiUser,
		"ClientIp":                config.ClientIp,
		"DomainName":              req.Domain,
		"Years":                   "1",
		"RegistrantFirstName":     req.FirstName,
		"RegistrantLastName":      req.LastName,
		"RegistrantAddress1":      req.Address,
		"RegistrantCity":          req.City,
		"RegistrantStateProvince": "MH",
		"RegistrantPostalCode":    req.PostalCode,
		"RegistrantCountry":       req.Country,
		"RegistrantPhone":         req.Phone,
		"RegistrantEmailAddress":  req.Email,
	}

	for _, role := range []string{"Tech", "Admin", "AuxBilling"} {
		params[role+"FirstName"] = req.FirstName
		params[role+"LastName"] = req.LastName
		params[role+"Address1"] = req.Address
		params[role+"City"] = req.City
		params[role+"StateProvince"] = "MH"
		params[role+"PostalCode"] = req.PostalCode
		params[role+"Country"] = req.Country
		params[role+"Phone"] = req.Phone
		params[role+"EmailAddress"] = req.Email
	}

	var result struct {
		XMLName         xml.Name `xml:"ApiResponse"`
		Status          string   `xml:"Status,attr"`
		CommandResponse struct {
			DomainCreateResult struct {
				Registered bool   `xml:"Registered,attr"`
				Domain     string `xml:"Domain,attr"`
				OrderID    string `xml:"OrderID,attr"`
			} `xml:"DomainCreateResult"`
		} `xml:"CommandResponse"`
		Errors struct {
			Error string `xml:"Error"`
		} `xml:"Errors"`
	}

	_, err := client.DoXML(params, &result)
	if err != nil {
		return false, "", fmt.Errorf("API error: %v", err)
	}

	if result.Status == "ERROR" {
		return false, "", fmt.Errorf("Namecheap Error: %s", result.Errors.Error)
	}

	domainName := result.CommandResponse.DomainCreateResult.Domain

	if result.CommandResponse.DomainCreateResult.Registered {
		aRecord := req.ARecord
		if aRecord == "" {
			aRecord = "82.25.106.75"
		}
		cName := req.CName
		if cName == "" {
			cName = "indigo-spoonbill-233511.hostingersite.com"
		}

		if err := SetDNSRecords(client, req.Domain, aRecord, cName); err != nil {
			return true, domainName, fmt.Errorf("domain registered but DNS setup failed: %v", err)
		}

		if err := SaveDomainWithDNS(domainName, req.FirstName+" "+req.LastName, aRecord, cName, req.Price, req.Tax, req.Total); err != nil {
			return true, domainName, fmt.Errorf("DNS set but DB save failed: %v", err)
		}
	}

	return result.CommandResponse.DomainCreateResult.Registered, domainName, nil
}

func SaveDomainWithDNS(domainName string, customer string, aRecord string, cname string, price float64, tax float64, total float64) error {
	// Create domain purchase object
	domain := model.DomainPurchase{
		Name:      domainName,
		Purchased: true,
		Customer:  customer,
		Price:     price,
		Tax:       tax,
		Total:     total,
	}

	// Create DNS record object
	dns := model.DNSRecord{
		ARecord: aRecord,
		CName:   cname,
	}

	// Save both in a transaction
	tx := database.DB.Begin()

	if err := tx.Create(&domain).Error; err != nil {
		tx.Rollback()
		return err
	}

	dns.DomainPurchaseID = domain.ID
	if err := tx.Create(&dns).Error; err != nil {
		tx.Rollback()
		return err
	}

	return tx.Commit().Error
}
