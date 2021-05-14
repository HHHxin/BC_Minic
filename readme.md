## windows下运行（命令行）123
项目根目录下  123
* go build -o bc.exe main.go
## 功能：
* bc.exe
    * 查看所有功能
* bc.exe createblockchain [--address Address]
    * 创建区块链，并创建coinbase交易，输出地址为Address
* bc.exe printchain
    * 打印所有区块链信息
* bc.exe getbalance -address Address
    * 查询指定地址Address的余额
* bc.exe send -from From -to TO -amount AMOUNT
    * FROM地址向TO地址转账金额AMOUNT，变量格式：from："[\"Alice\",\"Bob\",\"troytan\"]"，可进行多笔交易。
