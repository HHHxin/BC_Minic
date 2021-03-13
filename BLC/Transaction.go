package BLC

import (
	"bytes"
	"crypto/sha256"
	"encoding/gob"
	"encoding/hex"
	"fmt"
	"log"
)

// 交易管理文件

// 定义一个交易结构
type Transaction struct {
	TxHash []byte      // 交易哈希（标识）
	Vins   []*TxInput  // 输入列表
	Vouts  []*TxOutput // 输出列表
}

// 实现coinbase交易
func NewCoinbaseTransaction(address string) *Transaction {
	// 输入
	// coinbase特点
	// txHash:nil
	// vout:-1
	// ScriptSig:系统奖励
	txInput := &TxInput{
		TxHash:    []byte{},
		Vout:      -1,
		ScriptSig: "system reward",
	}
	// 输出
	// value：
	// address：
	txOutput := &TxOutput{
		Value:        10,
		ScriptPubkey: address,
	}
	txCoinbase := &Transaction{
		TxHash: nil,
		Vins:   []*TxInput{txInput},
		Vouts:  []*TxOutput{txOutput},
	}
	// 交易哈希生成
	txCoinbase.HashTransaction()

	return txCoinbase
}

// 生成交易哈希（交易序列化）,并不是真正的Merkle树的结构的生成哈希
func (tx *Transaction) HashTransaction() {
	var result bytes.Buffer
	// 设置编码对象
	encoder := gob.NewEncoder(&result)
	if err := encoder.Encode(tx); err != nil {
		log.Panicf("tx Hash encoded failed %v\n", err)
	}

	// 生成哈希值
	hash := sha256.Sum256(result.Bytes())
	tx.TxHash = hash[:]
}

// 生成普通转账交易
func NewSimpleTransaciton(from string, to string, amount int, blockchain *BlockChain, txs []*Transaction) *Transaction {
	var txInputs []*TxInput
	var txOutputs []*TxOutput

	// 获取UTXO
	money, utoxsDic := blockchain.FindSpendableUTXO(from, amount, txs)
	fmt.Printf("money:%v\n", money)
	// 输入
	for txHash, indexArry := range utoxsDic {
		txHashBytes, err := hex.DecodeString(txHash)
		if err != nil {
			log.Panicf("decode string to []byte failed! %v\n", err)
		}

		// 遍历索引列表
		for _, index := range indexArry {
			txInputs = append(txInputs, &TxInput{txHashBytes, index, from})
		}
	}

	// 输出（源）
	txOutput := &TxOutput{amount, to}
	txOutputs = append(txOutputs, txOutput)

	// 输出（找零）
	if money > amount {
		txOutput = &TxOutput{money - amount, from}
		txOutputs = append(txOutputs, txOutput)
	} else {
		log.Panicf("余额不足...\n")
	}

	tx := Transaction{nil, txInputs, txOutputs}
	tx.HashTransaction()
	return &tx
}

// 判断指定的交易是否是一个coinbase交易
func (tx *Transaction) IsCoinbaseTransaction() bool {
	return tx.Vins[0].Vout == -1 && len(tx.Vins[0].TxHash) == 0
}
