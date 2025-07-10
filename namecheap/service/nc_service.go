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
		"Command":    "namecheap.domains.check",
		"DomainList": domain,
	}

	fmt.Printf("👤 Loaded API User: '%s'\n", config.ApiUser)
	fmt.Printf("🌐 Checking domain: %s\n", domain)
	fmt.Printf("📤 Sending request with params: %+v\n", params)

	var result struct {
		XMLName         xml.Name `xml:"ApiResponse"`
		Status          string   `xml:"Status,attr"`
		CommandResponse struct {
			DomainCheckResult struct {
				Available     string `xml:"Available,attr"`
				Domain        string `xml:"Domain,attr"`
				IsPremiumName string `xml:"IsPremiumName,attr"`
				IsRestricted  string `xml:"IsRestricted,attr"`
			} `xml:"DomainCheckResult"`
		} `xml:"CommandResponse"`
		Errors struct {
			Error string `xml:"Error"`
		} `xml:"Errors"`
	}

	_, err := client.DoXML(params, &result)
	if err != nil {
		fmt.Printf("❌ API call failed: %v\n", err)
		return false, fmt.Errorf("API error: %v", err)
	}

	// Pretty print raw XML response
	rawResponse, _ := xml.MarshalIndent(result, "", "  ")
	fmt.Println("🧾 Raw API Response:\n", string(rawResponse))

	if result.Errors.Error != "" {
		fmt.Println("⚠️ Namecheap internal error:", result.Errors.Error)
		return false, fmt.Errorf("Namecheap error: %s", result.Errors.Error)
	}

	// Extract result
	d := result.CommandResponse.DomainCheckResult

	// Post-call status logs
	fmt.Printf("📥 API Response Status: %s\n", result.Status)
	fmt.Printf("✅ Domain Checked: %s\n", d.Domain)
	fmt.Printf("🟢 Available: %s | ⚠️ Premium: %s | ⛔ Restricted: %s\n",
		d.Available, d.IsPremiumName, d.IsRestricted,
	)

	// Handle errors from Namecheap
	if result.Status == "ERROR" && result.Errors.Error != "" {
		fmt.Printf("🚫 Namecheap Error: %s\n", result.Errors.Error)
		return false, fmt.Errorf("Namecheap error: %s", result.Errors.Error)
	}

	// Handle restricted domains
	if d.IsRestricted == "true" {
		fmt.Println("⛔ Domain is restricted or banned.")
		return false, fmt.Errorf("domain is restricted or banned by registry")
	}

	// Optional: warn about premium names
	if d.IsPremiumName == "true" {
		fmt.Println("⚠️ Note: This is a premium domain. Higher cost may apply.")
	}

	available := d.Available == "true"
	fmt.Printf("🔍 Final Availability: %v\n", available)

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
