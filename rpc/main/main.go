package main

import (
	"encoding/json"
	"fmt"
	"frpc"
	"frpc/codec"
	"log"
	"net"
	"time"
)

func startServer(addr chan string) {
	l, err := net.Listen("tcp", ":0")
	if err != nil {
		log.Fatal("net work error: ", err)
	}
	log.Println("start rpc server on: ", l.Addr())
	addr <- l.Addr().String()
	frpc.Accept(l)
}
func main() {
	addr := make(chan string)
	go startServer(addr)

	conn, _ := net.Dial("tcp", <-addr)
	defer func() {
		_ = conn.Close()
	}()

	time.Sleep(time.Second)

	// send options
	_ = json.NewEncoder(conn).Encode(frpc.DefaultOption)
	cc := codec.NewGobCodec(conn)
	for i := 0; i < 5; i++ {
		header := &codec.Header{
			ServiceMethod: "Foo.sum",
			Seq:           uint64(i),
		}
		_ = cc.Write(header, fmt.Sprintf("frpc req %d", header.Seq))
		_ = cc.ReadHeader(header)
		var reply string
		_ = cc.ReadBody(&reply)
		log.Println("reply", reply)

	}

}
