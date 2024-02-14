package utils

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
	"fmt"
	"net/http"
	"net/url"
	"strings"

	"github.com/pkg/errors"
)

// openSSL 을 사용해서 난수 생성
// ref. https://stackoverflow.com/questions/70254968/create-key-and-certificate-in-golang-same-as-openssl-do-for-local-host
func GenKeyID() (string, error) {
	key, err := rsa.GenerateKey(rand.Reader, 4096)
	if err != nil {
		return "", errors.Wrap(err, "GenerateKey failed")
	}

	keyBytes := x509.MarshalPKCS1PrivateKey(key)
	keyPEM := pem.EncodeToMemory(
		&pem.Block{
			Type:  "RSA PRIVATE KEY",
			Bytes: keyBytes,
		},
	)

	return string(keyPEM), nil
}

func CurlPOSTToVault(path string, vaultAuth string, data string) error {
	params := url.Values{}
	params.Add(data, ``)
	body := strings.NewReader(params.Encode())

	req, err := http.NewRequest("POST", path, body)
	if err != nil {
		fmt.Println(err)
		return errors.Wrap(err, "http request failed")
	}

	req.Header.Set("X-Vault-Token", vaultAuth)
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		fmt.Println(err)
		return errors.Wrap(err, "http client failed")
	}

	defer resp.Body.Close()
}
