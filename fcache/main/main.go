package main

import (
	"fcache"
	"fmt"
	"log"
	"net/http"
)

var db = map[string]string{
	"tom":  "tom1",
	"same": "same1",
	"luna": "luna",
}

func main() {
	log.Println("2 << 10 ", 2<<10)

	fcache.NewGroup("fina", 2<<10, fcache.GetterFunc(func(key string) ([]byte, error) {
		log.Println("slow db search: key: ", key)
		if value, ok := db[key]; ok {
			return []byte(value), nil
		} else {
			return nil, fmt.Errorf("key : %s is not exist  ", key)
		}
	}))

	addr := "localhost:9999"

	handlerPool := fcache.NewHTTPPool(addr)
	log.Printf("fcache is running at : %s \n", addr)
	log.Fatal(http.ListenAndServe(addr, handlerPool))

}
