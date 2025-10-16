package mocks

import (
	"context"
)

type Wallet struct {
}
type SignTxParams struct {
	AccountID      string `json:"accountId"`
	ChainID        string `json:"chainId"`
	UnsignedTxData any    `json:"unsignedTxData"`
}

type SignTxRes struct {
	SignedTx string `json:"signedTx"`
}

type GetAddressPublicKeyParams struct {
	AccountId string `json:"accountId" validate:"required"`
	ChainId   string `json:"chainId" validate:"required"`
}

type GetAddressPublicKeyResponse struct {
	Address   string `json:"address"`
	PublicKey string `json:"publicKey"`
}

func (w *Wallet) SignTx(ctx context.Context, params *SignTxParams) (*SignTxRes, error) {
	return &SignTxRes{SignedTx: "signed"}, nil
}
func (w *Wallet) GetAddressPublicKey(ctx context.Context, params *GetAddressPublicKeyParams) (*GetAddressPublicKeyResponse, error) {
	return &GetAddressPublicKeyResponse{
		Address:   "0xac65a40b5a798e84B550Aebc611c2bE922541A0B",
		PublicKey: "024933f63ed21782d72722a5d668881baeda5599fc981c759bdcceb6827d9f4676",
	}, nil
}
