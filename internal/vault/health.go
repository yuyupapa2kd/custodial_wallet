package vault

import (
	"fmt"
	"time"

	"github.com/pkg/errors"
)

// ToDo : vault client 헬스체크
// curl -s -k -X GET http://127.0.0.1:8200/v1/sys/health
func (vs *VaultService) VaultHealthCheck() error {

	var vaultRespDTO VaultHealthResp

	resp, err := vs.client.R().
		SetResult(&vaultRespDTO).
		Get(vs.healthCheck)
	if err != nil {
		fmt.Printf("Fail to Connection Vault. Status Code: %d\n", resp.StatusCode())
		return errors.Wrap(err, "Fail to Connection Vault")
	}

	_lg.VaultHealthCheck(resp, vs.healthCheck, vaultRespDTO.Initialized, vaultRespDTO.Sealed, vaultRespDTO.ServerTimeUtc)
	if resp.StatusCode() != 200 {
		return errors.Wrap(errors.New("health check bad"), "Vault Server Bad Health")
	}

	return nil
}

func (vs *VaultService) HealthCheckLoop() {

	for {
		vs.VaultHealthCheck()
		time.Sleep(180 * time.Second)
	}

}

type VaultHealthResp struct {
	Initialized                bool   `json:"initialized"`
	Sealed                     bool   `json:"sealed"`
	Standby                    bool   `json:"standby"`
	PerformanceStandby         bool   `json:"performance_standby"`
	ReplicationPerformanceMode string `json:"replication_performance_mode"`
	ReplicationDrMode          string `json:"replication_dr_mode"`
	ServerTimeUtc              int    `json:"server_time_utc"`
	Version                    string `json:"version"`
	ClusterName                string `json:"cluster_name"`
	ClusterId                  string `json:"cluster_id"`
}

// // ToDo : ReConn 기능
// func (vs *VaultService) ReConnectToVault() error {

// }

// func ReConnVault() *VaultService {
// 	_client := resty.New()
// 	return &VaultService{_client}
// }

// ToDo : client reconnect 로직

// ToDo : 계정 생성 및 키 저장
// curl -s -k -X PATCH -H "X-Vault-Token: hvs.xs9zscrrZUyQ2UeYRpmQTP8r" -H "Content-Type: application/merge-patch+json" -d '{"data": {"eth":"kkkk"}}' http://127.0.0.1:8200/v1/secret/data/joseph

// vault 에 privKey 저장
// func (vs *VaultService) StoreKeyToVaultForEVM(keyID string, privKey string) error {
// 	payload := []byte(fmt.Sprintf(`{"data": {"evm": "%s"}}`, privKey))

// 	req, err := http.NewRequest("POST", vs.basePath+keyID, bytes.NewBuffer(payload))
// 	if err != nil {
// 		return errors.Wrap(err, "new http request failed")
// 	}

// 	req.Header.Set("X-Vault-Token", vs.vaultAuthToken)
// 	req.Header.Set("Content-Type", "application/json")

// 	client := &http.Client{}
// 	resp, err := client.Do(req)
// 	if err != nil {
// 		return errors.Wrap(err, "http client do failed")
// 	}
// 	defer resp.Body.Close()

// 	body, err := ioutil.ReadAll(resp.Body)
// 	if err != nil {
// 		return errors.Wrap(err, "read http body failed")
// 	}
// 	fmt.Println("Response Status : ", resp.Status)
// 	fmt.Println("Response Body : ", string(body))

// 	return nil
// }
