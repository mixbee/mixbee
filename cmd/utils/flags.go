package utils

import (
	"fmt"
	"github.com/mixbee/mixbee/common/config"
	"github.com/mixbee/mixbee/common/constants"
	"github.com/mixbee/mixbee/smartcontract/service/neovm"
	"github.com/urfave/cli"
	"strings"
)

const (
	DEFAULT_EXPORT_FILE             = "./blocks.dat"
	DEFAULT_ABI_PATH                = "./abi"
	DEFAULT_CROSS_CHAIN_DEPLOY_TIME = 60 * 60 * 2
)

var (
	//Mixbee setting
	ConfigFlag = cli.StringFlag{
		Name:  "config",
		Usage: "Use `<filename>` to specifies the genesis block config file. If doesn't specifies the genesis block config, Mixbee will use Polaris config with VBFT consensus as default.",
	}
	LogLevelFlag = cli.UintFlag{
		Name:  "loglevel",
		Usage: "Set the log level to `<level>` (0~6). 0:Debug 1:Info 2:Warn 3:Error 4:Fatal 5:Trace 6:MaxLevel",
		Value: config.DEFAULT_LOG_LEVEL,
	}
	DisableEventLogFlag = cli.BoolFlag{
		Name:  "disableeventlog",
		Usage: "If set disableeventlog flag, Mixbee will not record event log output by smart contract",
	}
	WalletFileFlag = cli.StringFlag{
		Name:  "wallet,w",
		Value: config.DEFAULT_WALLET_FILE_NAME,
		Usage: "Use `<filename>` as the wallet",
	}
	ImportEnableFlag = cli.BoolFlag{
		Name:  "import",
		Usage: "Import blocks for file",
	}
	ImportFileFlag = cli.StringFlag{
		Name:  "importfile",
		Usage: "Path of import file",
		Value: DEFAULT_EXPORT_FILE,
	}
	ImportHeightFlag = cli.UintFlag{
		Name:  "importheight",
		Usage: "Using to specifies the height of the imported target block. If the block height specified by importheight is less than the maximum height of the block file, it will only be imported to the height specified by importheight and the rest blocks will stop importing. The default value is 0, which means import all the blocks",
	}
	DataDirFlag = cli.StringFlag{
		Name:  "datadir",
		Usage: "Using dir `<path>` to storage block data",
		Value: config.DEFAULT_DATA_DIR,
	}

	//Consensus setting
	EnableConsensusFlag = cli.BoolFlag{
		Name:  "enableconsensus",
		Usage: "If set enableconsensus, will start consensus module",
	}
	MaxTxInBlockFlag = cli.IntFlag{
		Name:  "maxtxinblock",
		Usage: "Using maxtxinblock to set the max transaction number in block",
		Value: config.DEFAULT_MAX_TX_IN_BLOCK,
	}
	GasLimitFlag = cli.Uint64Flag{
		Name:  "gaslimit",
		Usage: "Using to set the gaslimit of the current node transaction pool to accept transactions. Transactions below this gaslimit will be discarded",
		Value: neovm.MIN_TRANSACTION_GAS,
	}
	GasPriceFlag = cli.Uint64Flag{
		Name:  "gasprice",
		Usage: "Using to set the lowest gasprice of the current node transaction pool to accept transactions. Transactions below this gasprice will be discarded",
		Value: config.DEFAULT_GAS_PRICE,
	}

	//Test Mode setting
	EnableTestModeFlag = cli.BoolFlag{
		Name:  "testmode",
		Usage: "Using to start a single node test network for ease of development and debug. In testmode, Mixbee will start rpc, rest and web socket server",
	}
	TestModeGenBlockTimeFlag = cli.UintFlag{
		Name:  "testmodegenblocktime",
		Usage: "Using to set the block-out time in test mode. The time unit is in seconds, and the minimum block-out time is 2 seconds.",
		Value: config.DEFAULT_GEN_BLOCK_TIME,
	}
	ClearTestModeDataFlag = cli.BoolFlag{
		Name:  "cleartestmodedata",
		Usage: "Clear test mode block data",
	}

	//P2P setting
	ReservedPeersOnlyFlag = cli.BoolFlag{
		Name:  "reservedonly",
		Usage: "connect reserved peers only",
	}
	ReservedPeersFileFlag = cli.StringFlag{
		Name:  "reservedfile",
		Usage: "reserved peers file",
		Value: config.DEFAULT_RESERVED_FILE,
	}
	NetworkIdFlag = cli.UintFlag{
		Name:  "networkid",
		Usage: "Using to specify the network ID. Different networkids cannot connect to the blockchain network. 1=Mixbee main net, 2=polaris test net, 3=testmode, and other for custom network",
		Value: config.NETWORK_ID_MAIN_NET,
	}
	NodePortFlag = cli.UintFlag{
		Name:  "nodeport",
		Usage: "Using to specify the P2P network port number",
		Value: config.DEFAULT_NODE_PORT,
	}
	DualPortSupportFlag = cli.BoolFlag{
		Name:  "dualport",
		Usage: "Using to initiates a dual network, i.e. a P2P network for processing transaction messages and a consensus network for consensus messages. ",
	}
	ConsensusPortFlag = cli.UintFlag{
		Name:  "consensusport",
		Usage: "Using to specifies the consensus network port number. By default, the consensus network reuses the P2P network, so it is not necessary to specify a consensus network port. After the dual network is enabled with the --dualport parameter, the consensus network port number must be set separately.",
		Value: config.DEFAULT_CONSENSUS_PORT,
	}
	MaxConnInBoundFlag = cli.UintFlag{
		Name:  "maxconninbound",
		Usage: "Max connection in bound",
		Value: config.DEFAULT_MAX_CONN_IN_BOUND,
	}
	MaxConnOutBoundFlag = cli.UintFlag{
		Name:  "maxconnoutbound",
		Usage: "Max connection out bound",
		Value: config.DEFAULT_MAX_CONN_OUT_BOUND,
	}
	MaxConnInBoundForSingleIPFlag = cli.UintFlag{
		Name:  "maxconninboundforsingleip",
		Usage: "Max connection in bound for single ip",
		Value: config.DEFAULT_MAX_CONN_IN_BOUND_FOR_SINGLE_IP,
	}
	// RPC settings
	RPCDisabledFlag = cli.BoolFlag{
		Name:  "disablerpc",
		Usage: "Using to shut down the rpc server. The Mixbee node starts the rpc server by default at startup.",
	}
	RPCPortFlag = cli.UintFlag{
		Name:  "rpcport",
		Usage: "Json rpc server listening port",
		Value: config.DEFAULT_RPC_PORT,
	}
	RPCLocalEnableFlag = cli.BoolFlag{
		Name:  "localrpc",
		Usage: "Enable local rpc server",
	}
	RPCLocalProtFlag = cli.UintFlag{
		Name:  "rpclocalport",
		Usage: "Json rpc local server listening port",
		Value: config.DEFAULT_RPC_LOCAL_PORT,
	}

	//Websocket setting
	WsEnabledFlag = cli.BoolFlag{
		Name:  "ws",
		Usage: "Enable websocket server",
	}
	WsPortFlag = cli.UintFlag{
		Name:  "wsport",
		Usage: "Ws server listening port",
		Value: config.DEFAULT_WS_PORT,
	}

	//Restful setting
	RestfulEnableFlag = cli.BoolFlag{
		Name:  "rest",
		Usage: "Enable restful api server",
	}
	RestfulPortFlag = cli.UintFlag{
		Name:  "restport",
		Usage: "Restful server listening port",
		Value: config.DEFAULT_REST_PORT,
	}

	//Account setting
	AccountPassFlag = cli.StringFlag{
		Name:   "password,p",
		Hidden: true,
		Usage:  "Using to specify the account `<password>` when Mixbee node starts. Because the account password entered in the command line is saved in the log, it is easy to leak the password. Therefore, it is not recommended to use this parameter in a production environment.",
	}
	AccountAddressFlag = cli.StringFlag{
		Name:  "account,a",
		Usage: "Using to specify the account `<address|label|index>` when the Mixbee node starts. If the account is null, it uses the wallet default account",
	}
	AccountDefaultFlag = cli.BoolFlag{
		Name:  "default,d",
		Usage: "Use default settings to create a new account (equal to '-t ecdsa -b 256 -s SHA256withECDSA')",
	}
	AccountTypeFlag = cli.StringFlag{
		Name:  "type,t",
		Usage: "Specifies the `<key-type>` by signature algorithm.",
	}
	AccountKeylenFlag = cli.StringFlag{
		Name:  "bit-length,b",
		Usage: "Specifies the `<bit-length>` of key",
	}
	AccountSigSchemeFlag = cli.StringFlag{
		Name:  "signature-scheme,s",
		Usage: "Specifies the signature scheme `<scheme>`",
	}
	AccountSetDefaultFlag = cli.BoolFlag{
		Name:  "as-default,d",
		Usage: "Set the specified account to default account",
	}
	AccountQuantityFlag = cli.UintFlag{
		Name:  "number,n",
		Value: 1,
		Usage: "Specifies the `<quantity>` of account to create.",
	}
	AccountSourceFileFlag = cli.StringFlag{
		Name:  "source,s",
		Usage: "Use `<filename>` as the source wallet file to import",
	}
	AccountLabelFlag = cli.StringFlag{
		Name:  "label,l",
		Usage: "Use `<label>` for newly created accounts for easy and fast use of accounts. Note that duplicate label names cannot appear in the same wallet file. An account with no label is an empty string.",
	}
	AccountKeyFlag = cli.StringFlag{
		Name:  "key,k",
		Usage: "Use `<private key(hex encoding)>` of the account",
	}
	AccountVerboseFlag = cli.BoolFlag{
		Name:  "verbose,v",
		Usage: "Display accounts with details",
	}
	AccountChangePasswdFlag = cli.BoolFlag{
		Name:  "changepasswd",
		Usage: "Change account password",
	}
	AccountLowSecurityFlag = cli.BoolFlag{
		Name:  "low-security",
		Usage: "Change account to low protection strength for low performance devices",
	}
	AccountWIFFlag = cli.BoolFlag{
		Name:  "wif",
		Usage: "Import WIF keys from the source file specified by --source option",
	}
	AccountMultiMFlag = cli.UintFlag{
		Name:  "m",
		Usage: fmt.Sprintf("M of multi signature address. m must > 0 and <= %d, and m must <= number of pub key", constants.MULTI_SIG_MAX_PUBKEY_SIZE),
	}
	AccountMultiPubKeyFlag = cli.StringFlag{
		Name:  "pubkey",
		Usage: fmt.Sprintf("Pub key list of multi address, split pub key with `,`. Number of pub key must > 0 and <= %d", constants.MULTI_SIG_MAX_PUBKEY_SIZE),
	}
	IdentityFlag = cli.BoolFlag{
		Name:  "mixid",
		Usage: "create an MBC ID instead of account",
	}

	//SmartContract setting
	ContractAddrFlag = cli.StringFlag{
		Name:  "address",
		Usage: "Contract address",
	}
	ContractStorageFlag = cli.BoolFlag{
		Name:  "needstore",
		Usage: "Is need use storage in contract",
	}
	ContractCodeFileFlag = cli.StringFlag{
		Name:  "code",
		Usage: "File path of contract code `<path>`",
	}
	ContractNameFlag = cli.StringFlag{
		Name:  "name",
		Usage: "Specifies contract name to `<name>`",
	}
	ContractVersionFlag = cli.StringFlag{
		Name:  "version",
		Usage: "Specifies contract version to `<ver>`",
	}
	ContractAuthorFlag = cli.StringFlag{
		Name:  "author",
		Usage: "Set `<address>` as the contract owner",
		Value: "",
	}
	ContractEmailFlag = cli.StringFlag{
		Name:  "email",
		Usage: "Set `<email>` owner's email address",
		Value: "",
	}
	ContractDescFlag = cli.StringFlag{
		Name:  "desc",
		Usage: "Set `<text>` as the description of the contract",
		Value: "",
	}
	ContractParamsFlag = cli.StringFlag{
		Name:  "params",
		Usage: "Invoke contract parameters list. use comma ',' to split params, and must add type prefix to params. Param type support bytearray(hexstring), string, integer, boolean,For example: string:foo,int:0,bool:true; If parameter is an object array, enclose array with '[]'. For example:  string:foo,[int:0,bool:true]",
	}
	ContractPrepareDeployFlag = cli.BoolFlag{
		Name:  "prepare,p",
		Usage: "Prepare deploy contract without commit to ledger",
	}
	ContractPrepareInvokeFlag = cli.BoolFlag{
		Name:  "prepare,p",
		Usage: "Prepare invoke contract without commit to ledger",
	}
	ContractReturnTypeFlag = cli.StringFlag{
		Name:  "return",
		Usage: "Return type of contract.Return type support bytearray(hexstring), string, integer, boolean. If return type is object array, enclose array with '[]'. For example [string,int,bool,string]. Only prepare invoke need this flag.",
	}

	//information cmd settings
	BlockHashInfoFlag = cli.StringFlag{
		Name:  "hash",
		Usage: "Get block info by block hash",
	}
	BlockHeightInfoFlag = cli.UintFlag{
		Name:  "height",
		Usage: "Get block info by block height",
	}

	//Transfer setting
	TransactionAssetFlag = cli.StringFlag{
		Name:  "asset",
		Usage: "Using to specifies the transfer asset `<mbc|mbg>`",
		Value: ASSET_MBC,
	}
	TransactionFromFlag = cli.StringFlag{
		Name:  "from",
		Usage: "Using to specifies the transfer-out account `<address|label|index>`",
	}
	TransactionToFlag = cli.StringFlag{
		Name:  "to",
		Usage: "Using to specifies the transfer-in account `<address|label|index>`",
	}
	TransactionAmountFlag = cli.StringFlag{
		Name:  "amount",
		Usage: "Using to specifies the transfer amount",
	}
	TransactionHashFlag = cli.StringFlag{
		Name:  "hash",
		Usage: "Transaction <hash>",
	}
	TransactionGasPriceFlag = cli.Uint64Flag{
		Name:  "gasprice",
		Usage: "Using to specifies the gas price of transaction. The gas price of the transaction cannot be less than the lowest gas price set by node's transaction pool, otherwise the transaction will be rejected. When there are transactions that are queued for packing into the block in the transaction pool, the transaction pool will deal with transactions according to the gas price and transactions with high gas prices will be prioritized",
		Value: 0,
	}
	TransactionGasLimitFlag = cli.Uint64Flag{
		Name:  "gaslimit",
		Usage: "Using to specifies the gas limit of the transaction. The gas limit of the transaction cannot be less than the minimum gas limit set by the node's transaction pool, otherwise the transaction will be rejected. Gasprice * gaslimit is actual MBG costs.",
		Value: neovm.MIN_TRANSACTION_GAS,
	}

	//Asset setting
	ApproveAssetFromFlag = cli.StringFlag{
		Name:  "from",
		Usage: "Using to specifies the transfer-out account `<address|label|index>`",
	}
	ApproveAssetToFlag = cli.StringFlag{
		Name:  "to",
		Usage: "Using to specifies the transfer-in account `<address|label|index>`",
	}
	ApproveAssetFlag = cli.StringFlag{
		Name:  "asset",
		Usage: "Using to specifies the transfer asset <mbc|mbg> for approve",
		Value: "mbc",
	}
	ApproveAmountFlag = cli.StringFlag{
		Name:  "amount",
		Usage: "Using to specifies the transfer amount for approve",
	}
	TransferFromAmountFlag = cli.StringFlag{
		Name:  "amount",
		Usage: "Using to specifies the transfer from amount",
	}
	TransferFromSenderFlag = cli.StringFlag{
		Name:  "sender",
		Usage: "Using to specifies the sender account `<address|label|index>` of transfer from transaction, if empty sender is to account",
	}

	//Cli setting
	CliRpcPortFlag = cli.UintFlag{
		Name:  "cliport",
		Usage: "Cli rpc port",
		Value: config.DEFAULT_CLI_RPC_PORT,
	}
	CliABIPathFlag = cli.StringFlag{
		Name:  "abi",
		Usage: "Abi path",
		Value: DEFAULT_ABI_PATH,
	}

	//Export setting
	ExportFileFlag = cli.StringFlag{
		Name:  "file",
		Usage: "Path of export file",
		Value: DEFAULT_EXPORT_FILE,
	}
	ExportHeightFlag = cli.UintFlag{
		Name:  "height",
		Usage: "Using to specifies the height of the exported block. When height of the local node's current block is greater than the height required for export, the greater part will not be exported. Height is equal to 0, which means exporting all the blocks of the current node.",
		Value: 0,
	}
	ExportSpeedFlag = cli.StringFlag{
		Name:  "speed",
		Usage: "Export block speed, `<h|m|l>` h for high speed, m for middle speed and l for low speed",
		Value: "m",
	}

	//PreExecute switcher
	PreExecEnableFlag = cli.BoolFlag{
		Name:  "preexec",
		Usage: "Enable preExecute in tx pool",
	}

	NonOptionFlag = cli.StringFlag{
		Name:  "option",
		Usage: "this command does not need option, please run directly",
	}

	//EnableCrossChainVerify setting
	EnableCrossChainVerifyFlag = cli.BoolFlag{
		Name:  "enablecrosschainverify",
		Usage: "If set enablecrosschainverify, will start cross chain verify module on mainchain",
	}

	//EnableCrossChainInteractive setting
	EnableCrossChainInteractiveFlag = cli.BoolFlag{
		Name:  "enablecrosschaininter",
		Usage: "If set enablecrosschaininter, will start cross chain verify module in subchain",
	}

	//CrossChainVerifyNode setting
	CrossChainVerifyNode = cli.StringSliceFlag{
		Name:  "crosschainnode",
		Usage: "If set enablecrosschaininter true,crosschainnode must not nil, example:node1:port1",
	}

	//key value
	MixTestKeyFlag = cli.StringFlag{
		Name:  "key",
		Usage: "mix test contract setkey key",
	}

	MixTestValueFlag = cli.StringFlag{
		Name:  "value",
		Usage: "mix test contract setkey value",
	}

	CrossChainAValueFlag = cli.StringFlag{
		Name:  "aamount",
		Usage: "cross chain amount while from send to toAddress",
	}

	CrossChainBValueFlag = cli.StringFlag{
		Name:  "bamount",
		Usage: "cross chain amount while from want to obtain",
	}

	CrossChainBChainIdFlag = cli.StringFlag{
		Name:  "bchainid",
		Usage: "cross chain other chain id",
	}

	CrossChainAChainIdFlag = cli.StringFlag{
		Name:  "achainid",
		Usage: "cross chain from chain id",
	}

	CrossChainDelayTimeFlag = cli.UintFlag{
		Name:  "delaytime",
		Usage: "cross chain delay time >= 60*60*2 second",
		Value: DEFAULT_CROSS_CHAIN_DEPLOY_TIME,
	}

	CrossChainNonceFlag = cli.UintFlag{
		Name:  "nonce",
		Usage: "cross chain nonce default 0",
		Value: 0,
	}

	CrossChainSeqIdFlag = cli.StringFlag{
		Name:  "seqid",
		Usage: "cross chain seqId",
	}

	CrossChainVerifyPublicKeyFlag = cli.StringFlag{
		Name:  "cpbk",
		Usage: "cross chain, main chain verify node public key",
	}
)

//GetFlagName deal with short flag, and return the flag name whether flag name have short name
func GetFlagName(flag cli.Flag) string {
	name := flag.GetName()
	if name == "" {
		return ""
	}
	return strings.TrimSpace(strings.Split(name, ",")[0])
}
