package evm

import (
	"crypto/ecdsa"
	"custodial-vault/internal/resources"
	"fmt"
	"math/big"

	"github.com/holiman/uint256"
	"github.com/pkg/errors"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/rlp"
)

// go-ethereum 사용해서 키 생성. address 와 privKey 반환
func GenPrivKeyNAddrForEVM() (string, string, error) {
	privKey, err := crypto.GenerateKey()
	if err != nil {
		fmt.Println(err)
		return "", "", errors.Wrap(err, "GenerateKey() failed")
	}

	privKeyBytes := crypto.FromECDSA(privKey)
	privKeyStr := hexutil.Encode(privKeyBytes)
	// fmt.Println("privKeyBytes : ", hexutil.Encode(privKeyBytes[2:]))
	// fmt.Println("privKey : ", privKeyStr)

	pubKey := privKey.Public()
	pubKeyECDSA, ok := pubKey.(*ecdsa.PublicKey)
	if !ok {
		fmt.Println(err)
		return "", "", errors.Wrap(errors.New("publicKey is not of type *ecdsa.PublicKey"), "cannot assert type")
	}

	address := crypto.PubkeyToAddress(*pubKeyECDSA).Hex()
	fmt.Println("address : ", address)

	return privKeyStr, address, nil
}

// 직렬화된 트렌젝션 데이터를 type.Transaction Struct로 변환
func ParseTxnForEVM(serializedTxn string, chainID *big.Int) (*types.Transaction, error) {

	txBytes := common.FromHex(serializedTxn)

	// rlp decode 할때 서명값(r, s, v)이 필수이기 때문에 서명이 없는 serialized txn은 에러가 발생한다.
	// go-ethereum 의 types 패키지를 활용하여 typed Transaction을 생성하면 자동으로 default 서명이 들어가지만 (r:0, s:0, v:0)
	// 일반적으로 서명값에 디폴드값도 없는 트렌젝션을 받는 경우가 많을 것이기 때문에 한단계 더 거쳐서 rlp decode 해준다.

	// legacy Txn
	if len(txBytes) > 0 && txBytes[0] > 0x7f {
		var inner resources.LegacyTxnOptionalSig
		err := rlp.DecodeBytes(txBytes, &inner)
		if err != nil {
			return nil, err
		}
		return types.NewTx(&types.LegacyTx{
			Nonce:    inner.Nonce,
			GasPrice: inner.GasPrice,
			Gas:      inner.Gas,
			To:       inner.To,
			Value:    inner.Value,
			Data:     inner.Data,
		}), nil
	}

	// typed Txn
	if len(txBytes) <= 1 {
		return nil, fmt.Errorf("typed transaction too short")
	}
	switch txBytes[0] { // 0번째 인덱스에는 트렌젝션 타입에 대한 정보가 담겨있다.
	case types.AccessListTxType:
		var inner resources.AccessListTxnOptionalSig
		err := rlp.DecodeBytes(txBytes[1:], &inner)
		if err != nil {
			return nil, err
		}
		return types.NewTx(&types.AccessListTx{
			ChainID:    chainID,
			Nonce:      inner.Nonce,
			GasPrice:   inner.GasPrice,
			Gas:        inner.Gas,
			To:         inner.To,
			Value:      inner.Value,
			Data:       inner.Data,
			AccessList: inner.AccessList,
		}), nil

	case types.DynamicFeeTxType:
		var inner resources.DynamicFeeTxnOptionalSig
		err := rlp.DecodeBytes(txBytes[1:], &inner)
		if err != nil {
			return nil, err
		}
		return types.NewTx(&types.DynamicFeeTx{
			ChainID:    chainID,
			Nonce:      inner.Nonce,
			GasTipCap:  inner.GasTipCap,
			GasFeeCap:  inner.GasFeeCap,
			Gas:        inner.Gas,
			To:         inner.To,
			Value:      inner.Value,
			Data:       inner.Data,
			AccessList: inner.AccessList,
		}), nil

	case types.BlobTxType:
		var inner resources.BlobTxnOptionalSig
		err := rlp.DecodeBytes(txBytes[1:], &inner)
		if err != nil {
			return nil, err
		}
		blobChainID, _ := uint256.FromBig(chainID)
		return types.NewTx(&types.BlobTx{
			ChainID:    blobChainID,
			Nonce:      inner.Nonce,
			GasTipCap:  inner.GasTipCap,
			GasFeeCap:  inner.GasFeeCap,
			Gas:        inner.Gas,
			To:         inner.To,
			Value:      inner.Value,
			Data:       inner.Data,
			AccessList: inner.AccessList,
			BlobFeeCap: inner.BlobFeeCap,
			BlobHashes: inner.BlobHashes,
		}), nil

	default:
		return nil, fmt.Errorf("unsuppported transaction type")
	}

}
