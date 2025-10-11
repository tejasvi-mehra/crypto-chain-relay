package utils

type HDPathComponents struct {
	Purpose      uint32
	CoinType     uint32
	AccountIndex uint32
	ChangeIndex  uint32
	AddressIndex uint32
}

var HDPaths = map[string]HDPathComponents{
	"1":        {Purpose: 44, CoinType: 60, AccountIndex: 0, ChangeIndex: 0, AddressIndex: 0}, // Ethereum Mainnet
	"11155111": {Purpose: 44, CoinType: 60, AccountIndex: 0, ChangeIndex: 0, AddressIndex: 0}, // Sepolia testnet
	"signet":   {Purpose: 84, CoinType: 1, AccountIndex: 0, ChangeIndex: 0, AddressIndex: 0},  // Bitcoin signet testnet
	"bitcoin":  {Purpose: 84, CoinType: 0, AccountIndex: 0, ChangeIndex: 0, AddressIndex: 0},  // Bitcoin mainet
}

type ChainType string

var ETH ChainType = "eth"
var BTC ChainType = "btc"

var ChainMap = map[string]ChainType{
	"1":        ETH,
	"11155111": ETH,
	"signet":   BTC,
	"bitcoin":  BTC,
}
