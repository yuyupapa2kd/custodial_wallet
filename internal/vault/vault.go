package vault

import (
	"custodial-vault/configs"
	"custodial-vault/internal/logger"
	"fmt"
	"os"

	"github.com/go-resty/resty/v2"
)

type VaultService struct {
	client         *resty.Client
	basePath       string
	vaultAuthToken string
	healthCheck    string
}

var VaultSrv VaultService

var _lg *logger.Logger

func Init(lg *logger.Logger) *VaultService {
	_lg = lg

	VaultSrv.client = resty.New()
	VaultSrv.basePath = configs.RuntimeConf.RpcEndpoint + "/v1/secret/data/"
	VaultSrv.vaultAuthToken = os.Getenv("VaultAuthToken")
	VaultSrv.healthCheck = configs.RuntimeConf.RpcEndpoint + "/v1/sys/health"
	fmt.Println("VaultSrv.basePath : ", VaultSrv.basePath)
	fmt.Println("VaultSrv.vaultAuthToken : ", VaultSrv.vaultAuthToken)

	return &VaultSrv
}
