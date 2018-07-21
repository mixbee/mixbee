package main


import (
	"fmt"
	"github.com/mixbee/mixbee/cmd/abi"
	cmdcom "github.com/mixbee/mixbee/cmd/common"
	cmdsvr "github.com/mixbee/mixbee/cmd/sigsvr"
	cmdsvrcom "github.com/mixbee/mixbee/cmd/sigsvr/common"
	"github.com/mixbee/mixbee/cmd/utils"
	"github.com/mixbee/mixbee/common"
	"github.com/mixbee/mixbee/common/config"
	"github.com/mixbee/mixbee/common/log"
	"github.com/urfave/cli"
	"os"
	"os/signal"
	"runtime"
	"syscall"
)

func setupSigSvr() *cli.App {
	app := cli.NewApp()
	app.Usage = "Mixbee Sig server"
	app.Action = startSigSvr
	app.Version = config.Version
	app.Copyright = "Copyright in 2018 The Mixbee Authors"
	app.Flags = []cli.Flag{
		utils.LogLevelFlag,
		//account setting
		utils.WalletFileFlag,
		utils.AccountAddressFlag,
		utils.AccountPassFlag,
		//cli setting
		utils.CliRpcPortFlag,
		utils.CliABIPathFlag,
	}
	app.Before = func(context *cli.Context) error {
		runtime.GOMAXPROCS(runtime.NumCPU())
		return nil
	}
	return app
}

func startSigSvr(ctx *cli.Context) {
	logLevel := ctx.GlobalInt(utils.GetFlagName(utils.LogLevelFlag))
	log.InitLog(logLevel, log.PATH, log.Stdout)

	walletFile := ctx.GlobalString(utils.GetFlagName(utils.WalletFileFlag))
	if walletFile == "" {
		log.Infof("Please specificed wallet file using --wallet flag")
		return
	}
	if !common.FileExisted(walletFile) {
		log.Infof("Cannot find wallet file:%s. Please create wallet first", walletFile)
		return
	}
	acc, err := cmdcom.GetAccount(ctx)
	if err != nil {
		log.Infof("GetAccount error:%s", err)
		return
	}
	log.Infof("Using account:%s", acc.Address.ToBase58())

	rpcPort := ctx.Uint(utils.GetFlagName(utils.CliRpcPortFlag))
	if rpcPort == 0 {
		log.Infof("Please using sig server port by --%s flag", utils.GetFlagName(utils.CliRpcPortFlag))
		return
	}
	cmdsvrcom.DefAccount = acc
	go cmdsvr.DefCliRpcSvr.Start(rpcPort)

	abiPath := ctx.GlobalString(utils.GetFlagName(utils.CliABIPathFlag))
	abi.DefAbiMgr.Init(abiPath)

	log.Infof("Sig server init success")
	log.Infof("Sig server listing on: %d", rpcPort)

	exit := make(chan bool, 0)
	sc := make(chan os.Signal, 1)
	signal.Notify(sc, syscall.SIGINT, syscall.SIGTERM, syscall.SIGHUP)
	go func() {
		for sig := range sc {
			log.Infof("Sig server received exit signal:%v.", sig.String())
			close(exit)
			break
		}
	}()
	<-exit
}

func main() {
	if err := setupSigSvr().Run(os.Args); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

