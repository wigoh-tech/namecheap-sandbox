package service

import (
	"encoding/xml"
	"fmt"
	"namecheap-microservice/config"
	"namecheap-microservice/model"
	"strings"

	"github.com/namecheap/go-namecheap-sdk/v2/namecheap"
)

func TransferDomain(client *namecheap.Client, req model.DomainPurchaseRequest) (bool, string, error) {
	parts := strings.SplitN(req.Domain, ".", 2)
	if len(parts) != 2 {
		return false, "", fmt.Errorf("invalid domain format")
	}

	params := map[string]string{
		"Command":    "namecheap.domains.transfer.create",
		"ApiUser":    config.ApiUser,
		"ApiKey":     config.ApiKey,
		"Username":   config.ApiUser,
		"ClientIp":   config.ClientIp,
		"DomainName": req.Domain,
		"EPPCode":    req.EPPCode,
		"Years":      "1",
	}

	for _, role := range []string{"Registrant", "Admin", "Tech", "AuxBilling"} {
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
			TransferID string `xml:"TransferID,attr"`
		} `xml:"CommandResponse>DomainTransferCreateResult"`
		Errors struct {
			Error string `xml:"Error"`
		} `xml:"Errors"`
	}

	_, err := client.DoXML(params, &result)
	if err != nil {
		return false, "", err
	}
	if result.Status == "ERROR" {
		return false, "", fmt.Errorf("Namecheap error: %s", result.Errors.Error)
	}

	return true, result.CommandResponse.TransferID, nil
}
func GetTransferList(client *namecheap.Client) ([]string, error) {
	params := map[string]string{
		"Command":  "namecheap.domains.transfer.getList",
		"ApiUser":  config.ApiUser,
		"ApiKey":   config.ApiKey,
		"Username": config.ApiUser,
		"ClientIp": config.ClientIp,
	}

	var result struct {
		XMLName         xml.Name `xml:"ApiResponse"`
		Status          string   `xml:"Status,attr"`
		CommandResponse struct {
			Transfers []struct {
				DomainName string `xml:"DomainName,attr"`
			} `xml:"DomainTransferGetListResult>Transfer"`
		} `xml:"CommandResponse"`
		Errors struct {
			Error string `xml:"Error"`
		} `xml:"Errors"`
	}

	_, err := client.DoXML(params, &result)
	if err != nil {
		return nil, err
	}
	if result.Status == "ERROR" {
		return nil, fmt.Errorf("Namecheap error: %s", result.Errors.Error)
	}

	var domains []string
	for _, t := range result.CommandResponse.Transfers {
		domains = append(domains, t.DomainName)
	}
	return domains, nil
}
func GetTransferStatus(client *namecheap.Client, domain string) (string, error) {
	params := map[string]string{
		"Command":    "namecheap.domains.transfer.getStatus",
		"ApiUser":    config.ApiUser,
		"ApiKey":     config.ApiKey,
		"Username":   config.ApiUser,
		"ClientIp":   config.ClientIp,
		"DomainName": domain,
	}

	var result struct {
		XMLName         xml.Name `xml:"ApiResponse"`
		Status          string   `xml:"Status,attr"`
		CommandResponse struct {
			Status string `xml:"Status"`
		} `xml:"CommandResponse>DomainTransferGetStatusResult"`
		Errors struct {
			Error string `xml:"Error"`
		} `xml:"Errors"`
	}

	_, err := client.DoXML(params, &result)
	if err != nil {
		return "", err
	}
	if result.Status == "ERROR" {
		return "", fmt.Errorf("Namecheap error: %s", result.Errors.Error)
	}

	return result.CommandResponse.Status, nil
}
