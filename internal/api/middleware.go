package api

import (
	"custodial-vault/configs"
	"custodial-vault/internal/resources"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
)

func authCheck(c *gin.Context) {
	authToken := c.Request.Header.Get("auth-token")
	if authToken == "" {
		c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{})
		return
	}

	if authToken != os.Getenv("ServerAuthToken") {
		c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{})
	}
}

func clientIpCheck(c *gin.Context) {
	if c.ClientIP() != configs.RuntimeConf.WhiteIP {
		c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"result": nil, "errors": resources.ErrMsgMW.NOT_AUTHENTICATED_IP})
	}
}
