package BLC

import (
	"log"

	"github.com/boltdb/bolt"
)

// 区块链迭代器管理文件

// 迭代器基本结构
type BlockChainIterator struct {
	DB          *bolt.DB //迭代目标
	CurrentHash []byte   //单次迭代目标的哈希
}

// 创建迭代器对象
func (bc *BlockChain) Iterator() *BlockChainIterator {
	return &BlockChainIterator{DB: bc.DB, CurrentHash: bc.Tip}
}

// 实现迭代器函数next，获取到每一个区块
func (bcit *BlockChainIterator) Next() *Block {
	var block *Block
	err := bcit.DB.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(blockTableName))
		if b != nil {
			currentBlockBytes := b.Get(bcit.CurrentHash)
			block = DeserializeBlock(currentBlockBytes)
			// 更新迭代器中区块的哈希值
			bcit.CurrentHash = block.PrevBlockHash
		}
		return nil
	})
	if err != nil {
		log.Panicf("iterator the db failed! %v\n", err)
	}
	return block
}
