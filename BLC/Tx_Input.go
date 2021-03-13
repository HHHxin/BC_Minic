package BLC

// 交易输入管理

// 输入结构
type TxInput struct {
	TxHash    []byte // 交易哈希(未花费UTXO交易哈希)
	Vout      int    // 引用上一笔交易的输出的索引
	ScriptSig string // 签名（模拟）
}

// 验证引用的地址是否匹配
func (txInput *TxInput) CheckPubkeyWithAddress(address string) bool {
	return address == txInput.ScriptSig
}
