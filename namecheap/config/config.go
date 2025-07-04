// config/config.go
package config

import (
	"fmt"
	"log"
	"os"

	"github.com/joho/godotenv"
)

var (
	Env        string
	ApiUser    string
	ApiKey     string
	ClientIp   string
	UseSandbox bool
)

func LoadConfig(envFile string) {
	err := godotenv.Load(envFile)
	if err != nil {
		log.Fatalf("Error loading env file: %v", err)
	}

	Env = os.Getenv("ENV")
	ApiUser = os.Getenv("NAMECHEAP_API_USER")
	ApiKey = os.Getenv("NAMECHEAP_API_KEY")
	UseSandbox = os.Getenv("NAMECHEAP_USE_SANDBOX") == "true"

	ClientIp = os.Getenv("CLIENT_IP")

	fmt.Println("ENV:", Env, "ClientIp:", ClientIp) // Debugging line to check loaded values
}
