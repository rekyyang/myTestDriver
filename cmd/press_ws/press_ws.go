package main

import (
	"fmt"
	"log"
	"math/rand"
	"net/http"
	"sync"
	"time"

	"github.com/gorilla/websocket"
)

const (
	url0_ = `ws://10.179.227.13:8545/ws/v1/f381061f86f04e2a9490b0986be10a98`
	url1_ = `wss://bsc-mainnet-test.bk.nodereal.cc/ws/v1/f381061f86f04e2a9490b0986be10a98`
	url2_ = `ws://eth-goerli-test.bk.nodereal.cc/ws/v1/f381061f86f04e2a9490b0986be10a98`
	url3_ = `ws://localhost:8889/ws/v1/d470274753a04d2793b3dc747e421a49`
	req_  = "{\"jsonrpc\":\"2.0\",\"id\":1,\"method\":\"eth_getBlockByNumber\", \"params\":[\"latest\", false]}"
	//req_  = "{\"jsonrpc\":\"2.0\",\"id\":1,\"method\":\"eth_chainId\"}"
	subs_ = `{"jsonrpc":"2.0", "id": 1, "method": "eth_subscribe", "params": ["newHeads"]}`

	connectionEstablishInterval = 10 // ms
)

func main() {
	for {
		process()
	}
}

func process() {
	var wg sync.WaitGroup
	for i := 0; i < 20000; i++ {
		if i%1000 == 0 {
			fmt.Println(i)
		}
		wg.Add(2)
		url_ := url1_
		//if i%3 == 1 {
		//	url_ = url1_
		//} else {
		//	url_ = url2_
		//}
		idx := i
		go func() {
			defer wg.Done()
			createConnection(url_, idx)
		}()
		time.Sleep(connectionEstablishInterval * time.Millisecond)

		//go func() {
		// defer wg.Done()
		// createConnection("ws://bsc-mainnet-test.bk.nodereal.cc/ws/v1/beb06ff1ace649f4808094a7537124e7")
		//}()
	}
	wg.Wait()
}

func createConnection(url string, idx int) {
	//创建一个拨号器，也可以用默认的 websocket.DefaultDialer
	dialer := websocket.Dialer{}
	//向服务器发送连接请求，websocket 统一使用 ws://，默认端口和http一样都是80
	hdr := make(http.Header)
	hdr.Add("X-Forwarded-Host", "bsc-mainnet")
	connect, _, err := dialer.Dial(url, hdr)
	//ws://127.0.0.1:8889/ws/v1/d470274753a04d2793b3dc747e421a49
	//connect, _, err := dialer.Dial("ws://bsc-mainnet-test.bk.nodereal.cc/ws/v1/beb06ff1ace649f4808094a7537124e7", nil)
	if nil != err {
		//log.Println(err)
		//fmt.Println(err)
		return
	}
	//离开作用域关闭连接，go 的常规操作
	defer func() {
		log.Println("exit")
		connect.Close()
	}()

	//定时向客户端发送数据

	if idx%3 == 0 {
		_ = connect.WriteMessage(websocket.TextMessage, []byte(subs_))
	}
	go tickWriter(connect)

	//启动数据读取循环，读取客户端发送来的数据
	for {
		//从 websocket 中读取数据
		//messageType 消息类型，websocket 标准
		//messageData 消息数据
		messageType, _, err := connect.ReadMessage()
		if nil != err {
			//log.Println(err)
			//fmt.Println(err)
			break
		}
		switch messageType {
		case websocket.TextMessage: //文本数据
			//fmt.Println(string(messageData))
		case websocket.BinaryMessage: //二进制数据
			//fmt.Println(messageData)
		case websocket.CloseMessage: //关闭
		case websocket.PingMessage: //Ping
		case websocket.PongMessage: //Pong
		default:

		}
	}
}

func tickWriter(connect *websocket.Conn) {
	delay := rand.Float64() * 1000
	time.Sleep(time.Duration(delay) * time.Millisecond)
	for {
		//向客户端发送类型为文本的数据
		err := connect.WriteMessage(websocket.TextMessage, []byte(req_))
		//err := connect.WriteMessage(websocket.TextMessage, []byte("{\"jsonrpc\":\"2.0\",\"id\":1,\"method\":\"eth_blockNumber\"}"))
		if nil != err {
			//log.Println(err)
			//fmt.Println(err)
			break
		}
		//休息一秒
		time.Sleep(10000 * time.Millisecond)
	}
}
