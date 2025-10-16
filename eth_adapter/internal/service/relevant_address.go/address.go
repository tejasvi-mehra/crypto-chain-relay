package relevantaddress

import (
	"context"
	"strings"
)

type RelevantAddressMock struct {
	addressMap map[string]string
}

func NewRelevantAddressMock(addressMap map[string]string) *RelevantAddressMock {
	return &RelevantAddressMock{
		addressMap: addressMap,
	}
}

func (r *RelevantAddressMock) IsRelevantAddress(ctx context.Context, address string) (relevant bool, accountID string) {
	accountId, ok := r.addressMap[strings.ToLower(address)]
	return ok, accountId
}
