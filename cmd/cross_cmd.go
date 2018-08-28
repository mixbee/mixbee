package cmd

import (
	"fmt"
	"github.com/mixbee/mixbee/account"
	cmdcom "github.com/mixbee/mixbee/cmd/common"
	"github.com/mixbee/mixbee/cmd/utils"
	"github.com/urfave/cli"
	"strconv"
	"strings"
)

var CrossChainCommand = cli.Command{
	Name:        "cross",
	Usage:       "Handle cross chain",
	Description: "",
	Subcommands: []cli.Command{
		{
			Action:    crossQuery,
			Name:      "cquery",
			Usage:     "cross chain tx query by seqId",
			ArgsUsage: "<seqId>",
			Flags: []cli.Flag{
				utils.RPCPortFlag,
			},
		},
		{
			Action:    crossHistory,
			Name:      "chistory",
			Usage:     "cross chain tx query history by from",
			ArgsUsage: "<adress>",
			Flags: []cli.Flag{
				utils.RPCPortFlag,
			},
		},
		{
			Action: crossTranfer,
			Name:   "ctransfer",
			Usage:  "cross chain transfer asset",
			Flags: []cli.Flag{
				utils.RPCPortFlag,
				utils.TransactionGasPriceFlag,
				utils.TransactionGasLimitFlag,
				utils.TransactionFromFlag,
				utils.TransactionToFlag,
				utils.CrossChainAValueFlag,
				utils.CrossChainBValueFlag,
				utils.CrossChainBChainIdFlag,
				utils.CrossChainAChainIdFlag,
				utils.CrossChainDelayTimeFlag,
				utils.WalletFileFlag,
				utils.CrossChainNonceFlag,
				utils.CrossChainVerifyPublicKeyFlag,
				utils.AccountPassFlag,
			},
		},
		{
			Action: crossUnlock,
			Name:   "cunlock",
			Usage:  "cross chain lock expire tx",
			Flags: []cli.Flag{
				utils.RPCPortFlag,
				utils.WalletFileFlag,
				utils.TransactionGasPriceFlag,
				utils.TransactionGasLimitFlag,
				utils.TransactionFromFlag,
				utils.CrossChainSeqIdFlag,
			},
		},
		{
			Action:    crossPairQuery,
			Name:      "crossPairQuery",
			Usage:     "cross chain pair evidence tx query by seqId",
			ArgsUsage: "<seqId>",
			Flags: []cli.Flag{
				utils.RPCPortFlag,
			},
		},
		{
			Action: crossPaidDeposit,
			Name:   "cdeposit",
			Usage:  "cross chain paid deposit",
			Flags: []cli.Flag{
				utils.RPCPortFlag,
				utils.TransactionGasPriceFlag,
				utils.TransactionGasLimitFlag,
				utils.TransactionFromFlag,
				utils.WalletFileFlag,
				utils.CrossChainVerifyPublicKeyFlag,
				utils.AccountPassFlag,
				utils.TransactionAmountFlag,
			},
		},
	},
}

func crossTranfer(ctx *cli.Context) error {
	SetRpcPort(ctx)
	if !ctx.IsSet(utils.GetFlagName(utils.TransactionToFlag)) ||
		!ctx.IsSet(utils.GetFlagName(utils.TransactionFromFlag)) ||
		!ctx.IsSet(utils.GetFlagName(utils.CrossChainVerifyPublicKeyFlag)) ||
		!ctx.IsSet(utils.GetFlagName(utils.CrossChainBValueFlag)) ||
		!ctx.IsSet(utils.GetFlagName(utils.CrossChainBChainIdFlag)) ||
		!ctx.IsSet(utils.GetFlagName(utils.CrossChainAChainIdFlag)) ||
		!ctx.IsSet(utils.GetFlagName(utils.CrossChainAValueFlag)) {
		fmt.Println("Missing from,to,aValue,bValue,achainid or bChainId flag\n")
		cli.ShowSubcommandHelp(ctx)
		return nil
	}

	asset := ctx.String(utils.GetFlagName(utils.TransactionAssetFlag))
	if asset == "" {
		asset = utils.ASSET_MBC
	}
	from := ctx.String(utils.TransactionFromFlag.Name)
	fromAddr, err := cmdcom.ParseAddress(from, ctx)
	if err != nil {
		return fmt.Errorf("Parse from address:%s error:%s", from, err)
	}
	to := ctx.String(utils.TransactionToFlag.Name)
	toAddr, err := cmdcom.ParseAddress(to, ctx)
	if err != nil {
		return fmt.Errorf("Parse to address:%s error:%s", to, err)
	}

	var aAmount uint64
	aAmountStr := ctx.String(utils.CrossChainAValueFlag.Name)
	aAmount = utils.ParseMbc(aAmountStr)
	aAmountStr = utils.FormatMbc(aAmount)

	var bAmount uint64
	bAmountStr := ctx.String(utils.CrossChainBValueFlag.Name)
	bAmount = utils.ParseMbc(bAmountStr)
	bAmountStr = utils.FormatMbc(bAmount)

	err = utils.CheckAssetAmount(asset, aAmount)
	if err != nil {
		return err
	}

	bchainIdStr := ctx.String(utils.CrossChainBChainIdFlag.Name)
	bChainId, err := strconv.ParseUint(bchainIdStr, 10, 32)
	if err != nil {
		return fmt.Errorf("Parse bchainId:%s error:%s", bchainIdStr, err)
	}

	achainIdStr := ctx.String(utils.CrossChainAChainIdFlag.Name)
	aChainId, err := strconv.ParseUint(achainIdStr, 10, 32)
	if err != nil {
		return fmt.Errorf("Parse achainId:%s error:%s", bchainIdStr, err)
	}

	delayTime := ctx.Uint64(utils.CrossChainDelayTimeFlag.Name)

	nonce := ctx.Uint64(utils.CrossChainNonceFlag.Name)

	gasPrice := ctx.Uint64(utils.TransactionGasPriceFlag.Name)
	gasLimit := ctx.Uint64(utils.TransactionGasLimitFlag.Name)

	verifyPublicKey := ctx.String(utils.CrossChainVerifyPublicKeyFlag.Name)

	var signer *account.Account
	signer, err = cmdcom.GetAccount(ctx, fromAddr)
	if err != nil {
		return fmt.Errorf("GetAccount error:%s", err)
	}
	txHash, seqId, err := utils.CrossTransfer(gasPrice, gasLimit, signer, asset, toAddr, aAmount, bAmount, aChainId, bChainId, delayTime, nonce, verifyPublicKey)
	if err != nil {
		return fmt.Errorf("Transfer error:%s", err)
	}
	fmt.Printf("cross chain Transfer %s\n", strings.ToUpper(asset))
	fmt.Printf("  From:%s\n", fromAddr)
	fmt.Printf("  To:%s\n", toAddr)
	fmt.Printf("  a Amount:%s\n", aAmountStr)
	fmt.Printf("  b Amount:%s\n", bAmountStr)
	fmt.Printf("  achainID :%s\n", achainIdStr)
	fmt.Printf("  bchainID :%s\n", bchainIdStr)
	fmt.Printf("  seqId:%s\n", seqId)
	fmt.Printf("  TxHash:%s\n", txHash)
	fmt.Printf("\nTip:\n")
	fmt.Printf("  Using './mixbee info status %s' to query transaction status\n", txHash)
	return nil
}

func crossUnlock(ctx *cli.Context) error {
	SetRpcPort(ctx)

	if !ctx.IsSet(utils.GetFlagName(utils.TransactionFromFlag)) ||
		!ctx.IsSet(utils.GetFlagName(utils.CrossChainSeqIdFlag)) {
		fmt.Println("Missing from or seqid flag\n")
		cli.ShowSubcommandHelp(ctx)
		return nil
	}

	asset := ctx.String(utils.GetFlagName(utils.TransactionAssetFlag))
	if asset == "" {
		asset = utils.ASSET_MBC
	}
	from := ctx.String(utils.TransactionFromFlag.Name)
	fromAddr, err := cmdcom.ParseAddress(from, ctx)

	seqId := ctx.String(utils.CrossChainSeqIdFlag.Name)
	gasPrice := ctx.Uint64(utils.TransactionGasPriceFlag.Name)
	gasLimit := ctx.Uint64(utils.TransactionGasLimitFlag.Name)

	var signer *account.Account
	signer, err = cmdcom.GetAccount(ctx, fromAddr)
	if err != nil {
		return fmt.Errorf("GetAccount error:%s", err)
	}

	txHash, err := utils.CrossUnlockOrRelease(gasPrice, gasLimit, signer, seqId, utils.CONTRACT_CROSS_UNLOCK)
	if err != nil {
		return fmt.Errorf("crossUnlock error:%s", err)
	}

	fmt.Printf("unlock %s\n", strings.ToUpper("Cross chain"))
	fmt.Printf("TxHash:%s\n", txHash)
	fmt.Printf("  From:%s\n", fromAddr)
	fmt.Printf("  seqId:%s\n", seqId)
	fmt.Printf("\nTip:\n")
	fmt.Printf("  Using './mixbee info status %s' to query transaction status\n", txHash)
	return nil
}

func crossRelease(ctx *cli.Context) error {
	SetRpcPort(ctx)

	if !ctx.IsSet(utils.GetFlagName(utils.TransactionFromFlag)) ||
		!ctx.IsSet(utils.GetFlagName(utils.CrossChainSeqIdFlag)) {
		fmt.Println("Missing from or seqid flag\n")
		cli.ShowSubcommandHelp(ctx)
		return nil
	}

	asset := ctx.String(utils.GetFlagName(utils.TransactionAssetFlag))
	if asset == "" {
		asset = utils.ASSET_MBC
	}
	from := ctx.String(utils.TransactionFromFlag.Name)
	fromAddr, err := cmdcom.ParseAddress(from, ctx)

	seqId := ctx.String(utils.CrossChainSeqIdFlag.Name)
	gasPrice := ctx.Uint64(utils.TransactionGasPriceFlag.Name)
	gasLimit := ctx.Uint64(utils.TransactionGasLimitFlag.Name)

	var signer *account.Account
	signer, err = cmdcom.GetAccount(ctx, fromAddr)
	if err != nil {
		return fmt.Errorf("GetAccount error:%s", err)
	}

	txHash, err := utils.CrossUnlockOrRelease(gasPrice, gasLimit, signer, seqId, utils.CONTRACT_CROSS_RELEASE)
	if err != nil {
		return fmt.Errorf("crossRelease error:%s", err)
	}

	fmt.Printf("crossRelease %s\n", strings.ToUpper("Cross chain"))
	fmt.Printf("TxHash:%s\n", txHash)
	fmt.Printf("  From:%s\n", fromAddr)
	fmt.Printf("  seqId:%s\n", seqId)
	fmt.Printf("\nTip:\n")
	fmt.Printf("  Using './mixbee info status %s' to query transaction status\n", txHash)
	return nil
}

func crossQuery(ctx *cli.Context) error {
	SetRpcPort(ctx)
	if ctx.NArg() < 1 {
		fmt.Println("Missing argument. cross chain seqId.\n")
		cli.ShowSubcommandHelp(ctx)
		return nil
	}

	seqArg := ctx.Args().First()
	value, err := utils.CrossQuery(seqArg)
	if err != nil {
		return err
	}

	fmt.Printf("crossQuery seqId:%s\n", seqArg)
	fmt.Println("cross Info:", value.Value)
	return nil
}

func crossPairQuery(ctx *cli.Context) error {
	SetRpcPort(ctx)
	if ctx.NArg() < 1 {
		fmt.Println("Missing argument. cross chain seqId.\n")
		cli.ShowSubcommandHelp(ctx)
		return nil
	}

	seqArg := ctx.Args().First()
	value, err := utils.CrossPairEvidenceQuery(seqArg)
	if err != nil {
		return err
	}

	fmt.Printf("crossPairEvidenceQuery seqId:%s\n", seqArg)
	fmt.Println("cross Info:", value.Value)
	return nil
}

func crossHistory(ctx *cli.Context) error {
	SetRpcPort(ctx)
	if ctx.NArg() < 1 {
		fmt.Println("Missing argument. cross chain sender address.\n")
		cli.ShowSubcommandHelp(ctx)
		return nil
	}

	seqArg := ctx.Args().First()
	value, err := utils.CrossHistory(seqArg)
	if err != nil {
		return err
	}

	fmt.Printf("crossHistory seqId:%s\n", seqArg)
	fmt.Printf("cross chain seqIds:%+v\n", value.Value)
	return nil
}

func crossPaidDeposit(ctx *cli.Context) error {
	SetRpcPort(ctx)

	if !ctx.IsSet(utils.GetFlagName(utils.TransactionFromFlag)) ||
		!ctx.IsSet(utils.GetFlagName(utils.TransactionAmountFlag)) ||
		!ctx.IsSet(utils.GetFlagName(utils.CrossChainVerifyPublicKeyFlag)) {
		fmt.Println("Missing from or seqid flag\n")
		cli.ShowSubcommandHelp(ctx)
		return nil
	}

	from := ctx.String(utils.TransactionFromFlag.Name)
	fromAddr, err := cmdcom.ParseAddress(from, ctx)
	verifyPublicKey := ctx.String(utils.CrossChainVerifyPublicKeyFlag.Name)

	gasPrice := ctx.Uint64(utils.TransactionGasPriceFlag.Name)
	gasLimit := ctx.Uint64(utils.TransactionGasLimitFlag.Name)
	amountStr := ctx.String(utils.TransactionAmountFlag.Name)
	amount := utils.ParseMbc(amountStr)
	var signer *account.Account
	signer, err = cmdcom.GetAccount(ctx, fromAddr)
	if err != nil {
		return fmt.Errorf("GetAccount error:%s", err)
	}

	txHash, err := utils.CrossVerifyNodePaidDeposit(gasPrice, gasLimit, signer, verifyPublicKey, amount)
	if err != nil {
		return fmt.Errorf("crossUnlock error:%s", err)
	}

	fmt.Printf("crossPaidDeposit %s\n", strings.ToUpper("Cross chain"))
	fmt.Printf("TxHash:%s\n", txHash)
	fmt.Printf("  From:%s\n", fromAddr)
	fmt.Printf("  publicKey:%s\n", verifyPublicKey)
	fmt.Printf("\nTip:\n")
	fmt.Printf("  Using './mixbee info status %s' to query transaction status\n", txHash)
	return nil
}
