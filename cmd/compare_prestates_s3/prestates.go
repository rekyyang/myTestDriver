package main

import (
	"fmt"
	"io"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/ethereum/go-ethereum/rlp"
)

type Prestates struct {
	Accounts           []*AccountPrestate                  // this field is only used for encode/decode
	AccountPrestateMap map[common.Address]*AccountPrestate `rlp:"-"`
}

func (p *Prestates) Account(addr common.Address) (*AccountPrestate, bool) {
	v, ok := p.AccountPrestateMap[addr]
	return v, ok
}

func (p *Prestates) SetAccount(addr common.Address, acc *AccountPrestate) {
	p.AccountPrestateMap[addr] = acc
}

func NewPrestate() *Prestates {
	return &Prestates{
		Accounts:           []*AccountPrestate{},
		AccountPrestateMap: map[common.Address]*AccountPrestate{},
	}
}

func (p *Prestates) Copy() *Prestates {
	ret := NewPrestate()

	// `Accounts` field is skip
	for addr, acc := range p.AccountPrestateMap {
		newAcc := acc.Copy()
		ret.SetAccount(addr, newAcc)
	}

	return ret
}

type PrestatesRLP Prestates

func (p *Prestates) EncodeRLP(w io.Writer) error {
	p.Accounts = make([]*AccountPrestate, len(p.AccountPrestateMap))
	idx := 0
	for _, v := range p.AccountPrestateMap {
		p.Accounts[idx] = v
		idx++
	}

	return rlp.Encode(w, (*PrestatesRLP)(p))
}

func (p *Prestates) DecodeRLP(stream *rlp.Stream) error {
	pRLP := (*PrestatesRLP)(p)
	if err := stream.Decode(pRLP); err != nil {
		return err
	}

	p.AccountPrestateMap = make(map[common.Address]*AccountPrestate, len(p.Accounts))
	for _, acc := range p.Accounts {
		p.AccountPrestateMap[acc.Address] = acc
	}

	return nil
}

func ComparePrestates(pExp, pAct *Prestates) {
	r1 := make(map[string]struct{})
	r2 := make(map[string]struct{})

	for k, _ := range pExp.AccountPrestateMap {
		if _, ok := pAct.AccountPrestateMap[k]; !ok {
			r1[hexutil.Encode(k[:])] = struct{}{}
		}
	}

	for k, _ := range pAct.AccountPrestateMap {
		if _, ok := pExp.AccountPrestateMap[k]; !ok {
			r2[hexutil.Encode(k[:])] = struct{}{}
		}
	}

	fmt.Println("=========addr 1==========")
	fmt.Println(r1)
	fmt.Println("=========addr 2==========")
	fmt.Println(r2)
	fmt.Println("=========================")

	for k, v1 := range pExp.AccountPrestateMap {
		addr := hexutil.Encode(k[:])
		if v2, ok := pAct.AccountPrestateMap[k]; ok {
			for ks, vv1 := range v1.Storages {
				key := hexutil.Encode(ks[:])
				if vv2, ok := v2.Storages[ks]; ok {
					if vv1.String() != vv2.String() {
						fmt.Println(fmt.Sprintf("addr[%s], key[%s] storage mismatch, exp[%s], act[%s]", addr, key, vv1.String(), vv2.String()))
					}
				} else {
					fmt.Println(fmt.Sprintf("addr [%s], storage key[%s] missed", addr, key))
				}
			}
		}
	}
	fmt.Println("=========================")
}
