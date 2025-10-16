package repo

import (
	"context"
	db "transaction-history/internal/database/generated"
	"transaction-history/internal/framework/logger"
)

type MockRepo struct {
	txs []*db.Txs
}

func (m *MockRepo) GetAccountTransactions(ctx context.Context, accountID string, limit, offset int) ([]*db.Txs, error) {
	log := logger.GetLogger("GetAccountTransactions")
	log.With("accountID", accountID, "limit", limit, "offset", offset)
	res := make([]*db.Txs, 0)
	min := offset
	max := offset + limit
	for i, tx := range m.txs {
		if tx.AccountID == accountID && i >= min && i < max {
			res = append(res, tx)
		}
	}
	log.Debugw("got txs", "count", len(res))
	if len(res) == 0 {
		return nil, nil
	}
	return res, nil
}

func (m *MockRepo) GetAccountBalances(ctx context.Context, accountID string) ([]*AccountBalances, error) {
	log := logger.GetLogger("GetAccountTransactions")
	log.With("accountID", accountID)
	type amountEvent struct {
		Amount    float64
		EventType string
	}
	balanceMap := make(map[string]map[string]amountEvent)
	result := make([]*AccountBalances, 0)
	for _, tx := range m.txs {
		if tx.AccountID == accountID {
			if tx.Symbol == nil {
				continue
			}
			if tx.EventType == nil {
				continue
			}
			if tx.TxType == nil {
				continue
			}
			sybmol := *tx.Symbol
			eventType := *tx.EventType
			txType := *tx.TxType
			data, ok := balanceMap[sybmol]
			if !ok {
				b := map[string]amountEvent{
					txType: amountEvent{
						Amount:    tx.AmountDecimal,
						EventType: eventType,
					},
				}
				balanceMap[sybmol] = b
			} else {
				txTypeMap, ok := data[txType]
				if ok {
					amountData := amountEvent{
						Amount:    txTypeMap.Amount + tx.AmountDecimal,
						EventType: eventType,
					}
					data[txType] = amountData
					balanceMap[sybmol] = data
				} else {
					amountData := amountEvent{
						Amount:    tx.AmountDecimal,
						EventType: eventType,
					}
					data[txType] = amountData
					balanceMap[sybmol] = data
				}
			}
		}
	}
	for symbol, balanceInfo := range balanceMap {
		for txType, amount := range balanceInfo {
			res := &AccountBalances{
				Amount:    amount.Amount,
				TxType:    txType,
				EventType: amount.EventType,
				Symbol:    symbol,
			}
			result = append(result, res)
		}
	}
	log.Debugw("got result", "res", result)
	if len(result) == 0 {
		return nil, nil
	}
	return result, nil
}

func (m *MockRepo) SaveTx(ctx context.Context, tx []*db.Txs) error {
	m.txs = append(m.txs, tx...)
	return nil
}
