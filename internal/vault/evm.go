package vault

import (
	"fmt"

	"github.com/pkg/errors"
)

// type VaultService struct {
// 	client         *resty.Client
// 	basePath       string
// 	vaultAuthToken string
// 	healthCheck    string
// }

// var VaultSrv VaultService

// func Init() *VaultService {
// 	VaultSrv.client = resty.New()
// 	VaultSrv.basePath = configs.RuntimeConf.RpcEndpoint + "/v1/secret/data/"
// 	VaultSrv.vaultAuthToken = os.Getenv("VaultAuthToken")
// 	VaultSrv.healthCheck = "/v1/sys/health"
// 	fmt.Println("VaultSrv.basePath : ", VaultSrv.basePath)
// 	fmt.Println("VaultSrv.vaultAuthToken : ", VaultSrv.vaultAuthToken)

// 	return &VaultSrv
// }

func (vs *VaultService) StoreKeyToVaultForEVM(keyID string, privKey string) error {

	payload := map[string]interface{}{
		"data": map[string]string{
			"evm": privKey,
		},
	}

	resp, err := vs.client.R().
		SetHeader("X-Vault-Token", vs.vaultAuthToken).
		SetBody(payload).
		Post(vs.basePath + keyID)
	if err != nil {
		fmt.Printf("Fail to save privKey to vault. Status Code: %d\n", resp.StatusCode())
		return errors.Wrap(err, "POST to vault failed")
	}
	fmt.Println("vaultAuthToken : ", vs.vaultAuthToken)
	fmt.Println("used path : ", vs.basePath+keyID)

	fmt.Printf("Success to save privKey to vault. Status Code: %d\n", resp.StatusCode())
	return nil
}

// vault 에서 privKey 조회
func (vs *VaultService) GetPrivKeyFromVaultForEVM(keyID string) (string, error) {

	var vaultRespDTO map[string]interface{}

	resp, err := vs.client.R().
		SetHeader("X-Vault-Token", vs.vaultAuthToken).
		SetResult(&vaultRespDTO).
		Get(vs.basePath + keyID)

	if err != nil {
		fmt.Printf("Fail to get privKey from vault. Status Code: %d\n", resp.StatusCode())
		return "", errors.Wrap(err, "GET from vault failed")
	}

	fmt.Printf("Success to get privKey from vault. Status Code: %d\n", resp.StatusCode())

	value, ok := vaultRespDTO["data"].(map[string]interface{})["data"].(map[string]interface{})["evm"].(string)
	if !ok {
		fmt.Printf("Fail to parse privKey from vaultRespDTO")
		return "", errors.Wrap(errors.New("Fail to parse privKey from vaultRespDTO"), "parsing error")

	}

	return value, nil
}

// vault 에 privKey 추가 -> 추가 네트워크 확장 시에 고려할 기능
func (vs *VaultService) AddKeyToVaultForEVM(keyId string, privKey string) error {

	payload := map[string]interface{}{
		"data": map[string]string{
			"evm": privKey,
		},
	}

	resp, err := vs.client.R().
		SetHeader("X-Vault-Token", vs.vaultAuthToken).
		SetHeader("Content-Type", "application/merge-patch+json").
		SetBody(payload).
		Patch(vs.basePath + keyId)

	if err != nil {
		fmt.Printf("Fail to save privKey to vault. Status Code: %d\n", resp.StatusCode())
		return err
	}

	fmt.Printf("Success to save privKey to vault. Status Code: %d\n", resp.StatusCode())
	return nil

}
