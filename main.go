package main

import (
	"bitcoin/BLC"
)

func main() {

	// bc := BLC.CreateBlockChainWithGenesisBlock()

	// fmt.Printf("blockchain:%v\n", bc.Block[0])
	// bc.AddBlock(bc.Block[len(bc.Block)-1].Height+1, bc.Block[len(bc.Block)-1].Hash, []byte("alice send 10 btc to bob"))
	// bc.AddBlock(bc.Block[len(bc.Block)-1].Height+1, bc.Block[len(bc.Block)-1].Hash, []byte("alice send 5 btc to troytan"))

	// bc.AddBlock([]byte("a send 100 btc to b"))
	// bc.AddBlock([]byte("b send 100 btc to c"))

	// bc.PrintChain()

	cli := BLC.CLI{}
	cli.Run()

	// result := BLC.Base58Encode([]byte("this is the example"))
	// fmt.Printf("result: %s\n", result)
	// result = BLC.Base58Decode([]byte("1nj2SLMErZakmBni8xhSXtimREn"))
	// fmt.Printf("result: %s\n", result)
}
