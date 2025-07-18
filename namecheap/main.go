package main

import (
	"fmt"
	"log"
	"namecheap-microservice/service"
	"net/http"
	"time"

	"namecheap-microservice/config"
	"namecheap-microservice/database"
	"namecheap-microservice/routes"
)

func main() {
	// Always load .env.localdemo first
	config.LoadConfig(".env.localdemo")

	// If ENV was set to production, load .env.production next
	if config.Env == "production" {
		config.LoadConfig(".env.production")
	}

	database.ConnectDB()

	go func() {
		for {
			service.UnrevokeOldDomains()
			time.Sleep(1 * time.Minute)
		}
	}()

	routes.SetupRoutes()

	port := "8080"
	fmt.Println("Server running on port " + port)

	log.Fatal(http.ListenAndServe(":"+port, nil))

}
