package BLC

// UTXO结构
type UTXO struct {
	TxHash []byte    // UTXO对应的交易哈希
	Index  int       // UTXO及其所属交易的输出列表中的索引
	Output *TxOutput // Output本身
}
