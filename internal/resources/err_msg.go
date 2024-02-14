package resources

type ErrMsgForVaultSrv struct {
	FAILED_CREATE_KEYID  string
	FAILED_STORE_PRIVKEY string
	FAILED_QUERY_PRIVKEY string
}

type ErrMsgForEVM struct {
	FAILED_CREATE_PRIVKEY_EVM        string
	FAILED_HEX_TO_ECDSA              string
	FAILED_PARSE_TXN_EVM             string
	FAILED_SIGN_TXN_EVM              string
	FAILED_MARSHAL_BINARY_SIGNED_TXN string
}

type ErrMsgForMiddleWare struct {
	NOT_AUTHENTICATED_IP string
}

var ErrMsgVault ErrMsgForVaultSrv
var ErrMsgEVM ErrMsgForEVM
var ErrMsgMW ErrMsgForMiddleWare

func SetErrMsgVault() {
	ErrMsgVault.FAILED_CREATE_KEYID = "random uuid generate failed"
	ErrMsgVault.FAILED_STORE_PRIVKEY = "store privKey to vault failed"
	ErrMsgVault.FAILED_QUERY_PRIVKEY = "GetPrivKeyFromVaultForEVM() failed"
}

func SetErrMsgEVM() {
	ErrMsgEVM.FAILED_CREATE_PRIVKEY_EVM = "create privKey failed"
	ErrMsgEVM.FAILED_HEX_TO_ECDSA = "privKeyStr to ECDSA failed"
	ErrMsgEVM.FAILED_PARSE_TXN_EVM = "parse txn failed"
	ErrMsgEVM.FAILED_SIGN_TXN_EVM = "signing txn failed"
	ErrMsgEVM.FAILED_MARSHAL_BINARY_SIGNED_TXN = "marshal binary txn failed"
}

func SetErrMsgMiddleWare() {
	ErrMsgMW.NOT_AUTHENTICATED_IP = "not authenticated ip"
}
