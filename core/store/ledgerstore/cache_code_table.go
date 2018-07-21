

package ledgerstore

import (
	"fmt"

	"github.com/mixbee/mixbee/core/payload"
	scom "github.com/mixbee/mixbee/core/store/common"
)

type CacheCodeTable struct {
	store scom.StateStore
}

func (table *CacheCodeTable) GetCode(address []byte) ([]byte, error) {
	value, _ := table.store.TryGet(scom.ST_CONTRACT, address)
	if value == nil {
		return nil, fmt.Errorf("[GetCode] TryGet contract error! address:%x", address)
	}

	return value.Value.(*payload.DeployCode).Code, nil
}
