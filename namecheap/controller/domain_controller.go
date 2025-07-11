package controller

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"strings"

	"namecheap-microservice/database"
	"namecheap-microservice/utils"
	"net/http"
	"net/url"

	"namecheap-microservice/model"
	"namecheap-microservice/service"
)

func CheckDomainHandler(w http.ResponseWriter, r *http.Request) {
	raw := r.URL.Query().Get("domain")
	if raw == "" {
		http.Error(w, `{"error": "domain is required"}`, http.StatusBadRequest)
		return
	}

	// Safely decode the domain from URL encoding
	domain, err := url.QueryUnescape(raw)
	if err != nil {
		fmt.Println("Error decoding domain:", err) //test 1 to check if the error is logged
		// Return a 400 Bad Request error with a JSON response
		http.Error(w, `{"error": "invalid domain encoding"}`, http.StatusBadRequest)
		return
	}
	// Validate the domain format
	parsedDomain, err := utils.ParseDomain(domain)
	if err != nil {
		fmt.Println("Error parsing domain:", err) //test 2 to check if the error is logged
		http.Error(w, fmt.Sprintf(`{"error": "%s"}`, err.Error()), http.StatusBadRequest)
		return
	}

	// You can now use `parsedDomain`
	fmt.Println("Valid domain:", parsedDomain)

	// Call the service to get domain info

	client := service.NewNamecheapClient()
	available, err := service.CheckDomain(client, domain)
	if err != nil {
		fmt.Println("CheckDomain error:", err) //test 3 to check if the error is logged
		http.Error(w, fmt.Sprintf(`{"error": "%s"}`, err.Error()), http.StatusInternalServerError)
		return
	}

	fmt.Printf("🔍 Domain availability for %s: %v\n", domain, available)

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"available": available,
	})
}

func BuyDomainHandler(w http.ResponseWriter, r *http.Request) {
	var req model.DomainPurchaseRequest
	fmt.Println("🚀 /buy-domain endpoint hit")

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, `{"error":"invalid request body"}`, http.StatusBadRequest)
		return
	}
	fmt.Printf("📦 Request payload: %+v\n", req)

	client := service.NewNamecheapClient()

	// 🔁 Step 1: Use wholesale price in USD from Namecheap
	wholesaleUSD, err := service.GetWholesalePrice(client, req.Domain)
	if err != nil {
		http.Error(w, `{"error": "failed to fetch wholesale price"}`, http.StatusInternalServerError)
		return
	}

	// 🔁 Step 2: Convert USD to INR (you can later use a live exchange rate API)
	const exchangeRate = 83.0
	retailINR := wholesaleUSD * exchangeRate

	// 💰 Step 3: Add 18% GST
	tax := retailINR * 0.18
	total := retailINR + tax

	// 📦 Step 4: Fill price info into req for DB save
	req.Price = retailINR
	req.Tax = tax
	req.Total = total

	success, domainName, err := service.BuyDomain(client, req)
	if err != nil {
		fmt.Println("BuyDomain error:", err)
		http.Error(w, fmt.Sprintf(`{"error": "%s"}`, err.Error()), http.StatusInternalServerError)
		return
	}

	fmt.Println("✅ Domain purchase completed successfully:", domainName)

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"success": success,
		"domain":  domainName,
	})
}

func GetDomainPriceHandler(w http.ResponseWriter, r *http.Request) {
	domain := r.URL.Query().Get("domain")
	if domain == "" {
		http.Error(w, `{"error": "domain is required"}`, http.StatusBadRequest)
		return
	}

	parts := strings.Split(domain, ".")
	if len(parts) < 2 {
		http.Error(w, `{"error": "invalid domain"}`, http.StatusBadRequest)
		return
	}

	tld := parts[len(parts)-1]
	base, ok := service.TLDBasePrices[tld]
	if !ok {
		base = 1000.00 // fallback
	}

	tax := base * 0.3
	total := base + tax

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]float64{
		"base":  base,
		"tax":   tax,
		"total": total,
	})
}
func ListDomains(w http.ResponseWriter, r *http.Request) {
	var domains []model.DomainPurchase
	if err := database.DB.Preload("DNSRecord").Find(&domains).Error; err != nil {
		http.Error(w, `{"error": "failed to fetch domains"}`, http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"domains": domains,
	})
}
func RevokeDomainHandler(w http.ResponseWriter, r *http.Request) {
	// Read raw body
	bodyBytes, _ := io.ReadAll(r.Body)
	fmt.Println("Received raw body:", string(bodyBytes))

	// Rewind body for decoding
	r.Body = io.NopCloser(bytes.NewBuffer(bodyBytes))

	var body struct {
		Domain string `json:"domain"`
	}
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil || body.Domain == "" {
		fmt.Println("Decode error or empty domain:", err)
		http.Error(w, `{"error":"invalid domain"}`, http.StatusBadRequest)
		return
	}

	fmt.Println("Decoded domain:", body.Domain)

	if err := service.RevokeDomain(body.Domain); err != nil {
		http.Error(w, fmt.Sprintf(`{"error":"%s"}`, err.Error()), http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(map[string]interface{}{
		"success": true,
		"message": "domain revoked",
	})
}

func UpdateDNSHandler(w http.ResponseWriter, r *http.Request) {
	var body struct {
		Domain     string `json:"domain"`
		ARecord    string `json:"aRecord"`
		CName      string `json:"cName"`
		RecordType string `json:"recordType"`
	}
	fmt.Println("🚀 /update-dns endpoint hit")

	if err := json.NewDecoder(r.Body).Decode(&body); err != nil || body.Domain == "" || body.RecordType == "" {
		http.Error(w, `{"error": "Invalid input: missing domain or record type"}`, http.StatusBadRequest)
		return
	}

	client := service.NewNamecheapClient()

	switch body.RecordType {
	case "A":
		if body.ARecord == "" {
			http.Error(w, `{"error": "Missing A record"}`, http.StatusBadRequest)
			return
		}
		err := service.SetDNSRecords(client, body.Domain, body.ARecord, "")
		if err != nil {
			fmt.Println("SetDNSRecords A error:", err)
			http.Error(w, fmt.Sprintf(`{"error": "Failed to update A record: %s"}`, err.Error()), http.StatusInternalServerError)
			return
		}
	case "CNAME":
		if body.CName == "" {
			http.Error(w, `{"error": "Missing CNAME record"}`, http.StatusBadRequest)
			return
		}
		err := service.SetDNSRecords(client, body.Domain, "", body.CName)
		if err != nil {
			fmt.Println("SetDNSRecords CNAME error:", err)
			http.Error(w, fmt.Sprintf(`{"error": "Failed to update CNAME: %s"}`, err.Error()), http.StatusInternalServerError)
			return
		}
	default:
		http.Error(w, `{"error": "Unsupported record type"}`, http.StatusBadRequest)
		return
	}

	// Update database with both values, but only the valid one will be changed
	if err := database.UpdateDNSInDB(body.Domain, body.ARecord, body.CName); err != nil {
		http.Error(w, `{"error": "Failed to update DNS in DB"}`, http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(map[string]any{
		"success": true,
		"message": "DNS updated successfully",
	})
}
