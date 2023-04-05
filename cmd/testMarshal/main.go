package main

import (
	"fmt"

	"github.com/ethereum/go-ethereum/common"
	jsoniter "github.com/json-iterator/go"
)

type Person struct {
	Name  string       `json:"name,omitempty"`
	Age   string       `json:"age,omitempty"`
	Desc  []string     `json:"desc"`
	Phone *common.Hash `json:"phone"`
}

func (p Person) MarshalJSON() ([]byte, error) {
	type Alias Person
	if p.Phone == nil {
		return jsoniter.Marshal(&struct {
			Desc  []string     `json:"desc,omitempty"`
			Phone *common.Hash `json:"phone,omitempty"`
			*Alias
		}{
			Desc:  nil,
			Phone: nil,
			Alias: (*Alias)(&p),
		})
	} else {
		return jsoniter.Marshal(&struct {
			Desc  []string     `json:"desc"`
			Phone *common.Hash `json:"phone"`
			*Alias
		}{
			Desc:  p.Desc,
			Phone: p.Phone,
			Alias: (*Alias)(&p),
		})
	}
}

func main() {
	var desc0 []string = make([]string, 0)
	desc1 := []string{"123", "345"}
	phone := "11111111"
	_ = &desc0
	_ = &desc1
	_ = &phone
	p := Person{
		Name:  "aa",
		Age:   "12",
		Desc:  desc0,
		Phone: &common.Hash{},
	}
	bt, _ := jsoniter.Marshal(p)
	fmt.Println(string(bt))
}
