package main

import (
	"custodial-vault/configs"
	"custodial-vault/internal/api"
	"custodial-vault/internal/vault"
	"fmt"
)

func init() {

}

func main() {
	profile := "dev"
	configs.SetRuntimeConfig(profile)
	fmt.Println("======================================================")
	fmt.Println("connect to network type of ", profile)
	fmt.Println("======================================================")

	vault.Init()

	r := api.SetRouter()
	go r.Run(":" + configs.RuntimeConf.Server.Port)
	fmt.Println("server start")

	select {}
}
