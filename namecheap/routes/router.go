package routes

import (
	"namecheap-microservice/controller"
	"namecheap-microservice/middleware"
	"net/http"
)

func SetupRoutes() {
	http.HandleFunc("/check-domain", middleware.CORSMiddleware(controller.CheckDomainHandler))
	http.HandleFunc("/buy-domain", middleware.CORSMiddleware(controller.BuyDomainHandler))
	http.HandleFunc("/domains", middleware.CORSMiddleware(controller.ListDomains))
	http.HandleFunc("/revoke-domain", middleware.CORSMiddleware(controller.RevokeDomainHandler))
	http.HandleFunc("/move-domain", middleware.CORSMiddleware(controller.MoveDomainHandler))

}
