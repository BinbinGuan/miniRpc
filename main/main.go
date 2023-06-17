package main

import (
	"encoding/json"
	"fmt"
	"log"
	"miniRpc"
	"miniRpc/codec"
	"net"
	"time"
)

func startServer(addr chan string) {
	// pick a free port
	l, err := net.Listen("tcp", ":0")
	if err != nil {
		log.Fatal("network error:", err)
	}
	log.Println("start rpc server on", l.Addr())
	addr <- l.Addr().String() //在golang中可以使用channel作为同步的工具。 通过channel可以实现两个goroutine之间的通信。 创建一个channel， make(chan TYPE {, NUM}) , TYPE指的是channel中传输的数据类型，第二个参数是可选的，指的是channel的容量大小。 向channel传入数据， CHAN <- DATA , CHAN 指的是目的channel即收集数据的一方， DATA 则是要传的数据
	miniRpc.Accept(l)
}

func main() {
	log.SetFlags(0)
	addr := make(chan string)
	go startServer(addr)

	// in fact, following code is like a simple miniRpc client
	conn, _ := net.Dial("tcp", <-addr)
	defer func() { _ = conn.Close() }()

	time.Sleep(time.Second)
	// send options
	_ = json.NewEncoder(conn).Encode(miniRpc.DefaultOption)
	cc := codec.NewGobCodec(conn)
	// send request & receive response
	for i := 0; i < 5; i++ {
		h := &codec.Header{
			ServiceMethod: "Foo.Sum",
			Seq:           uint64(i),
		}
		_ = cc.Write(h, fmt.Sprintf("miniRpc req %d", h.Seq))
		_ = cc.ReadHeader(h)
		var reply string
		_ = cc.ReadBody(&reply)
		log.Println("reply:", reply)
	}
}
