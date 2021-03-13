package BLC

import (
	"flag"
	"fmt"
	"log"
	"os"
)

// 对blockchain的命令行操作进行管理

// client对象
type CLI struct {
}

// 用法展示
func PrintUsage() {
	fmt.Println("Usage:")
	// 初始化
	fmt.Printf("\tcreateblockchain --address Address -- 创建区块链\n")
	// 添加区块
	fmt.Printf("\taddblock --txs Transaction -- 添加区块\n")
	// 打印完整的区块信息
	fmt.Printf("\tprintchain -- 输出区块链信息\n")
	// 通过命令转账
	fmt.Printf("\tsend -from FROM -to TO -amount AMOUNT -- 发起转账\n")
	fmt.Printf("\t\t-from FROM -- 转账源地址\n")
	fmt.Printf("\t\t-to TO -- 转账目标地址\n")
	fmt.Printf("\t\t-amount AMOUNT -- 转账金额\n")
	fmt.Printf("\tgetbalance -address FROM -- 查询指定地址的余额\n")
	fmt.Printf("\t查询余额参数说明\n")
	fmt.Printf("\t\t-address --查询余额的地址\n")
}

// 查询余额
func (cli *CLI) getBalance(from string) {
	if !dbExist() {
		fmt.Printf("数据库不存在...")
		os.Exit(1)
	}
	blockchain := BlockchainObject()
	defer blockchain.DB.Close()
	amount := blockchain.getBalance(from)
	fmt.Printf("\t地址 [%s] 的余额：[%d]\n", from, amount)
}

// 发起交易
func (cli *CLI) send(from, to, amount []string) {
	if !dbExist() {
		fmt.Printf("数据库不存在...")
		os.Exit(1)
	}
	blockchain := BlockchainObject()
	defer blockchain.DB.Close()
	if len(from) != len(to) || len(from) != len(amount) {
		fmt.Printf("交易参数输入有误，请检查一致性...\n")
		os.Exit(1)
	}
	blockchain.MineNewBlock(from, to, amount)
}

// 初始化区块链
func (cli *CLI) createBlockchain(address string) {
	CreateBlockChainWithGenesisBlock(address)
}

// 添加区块
func (cli *CLI) addBlock(txs []*Transaction) {
	if !dbExist() {
		fmt.Printf("数据库不存在...")
		os.Exit(1)
	}
	// 获取bc对象
	blockchain := BlockchainObject()
	blockchain.AddBlock(txs)
}

// 打印完整的区块信息
func (cli *CLI) printChain() {
	if !dbExist() {
		fmt.Printf("数据库不存在...")
		os.Exit(1)
	}
	// 获取bc对象
	blockchain := BlockchainObject()
	blockchain.PrintChain()
}

// 命令行运行函数
func (cli *CLI) Run() {
	// 检测参数数量
	IsValidArgs()
	// 新建相关命令
	// 添加区块
	addBlockCmd := flag.NewFlagSet("addblock", flag.ExitOnError)
	// 输出区块链完整信息
	printChainCmd := flag.NewFlagSet("printchain", flag.ExitOnError)
	// 创建区块链
	createBLCWithGenesisBlockCmd := flag.NewFlagSet("createblockchain", flag.ExitOnError)
	// 发起交易
	sendCmd := flag.NewFlagSet("send", flag.ExitOnError)
	// 查询余额
	getBalanceCmd := flag.NewFlagSet("getbalance", flag.ExitOnError)

	// 数据参数处理
	// 添加区块
	flagAddBlockArg := addBlockCmd.String("data", "sent 100 btc to player", "添加区块")
	// 创建区块链是指定矿工地址
	flagCreateBlockchainArg := createBLCWithGenesisBlockCmd.String("address", "troytan", "指定接收系统奖励的矿工地址")
	// 发起交易
	flagSendFromArg := sendCmd.String("from", "[\"troytan\",\"troytan\",\"troytan\"]", "转账源地址")
	flagSendToArg := sendCmd.String("to", "[\"aaa\",\"bbb\",\"ccc\"]", "转账目标地址")
	flagSendAmountArg := sendCmd.String("amount", "[\"1\",\"2\",\"3\"]", "转账金额")
	// 查询余额
	flagGetBalanceArg := getBalanceCmd.String("address", "", "余额")

	// 判断命令
	switch os.Args[1] {
	case "getbalance":
		if err := getBalanceCmd.Parse(os.Args[2:]); err != nil {
			log.Panicf("parse getbalanceCmd failed! %v\n", err)
		}
	case "send":
		if err := sendCmd.Parse(os.Args[2:]); err != nil {
			log.Panicf("parse sendCmd failed! %v\n", err)
		}
	case "addblock":
		if err := addBlockCmd.Parse(os.Args[2:]); err != nil {
			log.Panicf("parse addBlockCmd failed! %v\n", err)
		}
	case "printchain":
		if err := printChainCmd.Parse(os.Args[2:]); err != nil {
			log.Panicf("parse printChainCmd failed! %v\n", err)
		}
	case "createblockchain":
		if err := createBLCWithGenesisBlockCmd.Parse(os.Args[2:]); err != nil {
			log.Panicf("parse createBLCWithGenesisBlockCmd failed! %v\n", err)
		}
	default:
		PrintUsage()
		os.Exit(1)
	}

	// 添加区块
	if addBlockCmd.Parsed() {
		if *flagAddBlockArg == "" {
			PrintUsage()
			os.Exit(1)
		}
		cli.addBlock([]*Transaction{})
	}
	// 输出区块链信息
	if printChainCmd.Parsed() {
		cli.printChain()
	}
	// 创建区块链命令
	if createBLCWithGenesisBlockCmd.Parsed() {
		if *flagCreateBlockchainArg == "" {
			PrintUsage()
			os.Exit(1)
		}
		cli.createBlockchain(*flagCreateBlockchainArg)
	}
	// 发起交易
	if sendCmd.Parsed() {
		if *flagSendFromArg == "" {
			fmt.Printf("源地址不能为空...\n")
			PrintUsage()
			os.Exit(1)
		}
		if *flagSendToArg == "" {
			fmt.Printf("目标地址不能为空...%s\n", *flagSendFromArg)
			PrintUsage()
			os.Exit(1)
		}
		if *flagSendAmountArg == "" {
			fmt.Printf("转账金额不能为空...\n")
			PrintUsage()
			os.Exit(1)
		}
		fmt.Printf("\tFROM:[%s]\n", JSONToSlice(*flagSendFromArg))
		fmt.Printf("\tTO:[%s]\n", JSONToSlice(*flagSendToArg))
		fmt.Printf("\tAMOUNT:[%s]\n", JSONToSlice(*flagSendAmountArg))
		cli.send(JSONToSlice(*flagSendFromArg), JSONToSlice(*flagSendToArg), JSONToSlice(*flagSendAmountArg))
	}
	// 查询余额
	if getBalanceCmd.Parsed() {
		if *flagGetBalanceArg == "" {
			fmt.Printf("查询地址不能为空\n")
			PrintUsage()
			os.Exit(1)
		}
		cli.getBalance(*flagGetBalanceArg)
	}
}
