package service

import (
	"encoding/xml"
	"fmt"
	"strings"

	"namecheap-microservice/config"

	"github.com/namecheap/go-namecheap-sdk/v2/namecheap"
)

// GetDomainInfo retrieves domain info using the Namecheap SDK
func CheckDomain(client *namecheap.Client, domain string) (bool, error) {
	params := map[string]string{
		"Command": "namecheap.domains.check",
		// "ApiUser":    config.ApiUser,
		// "ApiKey":     config.ApiKey,
		// "Username":   config.ApiUser,
		// "ClientIp":   config.ClientIp,
		"DomainList": domain,
	}
	fmt.Printf("👤 Loaded API User: '%s'\n", config.ApiUser)

	var result struct {
		XMLName         xml.Name `xml:"ApiResponse"`
		Status          string   `xml:"Status,attr"`
		CommandResponse struct {
			DomainCheckResult struct {
				Available string `xml:"Available,attr"` // "true" or "false"
				Domain    string `xml:"Domain,attr"`
			} `xml:"DomainCheckResult"`
		} `xml:"CommandResponse"`
		Errors struct {
			Error string `xml:"Error"`
		} `xml:"Errors"`
	}

	// Log: BEFORE API call
	fmt.Printf("🌐 Checking domain: %s\n", domain)              //test1
	fmt.Printf("📤 Sending request with params: %+v\n", params) //test 2

	_, err := client.DoXML(params, &result)
	if err != nil {
		fmt.Printf("❌ API call failed: %v\n", err) //test 3
		return false, fmt.Errorf("API error: %v", err)
	}
	// rawResponse, _ := xml.MarshalIndent(result, "", "  ")
	//test 4

	// Log: AFTER API call
	fmt.Printf("📥 Raw API response Status: %s\n", result.Status) //test 5
	fmt.Printf("✅ Domain: %s | Available: %s\n",                 //test 6
		result.CommandResponse.DomainCheckResult.Domain,
		result.CommandResponse.DomainCheckResult.Available,
	)

	if result.Status == "ERROR" {
		fmt.Printf("⚠️ Namecheap error: %s\n", result.Errors.Error)
		return false, fmt.Errorf("Namecheap error: %s", result.Errors.Error)
	}

	available := result.CommandResponse.DomainCheckResult.Available == "true"
	fmt.Printf("🔎 Final check: Is domain available? %v\n", available)

	return available, nil
}

func SetDNSRecords(client *namecheap.Client, domain string, ip string) error {
	parts := strings.SplitN(domain, ".", 2)
	if len(parts) != 2 {
		return fmt.Errorf("invalid domain format: %s", domain)
	}
	sld := parts[0]
	tld := parts[1]
	params := map[string]string{
		"Command":     "namecheap.domains.dns.setHosts",
		"ApiUser":     config.ApiUser,
		"ApiKey":      config.ApiKey,
		"Username":    config.ApiUser,
		"ClientIp":    config.ClientIp,
		"SLD":         sld,
		"TLD":         tld,
		"DomainName":  domain,
		"HostName1":   "@",
		"RecordType1": "A",
		"Address1":    "82.25.106.75",
		"TTL1":        "1800",
		"HostName2":   "www",
		"RecordType2": "CNAME",
		"Address2":    "indigo-spoonbill-233511.hostingersite.com",
		"TTL2":        "1800",
	}

	var response struct {
		XMLName xml.Name `xml:"ApiResponse"`
		Status  string   `xml:"Status,attr"`
		Errors  struct {
			Error string `xml:"Error"`
		} `xml:"Errors"`
		CommandResponse struct {
			IsSuccess bool `xml:"IsSuccess,attr"`
		} `xml:"CommandResponse>DomainDNSSetHostsResult"`
	}

	_, err := client.DoXML(params, &response)
	if err != nil {
		return fmt.Errorf("API call failed: %v", err)
	}

	if response.Status == "ERROR" {
		return fmt.Errorf("Namecheap error: %s", response.Errors.Error)
	}

	if !response.CommandResponse.IsSuccess {
		return fmt.Errorf("DNS update failed, unknown reason")
	}

	fmt.Println("✅ DNS records successfully set for", domain)
	return nil
}
