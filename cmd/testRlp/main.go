package main

import (
	"fmt"
	"math/big"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/node-real/go-pkg/rlp"
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
	s1 := Source1{
		Integer:     big.NewInt(123),
		String:      "aaaa",
		Float:       big.NewFloat(3.14),
		Withdrawals: make([]*types.Withdrawal, 0),
		//Withdrawals: make([]string, 0),
	}
	s1.Withdrawals = append(s1.Withdrawals, &types.Withdrawal{
		Index:     2,
		Validator: 0,
		Address:   common.Address{},
		Amount:    0,
	})
	//s1.Withdrawals = append(s1.Withdrawals, "1234")

	_, r1, err := rlp.EncodeToReader(s1)
	fmt.Println(err)
	var t1 Target1
	err = rlp.Decode(r1, &t1)
	fmt.Println(err)
	fmt.Println(t1)
}
