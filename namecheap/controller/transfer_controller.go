package controller

import (
	"encoding/json"

	"net/http"

	"namecheap-microservice/model"
	"namecheap-microservice/service"
)

func TransferDomainHandler(w http.ResponseWriter, r *http.Request) {
	var req model.DomainPurchaseRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, `{"error":"invalid input"}`, http.StatusBadRequest)
		return
	}

	client := service.NewNamecheapClient()
	ok, transferID, err := service.TransferDomain(client, req)
	if err != nil {
		http.Error(w, `{"error":"`+err.Error()+`"}`, http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(map[string]any{
		"success":     ok,
		"transfer_id": transferID,
	})
}
func GetTransferListHandler(w http.ResponseWriter, r *http.Request) {
	client := service.NewNamecheapClient()
	list, err := service.GetTransferList(client)
	if err != nil {
		http.Error(w, `{"error":"`+err.Error()+`"}`, http.StatusInternalServerError)
		return
	}
	json.NewEncoder(w).Encode(map[string]any{
		"success": true,
		"domains": list,
	})
}
func GetTransferStatusHandler(w http.ResponseWriter, r *http.Request) {
	domain := r.URL.Query().Get("domain")
	if domain == "" {
		http.Error(w, `{"error":"domain required"}`, http.StatusBadRequest)
		return
	}

	client := service.NewNamecheapClient()
	status, err := service.GetTransferStatus(client, domain)
	if err != nil {
		http.Error(w, `{"error":"`+err.Error()+`"}`, http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(map[string]string{
		"domain": domain,
		"status": status,
	})
}
