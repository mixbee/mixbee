package cmd

import (
	"fmt"
	"github.com/mixbee/mixbee/account"
	cmdcom "github.com/mixbee/mixbee/cmd/common"
	"github.com/mixbee/mixbee/cmd/utils"
	nutils "github.com/mixbee/mixbee/smartcontract/service/native/utils"
	"github.com/urfave/cli"
	"strconv"
	"strings"
)

var AssetCommand = cli.Command{
	Name:        "asset",
	Usage:       "Handle assets",
	Description: "Asset management commands can check account balance, MBC/MBG transfers, extract MBGs, and view unbound MBGs, and so on.",
	Subcommands: []cli.Command{
		{
			Action:      transfer,
			Name:        "transfer",
			Usage:       "Transfer mbc or mbg to another account",
			ArgsUsage:   " ",
			Description: "Transfer mbc or mbg to another account. If from address does not specified, using default account",
			Flags: []cli.Flag{
				utils.RPCPortFlag,
				utils.TransactionGasPriceFlag,
				utils.TransactionGasLimitFlag,
				utils.TransactionAssetFlag,
				utils.TransactionFromFlag,
				utils.TransactionToFlag,
				utils.TransactionAmountFlag,
				utils.WalletFileFlag,
			},
		},
		{
			Action:    approve,
			Name:      "approve",
			ArgsUsage: " ",
			Usage:     "Approve another user can transfer asset",
			Flags: []cli.Flag{
				utils.RPCPortFlag,
				utils.TransactionGasPriceFlag,
				utils.TransactionGasLimitFlag,
				utils.ApproveAssetFlag,
				utils.ApproveAssetFromFlag,
				utils.ApproveAssetToFlag,
				utils.ApproveAmountFlag,
				utils.WalletFileFlag,
			},
		},
		{
			Action:    transferFrom,
			Name:      "transferfrom",
			ArgsUsage: " ",
			Usage:     "Using to transfer asset after approve",
			Flags: []cli.Flag{
				utils.RPCPortFlag,
				utils.TransactionGasPriceFlag,
				utils.TransactionGasLimitFlag,
				utils.ApproveAssetFlag,
				utils.TransferFromSenderFlag,
				utils.ApproveAssetFromFlag,
				utils.ApproveAssetToFlag,
				utils.TransferFromAmountFlag,
				utils.WalletFileFlag,
			},
		},
		{
			Action:    getBalance,
			Name:      "balance",
			Usage:     "Show balance of mbc and mbg of specified account",
			ArgsUsage: "<address|label|index>",
			Flags: []cli.Flag{
				utils.RPCPortFlag,
				utils.WalletFileFlag,
			},
		},
		{
			Action: getAllowance,
			Name:   "allowance",
			Usage:  "Show approve balance of mbc or mbg of specified account",
			Flags: []cli.Flag{
				utils.RPCPortFlag,
				utils.ApproveAssetFlag,
				utils.ApproveAssetFromFlag,
				utils.ApproveAssetToFlag,
				utils.WalletFileFlag,
			},
		},
		{
			Action:    unboundMbg,
			Name:      "unboundmbg",
			Usage:     "Show the balance of unbound MBG",
			ArgsUsage: "<address|label|index>",
			Flags: []cli.Flag{
				utils.RPCPortFlag,
				utils.WalletFileFlag,
			},
		},
		{
			Action:    withdrawMbg,
			Name:      "withdrawMbg",
			Usage:     "Withdraw MBG",
			ArgsUsage: "<address|label|index>",
			Flags: []cli.Flag{
				utils.RPCPortFlag,
				utils.TransactionGasPriceFlag,
				utils.TransactionGasLimitFlag,
				utils.WalletFileFlag,
			},
		},
		{
			Action:    getKey,
			Name:      "getkey",
			Usage:     "get key from mixTest",
			ArgsUsage: "<key>",
			Flags: []cli.Flag{
				utils.RPCPortFlag,
				utils.WalletFileFlag,
			},
		},
		{
			Action: setKey,
			Name:   "setkey",
			Usage:  "set key from mixTest",
			Flags: []cli.Flag{
				utils.RPCPortFlag,
				utils.TransactionGasPriceFlag,
				utils.TransactionGasLimitFlag,
				utils.TransactionFromFlag,
				utils.MixTestKeyFlag,
				utils.MixTestValueFlag,
				utils.WalletFileFlag,
			},
		},
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
		//{
		//	Action: crossRelease,
		//	Name:   "crelease",
		//	Usage:  "cross chain release tx asset 2 toAddress",
		//	Flags: []cli.Flag{
		//		utils.RPCPortFlag,
		//		utils.WalletFileFlag,
		//		utils.TransactionGasPriceFlag,
		//		utils.TransactionGasLimitFlag,
		//		utils.TransactionFromFlag,
		//		utils.CrossChainSeqIdFlag,
		//	},
		//},
		{
			Action:    crossPairQuery,
			Name:      "crossPairQuery",
			Usage:     "cross chain pair evidence tx query by seqId",
			ArgsUsage: "<seqId>",
			Flags: []cli.Flag{
				utils.RPCPortFlag,
			},
		},
	},
}

func transfer(ctx *cli.Context) error {
	SetRpcPort(ctx)
	if !ctx.IsSet(utils.GetFlagName(utils.TransactionToFlag)) ||
		!ctx.IsSet(utils.GetFlagName(utils.TransactionFromFlag)) ||
		!ctx.IsSet(utils.GetFlagName(utils.TransactionAmountFlag)) {
		fmt.Println("Missing from, to or amount flag\n")
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

	var amount uint64
	amountStr := ctx.String(utils.TransactionAmountFlag.Name)
	switch strings.ToLower(asset) {
	case "mbc":
		amount = utils.ParseMbc(amountStr)
		amountStr = utils.FormatMbc(amount)
	case "mbg":
		amount = utils.ParseMbg(amountStr)
		amountStr = utils.FormatMbg(amount)
	default:
		return fmt.Errorf("unsupport asset:%s", asset)
	}

	err = utils.CheckAssetAmount(asset, amount)
	if err != nil {
		return err
	}

	gasPrice := ctx.Uint64(utils.TransactionGasPriceFlag.Name)
	gasLimit := ctx.Uint64(utils.TransactionGasLimitFlag.Name)

	var signer *account.Account
	signer, err = cmdcom.GetAccount(ctx, fromAddr)
	if err != nil {
		return fmt.Errorf("GetAccount error:%s", err)
	}
	txHash, err := utils.Transfer(gasPrice, gasLimit, signer, asset, fromAddr, toAddr, amount)
	if err != nil {
		return fmt.Errorf("Transfer error:%s", err)
	}
	fmt.Printf("Transfer %s\n", strings.ToUpper(asset))
	fmt.Printf("  From:%s\n", fromAddr)
	fmt.Printf("  To:%s\n", toAddr)
	fmt.Printf("  Amount:%s\n", amountStr)
	fmt.Printf("  TxHash:%s\n", txHash)
	fmt.Printf("\nTip:\n")
	fmt.Printf("  Using './mixbee info status %s' to query transaction status\n", txHash)
	return nil
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

func getKey(ctx *cli.Context) error {
	SetRpcPort(ctx)
	if ctx.NArg() < 2 {
		fmt.Println("Missing argument. mixTest key .\n")
		cli.ShowSubcommandHelp(ctx)
		return nil
	}

	addrArg := ctx.Args().First()
	keyArg := ctx.Args().Get(1)
	value, err := utils.GetKey(addrArg + keyArg)
	if err != nil {
		return err
	}

	fmt.Printf("GetKey:%s\n", keyArg)
	fmt.Printf("Value:%v\n", value.Value)
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

func setKey(ctx *cli.Context) error {
	SetRpcPort(ctx)

	if !ctx.IsSet(utils.GetFlagName(utils.MixTestKeyFlag)) ||
		!ctx.IsSet(utils.GetFlagName(utils.TransactionFromFlag)) ||
		!ctx.IsSet(utils.GetFlagName(utils.MixTestValueFlag)) {
		fmt.Println("Missing from, to or amount flag\n")
		cli.ShowSubcommandHelp(ctx)
		return nil
	}

	from := ctx.String(utils.TransactionFromFlag.Name)
	fromAddr, err := cmdcom.ParseAddress(from, ctx)
	if err != nil {
		return fmt.Errorf("Pxarse from address:%s error:%s", from, err)
	}
	key := ctx.String(utils.MixTestKeyFlag.Name)
	value := ctx.String(utils.MixTestValueFlag.Name)

	gasPrice := ctx.Uint64(utils.TransactionGasPriceFlag.Name)
	gasLimit := ctx.Uint64(utils.TransactionGasLimitFlag.Name)

	var signer *account.Account
	signer, err = cmdcom.GetAccount(ctx, fromAddr)
	if err != nil {
		return fmt.Errorf("GetAccount error:%s", err)
	}
	txHash, err := utils.SetKey(gasPrice, gasLimit, signer, key, value)
	if err != nil {
		return fmt.Errorf("Transfer error:%s", err)
	}
	fmt.Printf("setKey %s\n", strings.ToUpper("mixT"))
	fmt.Printf("  From:%s\n", fromAddr)
	fmt.Printf("  Key:%s\n", key)
	fmt.Printf("  Value:%s\n", value)
	fmt.Printf("  TxHash:%s\n", txHash)
	fmt.Printf("\nTip:\n")
	fmt.Printf("  Using './mixbee info status %s' to query transaction status\n", txHash)
	return nil
}

func getBalance(ctx *cli.Context) error {
	SetRpcPort(ctx)
	if ctx.NArg() < 1 {
		fmt.Println("Missing argument. Account address, label or index expected.\n")
		cli.ShowSubcommandHelp(ctx)
		return nil
	}

	addrArg := ctx.Args().First()
	accAddr, err := cmdcom.ParseAddress(addrArg, ctx)
	if err != nil {
		return err
	}
	balance, err := utils.GetBalance(accAddr)
	if err != nil {
		return err
	}

	mbg, err := strconv.ParseUint(balance.Mbg, 10, 64)
	if err != nil {
		return err
	}
	fmt.Printf("BalanceOf:%s\n", accAddr)
	fmt.Printf("  MBC:%s\n", balance.Mbc)
	fmt.Printf("  MBG:%s\n", utils.FormatMbg(mbg))
	return nil
}

func getAllowance(ctx *cli.Context) error {
	SetRpcPort(ctx)
	from := ctx.String(utils.GetFlagName(utils.ApproveAssetFromFlag))
	to := ctx.String(utils.GetFlagName(utils.ApproveAssetToFlag))
	if from == "" || to == "" {
		fmt.Printf("Missing approve from or to argument\n")
		cli.ShowSubcommandHelp(ctx)
		return nil
	}
	asset := ctx.String(utils.GetFlagName(utils.ApproveAssetFlag))
	if asset == "" {
		asset = utils.ASSET_MBC
	}
	fromAddr, err := cmdcom.ParseAddress(from, ctx)
	if err != nil {
		return err
	}
	toAddr, err := cmdcom.ParseAddress(to, ctx)
	if err != nil {
		return err
	}
	balanceStr, err := utils.GetAllowance(asset, fromAddr, toAddr)
	if err != nil {
		return err
	}
	switch strings.ToLower(asset) {
	case "mbc":
	case "mbg":
		balance, err := strconv.ParseUint(balanceStr, 10, 64)
		if err != nil {
			return err
		}
		balanceStr = utils.FormatMbg(balance)
	default:
		return fmt.Errorf("unsupport asset:%s", asset)
	}
	fmt.Printf("Allowance:%s\n", asset)
	fmt.Printf("  From:%s\n", fromAddr)
	fmt.Printf("  To:%s\n", toAddr)
	fmt.Printf("  Balance:%s\n", balanceStr)
	return nil
}

func approve(ctx *cli.Context) error {
	SetRpcPort(ctx)
	asset := ctx.String(utils.GetFlagName(utils.ApproveAssetFlag))
	from := ctx.String(utils.GetFlagName(utils.ApproveAssetFromFlag))
	to := ctx.String(utils.GetFlagName(utils.ApproveAssetToFlag))
	amountStr := ctx.String(utils.GetFlagName(utils.ApproveAmountFlag))
	if asset == "" ||
		from == "" ||
		to == "" ||
		amountStr == "" {
		fmt.Printf("Missing asset, from, to, or amount argument\n")
		cli.ShowSubcommandHelp(ctx)
		return nil
	}
	fromAddr, err := cmdcom.ParseAddress(from, ctx)
	if err != nil {
		return err
	}
	toAddr, err := cmdcom.ParseAddress(to, ctx)
	if err != nil {
		return err
	}
	var amount uint64
	switch strings.ToLower(asset) {
	case "mbc":
		amount = utils.ParseMbc(amountStr)
		amountStr = utils.FormatMbc(amount)
	case "mbg":
		amount = utils.ParseMbg(amountStr)
		amountStr = utils.FormatMbg(amount)
	default:
		return fmt.Errorf("unsupport asset:%s", asset)
	}

	err = utils.CheckAssetAmount(asset, amount)
	if err != nil {
		return err
	}

	gasPrice := ctx.Uint64(utils.TransactionGasPriceFlag.Name)
	gasLimit := ctx.Uint64(utils.TransactionGasLimitFlag.Name)

	var signer *account.Account
	signer, err = cmdcom.GetAccount(ctx, fromAddr)
	if err != nil {
		return fmt.Errorf("GetAccount error:%s", err)
	}

	txHash, err := utils.Approve(gasPrice, gasLimit, signer, asset, fromAddr, toAddr, amount)
	if err != nil {
		return fmt.Errorf("approve error:%s", err)
	}

	fmt.Printf("Approve:\n")
	fmt.Printf("  Asset:%s\n", asset)
	fmt.Printf("  From:%s\n", fromAddr)
	fmt.Printf("  To:%s\n", toAddr)
	fmt.Printf("  Amount:%s\n", amountStr)
	fmt.Printf("  TxHash:%s\n", txHash)
	fmt.Printf("\nTip:\n")
	fmt.Printf("  Using './mixbee info status %s' to query transaction status\n", txHash)
	return nil
}

func transferFrom(ctx *cli.Context) error {
	SetRpcPort(ctx)
	asset := ctx.String(utils.GetFlagName(utils.ApproveAssetFlag))
	from := ctx.String(utils.GetFlagName(utils.ApproveAssetFromFlag))
	to := ctx.String(utils.GetFlagName(utils.ApproveAssetToFlag))
	amountStr := ctx.String(utils.GetFlagName(utils.TransferFromAmountFlag))
	if asset == "" ||
		from == "" ||
		to == "" ||
		amountStr == "" {
		fmt.Printf("Missing asset, from, to, or amount argument\n")
		cli.ShowSubcommandHelp(ctx)
		return nil
	}
	fromAddr, err := cmdcom.ParseAddress(from, ctx)
	if err != nil {
		return err
	}
	toAddr, err := cmdcom.ParseAddress(to, ctx)
	if err != nil {
		return err
	}

	var sendAddr string
	sender := ctx.String(utils.GetFlagName(utils.TransferFromSenderFlag))
	if sender == "" {
		sendAddr = toAddr
	} else {
		sendAddr, err = cmdcom.ParseAddress(sender, ctx)
		if err != nil {
			return err
		}
	}

	var signer *account.Account
	signer, err = cmdcom.GetAccount(ctx, sendAddr)
	if err != nil {
		return fmt.Errorf("GetAccount error:%s", err)
	}

	var amount uint64
	switch strings.ToLower(asset) {
	case "mbc":
		amount = utils.ParseMbc(amountStr)
		amountStr = utils.FormatMbc(amount)
	case "mbg":
		amount = utils.ParseMbg(amountStr)
		amountStr = utils.FormatMbg(amount)
	default:
		return fmt.Errorf("unsupport asset:%s", asset)
	}

	err = utils.CheckAssetAmount(asset, amount)
	if err != nil {
		return err
	}

	gasPrice := ctx.Uint64(utils.TransactionGasPriceFlag.Name)
	gasLimit := ctx.Uint64(utils.TransactionGasLimitFlag.Name)

	txHash, err := utils.TransferFrom(gasPrice, gasLimit, signer, asset, sendAddr, fromAddr, toAddr, amount)
	if err != nil {
		return err
	}

	fmt.Printf("Transfer from:\n")
	fmt.Printf("  Asset:%s\n", asset)
	fmt.Printf("  Sender:%s\n", sendAddr)
	fmt.Printf("  From:%s\n", fromAddr)
	fmt.Printf("  To:%s\n", toAddr)
	fmt.Printf("  Amount:%s\n", amountStr)
	fmt.Printf("  TxHash:%s\n", txHash)
	fmt.Printf("\nTip:\n")
	fmt.Printf("  Using './mixbee info status %s' to query transaction status\n", txHash)
	return nil
}

func unboundMbg(ctx *cli.Context) error {
	SetRpcPort(ctx)
	if ctx.NArg() < 1 {
		fmt.Println("Missing argument. Account address, label or index expected.\n")
		cli.ShowSubcommandHelp(ctx)
		return nil
	}
	addrArg := ctx.Args().First()
	accAddr, err := cmdcom.ParseAddress(addrArg, ctx)
	if err != nil {
		return err
	}
	fromAddr := nutils.MbcContractAddress.ToBase58()
	balanceStr, err := utils.GetAllowance("mbg", fromAddr, accAddr)
	if err != nil {
		return err
	}
	balance, err := strconv.ParseUint(balanceStr, 10, 64)
	if err != nil {
		return err
	}
	balanceStr = utils.FormatMbg(balance)
	fmt.Printf("Unbound MBG:\n")
	fmt.Printf("  Account:%s\n", accAddr)
	fmt.Printf("  MBG:%s\n", balanceStr)
	return nil
}

func withdrawMbg(ctx *cli.Context) error {
	SetRpcPort(ctx)
	if ctx.NArg() < 1 {
		fmt.Println("Missing argument. Account address, label or index expected.\n")
		cli.ShowSubcommandHelp(ctx)
		return nil
	}
	addrArg := ctx.Args().First()
	accAddr, err := cmdcom.ParseAddress(addrArg, ctx)
	if err != nil {
		return err
	}
	fromAddr := nutils.MbcContractAddress.ToBase58()
	balance, err := utils.GetAllowance("mbg", fromAddr, accAddr)
	if err != nil {
		return err
	}

	amount, err := strconv.ParseUint(balance, 10, 64)
	if err != nil {
		return err
	}
	if amount <= 0 {
		return fmt.Errorf("Don't have unbound mbg\n")
	}

	var signer *account.Account
	signer, err = cmdcom.GetAccount(ctx, accAddr)
	if err != nil {
		return fmt.Errorf("GetAccount error:%s", err)
	}

	gasPrice := ctx.Uint64(utils.TransactionGasPriceFlag.Name)
	gasLimit := ctx.Uint64(utils.TransactionGasLimitFlag.Name)

	txHash, err := utils.TransferFrom(gasPrice, gasLimit, signer, "mbg", accAddr, fromAddr, accAddr, amount)
	if err != nil {
		return err
	}

	fmt.Printf("Withdraw MBG:\n")
	fmt.Printf("  Account:%s\n", accAddr)
	fmt.Printf("  Amount:%s\n", utils.FormatMbg(amount))
	fmt.Printf("  TxHash:%s\n", txHash)
	fmt.Printf("\nTip:\n")
	fmt.Printf("  Using './mixbee info status %s' to query transaction status\n", txHash)
	return nil
}
