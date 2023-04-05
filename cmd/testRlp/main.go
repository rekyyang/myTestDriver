package main

import (
	"fmt"
	"math/big"
	"strings"

	"github.com/ethereum/go-ethereum/core/types"
)

type Source1 struct {
	Integer     *big.Int
	String      string
	Float       *big.Float
	Withdrawals []*types.Withdrawal `rlp:"optional"`
	//Withdrawals common.Hash `rlp:"optional"`
}

type Target1 struct {
	Integer *big.Int
	String  string
	Float   *big.Float
	//Withdrawals []*types.Withdrawal `rlp:"optional"`
	//Hash3   common.Hash `rlp:"optional"`
}

func main() {
	//s1 := Source1{
	//	Integer:     big.NewInt(123),
	//	String:      "aaaa",
	//	Float:       big.NewFloat(3.14),
	//	Withdrawals: make([]*types.Withdrawal, 0),
	//	//Withdrawals: make([]string, 0),
	//}
	//s1.Withdrawals = append(s1.Withdrawals, &types.Withdrawal{
	//	Index:     2,
	//	Validator: 0,
	//	Address:   common.Address{},
	//	Amount:    0,
	//})
	////s1.Withdrawals = append(s1.Withdrawals, "1234")
	//
	//_, r1, err := rlp.EncodeToReader(s1)
	//fmt.Println(err)
	//var t1 Target1
	//err = rlp.Decode(r1, &t1)
	//fmt.Println(err)
	//fmt.Println(t1)
	k := "arbitrum-nitro-eth_call"
	index := strings.LastIndex(k, "-")
	if index == -1 {
		return
	}
	pkg := k[:index]
	method := k[index+1:]
	fmt.Println(1)
	fmt.Println(pkg)
	fmt.Println(2)
	fmt.Println(method)
}
