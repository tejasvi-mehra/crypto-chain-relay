package blockrepo

import (
	"context"
	db "eth-adapter/internal/database/generated"
	"eth-adapter/internal/framework/logger"
)

type MockBlockRepo struct{}

func (b *MockBlockRepo) SaveBlock(ctx context.Context, block *db.BlockData) error {
	log := logger.GetLogger("saveblock")
	log.Debugw("saved block", "data", block)
	return nil
}
