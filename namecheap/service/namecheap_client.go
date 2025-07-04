package service

import (
	"fmt"
	"namecheap-microservice/config"

	"github.com/namecheap/go-namecheap-sdk/v2/namecheap"
)

func NewNamecheapClient() *namecheap.Client {

	fmt.Printf("🌍 ENV: %s | Using Sandbox: %v\n", config.Env, config.UseSandbox)

	return namecheap.NewClient(&namecheap.ClientOptions{
		UserName:   config.ApiUser,
		ApiUser:    config.ApiUser,
		ApiKey:     config.ApiKey,
		ClientIp:   config.ClientIp,
		UseSandbox: config.UseSandbox,
	})
}
