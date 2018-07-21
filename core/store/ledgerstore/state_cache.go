

package ledgerstore

import (
	"github.com/hashicorp/golang-lru"
	"github.com/mixbee/mixbee/core/states"
)

const (
	STATE_CACHE_SIZE = 100000
)

type StateCache struct {
	stateCache *lru.ARCCache
}

func NewStateCache() (*StateCache, error) {
	stateCache, err := lru.NewARC(STATE_CACHE_SIZE)
	if err != nil {
		return nil, err
	}
	return &StateCache{
		stateCache: stateCache,
	}, nil
}

func (this *StateCache) GetState(key []byte) states.StateValue {
	state, ok := this.stateCache.Get(string(key))
	if !ok {
		return nil
	}
	return state.(states.StateValue)
}

func (this *StateCache) AddState(key []byte, state states.StateValue) {
	this.stateCache.Add(string(key), state)
}

func (this *StateCache) DeleteState(key []byte) {
	this.stateCache.Remove(string(key))
}
