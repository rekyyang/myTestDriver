package main

import (
	"fmt"

	"github.com/ethereum/go-ethereum/common"
)

func main() {
	fmt.Println(common.BytesToAddress([]byte{65536}))
}
