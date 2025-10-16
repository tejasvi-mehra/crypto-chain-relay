package abistore

import (
	"strings"

	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
)

// ideally real one should read ABIs from file
type MockAbiStore struct {
	erc20ABI abi.ABI
}

func NewMockAbiStore() *MockAbiStore {
	erc20ABI, _ := abi.JSON(strings.NewReader(`[{"anonymous":false,"inputs":[{"indexed":true,"name":"from","type":"address"},{"indexed":true,"name":"to","type":"address"},{"indexed":false,"name":"value","type":"uint256"}],"name":"Transfer","type":"event"},{"constant":true,"inputs":[],"name":"symbol","outputs":[{"name":"","type":"string"}],"type":"function"}]`))
	return &MockAbiStore{
		erc20ABI: erc20ABI,
	}
}

func (a *MockAbiStore) UnpackIntoInterface(v interface{}, name string, data []byte) error {
	return a.erc20ABI.UnpackIntoInterface(v, name, data)

}
func (a *MockAbiStore) GetSignatureHashes() []common.Hash {
	erc20TransferEventSignature := []byte("Transfer(address,address,uint256)")
	transferSigHash := crypto.Keccak256Hash(erc20TransferEventSignature)
	return []common.Hash{transferSigHash}
}

func (a *MockAbiStore) Pack(name string, args ...interface{}) ([]byte, error) {
	return a.erc20ABI.Pack(name, args...)
}
