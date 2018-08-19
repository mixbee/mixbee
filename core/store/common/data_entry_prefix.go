

package common

// DataEntryPrefix
type DataEntryPrefix byte

const (
	// DATA
	DATA_BLOCK       DataEntryPrefix = 0x00 //Block height => block hash key prefix
	DATA_HEADER                      = 0x01 //Block hash => block hash key prefix
	DATA_TRANSACTION                 = 0x02 //Transction hash = > transaction key prefix

	// Transaction
	ST_BOOKKEEPER DataEntryPrefix = 0x03 //BookKeeper state key prefix
	ST_CONTRACT   DataEntryPrefix = 0x04 //Smart contract state key prefix
	ST_STORAGE    DataEntryPrefix = 0x05 //Smart contract storage key prefix
	ST_VALIDATOR  DataEntryPrefix = 0x07 //no use
	ST_VOTE       DataEntryPrefix = 0x08 //Vote state key prefix

	IX_HEADER_HASH_LIST DataEntryPrefix = 0x09 //Block height => block hash key prefix

	//SYSTEM
	SYS_CURRENT_BLOCK      DataEntryPrefix = 0x10 //Current block key prefix
	SYS_VERSION            DataEntryPrefix = 0x11 //Store version key prefix
	SYS_CURRENT_STATE_ROOT DataEntryPrefix = 0x12 //no use
	SYS_BLOCK_MERKLE_TREE  DataEntryPrefix = 0x13 // Block merkle tree root key prefix

	EVENT_NOTIFY DataEntryPrefix = 0x14 //Event notify key prefix



	// ASSET
	ST_Info           DataEntryPrefix = 0xc0
	ST_QuantityIssued DataEntryPrefix = 0xc1
	ST_ACCOUNT        DataEntryPrefix = 0xc4
	ST_AssetState     DataEntryPrefix = 0xc5
	ST_Validator      DataEntryPrefix = 0xc6
)
