package resources

type ResPrivKeyGen struct {
	KeyID   string `json:"keyId"`
	Address string `json:"address"`
}

type ResSignedTxnEVM struct {
	SignedTxn string `json:"signedTxn"`
}
