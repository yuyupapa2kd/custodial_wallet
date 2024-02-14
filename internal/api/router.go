package api

import (
	"fmt"
	"net"
	"net/http"

	"github.com/gin-gonic/gin"
	swaggerFile "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"

	"custodial-vault/configs"
	"custodial-vault/docs"
	"custodial-vault/internal/logger"
	"custodial-vault/internal/resources"
	"custodial-vault/internal/vault"
)

var _lg *logger.Logger
var _vs *vault.VaultService

func Setup(r *gin.Engine, lg *logger.Logger, vs *vault.VaultService) {
	_lg = lg
	_vs = vs

	resources.SetHCode()
	resources.SetErrMsgVault()
	resources.SetErrMsgEVM()
	resources.SetErrMsgMiddleWare()

	setupSwagger(r)

	r.POST("/generateKey/evm", CreateKeyForEVM)
	r.GET("/signTxn/gnd/:keyId/:serializedTxn", GenSignedTxnForGND)

}

func getOutboundIP() net.IP {
	conn, err := net.Dial("udp", "8.8.8.8:80")
	if err != nil {
		panic(err)
	}
	defer conn.Close()

	localAddr := conn.LocalAddr().(*net.UDPAddr)

	return localAddr.IP
}

func setupSwagger(r *gin.Engine) {

	localAddr := getOutboundIP()
	IPNPortString := localAddr.String() + ":" + configs.RuntimeConf.Server.Port

	docs.SwaggerInfo.Title = "Custodial Wallet Service for Crypted.Inc."
	docs.SwaggerInfo.Description = "Create PrivKey and Signig Txn Service"
	docs.SwaggerInfo.Version = "0.1"
	docs.SwaggerInfo.Host = IPNPortString
	docs.SwaggerInfo.BasePath = ""
	docs.SwaggerInfo.Schemes = []string{"http", "https"}

	r.GET("/", func(c *gin.Context) {
		c.Redirect(http.StatusFound, "/swagger/index.html")
	})

	swaggerUrlString := "http://" + IPNPortString + "/swagger/doc.json"
	fmt.Println("swaggerUrlString : ", swaggerUrlString)
	// url := ginSwagger.URL("http://localhost:4000/swagger/doc.json")
	url := ginSwagger.URL(swaggerUrlString)
	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFile.Handler, url))
}
