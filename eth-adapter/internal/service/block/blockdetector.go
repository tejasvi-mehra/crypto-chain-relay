package block

import (
	"context"
	"eth-adapter/internal/config"
	db "eth-adapter/internal/database/generated"
	swagger "eth-adapter/internal/event/generated/1.0.0/transaction-events/go"
	"eth-adapter/internal/framework/logger"
	"eth-adapter/internal/service/trace_id"
	"fmt"
	"math/big"
	"time"

	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"golang.org/x/exp/slices"
)

type RpcClient interface {
	BlockByNumber(ctx context.Context, number *big.Int) (*types.Block, error)
	TransactionReceipt(ctx context.Context, txHash common.Hash) (*types.Receipt, error)
	CallContract(ctx context.Context, msg ethereum.CallMsg, blockNumber *big.Int) ([]byte, error)
}

type RelevantAddress interface {
	IsRelevantAddress(ctx context.Context, address string) (relevant bool, accountID string)
}

type BlockRepo interface {
	SaveBlock(ctx context.Context, block *db.BlockData) error
}

type PublishTx interface {
	PublishTxs(ctx context.Context, txs []swagger.TxEventData) error
}

type ABIStore interface {
	UnpackIntoInterface(v interface{}, name string, data []byte) error
	GetSignatureHashes() []common.Hash
	Pack(name string, args ...interface{}) ([]byte, error)
}

type BlockDetector struct {
	rpcClient RpcClient
	address   RelevantAddress
	publisher PublishTx
	blockRepo BlockRepo
	abiStore  ABIStore
}

func NewBlockDetector(
	rpcClient RpcClient,
	address RelevantAddress,
	publisher PublishTx,
	blockRepo BlockRepo,
	abiStore ABIStore,
) *BlockDetector {
	return &BlockDetector{
		rpcClient: rpcClient,
		address:   address,
		publisher: publisher,
		blockRepo: blockRepo,
		abiStore:  abiStore,
	}
}

func (b *BlockDetector) ProcessBlock(ctx context.Context, blockNo int64) error {
	log := logger.GetLogger("ProcessBlock")
	defer func() {
		log.Debugw("completed")
	}()
	log.With("blockNo", blockNo)
	// parse blockNo
	number := big.NewInt(blockNo)
	// fetch block from rpc client
	block, err := b.rpcClient.BlockByNumber(ctx, number)
	if err != nil {
		log.Errorw("failed to fetch block from rpc", "err", err)
		return err
	}
	// blockInfo
	blockHash := block.Hash().Hex()
	parentHash := block.ParentHash().Hex()
	blockTimeEpochSec := block.Time()
	blockTimestamp := time.Unix(int64(blockTimeEpochSec), 0)
	// start time for processing
	startTime := time.Now().UTC()
	// process block txs
	txs := make([]swagger.TxEventData, 0)
	for _, tx := range block.Transactions() {
		// init vars
		var fromAddress, toAddress, txHash, withdrawalAccountID, depositAccountID, symbol, amountDecimal string
		var amountInfo *swagger.TransferAmountInfo
		var isWithdrawal, isDeposit, isRelevant bool
		// Recover sender
		from, err := types.Sender(types.LatestSignerForChainID(tx.ChainId()), tx)
		if err != nil {
			log.Errorw("Failed to recover sender for tx", "txhash", tx.Hash().Hex(), "err", err)
			continue
		}
		fromAddress = from.Hex()
		// get txHash
		txHash = tx.Hash().Hex()
		// check need to process withdrawal
		isRelevant, withdrawalAccountID = b.address.IsRelevantAddress(ctx, fromAddress)
		if isRelevant {
			// process withdrawal
			isWithdrawal = true
		}
		// get tx receipt
		receipt, err := b.rpcClient.TransactionReceipt(ctx, tx.Hash())
		if err != nil {
			log.Errorw("Failed to get receipt for tx", txHash, err)
			continue
		}
		// ETH value > 0 -> native tx
		if tx.Value().Cmp(big.NewInt(0)) > 0 {
			// save to address and amount
			toAddress = tx.To().Hex()
			isRelevant, depositAccountID = b.address.IsRelevantAddress(ctx, toAddress)
			if isRelevant {
				// process deposit
				isDeposit = true
			}
			amountDecimal = new(big.Float).Quo(new(big.Float).SetInt(tx.Value()), big.NewFloat(1e18)).String()
			amountInfo = &swagger.TransferAmountInfo{
				Ticker: "ETH",
				// Atomic Amount in DEC
				ValueAtomic: tx.Value().String(),
				// Amount for UI
				ValueDecimal: amountDecimal,
				// Decimal of the asset
				Decimals: 18,
				// spec of asset transferred, native or token standard.
				Standard: "native",
			}
			symbol = "ETH"
		} else {
			// check erc20
			for _, vLog := range receipt.Logs {
				// Filter logs that have first topic Transfer event signature
				if len(vLog.Topics) > 0 && slices.Contains(b.abiStore.GetSignatureHashes(), vLog.Topics[0]) {
					// Decode indexed 'from' and 'to' from topics
					from := common.HexToAddress(vLog.Topics[1].Hex())
					to := common.HexToAddress(vLog.Topics[2].Hex())

					fromAddress = from.Hex()
					toAddress = to.Hex()
					log.Debugw("address data", "from", fromAddress, "to", toAddress)
					isRelevant, depositAccountID = b.address.IsRelevantAddress(ctx, toAddress)
					if isRelevant {
						// process deposit
						isDeposit = true
					}
					if !isDeposit && !isWithdrawal {
						// not withdraw or deposit just continue
						continue
					}

					// Decode value (amount) from data
					var transferEvent struct {
						Value *big.Int
					}
					err := b.abiStore.UnpackIntoInterface(&transferEvent, "Transfer", vLog.Data)
					if err != nil {
						log.Errorw("Failed to unpack transfer event data", "hash", txHash, "err", err)
						continue
					}
					amountDecimal = new(big.Float).Quo(new(big.Float).SetInt(transferEvent.Value), big.NewFloat(1e18)).String()
					// get symbol
					contractAddress := vLog.Address
					callData, err := b.abiStore.Pack("symbol")
					if err != nil {
						log.Errorw("Failed to pack symbol call", "hash", txHash, "err", err)
						continue
					}

					res, err := b.rpcClient.CallContract(ctx, ethereum.CallMsg{
						To:   &contractAddress,
						Data: callData,
					}, nil)
					if err != nil {
						log.Errorw("symbol() call failed for contract", "contractAddress", contractAddress.Hex(), "err", err)
						continue
					}

					err = b.abiStore.UnpackIntoInterface(&symbol, "symbol", res)
					if err != nil {
						log.Errorw("Failed to unpack symbol response", "err", err)
						continue
					}
					amountInfo = &swagger.TransferAmountInfo{
						Ticker: symbol,
						// Atomic Amount in DEC
						ValueAtomic: tx.Value().String(),
						// Amount for UI
						ValueDecimal: amountDecimal,
						// Decimal of the asset
						Decimals: 18,
						// spec of asset transferred, native or token standard.
						Standard: "native",
						// contract address
						ContractAddress: contractAddress.Hex(),
					}

				}
			}
		}

		// check if need to process tx
		if isWithdrawal {
			// make withdrawal object
			statusDetected := swagger.DETECTED_Status
			withdrawalTxType := swagger.WITHDRAWAL_TransactionType
			// make ulid
			traceID := trace_id.UlidHashData{
				ChainID:         config.CHAIN_ID,
				Block:           fmt.Sprintf("%d", blockNo),
				BlockTimestamp:  blockTimestamp.Format(time.RFC3339),
				Hash:            txHash,
				ActionAddress:   fromAddress,
				TransactionType: string(withdrawalTxType),
				Symbol:          symbol,
				Value:           amountDecimal,
			}.GenerateULID()
			specific := &swagger.CryptoSpecificInfo{
				Blockchaintype: "evm",
				ChainID:        config.CHAIN_ID,
				Block:          int32(blockNo),
				BlockTimestamp: blockTimestamp.Format(time.RFC3339),
				FromAddress:    fromAddress,
				ToAddress:      toAddress,
				AmountInfo:     amountInfo,
				TxHash:         txHash,
			}
			withdrawalTx := swagger.TxEventData{
				Specversion:        "1.0",
				Source:             "eth-adapter",
				Datacontenttype:    "application/json",
				EventType:          "crypto",
				TraceID:            traceID,
				Status:             &statusDetected,
				TransactionType:    &withdrawalTxType,
				From:               fromAddress,
				To:                 toAddress,
				AccountID:          withdrawalAccountID,
				AmountDecimal:      amountDecimal,
				Symbol:             symbol,
				TransactionTime:    blockTimestamp.Format(time.RFC3339),
				CreatedAt:          time.Now().UTC().String(),
				CryptoSpecificInfo: specific,
			}
			// append withdrawal
			txs = append(txs, withdrawalTx)
			// process network fee event -
			// Gas price extraction considering EIP-1559 transactions
			var gasPrice *big.Int
			if tx.Type() == types.DynamicFeeTxType {
				// For EIP-1559 tx, effective gas price = min(maxFeePerGas, baseFee + maxPriorityFeePerGas)
				// But baseFee is dynamic, so from receipt or block header can get it. For simplicity use tx.GasFeeCap()
				gasPrice = tx.GasFeeCap()
			} else {
				gasPrice = tx.GasPrice()
			}
			actualGasUsed := receipt.GasUsed
			gasCost := new(big.Int).Mul(new(big.Int).SetUint64(actualGasUsed), gasPrice)
			gasCostEth := new(big.Float).Quo(new(big.Float).SetInt(gasCost), big.NewFloat(1e18))
			gasAmount := &swagger.TransferAmountInfo{
				Ticker:       "ETH",
				ValueAtomic:  gasCost.String(),
				ValueDecimal: gasCostEth.String(),
				Decimals:     18,
				Standard:     "native",
			}
			networkFeeTxType := swagger.NETWORK_FEE_TransactionType
			// make ulid
			traceIDNetworkFee := trace_id.UlidHashData{
				ChainID:         config.CHAIN_ID,
				Block:           fmt.Sprintf("%d", blockNo),
				BlockTimestamp:  blockTimestamp.Format(time.RFC3339),
				Hash:            txHash,
				ActionAddress:   fromAddress,
				TransactionType: string(withdrawalTxType),
				Symbol:          "ETH",
				Value:           gasCost.String(),
			}.GenerateULID()
			specificNetworkFee := &swagger.CryptoSpecificInfo{
				Blockchaintype: "evm",
				ChainID:        config.CHAIN_ID,
				Block:          int32(blockNo),
				BlockTimestamp: blockTimestamp.Format(time.RFC3339),
				FromAddress:    fromAddress,
				ToAddress:      toAddress,
				AmountInfo:     gasAmount,
				TxHash:         txHash,
			}
			networkFeeTx := swagger.TxEventData{
				Specversion:        "1.0",
				Source:             "eth-adapter",
				Datacontenttype:    "application/json",
				EventType:          "crypto",
				TraceID:            traceIDNetworkFee,
				Status:             &statusDetected,
				TransactionType:    &networkFeeTxType,
				From:               fromAddress,
				To:                 toAddress,
				AccountID:          withdrawalAccountID,
				AmountDecimal:      gasCostEth.String(),
				Symbol:             "ETH",
				TransactionTime:    blockTimestamp.Format(time.RFC3339),
				CreatedAt:          time.Now().UTC().String(),
				CryptoSpecificInfo: specificNetworkFee,
			}
			txs = append(txs, networkFeeTx)
		}
		if isDeposit {
			// make withdrawal object
			statusDetected := swagger.DETECTED_Status
			depositTxType := swagger.DEPOSIT_TransactionType
			// make ulid
			traceID := trace_id.UlidHashData{
				ChainID:         config.CHAIN_ID,
				Block:           fmt.Sprintf("%d", blockNo),
				BlockTimestamp:  blockTimestamp.Format(time.RFC3339),
				Hash:            txHash,
				ActionAddress:   fromAddress,
				TransactionType: string(depositTxType),
				Symbol:          symbol,
				Value:           amountDecimal,
			}.GenerateULID()
			specific := &swagger.CryptoSpecificInfo{
				Blockchaintype: "evm",
				ChainID:        config.CHAIN_ID,
				Block:          int32(blockNo),
				BlockTimestamp: blockTimestamp.Format(time.RFC3339),
				FromAddress:    fromAddress,
				ToAddress:      toAddress,
				AmountInfo:     amountInfo,
				TxHash:         txHash,
			}
			depositTx := swagger.TxEventData{
				Specversion:        "1.0",
				Source:             "eth-adapter",
				Datacontenttype:    "application/json",
				EventType:          "crypto",
				TraceID:            traceID,
				Status:             &statusDetected,
				TransactionType:    &depositTxType,
				From:               fromAddress,
				To:                 toAddress,
				AccountID:          depositAccountID,
				AmountDecimal:      amountDecimal,
				Symbol:             symbol,
				TransactionTime:    blockTimestamp.Format(time.RFC3339),
				CreatedAt:          time.Now().UTC().String(),
				CryptoSpecificInfo: specific,
			}
			// append deposit
			txs = append(txs, depositTx)
		}
	}
	blockProcessTime := time.Since(startTime).Microseconds()
	log.Debugw("done processing txs", "count", len(txs), "txs", txs)

	dbBlock := &db.BlockData{
		BlockNo:        blockNo,
		BlockHash:      blockHash,
		ParentHash:     parentHash,
		BlockTimestamp: blockTimestamp,
		ProcessTimeMs:  blockProcessTime,
		Confirmed:      false,
		RetryCount:     0,
		CreatedAt:      time.Now().UTC(),
		UpdatedAt:      time.Now().UTC(),
	}
	err = b.blockRepo.SaveBlock(ctx, dbBlock)
	if err != nil {
		log.Errorw("failed to save block in DB", "err", err)
		return err
	}
	log.Debugw("saved block in db", "dbBlock", dbBlock)
	err = b.publisher.PublishTxs(ctx, txs)
	if err != nil {
		log.Errorw("failed to publish txs", "err", err)
		return err
	}
	log.Debugw("published txs", "published", len(txs))

	return nil
}
