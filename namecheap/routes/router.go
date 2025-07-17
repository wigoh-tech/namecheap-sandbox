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
	http.HandleFunc("/api/add-dns-record", controller.AddDNSRecordHandler)
	http.HandleFunc("/api/registrar-lock", controller.GetRegistrarLockHandler)
	http.HandleFunc("/api/set-registrar-lock", controller.SetRegistrarLockHandler)
	http.HandleFunc("/api/tld-list", controller.GetTLDListHandler)
	http.HandleFunc("/api/transfer-domain", controller.TransferDomainHandler)
	http.HandleFunc("/api/transfer-list", controller.GetTransferListHandler)
	http.HandleFunc("/api/transfer-status", controller.GetTransferStatusHandler)
	http.HandleFunc(("/api/reactivate-domain"), middleware.CORSMiddleware(controller.ReactivateDomainHandler))

}
