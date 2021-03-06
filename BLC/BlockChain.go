package BLC

import (
	"encoding/hex"
	"fmt"
	"log"
	"math/big"
	"os"
	"strconv"

	"github.com/boltdb/bolt"
)

// 数据库名称
const dbName = "block.db"

// 表名称
const blockTableName = "blocks"

type BlockChain struct {
	// Block []*Block //区块的切片
	DB  *bolt.DB // 数据库对象
	Tip []byte   // 保存最新区块的哈希值
}

// 判断数据库文件是否存在
func dbExist() bool {
	if _, err := os.Stat(dbName); os.IsNotExist(err) {
		return false
	}
	return true
}

// 初始化区块链
func CreateBlockChainWithGenesisBlock(address string) *BlockChain {
	if dbExist() {
		fmt.Printf("创世区块已存在...")
		os.Exit(1)
	}

	var latesetBlockHash []byte
	// 1. 创建或者打开一个数据库
	db, err := bolt.Open(dbName, 0600, nil)
	if err != nil {
		log.Panicf("create db [%s] failed %v\n", dbName, err)
	}
	// 2. 创建桶
	db.Update(func(tx *bolt.Tx) error {
		var err error
		b := tx.Bucket([]byte(blockTableName))
		if b == nil {
			// 没找到桶
			b, err = tx.CreateBucket([]byte(blockTableName))
			if err != nil {
				log.Panicf("create db [%s] failed %v\n", blockTableName, err)
			}
		}
		// 生成一个coinbase交易
		txCoinbase := NewCoinbaseTransaction(address)
		// 生成创世区块
		genesisBlock := CreateGenesisBlock([]*Transaction{txCoinbase})
		// 存储
		// 1. key,value分别以什么数据代表
		// 2. 如何把block结构存入到数据库中--序列化
		err = b.Put(genesisBlock.Hash, genesisBlock.Serialize())
		if err != nil {
			log.Panicf("insert the genesis block failed %v\n", err)
		}
		latesetBlockHash = genesisBlock.Hash
		// 存储最新区块的哈希
		// l：latest
		err = b.Put([]byte("l"), genesisBlock.Hash)
		if err != nil {
			log.Panicf("save the latest hash of genesis block failed %v\n", err)
		}
		return nil
	})

	return &BlockChain{DB: db, Tip: latesetBlockHash}
}

// 添加区块到区块链
func (bc *BlockChain) AddBlock(txs []*Transaction) {
	// 更新区块数据
	bc.DB.Update(func(tx *bolt.Tx) error {
		// 1. 获取数据库桶
		b := tx.Bucket([]byte(blockTableName))
		if b != nil {
			// 2. 获取最后插入的区块
			blockBytes := b.Get(bc.Tip)
			// 3. 区块数据的反序列化
			latestBlock := DeserializeBlock(blockBytes)
			// height int64, prevBlockHash []byte, data []byte
			// 4. 创造一个新区块
			newBlock := NewBlock(latestBlock.Height+1, latestBlock.Hash, txs)
			// 5. 存入数据库
			err := b.Put(newBlock.Hash, newBlock.Serialize())
			if err != nil {
				log.Panicf("insert the new block to db failed %v\n", err)
			}
			// 更新最新区块的哈希（数据库）
			err = b.Put([]byte("l"), newBlock.Hash)
			if err != nil {
				log.Panicf("update the latest block hash to db failed %v", err)
			}

			// 更新区块链对象的最新区块哈希
			bc.Tip = newBlock.Hash
		}
		return nil
	})
}

// 遍历数据库，输出所有区块信息
func (bc *BlockChain) PrintChain() {
	var currentBlock *Block
	bcit := bc.Iterator()

	fmt.Println("打印区块完整信息...")

	for {
		fmt.Println("---------------------------------")
		currentBlock = bcit.Next()
		// 输出区块详情
		fmt.Printf("\tHash：%x\n", currentBlock.Hash)
		fmt.Printf("\tPrevBlockHash：%x\n", currentBlock.PrevBlockHash)
		fmt.Printf("\tTimeStamp：%v\n", currentBlock.TimeStamp)
		fmt.Printf("\tHeight：%d\n", currentBlock.Height)
		fmt.Printf("\tNonce：%d\n", currentBlock.Nonce)
		fmt.Printf("\tTransaction：%v\n", currentBlock.Txs)
		for _, tx := range currentBlock.Txs {
			fmt.Printf("\t\t----------------------------\n")
			fmt.Printf("\t\ttx-hash: %x\n", tx.TxHash)
			fmt.Printf("\t\t输入...\n")
			for _, vin := range tx.Vins {
				fmt.Printf("\t\tvin-txHash: %x\n", vin.TxHash)
				fmt.Printf("\t\tvin-vout: %v\n", vin.Vout)
				fmt.Printf("\t\tvin-scriptSig: %s\n", vin.ScriptSig)
			}
			fmt.Printf("\t\t输出...\n")
			for _, vout := range tx.Vouts {
				fmt.Printf("\t\tvout-value: %d\n", vout.Value)
				fmt.Printf("\t\tvout-scriptPubkey: %s\n", vout.ScriptPubkey)
			}
		}

		// 退出条件
		var hashInt big.Int
		hashInt.SetBytes(currentBlock.PrevBlockHash)
		if big.NewInt(0).Cmp(&hashInt) == 0 {
			// 遍历到创世区块
			break
		}
	}
}

// 获取blockchain对象
func BlockchainObject() *BlockChain {
	// 获取DB
	db, err := bolt.Open(dbName, 0600, nil)
	if err != nil {
		log.Panicf("open the db [%s] failed! %v\n", dbName, err)
	}
	// 获取TIp
	var tip []byte
	err = db.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(blockTableName))
		if b != nil {
			tip = b.Get([]byte("l"))
		}
		return nil
	})
	if err != nil {
		log.Panicf("get the blockchain object failed %v\n", err)
	}
	return &BlockChain{DB: db, Tip: tip}
}

// 实现挖矿功能
// 通过接受交易，生成区块
func (blockchain *BlockChain) MineNewBlock(from, to, amount []string) {
	var txs []*Transaction
	var block *Block

	for index, address := range from {
		value, _ := strconv.Atoi(amount[index])
		tx := NewSimpleTransaciton(address, to[index], value, blockchain, txs)
		txs = append(txs, tx)
	}

	// 从数据库中获取最新一个区块
	blockchain.DB.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(blockTableName))
		if b != nil {
			tip := b.Get([]byte("l"))
			blockBytes := b.Get(tip)
			block = DeserializeBlock(blockBytes)
		}
		return nil
	})

	// 通过数据库中最新的区块去生成新区块
	block = NewBlock(block.Height+1, block.Hash, txs)

	blockchain.DB.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(blockTableName))
		if b != nil {
			err := b.Put(block.Hash, block.Serialize())
			if err != nil {
				log.Fatalf("update the new block to db failed %v\n", err)
			}

			err = b.Put([]byte("l"), block.Hash)
			if err != nil {
				log.Fatalf("update the latest block hash to db failed %v\n", err)
			}
			blockchain.Tip = block.Hash
		}
		return nil
	})
}

// 查找指定地址的所有UTXO集合
/*
1. 遍历查找每一个区块中的每一个交易
2. 查找每一个交易中的每一个输出
3. 查找每一个交易中是否满足下列条件
	1. 属于传入的地址
	2. 是否未被花费
*/
func (blockchain *BlockChain) UnUTXOs(address string, txs []*Transaction) []*UTXO {
	var unUTXOs []*UTXO
	spentTxOutputs := blockchain.SpentOutputs(address)
	// 遍历
	bcit := blockchain.Iterator()
	// 缓存迭代
	for _, tx := range txs {
		if !tx.IsCoinbaseTransaction() {
			for _, in := range tx.Vins {
				if in.CheckPubkeyWithAddress(address) {
					key := hex.EncodeToString(in.TxHash)
					spentTxOutputs[key] = append(spentTxOutputs[key], in.Vout)
				}
			}
		}
	}
	// 遍历缓存中的UTXO
	for _, tx := range txs {
	workCacheTx:
		for index, vout := range tx.Vouts {
			if vout.CheckPubkeyWithAddress(address) {
				if len(spentTxOutputs) != 0 {
					var isUtxoTx bool
					for txHash, indexArray := range spentTxOutputs {
						txHashStr := hex.EncodeToString(tx.TxHash)
						if txHash == txHashStr {
							isUtxoTx = true
							var isSpentUTXO bool
							for _, voutIndex := range indexArray {
								if index == voutIndex {
									isSpentUTXO = true
									continue workCacheTx
								}
							}
							if isSpentUTXO == false {
								utxo := &UTXO{tx.TxHash, index, vout}
								unUTXOs = append(unUTXOs, utxo)
							}
						}
					}
					if isUtxoTx == false {
						utxo := &UTXO{tx.TxHash, index, vout}
						unUTXOs = append(unUTXOs, utxo)
					}
				} else {
					utxo := &UTXO{tx.TxHash, index, vout}
					unUTXOs = append(unUTXOs, utxo)
				}
			}
		}
	}

	// 数据库迭代
	for {
		block := bcit.Next()
		for _, tx := range block.Txs {
		work:
			for index, vout := range tx.Vouts {
				// index:当前输出在当前交易中的索引位置
				// vout:当前输出
				if vout.CheckPubkeyWithAddress(address) {
					if len(spentTxOutputs) != 0 {
						var isSpentOutput bool
						for txHash, indexArray := range spentTxOutputs {
							for _, i := range indexArray {
								if txHash == hex.EncodeToString(tx.TxHash) && index == i {
									isSpentOutput = true
									continue work
								}
							}

						}
						if isSpentOutput == false {
							unUTXOs = append(unUTXOs, &UTXO{tx.TxHash, index, vout})
						}
					} else {
						unUTXOs = append(unUTXOs, &UTXO{tx.TxHash, index, vout})
					}
				}
			}
		}

		var hashInt big.Int
		hashInt.SetBytes(block.PrevBlockHash)
		if hashInt.Cmp(big.NewInt(0)) == 0 {
			break
		}
	}

	return unUTXOs
}

// 获取指定地址上的所有"已花费"的输出
func (blockchain *BlockChain) SpentOutputs(address string) map[string][]int {
	spentOutputs := make(map[string][]int)
	bcit := blockchain.Iterator()
	for {
		block := bcit.Next()
		for _, tx := range block.Txs {
			if !tx.IsCoinbaseTransaction() {
				for _, in := range tx.Vins {
					if in.CheckPubkeyWithAddress(address) {
						key := hex.EncodeToString(in.TxHash)
						spentOutputs[key] = append(spentOutputs[key], in.Vout)
					}
				}
			}
		}

		var hashInt big.Int
		hashInt.SetBytes(block.PrevBlockHash)
		if hashInt.Cmp(big.NewInt(0)) == 0 {
			break
		}
	}
	return spentOutputs
}

// 查询余额
func (blockchain *BlockChain) getBalance(address string) int {
	var amount int
	utxos := blockchain.UnUTXOs(address, []*Transaction{})
	for _, utxo := range utxos {
		amount += utxo.Output.Value
	}
	return amount
}

// 查找指定地址的可用UTXO，超过amount就中断查找
// 更新当前数据库中指定地址的UTXO数量
// txs:缓存中的交易列表
func (blockchain *BlockChain) FindSpendableUTXO(from string, amount int, txs []*Transaction) (int, map[string][]int) {
	spendableUTXO := make(map[string][]int)

	var value int
	utxos := blockchain.UnUTXOs(from, txs)

	// 遍历UTXO集合
	for _, utxo := range utxos {
		value += utxo.Output.Value
		hash := hex.EncodeToString(utxo.TxHash)
		spendableUTXO[hash] = append(spendableUTXO[hash], utxo.Index)
		if value >= amount {
			break
		}
	}

	if value < amount {
		fmt.Printf("地址 [%s] 余额不足，当前余额 [%d]，转账金额 [%d]\n", from, value, amount)
		os.Exit(1)
	}

	return value, spendableUTXO
}
