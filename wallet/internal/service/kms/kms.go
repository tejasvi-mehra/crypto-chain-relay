package kms

import (
	"crypto/ecdsa"
	"fmt"
	"math/big"
	"wallet/internal/framework/logger"
	"wallet/internal/service/utils"

	"strconv"

	"github.com/btcsuite/btcd/btcutil/hdkeychain"
	"github.com/btcsuite/btcd/btcutil/psbt"
	"github.com/btcsuite/btcd/chaincfg"
	"github.com/caarlos0/env/v6"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/google/uuid"
	"github.com/tyler-smith/go-bip32"
	"github.com/tyler-smith/go-bip39"
)

// ideally this should be separate and in its own secure machine

type KMSConfig struct {
	SeedNameSpace string `env:"SEED_NAMESPACE,required"`
}

func NewKMSConfig() (*KMSConfig, error) {
	c := &KMSConfig{}
	if err := env.Parse(c); err != nil {
		return nil, err
	}

	return c, nil
}

type KMS struct {
	seedNameSpace uuid.UUID
}

func NewKMS(c *KMSConfig) (*KMS, error) {
	seedNameSpace, err := uuid.Parse(c.SeedNameSpace)
	if err != nil {
		return nil, err
	}
	return &KMS{
		seedNameSpace: seedNameSpace,
	}, nil
}

func (k *KMS) DeriveSeedUUID(accountID uuid.UUID) uuid.UUID {
	return uuid.NewSHA1(k.seedNameSpace, accountID[:])
}

func (k *KMS) DerivePublicKeyAndAddress(seedID uuid.UUID, chainID string, hdPath utils.HDPathComponents) (string, string, error) {
	log := logger.GetLogger("DerivePublicKeyAndAddress")
	// check chain type
	chainType, ok := utils.ChainMap[chainID]
	if !ok {
		log.Errorw("invalid chainID, not found in map", "chainID", chainID, "map", utils.ChainMap)
		return "", "", fmt.Errorf("invalid chainID")
	}
	var address, pubKey string
	var err error
	switch chainType {
	case utils.ETH:
		// create eth address
		address, pubKey, err = k.generateETHPublicKeyAndAddress(seedID, hdPath)
	case utils.BTC:
		// create btc address
		address, pubKey, err = k.generateBTCPublicKeyAndAddress(seedID, hdPath, chainID)
	default:
		log.Errorw("invalid chainType, not found in map", "chainID", chainID, "map", utils.ChainMap)
		return "", "", fmt.Errorf("invalid chainType")
	}
	return address, pubKey, err
}

func (k *KMS) SignTx(seedID uuid.UUID, chainID string, hdPath utils.HDPathComponents, tx any) ([]byte, error) {
	log := logger.GetLogger("DerivePublicKeyAndAddress")
	// check chain type
	chainType, ok := utils.ChainMap[chainID]
	if !ok {
		log.Errorw("invalid chainID, not found in map", "chainID", chainID, "map", utils.ChainMap)
		return nil, fmt.Errorf("invalid chainID")
	}
	switch chainType {
	case utils.ETH:
		// eth signing
		ethTx, ok := tx.(*types.Transaction)
		if !ok {
			log.Errorw("invalid ETH tx for signing", "tx", tx)
			return nil, fmt.Errorf("invalid ETH tx for signing")
		}
		return k.signEthTx(seedID, chainID, hdPath, ethTx)
	case utils.BTC:
		// btc signing to;do
		btcPsbt, ok := tx.(*psbt.Packet)
		if !ok {
			log.Errorw("invalid BTC PSBT for signing", "tx", tx)
			return nil, fmt.Errorf("invalid BTC PSBT for signing")
		}
		return k.signBTCPsbt(seedID, chainID, hdPath, btcPsbt)
	default:
		log.Errorw("invalid chainType, not found in map", "chainID", chainID, "map", utils.ChainMap)
		return nil, fmt.Errorf("invalid chainType")
	}
}

func (k *KMS) generateETHPublicKeyAndAddress(seedID uuid.UUID, hdPath utils.HDPathComponents) (string, string, error) {
	privateKey, err := k.generateETHPrivateKeyFromSeed(seedID, hdPath)
	if err != nil {
		return "", "", err
	}
	publicKey := privateKey.Public().(*ecdsa.PublicKey)
	ethAddress := crypto.PubkeyToAddress(*publicKey)

	pubKeyBytes := crypto.FromECDSAPub(publicKey)
	return ethAddress.Hex(), fmt.Sprintf("0x%x", pubKeyBytes), nil
}

func (k *KMS) generateETHPrivateKeyFromSeed(seedUUID uuid.UUID, hdPath utils.HDPathComponents) (*ecdsa.PrivateKey, error) {
	seed, err := k.generateSeed(seedUUID)
	if err != nil {
		return nil, err
	}
	masterKey, err := bip32.NewMasterKey(seed)
	if err != nil {
		return nil, fmt.Errorf("failed to generate master key")
	}
	purpose, _ := masterKey.NewChildKey(bip32.FirstHardenedChild + hdPath.Purpose)
	coinType, _ := purpose.NewChildKey(bip32.FirstHardenedChild + hdPath.CoinType)
	account, _ := coinType.NewChildKey(bip32.FirstHardenedChild + hdPath.AccountIndex)
	change, _ := account.NewChildKey(hdPath.ChangeIndex)
	addressIndex, _ := change.NewChildKey(hdPath.AddressIndex)

	privateKey, err := crypto.ToECDSA(addressIndex.Key)
	if err != nil {
		return nil, fmt.Errorf("failed to make privateKey: %v", err)
	}
	return privateKey, nil
}

func (k *KMS) generateBTCPublicKeyAndAddress(seedID uuid.UUID, hdPath utils.HDPathComponents, chainID string) (string, string, error) {
	seed, err := k.generateSeed(seedID)
	if err != nil {
		return "", "", err
	}
	// Bitcoin master key using hdkeychain from btcsuite
	masterKey, err := hdkeychain.NewMaster(seed, &chaincfg.SigNetParams)
	if err != nil {
		return "", "", fmt.Errorf("failed new master: %w", err)
	}
	purpose, err := masterKey.Derive(hdkeychain.HardenedKeyStart + hdPath.Purpose)
	if err != nil {
		return "", "", err
	}
	coinType, err := purpose.Derive(hdkeychain.HardenedKeyStart + hdPath.CoinType)
	if err != nil {
		return "", "", err
	}
	account, err := coinType.Derive(hdkeychain.HardenedKeyStart + hdPath.AccountIndex)
	if err != nil {
		return "", "", err
	}
	change, err := account.Derive(hdPath.ChangeIndex)
	if err != nil {
		return "", "", err
	}
	addressIndex, err := change.Derive(0) // first address at index 0
	if err != nil {
		return "", "", err
	}

	// Get the associated public key serialized
	pubKey, err := addressIndex.ECPubKey()
	if err != nil {
		return "", "", err
	}
	pubKeyBytes := pubKey.SerializeCompressed()

	pubKeyHex := fmt.Sprintf("%x", pubKeyBytes)

	// Generate native segwit (Bech32) address
	chainCfg := chaincfg.MainNetParams
	if chainID == "signet" {
		chainCfg = chaincfg.SigNetParams
	}
	addr, err := addressIndex.Address(&chainCfg)
	if err != nil {
		return "", "", err
	}

	btcAddr := addr.EncodeAddress()
	return btcAddr, pubKeyHex, nil
}

func (k *KMS) generateSeed(seedUUID uuid.UUID) ([]byte, error) {
	entropy := seedUUID[:]
	// Optional: append a fixed salt for extra entropy variability (repeatable)
	// salt := []byte("fixed_salt")
	// seed = append(entropy, salt...)

	mnemonic, err := bip39.NewMnemonic(entropy)
	if err != nil {
		return nil, err
	}
	seed := bip39.NewSeed(mnemonic, "")
	return seed, nil
}

func (k *KMS) signEthTx(seedID uuid.UUID, chainID string, hdPath utils.HDPathComponents, tx *types.Transaction) ([]byte, error) {
	privateKey, err := k.generateETHPrivateKeyFromSeed(seedID, hdPath)
	if err != nil {
		return nil, err
	}

	chainIDInt, err := strconv.ParseInt(chainID, 0, 64)
	if err != nil {
		return nil, fmt.Errorf("invalid chainID cannot convert to int: %v", err)
	}
	chainIDBig := big.NewInt(chainIDInt)

	signedTx, err := types.SignTx(tx, types.NewEIP155Signer(chainIDBig), privateKey)
	if err != nil {
		return nil, fmt.Errorf("failed to sign tx: %v", err)
	}

	signedTxBytes, err := signedTx.MarshalBinary()
	if err != nil {
		return nil, fmt.Errorf("failed to marshal sign tx: %v", err)
	}
	return signedTxBytes, err
}

func (k *KMS) signBTCPsbt(seedID uuid.UUID, chainID string, hdPath utils.HDPathComponents, btcPsbt *psbt.Packet) ([]byte, error) {
	// essentially we will fetch the master key the same way as we get for address generation and use that to sign - however not sure if we need additional inputs over just the psbt here
	return []byte("signed psbt"), nil
}
