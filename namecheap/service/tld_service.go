package service

import (
	"encoding/xml"
	"fmt"
	"namecheap-microservice/config"

	"github.com/namecheap/go-namecheap-sdk/v2/namecheap"
)

func GetTLDList(client *namecheap.Client) ([]string, error) {
	params := map[string]string{
		"Command":  "namecheap.domains.getTldList",
		"ApiUser":  config.ApiUser,
		"ApiKey":   config.ApiKey,
		"Username": config.ApiUser,
		"ClientIp": config.ClientIp,
	}

	var result struct {
		XMLName         xml.Name `xml:"ApiResponse"`
		Status          string   `xml:"Status,attr"`
		CommandResponse struct {
			Tlds []struct {
				Name string `xml:"Name,attr"`
			} `xml:"Tlds>Tld"`
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

	var tlds []string
	for _, tld := range result.CommandResponse.Tlds {
		tlds = append(tlds, tld.Name)
	}

	return tlds, nil
}
