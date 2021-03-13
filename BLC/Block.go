package BLC

import (
	"bytes"
	"crypto/sha256"
	"encoding/gob"
	"log"
	"time"
)

//区块基本结构与功能文件

//实现一个最基本的区块结构
type Block struct {
	TimeStamp     int64          //区块时间戳
	Hash          []byte         //当前区块哈希
	PrevBlockHash []byte         //父区块哈希
	Height        int64          //区块高度
	Txs           []*Transaction //交易数据(交易列表)
	Nonce         int64          //运行pow是的修改值
}

//新建区块
func NewBlock(height int64, prevBlockHash []byte, txs []*Transaction) *Block {
	var block Block

	block = Block{
		TimeStamp:     time.Now().Unix(),
		PrevBlockHash: prevBlockHash,
		Height:        height,
		Txs:           txs,
	}
	// block.SetHash()
	// 通过POW生成新的哈希
	pow := NewProofOfWork(&block)
	// 执行pow算法
	hash, nonce := pow.Run()
	block.Hash = hash
	block.Nonce = int64(nonce)
	return &block
}

// 计算并设置区块哈希
// func (b *Block) SetHash() {
// 	//调用sha256生成哈希
// 	//实现int64-》hash
// 	timeStampBytes := IntoHex(b.TimeStamp)
// 	heightByte := IntoHex(b.Height)
// 	blockByte := bytes.Join([][]byte{
// 		timeStampBytes,
// 		heightByte,
// 		b.PrevBlockHash,
// 		b.Data,
// 	}, []byte{})
// 	hash := sha256.Sum256(blockByte)
// 	b.Hash = hash[:]
// }

// 生成创世区块
func CreateGenesisBlock(txs []*Transaction) *Block {
	return NewBlock(1, nil, txs)
}

// 区块结构序列化
func (block *Block) Serialize() []byte {
	var buffer bytes.Buffer
	// 新建编码对象
	encoder := gob.NewEncoder(&buffer)
	// 编码序列化
	if err := encoder.Encode(block); nil != err {
		log.Panicf("serialize the block to []byte failed %v\n", err)
	}
	return buffer.Bytes()
}

// 区块数据反序列化
func DeserializeBlock(blockBytes []byte) *Block {
	var block Block
	// 新建decoder对象
	decoder := gob.NewDecoder(bytes.NewBuffer(blockBytes))
	if err := decoder.Decode(&block); err != nil {
		log.Panicf("deserialize the block to []byte failed %v\n", err)
	}
	return &block
}

// 把指定区块的所有交易结构序列化(类似Merkle的哈希计算方法)
func (block *Block) HashTransaction() []byte {
	var txHashes [][]byte
	// 将指定区块中所有交易哈希进行拼接
	for _, tx := range block.Txs {
		txHashes = append(txHashes, tx.TxHash)
	}
	txHash := sha256.Sum256(bytes.Join(txHashes, []byte{}))
	return txHash[:]
}
