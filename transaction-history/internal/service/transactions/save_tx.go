package transactions

import (
	"context"
	"encoding/json"
	"fmt"
	"strconv"
	"time"
	db "transaction-history/internal/database/generated"
	swagger "transaction-history/internal/event/generated/1.0.0/transaction-events/go"
	"transaction-history/internal/framework/logger"
)

func (t *TxService) SaveTxs(ctx context.Context, params []swagger.TxEventData) error {
	log := logger.GetLogger("SaveTxs")
	defer func() {
		log.Debugw("completed")
	}()
	log.Debugw("started")
	// if err := t.v.Validate(params); err != nil {
	// 	log.Errorw("invalid params", "err", err)
	// 	return fmt.Errorf("invalid params: %v", err)
	// }
	log.With("params", params)
	dbTxs := make([]*db.Txs, 0)
	for _, tx := range params {
		if tx.EventType == "crypto" && tx.CryptoSpecificInfo == nil {
			log.Errorw("invalid params, need crypto specific info")
			return fmt.Errorf("invalid params, need crypto specific info")
		}
		if tx.EventType == "fiat" && tx.FiatSpecificInfo == nil {
			log.Errorw("invalid params, need fiat specific info")
			return fmt.Errorf("invalid params, need fiat specific info")
		}
		if tx.EventType == "crypto" {
			// parse as crypto
			amountDecimal, _ := strconv.ParseFloat(tx.AmountDecimal, 64)
			txTime, _ := time.Parse(time.RFC3339, tx.TransactionTime)
			marshal, _ := json.Marshal(tx.CryptoSpecificInfo)
			var cryptoInfo map[string]any
			json.Unmarshal(marshal, &cryptoInfo)
			now := time.Now().UTC()
			dbTx := &db.Txs{
				ID:                 int(time.Now().Unix()),
				TxType:             (*string)(tx.TransactionType),
				TraceID:            &tx.TraceID,
				AmountDecimal:      amountDecimal,
				Symbol:             &tx.Symbol,
				AccountID:          tx.AccountID,
				From:               &tx.From,
				To:                 &tx.To,
				Status:             (*string)(tx.Status),
				TxTime:             &txTime,
				EventType:          &tx.EventType,
				CryptoSpecificInfo: cryptoInfo,
				CreatedAt:          &now,
				UpdatedAt:          &now,
			}
			dbTxs = append(dbTxs, dbTx)

		}
		if tx.EventType == "fiat" {
			// parse as crypto
			amountDecimal, _ := strconv.ParseFloat(tx.AmountDecimal, 64)
			txTime, _ := time.Parse(time.RFC3339, tx.TransactionTime)
			marshal, _ := json.Marshal(tx.FiatSpecificInfo)
			var fiatInfo map[string]any
			json.Unmarshal(marshal, &fiatInfo)
			now := time.Now().UTC()
			dbTx := &db.Txs{
				ID:               int(time.Now().Unix()),
				TxType:           (*string)(tx.TransactionType),
				TraceID:          &tx.TraceID,
				AmountDecimal:    amountDecimal,
				Symbol:           &tx.Symbol,
				AccountID:        tx.AccountID,
				From:             &tx.From,
				To:               &tx.To,
				Status:           (*string)(tx.Status),
				TxTime:           &txTime,
				EventType:        &tx.EventType,
				FiatSpecificInfo: fiatInfo,
				CreatedAt:        &now,
				UpdatedAt:        &now,
			}
			dbTxs = append(dbTxs, dbTx)
		}
	}
	err := t.repo.SaveTx(ctx, dbTxs)
	if err != nil {
		log.Errorw("failed to save transactions", "err", err)
		return err
	}
	return nil
}
