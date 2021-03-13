package BLC

import (
	"bytes"
	"encoding/binary"
	"encoding/json"
	"log"
	"os"
)

//实现int64转成[]byte
func IntoHex(data int64) []byte {
	buffer := new(bytes.Buffer)
	err := binary.Write(buffer, binary.BigEndian, data)
	if err != nil {
		log.Panicf("int transact to []byte failed! %v\n", err)
	}
	return buffer.Bytes()
}

// 标准JSON格式转切片
func JSONToSlice(jsonString string) []string {
	var strSlice []string
	// json
	if err := json.Unmarshal([]byte(jsonString), &strSlice); nil != err {
		log.Panicf("json to []string failed! %v\n", err)
	}
	return strSlice
}

// 参数数量检测
func IsValidArgs() {
	if len(os.Args) < 2 {
		PrintUsage()

		os.Exit(1)
	}
}
