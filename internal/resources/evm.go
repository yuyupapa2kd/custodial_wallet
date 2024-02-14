package resources

import (
	"math/big"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/holiman/uint256"
)

type LegacyTxnOptionalSig struct {
	Nonce    uint64
	GasPrice *big.Int
	Gas      uint64
	To       *common.Address `rlp:"nil"`
	Value    *big.Int
	Data     []byte
	V, R, S  *big.Int `rlp:"optional"`
}

type AccessListTxnOptionalSig struct {
	ChainID    *big.Int
	Nonce      uint64
	GasPrice   *big.Int
	Gas        uint64
	To         *common.Address `rlp:"nil"`
	Value      *big.Int
	Data       []byte
	AccessList types.AccessList
	V, R, S    *big.Int `rlp:"optional"`
}

type DynamicFeeTxnOptionalSig struct {
	ChainID    *big.Int
	Nonce      uint64
	GasTipCap  *big.Int
	GasFeeCap  *big.Int
	Gas        uint64
	To         *common.Address `rlp:"nil"`
	Value      *big.Int
	Data       []byte
	AccessList types.AccessList
	V, R, S    *big.Int `rlp:"optional"`
}

type BlobTxnOptionalSig struct {
	ChainID    *uint256.Int
	Nonce      uint64
	GasTipCap  *uint256.Int
	GasFeeCap  *uint256.Int
	Gas        uint64
	To         common.Address
	Value      *uint256.Int
	Data       []byte
	AccessList types.AccessList
	BlobFeeCap *uint256.Int
	BlobHashes []common.Hash
	Sidecar    *types.BlobTxSidecar `rlp:"-"`
	V, R, S    *big.Int             `rlp:"optional"`
}
