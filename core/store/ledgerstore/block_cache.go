

package ledgerstore

import (
	"fmt"

	"github.com/hashicorp/golang-lru"
	"github.com/mixbee/mixbee/common"
	"github.com/mixbee/mixbee/core/types"
)

const (
	BLOCK_CAHE_SIZE        = 10    //Block cache size
	TRANSACTION_CACHE_SIZE = 10000 //Transaction cache size
)

//Value of transaction cache
type TransactionCacheaValue struct {
	Tx     *types.Transaction
	Height uint32
}

//BlockCache with block cache and transaction hash
type BlockCache struct {
	blockCache       *lru.ARCCache
	transactionCache *lru.ARCCache
}

//NewBlockCache return BlockCache instance
func NewBlockCache() (*BlockCache, error) {
	blockCache, err := lru.NewARC(BLOCK_CAHE_SIZE)
	if err != nil {
		return nil, fmt.Errorf("NewARC block error %s", err)
	}
	transactionCache, err := lru.NewARC(TRANSACTION_CACHE_SIZE)
	if err != nil {
		return nil, fmt.Errorf("NewARC header error %s", err)
	}
	return &BlockCache{
		blockCache:       blockCache,
		transactionCache: transactionCache,
	}, nil
}

//AddBlock to cache
func (this *BlockCache) AddBlock(block *types.Block) {
	blockHash := block.Hash()
	this.blockCache.Add(string(blockHash.ToArray()), block)
}

//GetBlock return block by block hash from cache
func (this *BlockCache) GetBlock(blockHash common.Uint256) *types.Block {
	block, ok := this.blockCache.Get(string(blockHash.ToArray()))
	if !ok {
		return nil
	}
	return block.(*types.Block)
}

//ContainBlock retuen whether block is in cache
func (this *BlockCache) ContainBlock(blockHash common.Uint256) bool {
	return this.blockCache.Contains(string(blockHash.ToArray()))
}

//AddTransaction add transaction to block cache
func (this *BlockCache) AddTransaction(tx *types.Transaction, height uint32) {
	txHash := tx.Hash()
	this.transactionCache.Add(string(txHash.ToArray()), &TransactionCacheaValue{
		Tx:     tx,
		Height: height,
	})
}

//GetTransaction return transaction by transaction hash from cache
func (this *BlockCache) GetTransaction(txHash common.Uint256) (*types.Transaction, uint32) {
	value, ok := this.transactionCache.Get(string(txHash.ToArray()))
	if !ok {
		return nil, 0
	}
	txValue := value.(*TransactionCacheaValue)
	return txValue.Tx, txValue.Height
}

//ContainTransaction return whether transaction is in cache
func (this *BlockCache) ContainTransaction(txHash common.Uint256) bool {
	return this.transactionCache.Contains(string(txHash.ToArray()))
}
