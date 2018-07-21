package main

import (
	"encoding/hex"
	"fmt"
	"os"
	"os/signal"
	"runtime"
	"strings"
	"syscall"
	"time"

	"github.com/mixbee/mixbee-crypto/keypair"
	"github.com/mixbee/mixbee/mixbee-eventbus/actor"
	alog "github.com/mixbee/mixbee/mixbee-eventbus/log"
	"github.com/mixbee/mixbee/account"
	"github.com/mixbee/mixbee/cmd"
	cmdcom "github.com/mixbee/mixbee/cmd/common"
	"github.com/mixbee/mixbee/cmd/utils"
	"github.com/mixbee/mixbee/common"
	"github.com/mixbee/mixbee/common/config"
	"github.com/mixbee/mixbee/common/log"
	"github.com/mixbee/mixbee/consensus"
	"github.com/mixbee/mixbee/core/genesis"
	"github.com/mixbee/mixbee/core/ledger"
	"github.com/mixbee/mixbee/events"
	hserver "github.com/mixbee/mixbee/http/base/actor"
	"github.com/mixbee/mixbee/http/jsonrpc"
	"github.com/mixbee/mixbee/http/localrpc"
	"github.com/mixbee/mixbee/http/nodeinfo"
	"github.com/mixbee/mixbee/http/restful"
	"github.com/mixbee/mixbee/http/websocket"
	"github.com/mixbee/mixbee/p2pserver"
	netreqactor "github.com/mixbee/mixbee/p2pserver/actor/req"
	p2pactor "github.com/mixbee/mixbee/p2pserver/actor/server"
	"github.com/mixbee/mixbee/txnpool"
	tc "github.com/mixbee/mixbee/txnpool/common"
	"github.com/mixbee/mixbee/txnpool/proc"
	"github.com/mixbee/mixbee/validator/stateful"
	"github.com/mixbee/mixbee/validator/stateless"
	"github.com/urfave/cli"
)

func setupAPP() *cli.App {
	// https://github.com/urfave/cli, 使用 urfave/cli 命令行工具包
	app := cli.NewApp()
	app.Usage = "Ontology CLI"
	app.Action = startOntology  // 服务的起点
	app.Version = config.Version
	app.Copyright = "Copyright in 2018 The Mixbee Authors"
	app.Commands = []cli.Command{
		cmd.AccountCommand,
		cmd.InfoCommand,
		cmd.AssetCommand,
		cmd.ContractCommand,
		cmd.ExportCommand,
	}
	app.Flags = []cli.Flag{
		//common setting
		utils.ConfigFlag,
		utils.LogLevelFlag,
		utils.DisableEventLogFlag,
		utils.DataDirFlag,
		utils.ImportEnableFlag,
		utils.ImportHeightFlag,
		utils.ImportFileFlag,
		//account setting
		utils.WalletFileFlag,
		utils.AccountAddressFlag,
		utils.AccountPassFlag,
		//consensus setting
		utils.EnableConsensusFlag,
		utils.MaxTxInBlockFlag,
		//txpool setting
		utils.GasPriceFlag,
		utils.GasLimitFlag,
		utils.PreExecEnableFlag,
		//p2p setting
		utils.ReservedPeersOnlyFlag,
		utils.ReservedPeersFileFlag,
		utils.NetworkIdFlag,
		utils.NodePortFlag,
		utils.ConsensusPortFlag,
		utils.DualPortSupportFlag,
		utils.MaxConnInBoundFlag,
		utils.MaxConnOutBoundFlag,
		utils.MaxConnInBoundForSingleIPFlag,
		//test mode setting
		utils.EnableTestModeFlag,
		utils.TestModeGenBlockTimeFlag,
		utils.ClearTestModeDataFlag,
		//rpc setting
		utils.RPCDisabledFlag,
		utils.RPCPortFlag,
		utils.RPCLocalEnableFlag,
		utils.RPCLocalProtFlag,
		//rest setting
		utils.RestfulEnableFlag,
		utils.RestfulPortFlag,
		//ws setting
		utils.WsEnabledFlag,
		utils.WsPortFlag,
	}
	app.Before = func(context *cli.Context) error {
		runtime.GOMAXPROCS(runtime.NumCPU())
		return nil
	}
	return app
}

func main() {
	if err := setupAPP().Run(os.Args); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

func startOntology(ctx *cli.Context) {
	initLog(ctx)

	// 配置初始化
	_, err := initConfig(ctx)
	if err != nil {
		log.Errorf("initConfig error:%s", err)
		return
	}

	// 账户信息初始化
	acc, err := initAccount(ctx)
	if err != nil {
		log.Errorf("initWallet error:%s", err)
		return
	}
	// 账本信息初始化
	ldg, err := initLedger(ctx)
	if err != nil {
		log.Errorf("%s", err)
		return
	}
	defer ldg.Close()
	// 区块导入
	err = importBlocks(ctx)
	if err != nil {
		log.Errorf("importBlocks error:%s", err)
		return
	}
	// 交易池初始化
	txpool, err := initTxPool(ctx)
	if err != nil {
		log.Errorf("initTxPool error:%s", err)
		return
	}
	// p2p网络节点初始化
	p2pSvr, p2pPid, err := initP2PNode(ctx, txpool)
	if err != nil {
		log.Errorf("initP2PNode error:%s", err)
		return
	}
	// 共识服务初始化
	_, err = initConsensus(ctx, p2pPid, txpool, acc)
	if err != nil {
		log.Errorf("initConsensus error:%s", err)
		return
	}
	// rpc服务初始化
	err = initRpc(ctx)
	if err != nil {
		log.Errorf("initRpc error:%s", err)
		return
	}
	// 本地rpc初始化
	err = initLocalRpc(ctx)
	if err != nil {
		log.Errorf("initLocalRpc error:%s", err)
		return
	}
	// httpRestful, Ws, Node 等初始化
	initRestful(ctx)
	initWs(ctx)
	initNodeInfo(ctx, p2pSvr)

	// 一直打印当前区块高度
	go logCurrBlockHeight()
	waitToExit()
}

func initLog(ctx *cli.Context) {
	//init log module
	logLevel := ctx.GlobalInt(utils.GetFlagName(utils.LogLevelFlag))
	alog.InitLog(log.PATH)
	log.InitLog(logLevel, log.PATH, log.Stdout)
}

func initConfig(ctx *cli.Context) (*config.MixbeeConfig, error) {
	//init ontology config from cli
	cfg, err := cmd.SetMixbeeConfig(ctx)
	if err != nil {
		return nil, err
	}
	log.Infof("Config init success")
	return cfg, nil
}

func initAccount(ctx *cli.Context) (*account.Account, error) {
	if !config.DefConfig.Consensus.EnableConsensus {
		return nil, nil
	}
	// 读取钱包文件命令
	walletFile := ctx.GlobalString(utils.GetFlagName(utils.WalletFileFlag))
	if walletFile == "" {
		return nil, fmt.Errorf("Please config wallet file using --wallet flag")
	}
	// 文件不存在
	if !common.FileExisted(walletFile) {
		return nil, fmt.Errorf("Cannot find wallet file:%s. Please create wallet first", walletFile)
	}

	// 获取账户
	acc, err := cmdcom.GetAccount(ctx)
	if err != nil {
		return nil, fmt.Errorf("get account error:%s", err)
	}
	log.Infof("Using account:%s", acc.Address.ToBase58())

	if config.DefConfig.Genesis.ConsensusType == config.CONSENSUS_TYPE_SOLO {
		// solo 模式下的处理
		curPk := hex.EncodeToString(keypair.SerializePublicKey(acc.PublicKey))
		config.DefConfig.Genesis.SOLO.Bookkeepers = []string{curPk}
	}

	log.Infof("Account init success")
	return acc, nil
}

func initLedger(ctx *cli.Context) (*ledger.Ledger, error) {
	events.Init() //Init event hub

	var err error
	dbDir := config.DefConfig.Common.DataDir + string(os.PathSeparator) + config.DefConfig.P2PNode.NetworkName

	if ctx.GlobalBool(utils.GetFlagName(utils.EnableTestModeFlag)) && ctx.GlobalBool(utils.GetFlagName(utils.ClearTestModeDataFlag)) {
		// 清除账本数据操作
		err = os.RemoveAll(dbDir)
		if err != nil {
			log.Warnf("InitLedger remove:%s error:%s", dbDir, err)
		}
	}
	// 账本数据创建
	ledger.DefLedger, err = ledger.NewLedger(dbDir)
	if err != nil {
		return nil, fmt.Errorf("NewLedger error:%s", err)
	}
	bookKeepers, err := config.DefConfig.GetBookkeepers()
	if err != nil {
		return nil, fmt.Errorf("GetBookkeepers error:%s", err)
	}
	genesisConfig := config.DefConfig.Genesis
	// 根据共识算法，创建初始区块
	genesisBlock, err := genesis.BuildGenesisBlock(bookKeepers, genesisConfig)
	if err != nil {
		return nil, fmt.Errorf("genesisBlock error %s", err)
	}
	err = ledger.DefLedger.Init(bookKeepers, genesisBlock)
	if err != nil {
		return nil, fmt.Errorf("Init ledger error:%s", err)
	}

	log.Infof("Ledger init success")
	return ledger.DefLedger, nil
}

func initTxPool(ctx *cli.Context) (*proc.TXPoolServer, error) {
	preExec := ctx.GlobalBool(utils.GetFlagName(utils.PreExecEnableFlag))
	// 启动交易池服务
	txPoolServer, err := txnpool.StartTxnPoolServer(preExec)
	if err != nil {
		return nil, fmt.Errorf("Init txpool error:%s", err)
	}
	stlValidator, _ := stateless.NewValidator("stateless_validator")
	stlValidator.Register(txPoolServer.GetPID(tc.VerifyRspActor))
	stlValidator2, _ := stateless.NewValidator("stateless_validator2")
	stlValidator2.Register(txPoolServer.GetPID(tc.VerifyRspActor))
	stfValidator, _ := stateful.NewValidator("stateful_validator")
	stfValidator.Register(txPoolServer.GetPID(tc.VerifyRspActor))

	hserver.SetTxnPoolPid(txPoolServer.GetPID(tc.TxPoolActor))
	hserver.SetTxPid(txPoolServer.GetPID(tc.TxActor))

	log.Infof("TxPool init success")
	return txPoolServer, nil
}

func initP2PNode(ctx *cli.Context, txpoolSvr *proc.TXPoolServer) (*p2pserver.P2PServer, *actor.PID, error) {
	if config.DefConfig.Genesis.ConsensusType == config.CONSENSUS_TYPE_SOLO {
		// solo 模式，不需要初始化
		return nil, nil, nil
	}
	p2p := p2pserver.NewServer()

	p2pActor := p2pactor.NewP2PActor(p2p)
	p2pPID, err := p2pActor.Start()
	if err != nil {
		return nil, nil, fmt.Errorf("p2pActor init error %s", err)
	}
	p2p.SetPID(p2pPID)
	err = p2p.Start()
	if err != nil {
		return nil, nil, fmt.Errorf("p2p service start error %s", err)
	}
	netreqactor.SetTxnPoolPid(txpoolSvr.GetPID(tc.TxActor))
	txpoolSvr.RegisterActor(tc.NetActor, p2pPID)
	hserver.SetNetServerPID(p2pPID)
	p2p.WaitForPeersStart()
	log.Infof("P2P node init success")
	return p2p, p2pPID, nil
}

func initConsensus(ctx *cli.Context, p2pPid *actor.PID, txpoolSvr *proc.TXPoolServer, acc *account.Account) (consensus.ConsensusService, error) {
	if !config.DefConfig.Consensus.EnableConsensus {
		return nil, nil
	}
	pool := txpoolSvr.GetPID(tc.TxPoolActor)

	consensusType := strings.ToLower(config.DefConfig.Genesis.ConsensusType)
	// 新建共识服务
	consensusService, err := consensus.NewConsensusService(consensusType, acc, pool, nil, p2pPid)
	if err != nil {
		return nil, fmt.Errorf("NewConsensusService:%s error:%s", consensusType, err)
	}
	// 启动共识服务
	consensusService.Start()

	netreqactor.SetConsensusPid(consensusService.GetPID())
	hserver.SetConsensusPid(consensusService.GetPID())

	log.Infof("Consensus init success")
	return consensusService, nil
}

func initRpc(ctx *cli.Context) error {
	// todo 这边不是很懂
	if !config.DefConfig.Rpc.EnableHttpJsonRpc {
		return nil
	}
	var err error
	exitCh := make(chan interface{}, 0)
	go func() {
		err = jsonrpc.StartRPCServer()
		close(exitCh)
	}()

	flag := false
	select {
	case <-exitCh:
		if !flag {
			return err
		}
	case <-time.After(time.Millisecond * 5):
		flag = true
	}
	log.Infof("Rpc init success")
	return nil
}

func initLocalRpc(ctx *cli.Context) error {
	if !ctx.GlobalBool(utils.GetFlagName(utils.RPCLocalEnableFlag)) {
		return nil
	}
	var err error
	exitCh := make(chan interface{}, 0)
	go func() {
		err = localrpc.StartLocalServer()
		close(exitCh)
	}()

	flag := false
	select {
	case <-exitCh:
		if !flag {
			return err
		}
	case <-time.After(time.Millisecond * 5):
		flag = true
	}

	log.Infof("Local rpc init success")
	return nil
}

func initRestful(ctx *cli.Context) {
	if !config.DefConfig.Restful.EnableHttpRestful {
		return
	}
	go restful.StartServer()

	log.Infof("Restful init success")
}

func initWs(ctx *cli.Context) {
	if !config.DefConfig.Ws.EnableHttpWs {
		return
	}
	websocket.StartServer()

	log.Infof("Ws init success")
}

func initNodeInfo(ctx *cli.Context, p2pSvr *p2pserver.P2PServer) {
	if config.DefConfig.P2PNode.HttpInfoPort == 0 {
		return
	}
	go nodeinfo.StartServer(p2pSvr.GetNetWork())

	log.Infof("Nodeinfo init success")
}

func importBlocks(ctx *cli.Context) error {
	if !ctx.GlobalBool(utils.GetFlagName(utils.ImportEnableFlag)) {
		return nil
	}
	importFile := ctx.GlobalString(utils.GetFlagName(utils.ImportFileFlag))
	if importFile == "" {
		return fmt.Errorf("missing import file argument")
	}
	height := ctx.GlobalUint(utils.GetFlagName(utils.ImportHeightFlag))
	return utils.ImportBlocks(importFile, uint32(height))
}

func logCurrBlockHeight() {
	ticker := time.NewTicker(config.DEFAULT_GEN_BLOCK_TIME * time.Second)
	for {
		select {
		case <-ticker.C:
			log.Infof("CurrentBlockHeight = %d", ledger.DefLedger.GetCurrentBlockHeight())
			isNeedNewFile := log.CheckIfNeedNewFile()
			if isNeedNewFile {
				log.ClosePrintLog()
				log.InitLog(int(config.DefConfig.Common.LogLevel), log.PATH, log.Stdout)
			}
		}
	}
}

func waitToExit() {
	exit := make(chan bool, 0)
	sc := make(chan os.Signal, 1)
	signal.Notify(sc, syscall.SIGINT, syscall.SIGTERM, syscall.SIGHUP)
	go func() {
		for sig := range sc {
			log.Infof("Ontology received exit signal:%v.", sig.String())
			close(exit)
			break
		}
	}()
	<-exit
}
