package controller

import (
	"encoding/json"
	"fmt"

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

	success, domainName, err := service.BuyDomain(client, req)

	if err != nil {
		fmt.Println("BuyDomain error:", err)
		http.Error(w, fmt.Sprintf(`{"error": "%s"}`, err.Error()), http.StatusInternalServerError)
		return
	}
	fmt.Println("✅ Domain purchase completed successfully:", domainName)
	if success {
		baseprice := 13.9
		tax := baseprice * 0.3
		total := baseprice + tax

		domain := model.DomainPurchase{
			Name:      domainName,
			Purchased: true,
			Customer:  req.Email,
			Price:     baseprice,
			Tax:       tax,
			Total:     total,
		}

		dns := model.DNSRecord{
			ARecord: "82.25.106.75",
			CName:   "indigo-spoonbill-233511.hostingersite.com",
		}

		if err := database.SaveDomainWithDNS(domain, dns); err != nil {
			fmt.Println("DB save error:", err)
			http.Error(w, `{"error": "failed to save domain and DNS info"}`, http.StatusInternalServerError)
			return
		}

		fmt.Println("✅ Domain saved to database successfully:", domain.Name)
	}

	json.NewEncoder(w).Encode(map[string]interface{}{
		"success": success,
		"domain":  domainName,
	})
}
func ListDomains(w http.ResponseWriter, r *http.Request) {
	var domains []model.DomainPurchase
	if err := database.DB.Find(&domains).Error; err != nil {
		http.Error(w, `{"error": "failed to fetch domains"}`, http.StatusInternalServerError)
		return
	}
	json.NewEncoder(w).Encode(domains)
}
func RevokeDomainHandler(w http.ResponseWriter, r *http.Request) {
	var body struct {
		Domain string `json:"domain"`
	}
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil || body.Domain == "" {
		http.Error(w, `{"error":"invalid domain"}`, http.StatusBadRequest)
		return
	}

	if err := service.RevokeDomain(body.Domain); err != nil {
		http.Error(w, fmt.Sprintf(`{"error":"%s"}`, err.Error()), http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(map[string]interface{}{
		"success": true,
		"message": "domain revoked",
	})
}
func MoveDomainHandler(w http.ResponseWriter, r *http.Request) {
	var body struct {
		Domain  string `json:"domain"`
		ARecord string `json:"aRecord"`
		CName   string `json:"cname"`
	}
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		http.Error(w, `{"error":"invalid body"}`, http.StatusBadRequest)
		return
	}

	client := service.NewNamecheapClient()
	if err := service.MoveDomain(client, body.Domain, body.ARecord, body.CName); err != nil {
		http.Error(w, fmt.Sprintf(`{"error":"%s"}`, err.Error()), http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(map[string]interface{}{
		"success": true,
		"message": "DNS updated successfully",
	})
}
