package trace_id

import (
	"math/rand"
	"time"

	"eth-adapter/internal/framework/logger"

	"github.com/mitchellh/hashstructure"
	"github.com/oklog/ulid/v2"
)

type UlidHashData struct {
	ChainID         string
	Block           string
	BlockTimestamp  string
	Hash            string
	ActionAddress   string
	TransactionType string
	Symbol          string
	Value           string
}

func (u UlidHashData) GenerateULID() string {
	log := logger.GetLogger("GenerateULID")
	defer func() {
		log.Debugw("completed")
		log.Sync()
	}()
	// generate hash based on TransactionMsgData for entropy - hashstructure.Hash has complexity O(n) where n is number of elements in struct. For us its always constant time
	objectHash, err := hashstructure.Hash(u, nil)
	if err != nil {
		log.Warnw("cannot generate object hash, generate random ULID", "error", err)
		return ulid.Make().String()
	}
	// generate entropy based on objectHash
	entropy := ulid.Monotonic(rand.New(rand.NewSource(int64(objectHash))), 0)
	blockTime, err := time.Parse(time.RFC3339, u.BlockTimestamp)
	if err != nil {
		log.Warnw("cannot set blockTime, generate random ULID", "error", err)
		return ulid.Make().String()
	}
	// Make ULID - it has 2 parts, timestamp and entropy - since timestamp is always set for block, if the transactionMsgData is same, the same ULID will be generated each time - this allows block replays
	listenerULID, err := ulid.New(ulid.Timestamp(blockTime), entropy)
	if err != nil {
		log.Warnw("cannot make ulid, generate random ULID", "error", err)
		return ulid.Make().String()
	}
	return listenerULID.String()
}
