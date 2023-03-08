package main

import (
	"context"
	"fmt"
	"time"

	"github.com/ethereum/go-ethereum/common/hexutil"
	jsonrpc "github.com/node-real/go-pkg/jsonrpc2"
)

const (
	ep4x = "http://10.179.195.166:8545"
	ep8x = "10.179.224.122:8545"

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

func main() {
	startBn := 24281658
	endBn := startBn + 2000
	ep4s_url := []string{ep4x}
	client4x, _ := jsonrpc.NewClient(jsonrpc.WithURLEndpoint("states4x", ep4s_url))

	//ep8s_url := []string{ep8x}
	//client8x, _ := jsonrpc.NewClient(jsonrpc.WithURLEndpoint("states8x", ep8s_url))

	param_ := param{
		Data:     data,
		From:     from,
		To:       to,
		GasPrice: gasPrice,
	}

	start1 := time.Now()
	for bn := startBn; bn < endBn; bn++ {
		bn_ := hexutil.EncodeUint64(uint64(bn))
		resp, err := client4x.Call(context.Background(), jsonrpc.NewRequest(0, method, param_, bn_))
		if err != nil {
			fmt.Println(err.Error())
			return
		}
		fmt.Println(resp)
	}
	end1 := time.Now()
	fmt.Println(end1.Sub(start1))
}
