package routes

import (
	"namecheap-microservice/controller"
	"namecheap-microservice/middleware"
	"net/http"
)

func SetupRoutes() {
	http.HandleFunc("/check-domain", middleware.CORSMiddleware(controller.CheckDomainHandler))
	http.HandleFunc("/buy-domain", middleware.CORSMiddleware(controller.BuyDomainHandler))
	http.HandleFunc(("/domain-price"), middleware.CORSMiddleware(controller.GetDomainPriceHandler))
	http.HandleFunc("/domains", middleware.CORSMiddleware(controller.GetLiveDomainsHandler))
	http.HandleFunc("/revoke-domain", middleware.CORSMiddleware(controller.RevokeDomainHandler))

	http.HandleFunc(("/api/update-dns"), middleware.CORSMiddleware(controller.UpdateDNSHandler))
	http.HandleFunc("/api/add-dns-record", middleware.CORSMiddleware(controller.AddDNSRecordHandler))
	http.HandleFunc("/api/registrar-lock", middleware.CORSMiddleware(controller.GetRegistrarLockHandler))
	http.HandleFunc("/api/set-registrar-lock", middleware.CORSMiddleware(controller.SetRegistrarLockHandler))
	http.HandleFunc("/api/tld-list", middleware.CORSMiddleware(controller.GetTLDListHandler))
	http.HandleFunc("/api/transfer-domain", middleware.CORSMiddleware(controller.TransferDomainHandler))
	http.HandleFunc("/api/transfer-list", middleware.CORSMiddleware(controller.GetTransferListHandler))
	http.HandleFunc("/api/transfer-status", controller.GetTransferStatusHandler)
	http.HandleFunc(("/api/reactivate-domain"), middleware.CORSMiddleware(controller.ReactivateDomainHandler))

}
