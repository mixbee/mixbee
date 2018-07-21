

package statestore

import (
	"github.com/mixbee/mixbee/core/states"
	"github.com/mixbee/mixbee/core/store/common"
)

type MemoryStore struct {
	memory map[string]*common.StateItem
}

func NewMemDatabase() *MemoryStore {
	return &MemoryStore{
		memory: make(map[string]*common.StateItem),
	}
}

func (db *MemoryStore) Put(prefix byte, key []byte, value states.StateValue, state common.ItemState) {
	db.memory[string(append([]byte{prefix}, key...))] = &common.StateItem{
		Key:   string(key),
		Value: value,
		State: state,
	}
}

func (db *MemoryStore) Get(prefix byte, key []byte) *common.StateItem {
	if entry, ok := db.memory[string(append([]byte{prefix}, key...))]; ok {
		return entry
	}
	return nil
}

func (db *MemoryStore) Delete(prefix byte, key []byte) {
	if v, ok := db.memory[string(append([]byte{prefix}, key...))]; ok {
		v.State = common.Deleted
	} else {
		db.memory[string(append([]byte{prefix}, key...))] = &common.StateItem{
			Key:   string(key),
			State: common.Deleted,
		}
	}

}

func (db *MemoryStore) Find() []*common.StateItem {
	var memory []*common.StateItem
	for _, v := range db.memory {
		memory = append(memory, v)
	}
	return memory
}

func (db *MemoryStore) GetChangeSet() map[string]*common.StateItem {
	m := make(map[string]*common.StateItem)
	for k, v := range db.memory {
		if v.State != common.None {
			m[k] = v
		}
	}
	return m
}
