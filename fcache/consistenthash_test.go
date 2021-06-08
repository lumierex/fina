package fcache

import (
	"log"
	"strconv"
	"testing"
)

func TestMap_Get(t *testing.T) {

	hash := New(3, func(key []byte) uint32 {
		log.Println("string: ", string(key))
		i, _ := strconv.Atoi(string(key))
		log.Println("i: ", i)
		return uint32(i)
	})
	hash.Add("6", "4", "2")

	ret := hash.Get("11")
	log.Println("ret: ", ret)

	testCases := map[string]string{
		"1":  "2",
		"12": "2",
		"2":  "2",
	}
	for key, value := range testCases {
		if hash.Get(key) != value {
			t.Errorf("get for %s, should yield to  %s", key, value)
		}
	}

	// 找key的方式实际上是key大于等于hash值
	hash.Add("8")
	testCases["27"] = "8"
	for key, value := range testCases {
		if hash.Get(key) != value {
			t.Errorf("get for %s, should yield to  %s", key, value)
		}
	}

}
