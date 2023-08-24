package main

import (
	"io"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/rlp"
)

type Prestates struct {
	Accounts           []*AccountPrestate                  // this field is only used for encode/decode
	accountPrestateMap map[common.Address]*AccountPrestate `rlp:"-"`
}

func (p *Prestates) Account(addr common.Address) (*AccountPrestate, bool) {
	v, ok := p.accountPrestateMap[addr]
	return v, ok
}

func (p *Prestates) SetAccount(addr common.Address, acc *AccountPrestate) {
	p.accountPrestateMap[addr] = acc
}

func NewPrestate() *Prestates {
	return &Prestates{
		Accounts:           []*AccountPrestate{},
		accountPrestateMap: map[common.Address]*AccountPrestate{},
	}
}

func (p *Prestates) Copy() *Prestates {
	ret := NewPrestate()

	// `Accounts` field is skip
	for addr, acc := range p.accountPrestateMap {
		newAcc := acc.Copy()
		ret.SetAccount(addr, newAcc)
	}

	return ret
}

type PrestatesRLP Prestates

func (p *Prestates) EncodeRLP(w io.Writer) error {
	p.Accounts = make([]*AccountPrestate, len(p.accountPrestateMap))
	idx := 0
	for _, v := range p.accountPrestateMap {
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

	p.accountPrestateMap = make(map[common.Address]*AccountPrestate, len(p.Accounts))
	for _, acc := range p.Accounts {
		p.accountPrestateMap[acc.Address] = acc
	}

	return nil
}
