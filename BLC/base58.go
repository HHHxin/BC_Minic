package BLC

import (
	"bytes"
	"math/big"
)

// base58编码实现

var b58Alphabet = []byte("" +
	"123456789" +
	"abcdefghijkmnopqrstuvwxyz" +
	"ABCDEFGHJKLMNPQRSTUVWXYZ")

// 编码函数
func Base58Encode(input []byte) []byte {
	var result []byte
	x := big.NewInt(0).SetBytes(input)
	// 求余的基本长度
	base := big.NewInt(int64(len(b58Alphabet)))
	// 求余和商
	// 判断条件，出掉的最终结果是否为0
	mod := &big.Int{}
	zero := big.NewInt(0)
	for x.Cmp(zero) != 0 {
		x.DivMod(x, base, mod)
		// 倒序
		result = append(result, b58Alphabet[mod.Int64()])
	}
	Reverse(result)
	result = append([]byte{b58Alphabet[0]}, result...)
	return result
}

func Reverse(data []byte) {
	for i, j := 0, len(data)-1; i < j; i, j = i+1, j-1 {
		data[i], data[j] = data[j], data[i]
	}
}

// 解码函数
func Base58Decode(input []byte) []byte {
	result := big.NewInt(0)
	zeroBytes := 1
	data := input[zeroBytes:]
	for _, b := range data {
		charIndex := bytes.IndexByte(b58Alphabet, b)
		result.Mul(result, big.NewInt(58))
		result.Add(result, big.NewInt(int64(charIndex)))
	}

	decoded := result.Bytes()
	return decoded
}
