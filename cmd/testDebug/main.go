package main

import (
	"bytes"
	"context"
	"fmt"
	"os"
	"sync"
	"time"

	"github.com/ethereum/go-ethereum/common/hexutil"
	jsoniter "github.com/json-iterator/go"
	jsonrpc "github.com/node-real/go-pkg/jsonrpc2"
)

const (
	QuickNodeUrl = "https://blissful-palpable-brook.bsc-testnet.quiknode.pro/db72f088982a0655c592f4ea07773317730c969b/"
	NodeRealUrl  = "https://bsc-testnet.nodereal.cc/v1/d9fa5c156a3c4102a55af9e97f6eb88f"
)

func main() {
	clientQ, _ := jsonrpc.NewClient(jsonrpc.WithURLEndpoint("quicknode", []string{QuickNodeUrl}))
	clientN, _ := jsonrpc.NewClient(jsonrpc.WithURLEndpoint("nodereal", []string{NodeRealUrl}))

	_, _ = clientQ, clientN

	tmr := time.NewTicker(3 * time.Second)
	for {
		select {
		case <-tmr.C:
			resp, _ := clientN.Call(context.Background(), jsonrpc.NewRequest(114514, "eth_blockNumber"))
			currBn_ := ""
			jsoniter.Unmarshal(resp.Result, &currBn_)
			bnNum, _ := hexutil.DecodeUint64(currBn_)
			//bnNum = bnNum - 10000
			currBn := hexutil.EncodeUint64(bnNum)
			fmt.Println(currBn)

			respQ := &jsonrpc.Response{}
			respN := &jsonrpc.Response{}

			req := jsonrpc.NewRequest(114514, "debug_traceBlockByNumber", currBn, struct {
				Tracer string
			}{
				Tracer: "callTracer",
			})

			wg := sync.WaitGroup{}
			wg.Add(2)
			go func() {
				defer wg.Done()
				respN, _ = clientN.Call(context.Background(), req)
			}()

			go func() {
				defer wg.Done()
				respQ, _ = clientQ.Call(context.Background(), req)
			}()
			wg.Wait()

			content := struct {
				Request   *jsonrpc.Request
				NodeReal  *jsonrpc.Response
				QuickNode *jsonrpc.Response
			}{
				Request:   req,
				NodeReal:  respN,
				QuickNode: respQ,
			}

			raw, _ := jsoniter.Marshal(content)

			if bytes.Compare(respQ.Result, respN.Result) != 0 {
				fmt.Printf("not equal [%s] \n", currBn)
			}
			fn := fmt.Sprintf("blk-%s.json", currBn)
			fp, _ := os.Create(fn)
			fp.Write(raw)
			fp.Close()
		}
	}
}
