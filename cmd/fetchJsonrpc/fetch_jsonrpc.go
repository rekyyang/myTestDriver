package main

import (
	"context"
	"fmt"
	"os"
	"strconv"
	"time"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/ethereum/go-ethereum/core/types"
	jsoniter "github.com/json-iterator/go"
	jsonrpc "github.com/node-real/go-pkg/jsonrpc2"
)

const (
	StartBlkNo1 = 0x10
	StartBlkNo2 = 0x10000
	StartBlkNo3 = 0x832087

	BlkRange = 100
)

var (
	clientAlchemy, _  = jsonrpc.NewClient(jsonrpc.WithURLEndpoint("alchemy_goerli", []string{"https://eth-goerli.g.alchemy.com/v2/docs-demo"}))
	clientNodeReal, _ = jsonrpc.NewClient(jsonrpc.WithURLEndpoint("nodereal_goerli", []string{"https://eth-goerli.nodereal.io/v1/d32dc1e5d7554d04832cbf8dbda2c0ff"}))
)

// RPCTransaction represents a transaction that will serialize to the RPC representation of a transaction
type RPCTransaction struct {
	BlockHash        *common.Hash      `json:"blockHash"`
	BlockNumber      *hexutil.Big      `json:"blockNumber"`
	From             common.Address    `json:"from"`
	Gas              hexutil.Uint64    `json:"gas"`
	GasPrice         *hexutil.Big      `json:"gasPrice"`
	GasFeeCap        *hexutil.Big      `json:"maxFeePerGas,omitempty"`
	GasTipCap        *hexutil.Big      `json:"maxPriorityFeePerGas,omitempty"`
	Hash             common.Hash       `json:"hash"`
	Input            hexutil.Bytes     `json:"input"`
	Nonce            hexutil.Uint64    `json:"nonce"`
	To               *common.Address   `json:"to"`
	TransactionIndex *hexutil.Uint64   `json:"transactionIndex"`
	Value            *hexutil.Big      `json:"value"`
	Type             hexutil.Uint64    `json:"type"`
	Accesses         *types.AccessList `json:"accessList,omitempty"`
	ChainID          *hexutil.Big      `json:"chainId,omitempty"`
	V                *hexutil.Big      `json:"v"`
	R                *hexutil.Big      `json:"r"`
	S                *hexutil.Big      `json:"s"`
}

type Content struct {
	Req *jsonrpc.Request
	Rsp *jsonrpc.Response
}

type Block struct {
	Header       types.Header
	Transactions []*common.Hash
}

type TraceFilterMode string

type TraceFilterRequest struct {
	FromBlock   *hexutil.Uint64   `json:"fromBlock"`
	ToBlock     *hexutil.Uint64   `json:"toBlock"`
	FromAddress []*common.Address `json:"fromAddress"`
	ToAddress   []*common.Address `json:"toAddress"`
	Mode        TraceFilterMode   `json:"mode"`
	After       *uint64           `json:"after"`
	Count       *uint64           `json:"count"`
}

func main() {
	os.Mkdir("trace_block", os.ModePerm)
	fetchTraceBlock(StartBlkNo3, BlkRange)

	os.Mkdir("trace_replayBlockTransactions", os.ModePerm)
	fetchTraceReplayBlock(StartBlkNo3, 10)

	os.Mkdir("txs", os.ModePerm)
	fetchTransaction(StartBlkNo3, 10)

	os.Mkdir("trace_transaction", os.ModePerm)
	fetchTraceTransaction()

	os.Mkdir("trace_replayTransaction", os.ModePerm)
	fetchTraceReplayTransaction()

	os.Mkdir("trace_call", os.ModePerm)

	os.Mkdir("trace_get", os.ModePerm)
	fetchTraceGet(100, 10)

	os.Mkdir("trace_filter", os.ModePerm)
	fetchTraceFilter(StartBlkNo3, 10)

}

// trace_block
func fetchTraceBlock(bnStart, bnRange int) {
	method := "trace_block"
	var hdr = make(map[string]string)
	hdr["Origin"] = "https://docs.alchemy.com"
	//clientAlchemy, _ := jsonrpc.NewClient(jsonrpc.WithURLEndpoint("alchemy_goerli", []string{"https://eth-goerli.g.alchemy.com/v2/docs-demo"}))
	// block
	// from 0x10 to 0x1000
	for bn := bnStart; bn < bnStart+bnRange; bn++ {
		hx := "0x" + strconv.FormatInt(int64(bn), 16)
		req := jsonrpc.NewRequest(bn, method, hx)
		fn := generateFileName(req)
		resp, err := clientAlchemy.Call(context.Background(), req, jsonrpc.CallWithHeader(hdr))
		if err != nil {
			fmt.Println(err.Error())
			continue
		}
		writeReqResp(fn, req, resp)
		time.Sleep(200 * time.Millisecond)
	}
}

// trace_replayBlock
func fetchTraceReplayBlock(bnStart, bnRange int) {
	method := "trace_replayBlockTransactions"
	var hdr = make(map[string]string)
	hdr["Origin"] = "https://docs.alchemy.com"
	//clientAlchemy, _ := jsonrpc.NewClient(jsonrpc.WithURLEndpoint("alchemy_goerli", []string{"https://eth-goerli.g.alchemy.com/v2/docs-demo"}))
	// block
	// from 0x10 to 0x1000
	for bn := bnStart; bn < bnStart+bnRange; bn++ {
		hx := "0x" + strconv.FormatInt(int64(bn), 16)
		fmt.Println(hx)
		req := jsonrpc.NewRequest(bn, method, hx, []string{"trace"})
		fn := generateFileName(req)
		resp, err := clientAlchemy.Call(context.Background(), req, jsonrpc.CallWithHeader(hdr))
		if err != nil {
			fmt.Println(err.Error())
			continue
		}
		writeReqResp(fn, req, resp)
		time.Sleep(200 * time.Millisecond)
	}

	for bn := bnStart; bn < bnStart+bnRange; bn++ {
		hx := "0x" + strconv.FormatInt(int64(bn), 16)
		fmt.Println(hx)
		req := jsonrpc.NewRequest(bn, method, hx, []string{"stateDiff"})
		fn := generateFileName(req)
		resp, err := clientAlchemy.Call(context.Background(), req, jsonrpc.CallWithHeader(hdr))
		if err != nil {
			fmt.Println(err.Error())
			continue
		}
		writeReqResp(fn, req, resp)
		time.Sleep(200 * time.Millisecond)
	}
}

// get transactions
func fetchTransaction(bnStart, bnRange int) {
	noderealUrl := "https://eth-goerli.nodereal.cc/v1/f381061f86f04e2a9490b0986be10a98"
	//clientNodereal, _ := jsonrpc.NewClient(jsonrpc.WithURLEndpoint("nodereal_goerli", []string{noderealUrl}))

	txs := make([]*common.Hash, 0)
	methodGetBlockByNumber := "eth_getBlockByNumber"

	// select some blocks to get transactions
	for bn := bnStart; bn < bnStart+bnRange; bn++ {
		hx := "0x" + strconv.FormatInt(int64(bn), 16)
		fmt.Println(hx)
		req := jsonrpc.NewRequest(bn, methodGetBlockByNumber, hx, false)
		resp, err := clientNodeReal.Call(context.Background(), req)
		if err != nil {
			fmt.Println(err.Error())
			continue
		}
		//writeJson(fn, resp)
		blk := &Block{}
		err = jsoniter.Unmarshal(resp.Result, blk)
		if err != nil {
			fmt.Println(err.Error())
			continue
		}
		for _, tx := range blk.Transactions {
			txs = append(txs, tx)
		}
		time.Sleep(200 * time.Millisecond)
	}

	fn := fmt.Sprintf("txs/txs.json")
	writeJson(fn, txs)
}

func fetchTraceTransaction() {
	method := "trace_transaction"
	var hdr = make(map[string]string)
	hdr["Origin"] = "https://docs.alchemy.com"
	//clientAlchemy, _ := jsonrpc.NewClient(jsonrpc.WithURLEndpoint("alchemy_goerli", []string{"https://eth-goerli.g.alchemy.com/v2/docs-demo"}))

	txs := make([]common.Hash, 0)

	filePtr, err := os.Open("txs/txs.json")
	if err != nil {
		panic(err)
	}
	decoder := jsoniter.NewDecoder(filePtr)
	decoder.Decode(&txs)

	// from 0x10 to 0x1000
	for _, tx := range txs {
		req := jsonrpc.NewRequest("114514", method, tx)
		fn := generateFileName(req)
		resp, err := clientAlchemy.Call(context.Background(), req, jsonrpc.CallWithHeader(hdr))
		if err != nil {
			fmt.Println(err.Error())
			continue
		}
		writeReqResp(fn, req, resp)
		time.Sleep(200 * time.Millisecond)
	}

	for _, tx := range txs {
		req := jsonrpc.NewRequest("114514", method, tx)
		fn := generateFileName(req)
		resp, err := clientAlchemy.Call(context.Background(), req, jsonrpc.CallWithHeader(hdr))
		if err != nil {
			fmt.Println(err.Error())
			continue
		}
		writeReqResp(fn, req, resp)
		time.Sleep(200 * time.Millisecond)
	}
}

func fetchTraceReplayTransaction() {
	method := "trace_replayTransaction"
	var hdr = make(map[string]string)
	hdr["Origin"] = "https://docs.alchemy.com"
	//clientAlchemy, _ := jsonrpc.NewClient(jsonrpc.WithURLEndpoint("alchemy_goerli", []string{"https://eth-goerli.g.alchemy.com/v2/docs-demo"}))

	txs := make([]common.Hash, 0)

	filePtr, err := os.Open("txs/txs.json")
	if err != nil {
		panic(err)
	}
	decoder := jsoniter.NewDecoder(filePtr)
	decoder.Decode(&txs)

	// from 0x10 to 0x1000
	for _, tx := range txs {
		req := jsonrpc.NewRequest("114514", method, tx, []string{"trace"})
		fn := generateFileName(req)
		resp, err := clientAlchemy.Call(context.Background(), req, jsonrpc.CallWithHeader(hdr))
		if err != nil {
			fmt.Println(err.Error())
			continue
		}
		writeReqResp(fn, req, resp)
		time.Sleep(200 * time.Millisecond)
	}

	for _, tx := range txs {
		req := jsonrpc.NewRequest("114514", method, tx, []string{"stateDiff"})
		fn := generateFileName(req)
		resp, err := clientAlchemy.Call(context.Background(), req, jsonrpc.CallWithHeader(hdr))
		if err != nil {
			fmt.Println(err.Error())
			continue
		}
		writeReqResp(fn, req, resp)
		time.Sleep(200 * time.Millisecond)
	}
}

func fetchTraceGet(txCount, idxCount uint64) {
	method := "trace_get"
	var hdr = make(map[string]string)
	hdr["Origin"] = "https://docs.alchemy.com"
	//clientAlchemy, _ := jsonrpc.NewClient(jsonrpc.WithURLEndpoint("alchemy_goerli", []string{"https://eth-goerli.g.alchemy.com/v2/docs-demo"}))

	txs := make([]common.Hash, 0)

	filePtr, err := os.Open("txs/txs.json")
	if err != nil {
		panic(err)
	}
	decoder := jsoniter.NewDecoder(filePtr)
	decoder.Decode(&txs)

	// from 0x10 to 0x1000
	for _, tx := range txs {
		txCount--
		if txCount < 0 {
			break
		}
		for idx := uint64(0); idx < idxCount; idx++ {
			req := jsonrpc.NewRequest("114514", method, tx, hexutil.EncodeUint64(idx))
			fn := generateFileName(req)
			resp, err := clientAlchemy.Call(context.Background(), req, jsonrpc.CallWithHeader(hdr))
			if err != nil {
				fmt.Println(err.Error())
				continue
			}
			writeReqResp(fn, req, resp)
			time.Sleep(200 * time.Millisecond)
		}
	}
}

func fetchTraceFilter(bnStart, bnRange int) {
	method := "trace_filter"
	var hdr = make(map[string]string)
	hdr["Origin"] = "https://docs.alchemy.com"

	//blks := make([]*types.Block, 0)
	txs := make([]*RPCTransaction, 0)
	for bn := bnStart; bn <= bnStart+bnRange; bn++ {
		bn := hexutil.Uint64(bn).String()
		blk := types.Block{}
		req := jsonrpc.NewRequest(bn, "eth_getBlockByNumber", bn, true)
		resp, err := clientAlchemy.Call(context.Background(), req, jsonrpc.CallWithHeader(hdr))
		if err != nil {
			fmt.Println(err.Error())
			continue
		}
		jsoniter.Unmarshal(resp.Result, &blk)
		for _, tx := range blk.Transactions() {
			reqTx := jsonrpc.NewRequest(bn, "eth_getTransactionByHash", tx.Hash())
			respTx, err := clientAlchemy.Call(context.Background(), reqTx, jsonrpc.CallWithHeader(hdr))
			if err != nil {
				fmt.Println(err.Error())
				continue
			}
			rpcTx := RPCTransaction{}
			jsoniter.Unmarshal(respTx.Result, &rpcTx)
			txs = append(txs, &rpcTx)
		}
	}

	// full
	{
		fromBlk := hexutil.Uint64(bnStart)
		toBlk := hexutil.Uint64(bnStart + bnRange)
		filterReq := &TraceFilterRequest{
			FromBlock:   &fromBlk,
			ToBlock:     &toBlk,
			FromAddress: nil,
			ToAddress:   nil,
			After:       nil,
			Count:       nil,
		}
		req := jsonrpc.NewRequest(114514, method, filterReq)
		fn := generateFileName(req)
		resp, err := clientAlchemy.Call(context.Background(), req, jsonrpc.CallWithHeader(hdr))
		if err != nil {
			fmt.Println(err.Error())
		}
		writeReqResp(fn, req, resp)
		time.Sleep(200 * time.Millisecond)
	}

	// count
	{
		count := uint64(15)
		fromBlk := hexutil.Uint64(bnStart)
		toBlk := hexutil.Uint64(bnStart + bnRange)
		filterReq := &TraceFilterRequest{
			FromBlock:   &fromBlk,
			ToBlock:     &toBlk,
			FromAddress: nil,
			ToAddress:   nil,
			After:       nil,
			Count:       &count,
		}
		req := jsonrpc.NewRequest(114514, method, filterReq)
		fn := generateFileName(req)
		resp, err := clientAlchemy.Call(context.Background(), req, jsonrpc.CallWithHeader(hdr))
		if err != nil {
			fmt.Println(err.Error())
		}
		writeReqResp(fn, req, resp)
		time.Sleep(200 * time.Millisecond)
	}

	// after
	{
		after := uint64(15)
		fromBlk := hexutil.Uint64(bnStart)
		toBlk := hexutil.Uint64(bnStart + bnRange)
		filterReq := &TraceFilterRequest{
			FromBlock:   &fromBlk,
			ToBlock:     &toBlk,
			FromAddress: nil,
			ToAddress:   nil,
			After:       &after,
			Count:       nil,
		}
		req := jsonrpc.NewRequest(114514, method, filterReq)
		fn := generateFileName(req)
		resp, err := clientAlchemy.Call(context.Background(), req, jsonrpc.CallWithHeader(hdr))
		if err != nil {
			fmt.Println(err.Error())
		}
		writeReqResp(fn, req, resp)
		time.Sleep(200 * time.Millisecond)
	}

	// from to union
	{
		fromBlk := hexutil.Uint64(bnStart)
		toBlk := hexutil.Uint64(bnStart + bnRange)
		filterReq := &TraceFilterRequest{
			FromBlock:   &fromBlk,
			ToBlock:     &toBlk,
			FromAddress: nil,
			ToAddress:   nil,
			After:       nil,
			Count:       nil,
		}
		filterReq.FromAddress = append(filterReq.FromAddress, &txs[0].From)
		filterReq.FromAddress = append(filterReq.FromAddress, &txs[1].From)
		filterReq.FromAddress = append(filterReq.FromAddress, &txs[4].From)

		filterReq.ToAddress = append(filterReq.ToAddress, txs[0].To)
		filterReq.ToAddress = append(filterReq.ToAddress, txs[2].To)
		filterReq.ToAddress = append(filterReq.ToAddress, txs[8].To)
		req := jsonrpc.NewRequest(114514, method, filterReq)
		fn := generateFileName(req)
		resp, err := clientAlchemy.Call(context.Background(), req, jsonrpc.CallWithHeader(hdr))
		if err != nil {
			fmt.Println(err.Error())
		}
		writeReqResp(fn, req, resp)
		time.Sleep(200 * time.Millisecond)
	}

	// from to intersection
	{
		fromBlk := hexutil.Uint64(bnStart)
		toBlk := hexutil.Uint64(bnStart + bnRange)
		filterReq := &TraceFilterRequest{
			FromBlock:   &fromBlk,
			ToBlock:     &toBlk,
			FromAddress: nil,
			ToAddress:   nil,
			After:       nil,
			Count:       nil,
			Mode:        "intersection",
		}
		filterReq.FromAddress = append(filterReq.FromAddress, &txs[0].From)
		filterReq.FromAddress = append(filterReq.FromAddress, &txs[1].From)
		filterReq.FromAddress = append(filterReq.FromAddress, &txs[4].From)
		filterReq.FromAddress = append(filterReq.FromAddress, &txs[7].From)
		filterReq.FromAddress = append(filterReq.FromAddress, &txs[12].From)

		filterReq.ToAddress = append(filterReq.ToAddress, txs[0].To)
		filterReq.ToAddress = append(filterReq.ToAddress, txs[2].To)
		filterReq.ToAddress = append(filterReq.ToAddress, txs[8].To)
		filterReq.ToAddress = append(filterReq.ToAddress, txs[11].To)
		filterReq.ToAddress = append(filterReq.ToAddress, txs[15].To)
		req := jsonrpc.NewRequest(114514, method, filterReq)
		fn := generateFileName(req)
		resp, err := clientAlchemy.Call(context.Background(), req, jsonrpc.CallWithHeader(hdr))
		if err != nil {
			fmt.Println(err.Error())
		}
		writeReqResp(fn, req, resp)
		time.Sleep(200 * time.Millisecond)
	}
}

func generateFileName(req *jsonrpc.Request) string {
	return fmt.Sprintf("%s/%s_%s.json", req.Method, req.Method, string(req.Params))
}

func loadRespFile(req *jsonrpc.Request) *jsonrpc.Response {
	resp := &jsonrpc.Response{}
	filePtr, err := os.Open(generateFileName(req))
	if err != nil {
		return nil
	}
	defer filePtr.Close()
	decoder := jsoniter.NewDecoder(filePtr)
	err = decoder.Decode(resp)
	if err != nil {
		fmt.Println("解码错误", err.Error())
	} else {
		fmt.Println("解码成功")
	}
	return resp
}

func writeReqResp(fileName string, req *jsonrpc.Request, resp *jsonrpc.Response) {
	writeJson(fileName, Content{
		Req: req,
		Rsp: resp,
	})
}

func writeJson(fileName string, content interface{}) {
	// 创建文件
	filePtr, err := os.Create(fileName)
	if err != nil {
		fmt.Println("文件创建失败", err.Error())
		return
	}
	defer filePtr.Close()
	// 创建Json编码器
	encoder := jsoniter.NewEncoder(filePtr)
	err = encoder.Encode(content)
	if err != nil {
		fmt.Println("编码错误", err.Error())
	} else {
		fmt.Println(fileName)
	}
}
