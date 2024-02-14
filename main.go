package main

import (
	"custodial-vault/configs"
	"custodial-vault/internal/logger"
	"custodial-vault/internal/server"
	"custodial-vault/internal/vault"
	"fmt"
)

func main() {
	profile := "dev"
	configs.SetRuntimeConfig(profile)
	fmt.Println("======================================================")
	fmt.Println("custodial vault server set configs for : ", profile)
	fmt.Println("======================================================")

	lg := logger.Init()
	vs := vault.Init(lg)

	go vs.HealthCheckLoop()

	r := server.Init(lg, vs)
	r.Run(":" + configs.RuntimeConf.Server.Port)

	fmt.Println("======================================================")
	fmt.Println("custodial vault server start")
	fmt.Println("======================================================")
}

// export VaultAuthToken=hvs.0jDfm1DP87zU7qaX12FEWXrm
