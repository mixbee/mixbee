

package db

import (
	pool "github.com/valyala/bytebufferpool"

	"github.com/mixbee/mixbee/common"
)

// DataEntryPrefix
type KeyPrefix byte

const (
	//SYSTEM
	SYS_VERSION       KeyPrefix = 0
	SYS_GENESIS_BLOCK KeyPrefix = 1 // key: prefix, value: gensisBlock

	SYS_BEST_BLOCK        KeyPrefix = 2 // key : prefix, value: bestblock
	SYS_BEST_BLOCK_HEADER KeyPrefix = 3 // key: prefix, value: BlockHeader

	// DATA
	//DATA_Block KeyPrefix = iota
	//DATA_Header
	DATA_TRANSACTION KeyPrefix = 10 // key: prefix+txid, value: height + tx
)

func GenGenesisBlockKey() *pool.ByteBuffer {
	key := keyPool.Get()
	key.WriteByte(byte(SYS_GENESIS_BLOCK))
	return key
}

func GenBestBlockHeaderKey() *pool.ByteBuffer {
	key := keyPool.Get()
	key.WriteByte(byte(SYS_BEST_BLOCK_HEADER))
	return key
}

func GenDataTransactionKey(hash common.Uint256) *pool.ByteBuffer {
	key := keyPool.Get()
	key.WriteByte(byte(DATA_TRANSACTION))
	key.Write(hash.ToArray())
	return key
}
