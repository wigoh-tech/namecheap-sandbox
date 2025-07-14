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

	// Fill other roles: Tech, Admin, AuxBilling
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

	// Call Namecheap API to register domain
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
	if domainName == "" {
		domainName = req.Domain
	}

	// After domain is registered
	if result.CommandResponse.DomainCreateResult.Registered {
		// Step 1: Prepare DNS records
		records := req.DNSRecords
		if len(records) == 0 {
			records = []model.DNSInputRecord{
				{Type: "A", Host: "@", Value: "82.25.106.75", TTL: 1800},
				{Type: "CNAME", Host: "www", Value: "indigo-spoonbill-233511.hostingersite.com", TTL: 1800},
			}
		}

		// Step 2: Setup DNS records in Namecheap
		if err := SetDNSRecordsAdvanced(client, domainName, records); err != nil {
			return true, domainName, fmt.Errorf("domain registered but DNS setup failed: %v", err)
		}
		var dbRecords []model.DNSRecord
		for _, r := range records {
			dbRecords = append(dbRecords, model.DNSRecord{
				Type:   r.Type,
				Host:   r.Host,
				Value:  r.Value,
				TTL:    r.TTL,
				MXPref: r.MXPref,
				Flag:   r.Flag,
				Tag:    r.Tag,
			})
		}

		// Step 3: Save to PostgreSQL
		if err := SaveDomainWithDNS(domainName, req.FirstName+" "+req.LastName, dbRecords, req.Price, req.Tax, req.Total); err != nil {
			return true, domainName, fmt.Errorf("DNS set but DB save failed: %v", err)
		}
	}

	// Final result
	return result.CommandResponse.DomainCreateResult.Registered, domainName, nil
}

func SaveDomainWithDNS(domainName string, customer string, records []model.DNSRecord, price float64, tax float64, total float64) error {
	// Create domain purchase object
	domain := model.DomainPurchase{
		Name:      domainName,
		Purchased: true,
		Customer:  customer,
		Price:     price,
		Tax:       tax,
		Total:     total,
	}

	// Save both in a transaction
	tx := database.DB.Begin()

	if err := tx.Create(&domain).Error; err != nil {
		tx.Rollback()
		return err
	}

	// Set the foreign key and insert records
	for i := range records {
		records[i].DomainPurchaseID = domain.ID
		if err := tx.Create(&records[i]).Error; err != nil {
			tx.Rollback()
			return err
		}
	}

	return tx.Commit().Error
}
