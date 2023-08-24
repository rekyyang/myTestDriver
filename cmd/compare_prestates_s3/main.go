package main

import (
	"context"
	"fmt"
	"time"

	awsconfig "github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/ethereum/go-ethereum/rlp"
	jsoniter "github.com/json-iterator/go"
	jsonrpc "github.com/node-real/go-pkg/jsonrpc2"
)

const URL1 = "https://bsc-mainnet.nodereal.io/v1/4e7dd5235d434c5a837f7e48e9af9b4d"
const URL2 = "http://10.73.228.48:8545"

var _ = s3.Client{}

type Comparer struct {
	S3Cli     *s3.Client
	Bucket    string
	ExpFolder string
	ActFolder string

	RpcCli jsonrpc.Client
}

type BlkWithHash struct {
	Hash        string
	BlockNumber string
}

func (c *Comparer) Compare(key string) error {
	keyAct := fmt.Sprintf("%s/%s", c.ActFolder, key)
	keyExp := fmt.Sprintf("%s/%s", c.ExpFolder, key)

	objAct, err := c.S3Cli.GetObject(context.Background(), &s3.GetObjectInput{
		Bucket: &c.Bucket,
		Key:    &keyAct,
	})

	if err != nil {
		return err
	}

	objExp, err := c.S3Cli.GetObject(context.Background(), &s3.GetObjectInput{
		Bucket: &c.Bucket,
		Key:    &keyExp,
	})

	if err != nil {
		return err
	}

	pAct := Prestates{}
	pExp := Prestates{}

	err = rlp.Decode(objAct.Body, &pAct)
	if err != nil {
		return err
	}
	err = rlp.Decode(objExp.Body, &pExp)
	if err != nil {
		return err
	}
	ComparePrestates(&pExp, &pAct)
	return nil
}

func (c *Comparer) GetLatestBlockKey() string {
	resp, err := c.RpcCli.Call(context.Background(), jsonrpc.NewRequest("23333", "eth_blockNumber"))
	if err != nil {
		fmt.Println(fmt.Sprintf("failed to get latest blk: %s", err.Error()))
		return ""
	}

	if resp.Error != nil {
		fmt.Println(fmt.Sprintf("failed to get latest blk: %s", resp.Error.Error()))
		return ""
	}

	bnStr := ""
	err = jsoniter.Unmarshal(resp.Result, &bnStr)
	fmt.Println(bnStr)

	bnStr = hexutil.EncodeUint64(hexutil.MustDecodeUint64(bnStr) - 10)

	blk := BlkWithHash{}

	resp, err = c.RpcCli.Call(context.Background(), jsonrpc.NewRequest("23334", "eth_getBlockByNumber", bnStr, false))
	if err != nil {
		fmt.Println(fmt.Sprintf("failed to get latest blk: %s", err.Error()))
		return ""
	}

	if resp.Error != nil {
		fmt.Println(fmt.Sprintf("failed to get latest blk: %s", resp.Error.Error()))
		return ""
	}

	if err != nil {
		fmt.Println(fmt.Sprintf("failed to unmarshal blk: %s", err.Error()))
		return ""
	}
	bn, _ := hexutil.DecodeUint64(blk.BlockNumber)
	return fmt.Sprintf("%010d_%s", bn, blk.Hash)
}

func (c *Comparer) Run() {
	// go func() {
	for {
		key := c.GetLatestBlockKey()
		fmt.Println(key)
		if key == "" {
			continue
		}

		err := c.Compare(key)
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
		S3Cli:     s3.NewFromConfig(s3cfg),
		Bucket:    "tf-nodereal-prod-meganode-tracer-cache-ap",
		ExpFolder: "bsc-mainnet-prestate",
		ActFolder: "bsc-mainnet-prestate-test",
		RpcCli:    rpcCli,
	}
	cmper.Run()
}
