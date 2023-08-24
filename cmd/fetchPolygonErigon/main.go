package main

import (
	"context"
	"fmt"
	"os"

	"github.com/ethereum/go-ethereum/common/hexutil"
	jsoniter "github.com/json-iterator/go"
	jsonrpc "github.com/node-real/go-pkg/jsonrpc2"
)

const ErigonURL = "https://rpc-mainnet.matic.quiknode.pro"

func main() {
	client, err := jsonrpc.NewClient(jsonrpc.WithURLEndpoint("Erigon", []string{ErigonURL}))
	if err != nil {
		fmt.Println(err.Error())
		return
	}
	latest := hexutil.Big{}
	for {
		resp, err := client.Call(
			context.Background(),
			jsonrpc.NewRequest(114514, "eth_blockNumber"),
		)
		if err != nil {
			fmt.Println(err.Error())
			continue
		}
		bn := hexutil.Big{}
		err = jsoniter.Unmarshal(resp.Result, &bn)
		fmt.Println(bn)
		bn.ToInt().SetUint64(bn.ToInt().Uint64() - 10)
		if bn.String() != latest.String() {
			latest = bn
		} else {
			continue
		}
		if err != nil {
			fmt.Println(err.Error())
			continue
		}
		req := jsonrpc.NewRequest(114514, "trace_replayBlockTransactions", bn, []string{"trace", "stateDiff"})
		resp, err = client.Call(
			context.Background(),
			req,
		)
		if err != nil {
			fmt.Println(err.Error())
			continue
		}
		f, err := os.Create(fmt.Sprintf("replay-%s.json", bn.String()))
		content := struct {
			Req  *jsonrpc.Request
			Resp *jsonrpc.Response
		}{
			Req:  req,
			Resp: resp,
		}
		raw, _ := jsoniter.Marshal(content)
		f.Write(raw)
	}
}
