package main

import (
	"fmt"
	"io"
	"math/big"
	"reflect"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/ethereum/go-ethereum/rlp"
	jsoniter "github.com/json-iterator/go"
)

type MyBig struct {
	hexutil.Big
}

type bg struct {
	big.Int
	Nil bool
}

func (b MyBig) EncodeRLP(w io.Writer) error {
	return rlp.Encode(w, *b.Big.ToInt())
}

func (b *MyBig) DecodeRLP(s *rlp.Stream) error {
	var b_ big.Int
	err := s.Decode(&b_)
	if err != nil {
		return err
	}
	*b = MyBig{
		hexutil.Big(b_),
	}
	return nil
}

type ParityTrace struct {
	// Do not change the ordering of these fields -- allows for easier comparison with other clients
	Action      interface{}  `json:"action"` // Can be either CallTraceAction or CreateTraceAction
	BlockHash   *common.Hash `json:"blockHash,omitempty" rlp:"nil"`
	BlockNumber *uint64      `json:"blockNumber,omitempty" rlp:"nil"`
	Error       string       `json:"error,omitempty"`
	//Result              interface{}     `json:"result"`
	Subtraces           uint64          `json:"subtraces"`
	TraceAddress        []uint64        `json:"traceAddress"`
	TransactionHash     *common.Hash    `json:"transactionHash,omitempty" rlp:"nil"`
	TransactionPosition MyBig           `json:"transactionPosition,omitempty"`
	Type                string          `json:"type"`
	Address             *common.Address `rlp:"nil"`
	Big                 big.Int
	Big2                hexutil.Big
	Big3                *MyBig `rlp:"nil"`
	Bts                 hexutil.Bytes
}

type ParityTraceAlias ParityTrace

type ParityTraceWrapper struct {
	ParityTraceAlias
	Action ActionWrapper
	//Result TraceResultWrapper
}

type Action1 struct {
	Act string
	Big *MyBig
	Num uint64
}

type Action2 struct {
	Act    string
	Num    uint64
	SubAct *common.Hash `rlp:"nil"`
}

type ActionWrapper struct {
	Type string
	Raw  *rlp.RawValue `rlp:"nil"`
}

func NewActionWrapper(action interface{}) ActionWrapper {
	if action == nil {
		return ActionWrapper{
			Type: "",
			Raw:  nil,
		}
	}
	raw, _ := rlp.EncodeToBytes(action)
	switch action.(type) {
	case *Action1:
	case *Action2:
	default:
		return ActionWrapper{
			Type: "",
			Raw:  nil,
		}
	}
	return ActionWrapper{
		Type: reflect.TypeOf(action).String(),
		Raw:  (*rlp.RawValue)(&raw),
	}
}

func (a *ActionWrapper) GetAction() interface{} {
	if a.Raw == nil {
		return nil
	}
	switch a.Type {
	case reflect.TypeOf(&Action1{}).String():
		action := Action1{}
		rlp.DecodeBytes(*a.Raw, &action)
		return &action
	case reflect.TypeOf(&Action2{}).String():
		action := Action2{}
		rlp.DecodeBytes(*a.Raw, &action)
		return &action
	default:
		return nil
	}
}

func (p *ParityTrace) EncodeRLP(w io.Writer) error {
	var p_ ParityTraceWrapper
	p_.ParityTraceAlias = ParityTraceAlias(*p)
	//var err error
	p_.Action = NewActionWrapper(p.Action)
	//p_.Result, err = NewTraceResultWrapper(p.Result)
	//if err != nil {
	//	return err
	//}
	return rlp.Encode(w, p_)
}

func (p *ParityTrace) DecodeRLP(stream *rlp.Stream) error {
	var p_ ParityTraceWrapper
	err := stream.Decode(&p_)
	if err != nil {
		return err
	}
	p_.ParityTraceAlias.Action = p_.Action.GetAction()
	if err != nil {
		return err
	}
	//p_.ParityTraceAlias.Result, err = p_.Result.GetTraceResult()
	//if err != nil {
	//	return err
	//}

	*p = ParityTrace(p_.ParityTraceAlias)
	return nil
}

func main1() {
	p1 := make([]ParityTrace, 0)
	hs1 := common.BytesToHash([]byte("1512"))
	hs2 := common.BytesToHash([]byte("2451"))
	pos := uint64(1)
	pos0 := uint64(0)
	_, _ = pos, pos0
	p1 = append(p1, ParityTrace{
		Action: &Action1{
			Act: "p1_act",
			Num: 114514,
			Big: &MyBig{},
		},
		BlockHash:   &hs1,
		BlockNumber: nil,
		Error:       "brb3h",
		//Result:              nil,
		Subtraces:       2,
		TraceAddress:    nil,
		TransactionHash: nil,
		//TransactionPosition: MyBig{},
		Type: "brr",
		Bts:  hexutil.Bytes("tetrameter"),
		Big3: &MyBig{},
	})
	hs := common.BytesToHash([]byte("homo"))
	_ = hs
	bg := big.Int{}
	bg.SetUint64(1453)
	bg2 := big.Int{}
	bg2.SetUint64(1566)
	p1 = append(p1, ParityTrace{
		//Action: &Action2{
		//	Act:    "p2_act",
		//	Num:    1919810,
		//	SubAct: &hs,
		//},
		BlockHash:   nil,
		BlockNumber: nil,
		Error:       "agegegeg",
		//Result:              nil,
		Subtraces:       2,
		TraceAddress:    nil,
		TransactionHash: &hs2,
		//TransactionPosition: nil,
		Type: "br1sr1112ffdfasdfasdfewfwewerwerq",
		Big:  bg,
		Big2: hexutil.Big(bg),
		Big3: &MyBig{hexutil.Big(bg2)},
	})
	raw, err := rlp.EncodeToBytes(p1)
	fmt.Println(err)
	fmt.Println(string(raw))
	//fmt.Println(raw)
	var p2 []ParityTrace
	err = rlp.DecodeBytes(raw, &p2)
	fmt.Println(err)
	fmt.Println(p1)
	fmt.Println("\n")
	fmt.Println(p2)
	fmt.Println("\n")
	fmt.Println(p2[1].Big.String())

	js, _ := jsoniter.Marshal(p1)
	fmt.Println(string(js))
}

type TestStruct struct {
	Name      string `rlp:"-"`
	Position1 *uint64
}

func main() {
	var t1 TestStruct
	t1.Name = "a"
	t1.Position1 = nil

	var t2 TestStruct
	t2.Name = "a"
	p := uint64(0)
	t2.Position1 = &p

	raw1, err := rlp.EncodeToBytes(t1)
	fmt.Println(err)
	raw2, err := rlp.EncodeToBytes(t2)
	fmt.Println(err)

	fmt.Println(raw1)
	fmt.Println(raw2)
}
