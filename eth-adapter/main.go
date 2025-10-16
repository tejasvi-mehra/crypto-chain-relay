package main

import (
	"context"
	"eth-adapter/internal/framework/logger"
	abistore "eth-adapter/internal/service/abi_store"
	"eth-adapter/internal/service/block"
	blockrepo "eth-adapter/internal/service/block_repo"
	"eth-adapter/internal/service/publisher"
	relevantaddress "eth-adapter/internal/service/relevant_address.go"

	"github.com/ethereum/go-ethereum/ethclient"
)

func main() {
	log := logger.GetLogger("MAIN")
	blockNo := 9387584
	addressMap := map[string]string{
		"0x80b724ba8005e91a3a2bb9b28719f63ea886dd5c": "ce1268c5-9d05-46b4-8e69-59221d4e26f5",
		"0xac65a40b5a798e84b550aebc611c2be922541a0b": "ce1268c5-9d05-46b4-8e69-59221d4e26f3",
	}
	node := "https://ethereum-sepolia-rpc.publicnode.com"
	log.Debugw("processing block", "blockNo", blockNo, "addresses", addressMap, "node", node)
	// setup rpc
	rpcClient, err := ethclient.Dial(node)
	if err != nil {
		log.Fatal(err)
	}
	// setup abi store
	abiStore := abistore.NewMockAbiStore()
	// setup relevant addresses
	relevantAddress := relevantaddress.NewRelevantAddressMock(addressMap)
	// publisher
	publisher := &publisher.MockPublisher{}
	// blockrepo
	blockRepo := &blockrepo.MockBlockRepo{}
	// block detector
	blockDetector := block.NewBlockDetector(rpcClient, relevantAddress, publisher, blockRepo, abiStore)
	// processes block
	err = blockDetector.ProcessBlock(context.Background(), int64(blockNo))
	if err != nil {
		log.Panicw("error in process block", "err", err)
	}
}
