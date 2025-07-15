package service

import (
	"encoding/xml"
	"fmt"
	"strings"

	"namecheap-microservice/config"
	"namecheap-microservice/model"
	"strconv"

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
func GetRegistrarLockStatus(client *namecheap.Client, domain string) (bool, error) {
	params := map[string]string{
		"Command":    "namecheap.domains.getRegistrarLock",
		"ApiUser":    config.ApiUser,
		"ApiKey":     config.ApiKey,
		"Username":   config.ApiUser,
		"ClientIp":   config.ClientIp,
		"DomainName": domain, // ✅ Namecheap requires this
	}

	// 🔍 Debug logging — check if DomainName is being passed
	fmt.Println("📦 Params being sent to Namecheap:")
	for k, v := range params {
		if k == "ApiKey" {
			v = "*****"
		}
		fmt.Printf("   %s: %s\n", k, v)
	}

	var result struct {
		XMLName         xml.Name `xml:"ApiResponse"`
		Status          string   `xml:"Status,attr"`
		CommandResponse struct {
			RegistrarLockStatus string `xml:"RegistrarLockStatus"`
		} `xml:"CommandResponse>DomainGetRegistrarLockResult"`
		Errors struct {
			Error string `xml:"Error"`
		} `xml:"Errors"`
	}

	// 📡 Call Namecheap with those params
	_, err := client.DoXML(params, &result)
	if err != nil {
		return false, fmt.Errorf("API call failed: %v", err)
	}
	if result.Status == "ERROR" {
		return false, fmt.Errorf("Namecheap error: %s", result.Errors.Error)
	}

	// ✅ Parse status
	locked := strings.ToLower(result.CommandResponse.RegistrarLockStatus) == "true"
	return locked, nil
}

func SetRegistrarLock(client *namecheap.Client, domain string, lock bool) error {
	parts := strings.SplitN(domain, ".", 2)
	if len(parts) != 2 {
		return fmt.Errorf("invalid domain format")
	}
	sld, tld := parts[0], parts[1]

	value := "UNLOCK"
	if lock {
		value = "LOCK"
	}
	params := map[string]string{
		"Command":    "namecheap.domains.setRegistrarLock",
		"ApiUser":    config.ApiUser,
		"ApiKey":     config.ApiKey,
		"Username":   config.ApiUser,
		"ClientIp":   config.ClientIp,
		"DomainName": domain,
		"SLD":        sld,
		"TLD":        tld,
		"LockAction": value,
	}

	if lock {
		params["LockAction"] = "LOCK"
	}

	// Debug log
	fmt.Println("🔐 LockAction being sent:", params["LockAction"])

	var result struct {
		XMLName xml.Name `xml:"ApiResponse"`
		Status  string   `xml:"Status,attr"`
		Errors  struct {
			Error string `xml:"Error"`
		} `xml:"Errors"`
		CommandResponse struct {
			IsSuccess bool `xml:"IsSuccess,attr"`
		} `xml:"CommandResponse>DomainSetRegistrarLockResult"`
	}

	_, err := client.DoXML(params, &result)
	if err != nil {
		return err
	}
	if result.Status == "ERROR" {
		return fmt.Errorf("Namecheap error: %s", result.Errors.Error)
	}
	if !result.CommandResponse.IsSuccess {
		return fmt.Errorf("Failed to change lock state")
	}
	return nil
}

func SetDNSRecordsAdvanced(client *namecheap.Client, domain string, records []model.DNSInputRecord) error {
	parts := strings.SplitN(domain, ".", 2)
	if len(parts) != 2 {
		return fmt.Errorf("invalid domain format: %s", domain)
	}
	sld := parts[0]
	tld := parts[1]

	params := map[string]string{
		"Command":  "namecheap.domains.dns.setHosts",
		"ApiUser":  config.ApiUser,
		"ApiKey":   config.ApiKey,
		"Username": config.ApiUser,
		"ClientIp": config.ClientIp,
		"SLD":      sld,
		"TLD":      tld,
	}

	// 🧾 Log the incoming DNS records
	fmt.Println("🟡 Preparing DNS Records for:", domain)
	for i, rec := range records {
		n := i + 1
		params[fmt.Sprintf("HostName%d", n)] = rec.Host
		params[fmt.Sprintf("RecordType%d", n)] = rec.Type
		params[fmt.Sprintf("Address%d", n)] = rec.Value
		params[fmt.Sprintf("TTL%d", n)] = strconv.Itoa(rec.TTL)

		if rec.Type == "MX" {
			params[fmt.Sprintf("MXPref%d", n)] = strconv.Itoa(rec.MXPref)
		}
		if rec.Type == "CAA" {
			params[fmt.Sprintf("Flag%d", n)] = strconv.Itoa(rec.Flag)
			params[fmt.Sprintf("Tag%d", n)] = rec.Tag
		}

		fmt.Printf("🔧 Record %d: Host=%s | Type=%s | Value=%s | TTL=%d\n", n, rec.Host, rec.Type, rec.Value, rec.TTL)
	}

	// ✅ Log final API parameters
	fmt.Println("🌐 DNS API Params being sent to Namecheap:")
	for k, v := range params {
		if strings.HasPrefix(k, "ApiKey") {
			v = "*****" // hide secret
		}
		fmt.Printf("   %s: %s\n", k, v)
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
		fmt.Println("❌ API call failed:", err)
		return fmt.Errorf("API call failed: %v", err)
	}

	if response.Status == "ERROR" {
		fmt.Println("🟥 Namecheap API Error:", response.Errors.Error)
		return fmt.Errorf("Namecheap error: %s", response.Errors.Error)
	}
	if !response.CommandResponse.IsSuccess {
		fmt.Println("❌ DNS update failed (Namecheap returned IsSuccess=false)")
		return fmt.Errorf("DNS update failed")
	}

	fmt.Println("✅ DNS records successfully set for", domain)
	return nil
}
