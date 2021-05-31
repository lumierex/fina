package main

import (
	"fmt"
	"frpc"
	"log"
	"net"
	"sync"
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

	// day1
	//conn, _ := net.Dial("tcp", <-addr)
	//defer func() {
	//	_ = conn.Close()
	//}()

	//time.Sleep(time.Second)

	// day 1
	// send options
	//_ = json.NewEncoder(conn).Encode(frpc.DefaultOption)
	//cc := codec.NewGobCodec(conn)
	//for i := 0; i < 5; i++ {
	//	header := &codec.Header{
	//		ServiceMethod: "Foo.sum",
	//		Seq:           uint64(i),
	//	}
	//	_ = cc.Write(header, fmt.Sprintf("frpc req %d", header.Seq))
	//	_ = cc.ReadHeader(header)
	//	var reply string
	//	_ = cc.ReadBody(&reply)
	//	log.Println("reply", reply)
	//
	//}

	// day2
	log.SetFlags(0)
	client, _ := frpc.Dial("tcp", <-addr)
	var wg sync.WaitGroup
	for i := 0; i < 5; i++ {
		wg.Add(1)
		go func(i int) {
			defer wg.Done()
			args := fmt.Sprintf("frpc req %d", i)
			var reply string
			if err := client.Call("Foo.Sun", args, &reply); err != nil {
				log.Fatal("call Foo.Sum error : ", err)
			}
			log.Println("reply: ", reply)
		}(i)
	}
	wg.Wait()
}
