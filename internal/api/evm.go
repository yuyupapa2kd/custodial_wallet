package api

import (
	"custodial-vault/configs"
	"custodial-vault/internal/evm"
	"custodial-vault/internal/resources"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/pkg/errors"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/crypto"
)

// CustodialWallet godoc
// @Tags PrivKey
// @Summary Generate PrivKey for EVM
// @Description Generate PrivKey for EVM
// @Produce       json
// @Success 200 {object} resources.ResJSON{data=resources.ResPrivKeyGen}
// @Failure 400 {object} resources.ResJSON{data=resources.ResErr}
// @Router /generateKey/evm [post]
func CreateKeyForEVM(c *gin.Context) {
	var res resources.ResPrivKeyGen
	// keyID 생성
	u, err := uuid.NewRandom()
	if err != nil {
		errW := errors.Wrap(err, resources.ErrMsgVault.FAILED_CREATE_KEYID)
		c.JSON(resources.HCode.FAILED_CREATE_KEYID, gin.H{"result": res, "errors": errW.Error()})
		_lg.LevelError(errW, c)
		return
	}
	keyID := u.String()
	// fmt.Println("keyID : ", keyID)

	// privKey 와 address 생성
	privKey, addr, err := evm.GenPrivKeyNAddrForEVM()
	if err != nil {
		errW := errors.Wrap(err, resources.ErrMsgEVM.FAILED_CREATE_PRIVKEY_EVM)
		c.JSON(resources.HCode.FAILED_CREATE_PRIVKEY, gin.H{"result": res, "errors": errW.Error()})
		_lg.LevelError(errW, c)
		return
	}

	// privKey 를 vault 에 저장
	err = _vs.StoreKeyToVaultForEVM(keyID, privKey)
	if err != nil {
		errW := errors.Wrap(err, resources.ErrMsgVault.FAILED_STORE_PRIVKEY)
		c.JSON(resources.HCode.FAILED_STORE_PRIVKEY, gin.H{"result": res, "errors": errW.Error()})
		_lg.LevelError(errW, c)
		return
	}

	res.KeyID = keyID
	res.Address = addr

	c.JSON(resources.HCode.CREATED, gin.H{"result": res, "errors": ""})
	_lg.LevelInfo(c)
	return
}

// CustodialWallet godoc
// @Tags SigingTxn
// @Summary Signing Txn for GNDChain
// @Description Signing Txn for GNDChain
// @Produce json
// @Param keyId path string true "keyId"
// @Param serializedTxn path string true "serializedTxn"
// @Success 200 {object} resources.ResJSON{data=resources.ResSignedTxnEVM}
// @Failure 400 {object} resources.ResJSON{data=resources.ResErr}
// @Router /signTxn/gnd/{keyId}/{serializedTxn} [get]
func GenSignedTxnForGND(c *gin.Context) {

	// keyID := c.MustGet("keyId").(string)
	// serializedTxn := c.MustGet("serializedTxn").(string)
	keyID := c.Param("keyId")
	serializedTxn := c.Param("serializedTxn")

	var res resources.ResSignedTxnEVM

	privKeyStr, err := _vs.GetPrivKeyFromVaultForEVM(keyID)
	if err != nil {
		errW := errors.Wrap(err, resources.ErrMsgVault.FAILED_QUERY_PRIVKEY)
		c.JSON(resources.HCode.FAILED_QUERY_PRIVKEY, gin.H{"result": res, "errors": errW.Error()})
		_lg.LevelError(errW, c)
		return
	}

	privKey, err := crypto.HexToECDSA(privKeyStr[2:])
	if err != nil {
		errW := errors.Wrap(err, resources.ErrMsgEVM.FAILED_HEX_TO_ECDSA)
		c.JSON(resources.HCode.FAILED_HEX_TO_ECDSA, gin.H{"result": res, "errors": errW.Error()})
		_lg.LevelError(errW, c)
		return
	}

	chainID := configs.RuntimeConf.EvmNetID.Gnd
	tx, err := evm.ParseTxnForEVM(serializedTxn, chainID)
	if err != nil {
		errW := errors.Wrap(err, resources.ErrMsgEVM.FAILED_PARSE_TXN_EVM)
		c.JSON(resources.HCode.FAILED_PARSE_TXN_EVM, gin.H{"result": res, "errors": errW.Error()})
		_lg.LevelError(errW, c)
		return
	}

	signedTxn, err := types.SignTx(tx, types.NewEIP155Signer(chainID), privKey)
	if err != nil {
		errW := errors.Wrap(err, resources.ErrMsgEVM.FAILED_SIGN_TXN_EVM)
		c.JSON(resources.HCode.FAILED_SIGN_TXN_EVM, gin.H{"result": res, "errors": errW.Error()})
		_lg.LevelError(errW, c)
		return
	}

	byteSignedTxn, err := signedTxn.MarshalBinary()
	if err != nil {
		errW := errors.Wrap(err, resources.ErrMsgEVM.FAILED_MARSHAL_BINARY_SIGNED_TXN)
		c.JSON(resources.HCode.FAILED_MARSHAL_BINARY_SIGNED_TXN, gin.H{"result": res, "errors": errW.Error()})
		_lg.LevelError(errW, c)
		return
	}

	res.SignedTxn = "0x" + common.Bytes2Hex(byteSignedTxn)
	c.JSON(resources.HCode.OK, gin.H{"result": res, "errors": ""})
	_lg.LevelInfo(c)
	return

}

// // 직렬화된 트렌젝션 데이터를 type.Transaction Struct로 변환
// func parseTxnForEVM(serializedTxn string, chainID *big.Int) (*types.Transaction, error) {

// 	txBytes := common.FromHex(serializedTxn)

// 	// rlp decode 할때 서명값(r, s, v)이 필수이기 때문에 서명이 없는 serialized txn은 에러가 발생한다.
// 	// go-ethereum 의 types 패키지를 활용하여 typed Transaction을 생성하면 자동으로 default 서명이 들어가지만 (r:0, s:0, v:0)
// 	// 일반적으로 서명값에 디폴드값도 없는 트렌젝션을 받는 경우가 많을 것이기 때문에 한단계 더 거쳐서 rlp decode 해준다.

// 	// legacy Txn
// 	if len(txBytes) > 0 && txBytes[0] > 0x7f {
// 		var inner resources.LegacyTxnOptionalSig
// 		err := rlp.DecodeBytes(txBytes, &inner)
// 		if err != nil {
// 			return nil, err
// 		}
// 		return types.NewTx(&types.LegacyTx{
// 			Nonce:    inner.Nonce,
// 			GasPrice: inner.GasPrice,
// 			Gas:      inner.Gas,
// 			To:       inner.To,
// 			Value:    inner.V,
// 			Data:     inner.Data,
// 		}), nil
// 	}

// 	// typed Txn
// 	if len(txBytes) <= 1 {
// 		return nil, fmt.Errorf("typed transaction too short")
// 	}
// 	switch txBytes[0] { // 0번째 인덱스에는 트렌젝션 타입에 대한 정보가 담겨있다.
// 	case types.AccessListTxType:
// 		var inner resources.AccessListTxnOptionalSig
// 		err := rlp.DecodeBytes(txBytes[1:], &inner)
// 		if err != nil {
// 			return nil, err
// 		}
// 		return types.NewTx(&types.AccessListTx{
// 			ChainID:    chainID,
// 			Nonce:      inner.Nonce,
// 			GasPrice:   inner.GasPrice,
// 			Gas:        inner.Gas,
// 			To:         inner.To,
// 			Value:      inner.Value,
// 			Data:       inner.Data,
// 			AccessList: inner.AccessList,
// 		}), nil

// 	case types.DynamicFeeTxType:
// 		var inner resources.DynamicFeeTxnOptionalSig
// 		err := rlp.DecodeBytes(txBytes[1:], &inner)
// 		if err != nil {
// 			return nil, err
// 		}
// 		return types.NewTx(&types.DynamicFeeTx{
// 			ChainID:    chainID,
// 			Nonce:      inner.Nonce,
// 			GasTipCap:  inner.GasTipCap,
// 			GasFeeCap:  inner.GasFeeCap,
// 			Gas:        inner.Gas,
// 			To:         inner.To,
// 			Value:      inner.Value,
// 			Data:       inner.Data,
// 			AccessList: inner.AccessList,
// 		}), nil

// 	case types.BlobTxType:
// 		var inner resources.BlobTxnOptionalSig
// 		err := rlp.DecodeBytes(txBytes[1:], &inner)
// 		if err != nil {
// 			return nil, err
// 		}
// 		blobChainID, _ := uint256.FromBig(chainID)
// 		return types.NewTx(&types.BlobTx{
// 			ChainID:    blobChainID,
// 			Nonce:      inner.Nonce,
// 			GasTipCap:  inner.GasTipCap,
// 			GasFeeCap:  inner.GasFeeCap,
// 			Gas:        inner.Gas,
// 			To:         inner.To,
// 			Value:      inner.Value,
// 			Data:       inner.Data,
// 			AccessList: inner.AccessList,
// 			BlobFeeCap: inner.BlobFeeCap,
// 			BlobHashes: inner.BlobHashes,
// 		}), nil

// 	default:
// 		return nil, fmt.Errorf("unsuppported transaction type")
// 	}

// }

// func QueryPrivKeyForEVM(keyID string) (string, error) {
// 	privKey, err := vault.VaultSrv.GetPrivKeyFromVaultForEVM(keyID)
// 	if err != nil {
// 		return "", err
// 	}

// 	return privKey, nil
// }
