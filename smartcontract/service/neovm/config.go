

package neovm

import "sync"

var (
	//Gas Limit
	MIN_TRANSACTION_GAS           uint64 = 20000 // Per transaction base cost.
	BLOCKCHAIN_GETHEADER_GAS      uint64 = 100
	BLOCKCHAIN_GETBLOCK_GAS       uint64 = 200
	BLOCKCHAIN_GETTRANSACTION_GAS uint64 = 100
	BLOCKCHAIN_GETCONTRACT_GAS    uint64 = 100
	CONTRACT_CREATE_GAS           uint64 = 20000000
	CONTRACT_MIGRATE_GAS          uint64 = 20000000
	UINT_DEPLOY_CODE_LEN_GAS      uint64 = 200000
	UINT_INVOKE_CODE_LEN_GAS      uint64 = 20000
	NATIVE_INVOKE_GAS             uint64 = 1000
	STORAGE_GET_GAS               uint64 = 200
	STORAGE_PUT_GAS               uint64 = 4000
	STORAGE_DELETE_GAS            uint64 = 100
	RUNTIME_CHECKWITNESS_GAS      uint64 = 200
	APPCALL_GAS                   uint64 = 10
	TAILCALL_GAS                  uint64 = 10
	SHA1_GAS                      uint64 = 10
	SHA256_GAS                    uint64 = 10
	HASH160_GAS                   uint64 = 20
	HASH256_GAS                   uint64 = 20
	OPCODE_GAS                    uint64 = 1

	PER_UNIT_CODE_LEN   int = 1024
	METHOD_LENGTH_LIMIT int = 1024
	MAX_STACK_SIZE      int = 1024
	VM_STEP_LIMIT       int = 400000

	// API Name
	ATTRIBUTE_GETUSAGE_NAME = "Neo.Attribute.GetUsage"
	ATTRIBUTE_GETDATA_NAME  = "Neo.Attribute.GetData"

	BLOCK_GETTRANSACTIONCOUNT_NAME       = "Neo.Block.GetTransactionCount"
	BLOCK_GETTRANSACTIONS_NAME           = "Neo.Block.GetTransactions"
	BLOCK_GETTRANSACTION_NAME            = "Neo.Block.GetTransaction"


	BLOCKCHAIN_GETHEIGHT_NAME            = "Neo.Blockchain.GetHeight"
	BLOCKCHAIN_GETHEADER_NAME            = "Neo.Blockchain.GetHeader"
	BLOCKCHAIN_GETBLOCK_NAME             = "Neo.Blockchain.GetBlock"
	BLOCKCHAIN_GETTRANSACTION_NAME       = "Neo.Blockchain.GetTransaction"
	BLOCKCHAIN_GETCONTRACT_NAME          = "Neo.Blockchain.GetContract"
	BLOCKCHAIN_GETTRANSACTIONHEIGHT_NAME = "Neo.Blockchain.GetTransactionHeight"


	HEADER_GETINDEX_NAME         = "Neo.Header.GetIndex"
	HEADER_GETHASH_NAME          = "Neo.Header.GetHash"
	HEADER_GETVERSION_NAME       = "Neo.Header.GetVersion"
	HEADER_GETPREVHASH_NAME      = "Neo.Header.GetPrevHash"
	HEADER_GETTIMESTAMP_NAME     = "Neo.Header.GetTimestamp"
	HEADER_GETCONSENSUSDATA_NAME = "Neo.Header.GetConsensusData"
	HEADER_GETNEXTCONSENSUS_NAME = "Neo.Header.GetNextConsensus"
	HEADER_GETMERKLEROOT_NAME    = "Neo.Header.GetMerkleRoot"


	TRANSACTION_GETHASH_NAME       = "Neo.Transaction.GetHash"
	TRANSACTION_GETTYPE_NAME       = "Neo.Transaction.GetType"
	TRANSACTION_GETATTRIBUTES_NAME = "Neo.Transaction.GetAttributes"


	CONTRACT_CREATE_NAME            = "Neo.Contract.Create"
	CONTRACT_MIGRATE_NAME           = "Neo.Contract.Migrate"
	CONTRACT_GETSTORAGECONTEXT_NAME = "Neo.Contract.GetStorageContext"
	CONTRACT_DESTROY_NAME           = "Neo.Contract.Destroy"
	CONTRACT_GETSCRIPT_NAME         = "Neo.Contract.GetScript"


	STORAGE_GET_NAME                = "Neo.Storage.Get"
	STORAGE_PUT_NAME                = "Neo.Storage.Put"
	STORAGE_DELETE_NAME             = "Neo.Storage.Delete"
	STORAGE_GETCONTEXT_NAME         = "Neo.Storage.GetContext"
	STORAGE_GETREADONLYCONTEXT_NAME = "Neo.Storage.GetReadOnlyContext"

	STORAGECONTEXT_ASREADONLY_NAME = "System.StorageContext.AsReadOnly"

	RUNTIME_GETTIME_NAME      = "Neo.Runtime.GetTime"
	RUNTIME_CHECKWITNESS_NAME = "Neo.Runtime.CheckWitness"
	RUNTIME_NOTIFY_NAME       = "Neo.Runtime.Notify"
	RUNTIME_LOG_NAME          = "Neo.Runtime.Log"
	RUNTIME_GETTRIGGER_NAME   = "Neo.Runtime.GetTrigger"
	RUNTIME_SERIALIZE_NAME    = "Neo.Runtime.Serialize"
	RUNTIME_DESERIALIZE_NAME  = "Neo.Runtime.Deserialize"


	// asset
	ASSET_GETASSETID_NAME   = "Neo.Asset.GetAssetId"
	ASSET_GETASSETTYPE_NAME = "Neo.Asset.GetAssetType"
	ASSET_GETAMOUNT_NAME    = "Neo.Asset.GetAmount"
	ASSET_GETAVAILABLE_NAME = "Neo.Asset.GetAvailable"
	ASSET_GETPRECISION_NAME = "Neo.Asset.GetPrecision"
	ASSET_GETOWNER_NAME     = "Neo.Asset.GetOwner"
	ASSET_GETADMIN_NAME     = "Neo.Asset.GetAdmin"
	ASSET_GETISSUER_NAME    = "Neo.Asset.GetIssuer"
	ASSET_CREATE_NAME       = "Neo.Asset.Create"
	ASSET_RENEW_NAME        = "Neo.Asset.Renew"

	//account api
	ACCOUNT_GETSCRIPTHASH_NAME 	= "Neo.Account.GetScriptHash"
	ACCOUNT_GETVOTES_NAME 		= "Neo.Account.GetVotes"
	ACCOUNT_GETBALANCE_NAME 	= "Neo.Account.GetBalance"




	NATIVE_INVOKE_NAME = "Mixbee.Native.Invoke"


	GETSCRIPTCONTAINER_NAME     = "System.ExecutionEngine.GetScriptContainer"
	GETEXECUTINGSCRIPTHASH_NAME = "System.ExecutionEngine.GetExecutingScriptHash"
	GETCALLINGSCRIPTHASH_NAME   = "System.ExecutionEngine.GetCallingScriptHash"
	GETENTRYSCRIPTHASH_NAME     = "System.ExecutionEngine.GetEntryScriptHash"

	APPCALL_NAME              = "APPCALL"
	TAILCALL_NAME             = "TAILCALL"
	SHA1_NAME                 = "SHA1"
	SHA256_NAME               = "SHA256"
	HASH160_NAME              = "HASH160"
	HASH256_NAME              = "HASH256"
	UINT_DEPLOY_CODE_LEN_NAME = "Deploy.Code.Gas"
	UINT_INVOKE_CODE_LEN_NAME = "Invoke.Code.Gas"

	GAS_TABLE = initGAS_TABLE()

	GAS_TABLE_KEYS = []string{
		BLOCKCHAIN_GETHEADER_NAME,
		BLOCKCHAIN_GETBLOCK_NAME,
		BLOCKCHAIN_GETTRANSACTION_NAME,
		BLOCKCHAIN_GETCONTRACT_NAME,
		CONTRACT_CREATE_NAME,
		CONTRACT_MIGRATE_NAME,
		STORAGE_GET_NAME,
		STORAGE_PUT_NAME,
		STORAGE_DELETE_NAME,
		RUNTIME_CHECKWITNESS_NAME,
		NATIVE_INVOKE_NAME,
		APPCALL_NAME,
		TAILCALL_NAME,
		SHA1_NAME,
		SHA256_NAME,
		HASH160_NAME,
		HASH256_NAME,
		UINT_DEPLOY_CODE_LEN_NAME,
		UINT_INVOKE_CODE_LEN_NAME,
	}
)

func initGAS_TABLE() *sync.Map {
	m := sync.Map{}
	m.Store(BLOCKCHAIN_GETHEADER_NAME, BLOCKCHAIN_GETHEADER_GAS)
	m.Store(BLOCKCHAIN_GETBLOCK_NAME, BLOCKCHAIN_GETBLOCK_GAS)
	m.Store(BLOCKCHAIN_GETTRANSACTION_NAME, BLOCKCHAIN_GETTRANSACTION_GAS)
	m.Store(BLOCKCHAIN_GETCONTRACT_NAME, BLOCKCHAIN_GETCONTRACT_GAS)
	m.Store(CONTRACT_CREATE_NAME, CONTRACT_CREATE_GAS)
	m.Store(CONTRACT_MIGRATE_NAME, CONTRACT_MIGRATE_GAS)
	m.Store(STORAGE_GET_NAME, STORAGE_GET_GAS)
	m.Store(STORAGE_PUT_NAME, STORAGE_PUT_GAS)
	m.Store(STORAGE_DELETE_NAME, STORAGE_DELETE_GAS)
	m.Store(RUNTIME_CHECKWITNESS_NAME, RUNTIME_CHECKWITNESS_GAS)
	m.Store(NATIVE_INVOKE_NAME, NATIVE_INVOKE_GAS)
	m.Store(APPCALL_NAME, APPCALL_GAS)
	m.Store(TAILCALL_NAME, TAILCALL_GAS)
	m.Store(SHA1_NAME, SHA1_GAS)
	m.Store(SHA256_NAME, SHA256_GAS)
	m.Store(HASH160_NAME, HASH160_GAS)
	m.Store(HASH256_NAME, HASH256_GAS)
	m.Store(UINT_DEPLOY_CODE_LEN_NAME, UINT_DEPLOY_CODE_LEN_GAS)
	m.Store(UINT_INVOKE_CODE_LEN_NAME, UINT_INVOKE_CODE_LEN_GAS)

	return &m
}
