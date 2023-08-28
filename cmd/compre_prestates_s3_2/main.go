package main

import (
	"context"
	"fmt"
	"myTestDriver/cmd/compre_prestates_s3_2/new_models"
	"myTestDriver/cmd/compre_prestates_s3_2/old_models"
	"time"

	awsconfig "github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/ethereum/go-ethereum/rlp"
	jsoniter "github.com/json-iterator/go"
	jsonrpc "github.com/node-real/go-pkg/jsonrpc2"
)

const URL1 = "https://bsc-mainnet.nodereal.io/v1/4e7dd5235d434c5a837f7e48e9af9b4d"
const URL2 = "http://bsc-mainnet-coordinator-ap.nodereal.internal"

var _ = s3.Client{}

type Comparer struct {
	S3Cli        *s3.Client
	ExpBucket    string
	ActBucket    string
	ExpFolder    string
	ActFolderEvm string
	ActFolderGat string

	RpcCli jsonrpc.Client
}

type BlkWithHash struct {
	Hash        string
	BlockNumber string
}

func (c *Comparer) Compare(key, keyOld string, bn uint64) error {
	keyActEvm := fmt.Sprintf("%s/%s", c.ActFolderEvm, key)
	keyActGat := fmt.Sprintf("%s/%s", c.ActFolderGat, key)
	keyExp := fmt.Sprintf("%s/%s", c.ExpFolder, keyOld)

	objActEvm, err := c.S3Cli.GetObject(context.Background(), &s3.GetObjectInput{
		Bucket: &c.ActBucket,
		Key:    &keyActEvm,
	})

	if err != nil {
		return err
	}

	objActGat, err := c.S3Cli.GetObject(context.Background(), &s3.GetObjectInput{
		Bucket: &c.ActBucket,
		Key:    &keyActGat,
	})

	if err != nil {
		return err
	}

	objExp, err := c.S3Cli.GetObject(context.Background(), &s3.GetObjectInput{
		Bucket: &c.ExpBucket,
		Key:    &keyExp,
	})

	if err != nil {
		return err
	}

	pActEvm := new_models.Prestates{}
	pActGat := new_models.Prestates{}
	pExp := old_models.Prestates{}

	err = rlp.Decode(objActEvm.Body, &pActEvm)
	if err != nil {
		return err
	}
	err = rlp.Decode(objActGat.Body, &pActGat)
	if err != nil {
		return err
	}
	err = rlp.Decode(objExp.Body, &pExp)
	if err != nil {
		return err
	}

	fmt.Println()
	fmt.Println("=========1 pActGat pActEvm===========")
	new_models.ComparePrestates(&pActGat, &pActEvm)
	fmt.Println()
	fmt.Println("=========2 pExp pActEvm==============")
	ComparePrestatesV2(&pExp, &pActEvm, bn)
	fmt.Println()
	fmt.Println("=========3 pExp pActGat==============")
	ComparePrestatesV2(&pExp, &pActGat, bn)
	fmt.Println("=====================================")
	fmt.Println()
	fmt.Println()
	fmt.Println()
	return nil
}

func (c *Comparer) GetLatestBlockKey() (string, string, uint64) {
	resp, err := c.RpcCli.Call(context.Background(), jsonrpc.NewRequest("23333", "eth_blockNumber"))
	if err != nil {
		fmt.Println(fmt.Sprintf("failed to get latest blk: %s", err.Error()))
		return "", "", 0
	}

	if resp.Error != nil {
		fmt.Println(fmt.Sprintf("failed to get latest blk: %s", resp.Error.Error()))
		return "", "", 0
	}

	bnStr := ""
	err = jsoniter.Unmarshal(resp.Result, &bnStr)

	bnStr = hexutil.EncodeUint64(hexutil.MustDecodeUint64(bnStr) - 10)
	fmt.Println()
	fmt.Println()
	fmt.Println()
	fmt.Println()
	fmt.Println()
	fmt.Println()
	fmt.Println(fmt.Sprintf("blockNumber : %s", bnStr))

	blk := BlkWithHash{}

	resp, err = c.RpcCli.Call(context.Background(), jsonrpc.NewRequest("23334", "eth_getBlockByNumber", bnStr, false))
	if err != nil {
		fmt.Println(fmt.Sprintf("failed to get latest blk: %s", err.Error()))
		return "", "", 0
	}

	if resp.Error != nil {
		fmt.Println(fmt.Sprintf("failed to get latest blk: %s", resp.Error.Error()))
		return "", "", 0
	}

	err = jsoniter.Unmarshal(resp.Result, &blk)
	if err != nil {
		fmt.Println(fmt.Sprintf("failed to unmarshal blk: %s", err.Error()))
		return "", "", 0
	}

	fmt.Println(blk)
	bn, _ := hexutil.DecodeUint64(bnStr)
	fmt.Println(fmt.Sprintf("key_new %s", fmt.Sprintf("%010d_%s", bn, blk.Hash)))
	fmt.Println(fmt.Sprintf("key_old %s", fmt.Sprintf("%d_%s", bn, blk.Hash)))
	return fmt.Sprintf("%010d_%s", bn, blk.Hash), fmt.Sprintf("%d_%s", bn, blk.Hash), bn
}

func (c *Comparer) Run() {
	// go func() {
	for {
		key, keyOld, bn := c.GetLatestBlockKey()
		fmt.Println(key)
		if key == "" {
			continue
		}

		err := c.Compare(key, keyOld, bn)
		if err != nil {
			fmt.Println(err.Error())
			continue
		}
		time.Sleep(3 * time.Second)
	}
	// }()
}

func main() {
	s3cfg, err := awsconfig.LoadDefaultConfig(context.Background(), awsconfig.WithRegion("ap-northeast-1"))
	if err != nil {
		panic(fmt.Errorf("load s3 config failed, err:%s", err.Error()))
	}
	rpcCli, err := jsonrpc.NewClient(jsonrpc.WithURLEndpoint("coord", []string{URL2}))
	if err != nil {
		panic(fmt.Errorf("init jsonrpc cli failed, err:%s", err.Error()))
	}
	cmper := Comparer{
		S3Cli:        s3.NewFromConfig(s3cfg),
		ActBucket:    "tf-nodereal-prod-meganode-tracer-cache-ap",
		ExpBucket:    "tf-nodereal-prod-bsc-states-tracer-ap",
		ExpFolder:    "bsc_states_tracer_prestate",
		ActFolderEvm: "bsc-mainnet-prestate-test",
		ActFolderGat: "bsc-mainnet-prestate",
		RpcCli:       rpcCli,
	}
	cmper.Run()
}

func ComparePrestatesV2(pExp *old_models.Prestates, pAct *new_models.Prestates, bn uint64) {
	r1 := make(map[string]struct{})
	r2 := make(map[string]struct{})

	for k, v := range pExp.AccountPrestateMap {
		if _, ok := pAct.AccountPrestateMap[k]; !ok {
			r1[hexutil.Encode(k[:])] = struct{}{}
			Storage := "no Storage"
			if v.Storages != nil && len(v.Storages) != 0 {
				Storage = "has Storage"
			}
			fmt.Println(fmt.Sprintf("dbug bn[%d], addr[%s] [%s], pAct miss content[%v]", bn, hexutil.Encode(k[:]), Storage, v))
		}
	}

	for k, v := range pAct.AccountPrestateMap {
		if _, ok := pExp.AccountPrestateMap[k]; !ok {
			r2[hexutil.Encode(k[:])] = struct{}{}
			Storage := "no Storage"
			if v.Storages != nil && len(v.Storages) != 0 {
				Storage = "has Storage"
			}
			fmt.Println(fmt.Sprintf("dbug bn[%d], addr[%s] [%s], pExp miss content[%v]", bn, hexutil.Encode(k[:]), Storage, v))
		}
	}

	fmt.Println("=========addr 1==========")
	fmt.Println(r1)
	fmt.Println("=========addr 2==========")
	fmt.Println(r2)
	fmt.Println("======storage diff=======")

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
					fmt.Println(fmt.Sprintf("pAct addr [%s], storage key[%s] missed", addr, key))
				}
			}
		}
	}

	for k, v1 := range pAct.AccountPrestateMap {
		addr := hexutil.Encode(k[:])
		if v2, ok := pExp.AccountPrestateMap[k]; ok {
			for ks, _ := range v1.Storages {
				key := hexutil.Encode(ks[:])
				if _, ok := v2.Storages[ks]; !ok {
					fmt.Println(fmt.Sprintf("pExp addr [%s], storage key[%s] missed", addr, key))
				}
			}
		}
	}
	fmt.Println("=========================")
}

//func (c *Comparer) GetAddrInfo(addr string, bn string) {
//	req := jsonrpc.NewRequest()
//}
