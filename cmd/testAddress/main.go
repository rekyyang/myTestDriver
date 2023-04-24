package main

import (
	"fmt"

	"github.com/ethereum/go-ethereum/common"
)

func main() {
	addr := common.HexToAddress("0x9af768b63d6cc359729d33759f7cbd7cf1526105")
	addr2_ := []byte("\232\367h\266=l\303Yr\235\063u\237|\275|\361Ra\005")
	addr2 := common.BytesToAddress(addr2_)

	fmt.Println(string(addr[:]))
	fmt.Println(addr2)
}
