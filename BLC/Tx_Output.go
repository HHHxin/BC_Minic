package BLC

// 交易输出管理

// 输出结构
type TxOutput struct {
	Value        int    // 金额
	ScriptPubkey string // 用户名（地址）
}

// 验证当前输出是否属于指定地址
func (txOutput *TxOutput) CheckPubkeyWithAddress(address string) bool {
	return address == txOutput.ScriptPubkey
}
