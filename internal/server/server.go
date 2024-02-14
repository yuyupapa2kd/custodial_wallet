package server

import (
	"custodial-vault/internal/api"
	"custodial-vault/internal/logger"
	"custodial-vault/internal/vault"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func Init(logger *logger.Logger, vaultService *vault.VaultService) *gin.Engine {
	router := gin.Default()

	router.Use(corsMiddleware())

	router.MaxMultipartMemory = 8 << 20 // 8MiB

	api.Setup(router, logger, vaultService)

	return router
}

func corsMiddleware() gin.HandlerFunc {
	config := cors.DefaultConfig()
	config.AddAllowHeaders("Authorization")
	config.AllowOriginFunc = func(origin string) bool {
		return true
	}

	return cors.New(config)
}
