package publisher

import (
	"context"
	swagger "eth-adapter/internal/event/generated/1.0.0/transaction-events/go"
	"eth-adapter/internal/framework/logger"
)

type MockPublisher struct{}

func (p *MockPublisher) PublishTxs(ctx context.Context, txs []swagger.TxEventData) error {
	log := logger.GetLogger("PublishTxs")
	log.Debugw("publish txs", "count", len(txs))
	log.Debugw("txs", "txs", txs)
	return nil
}
