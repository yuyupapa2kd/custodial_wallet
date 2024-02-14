package resources

type ReqCoin struct {
	Address string `json:"address"`
}

type ReqToken struct {
	Symbol  string `json:"symbol"`
	Address string `json:"address"`
}
