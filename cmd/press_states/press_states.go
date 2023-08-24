package main

import (
	"context"
	"fmt"
	"math/rand"
	"sync"
	"time"

	"github.com/ethereum/go-ethereum/common/hexutil"
	jsonrpc "github.com/node-real/go-pkg/jsonrpc2"
)

const (
	ep4x           = "http://10.179.208.27:8545"
	ep8x           = "http://10.179.223.42:8545"
	epErigonTracer = "https://eth-goerli.nodereal.cc/v1/d9fa5c156a3c4102a55af9e97f6eb88f"

	/*
		4x 10.179.195.166
		8x 10.179.224.122
	*/
	body = `
	{
    "jsonrpc": "2.0",
    "id": 1,
    "method": "eth_estimateGas",
    "params": [
        {
            "data": "0xa9059cbb00000000000000000000000010ef50d7c281f27767ee5e65c6ceb32bd4a200fe0000000000000000000000000000000000000000000000000000000000000001",
            "from": "0xd9912f744b3a888144d2f6a2b7fd90b5873d80c5",
            "to": "0x921b9ce5433c9389dd87814cbf0d8d7d91e1ef2c",
            "gasPrice": "0x0"
        },
        "0x1410f48"
    ]
}
	`

	body2 = `
{
    "jsonrpc": "2.0",
    "method": "trace_block",
    "params": [
        "0x86b85c"
    ],
    "id": 8778611
}
`

	data     = "0xa9059cbb00000000000000000000000010ef50d7c281f27767ee5e65c6ceb32bd4a200fe0000000000000000000000000000000000000000000000000000000000000001"
	from     = "0xd9912f744b3a888144d2f6a2b7fd90b5873d80c5"
	to       = "0x921b9ce5433c9389dd87814cbf0d8d7d91e1ef2c"
	gasPrice = "0x0"
	method   = "eth_estimateGas"
)

type param struct {
	Data     string
	From     string
	To       string
	GasPrice string
}

func main1() {
	startBn := 24281658
	endBn := startBn + 2000
	ep4s_url := []string{ep4x}
	client4x, _ := jsonrpc.NewClient(jsonrpc.WithURLEndpoint("states4x", ep4s_url))

	ep8s_url := []string{ep8x}
	client8x, _ := jsonrpc.NewClient(jsonrpc.WithURLEndpoint("states8x", ep8s_url))

	param_ := param{
		Data:     data,
		From:     from,
		To:       to,
		GasPrice: gasPrice,
	}

	start1 := time.Now()
	for bn := startBn; bn < endBn; bn++ {
		bn_ := hexutil.EncodeUint64(uint64(bn))
		_, err := client4x.Call(context.Background(), jsonrpc.NewRequest(0, method, param_, bn_))
		if err != nil {
			fmt.Println(err.Error())
			return
		}
		//fmt.Println(resp)
	}
	end1 := time.Now()
	fmt.Println(end1.Sub(start1))

	start2 := time.Now()
	for bn := startBn; bn < endBn; bn++ {
		bn_ := hexutil.EncodeUint64(uint64(bn))
		_, err := client8x.Call(context.Background(), jsonrpc.NewRequest(0, method, param_, bn_))
		if err != nil {
			fmt.Println(err.Error())
			return
		}
		//fmt.Println(resp)
	}
	end2 := time.Now()
	fmt.Println(end2.Sub(start2))
}

func main() {
	client, _ := jsonrpc.NewClient(jsonrpc.WithURLEndpoint("test", []string{epErigonTracer}))
	//req := jsonrpc.NewRequest(114514, "trace_block", "0x85f355")
	//req2 := jsonrpc.NewRequest(114514, "debug_traceBlockByNumber", "0x85f355", struct {
	//	Tracer string
	//}{
	//	Tracer: "callTracer",
	//})
	wg := sync.WaitGroup{}
	rand.Seed(time.Now().Unix())
	for i := 0; i < 1000000; i++ {
		wg.Add(2)
		go func() {
			bn := 0x87cc53 - rand.Uint64()%5000
			fmt.Println(hexutil.EncodeUint64(uint64(bn)))
			req1 := jsonrpc.NewRequest(114514,
				"trace_replayBlockTransactions",
				hexutil.EncodeUint64(uint64(bn)), []string{"trace"})
			_ = req1

			req2 := jsonrpc.NewRequest(114514,
				"trace_replayBlockTransactions",
				hexutil.EncodeUint64(uint64(bn)), []string{"trace", "stateDiff"})
			_ = req2

			req := jsonrpc.NewRequest(114514, "debug_traceBlockByNumber", hexutil.EncodeUint64(uint64(bn)), struct {
				Tracer string
			}{
				Tracer: "callTracer",
			})
			_ = req
			_, err := client.Call(context.Background(), req2)
			_ = err
			//fmt.Println(resp)
			wg.Done()
		}()

		go func() {
			//resp, err := client.Call(context.Background(), req2)
			//raw, _ := resp.MarshalJSON()
			//fmt.Printf("resp %s\n", raw)
			//fmt.Println(err)
			//_, _ = resp, err
			wg.Done()
		}()
		time.Sleep(50 * time.Millisecond)
	}
	wg.Wait()
}
