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
	http.HandleFunc("/domains", middleware.CORSMiddleware(controller.ListDomains))
	http.HandleFunc("/revoke-domain", middleware.CORSMiddleware(controller.RevokeDomainHandler))

	http.HandleFunc(("/api/update-dns"), middleware.CORSMiddleware(controller.UpdateDNSHandler))
}
