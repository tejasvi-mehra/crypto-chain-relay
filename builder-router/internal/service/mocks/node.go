package mocks

import "math/big"

// this is just for reference, this package should be replaced with RPC
type Node struct {
}

func (n *Node) GetNonce() (uint64, error) {
	return 10, nil
}
func (n *Node) GetGasInfo() (gasLimit uint64, gasPrice *big.Int, err error) {
	return 30000, big.NewInt(20000000), nil
}
func (n *Node) BroadcastTx(signedTx string) (txHash string, err error) {
	return "0x8b55d12417775983294949b74058b3a468280b015ed2be6bb358bb99e0c256dc", nil
}

func (n *Node) GetFeeInfo() (feeRete string, err error) {
	return "20000000", nil
}
