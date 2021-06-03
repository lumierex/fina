package fcache

import (
	"fmt"
	"log"
	"reflect"
	"testing"
)

func TestGetterFunc_Get(t *testing.T) {
	getter := GetterFunc(func(key string) ([]byte, error) {
		return []byte(key), nil
	})

	expect := []byte("key")
	if v, _ := getter.Get("key"); !reflect.DeepEqual(v, expect) {
		t.Fatalf("cb error")
	}

}

func TestGroup_Get(t *testing.T) {
	var db = map[string]string{
		"tom": "10000",
		"sam": "111111",
		"hi":  "hello world",
	}
	loadCounts := make(map[string]int, len(db))

	fina := NewGroup("fina", 10<<2, GetterFunc(func(key string) ([]byte, error) {
		// search
		log.Println("slow DB: search: key ", key)
		if value, ok := db[key]; ok {
			log.Println("[db] cache key : ", key)
			// make sum of load from db
			if _, ok := loadCounts[key]; !ok {
				loadCounts[key] = 0
			}
			loadCounts[key] += 1
			return []byte(value), nil
		}
		return nil, fmt.Errorf("%s not exist", key)
	}))

	for k, v := range db {
		if v1, err := fina.Get(k); err != nil || v1.String() != v {
			t.Fatalf("failed to hit tom, sam, hi")
		}

		// only be cached once
		// if loadCounts more than 1 than program has error logic
		if _, err := fina.Get(k); err != nil || loadCounts[k] > 1 {
			t.Fatalf("cache miss key %s ", k)
		}
	}

	if value, err := fina.Get("unknown"); err == nil {
		t.Fatalf("the value of %s showld be empty but return. ", value)
	}

}
