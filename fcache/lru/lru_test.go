package lru

import (
	"reflect"
	"testing"
)

type String string

func (data String) Len() int {
	return len(data)
}

func TestNew(t *testing.T) {

	// maxBytes equal to 0 we dont limit
	lru := New(int64(0), nil)
	lru.Add("key1", String("value1"))
	if v, ok := lru.Get("key1"); !ok && string(v.(String)) != "value1" {
		t.Fatalf("hit key %s failed", "key1")
	}
	if _, ok := lru.Get("key2"); ok {
		t.Fatalf("cache not exist")
	}
}

func TestRemoveOldest(t *testing.T) {
	k1, k2, k3 := "key1", "key2", "key3"

	v1, v2, v3 := "v1", "v2", "v3"

	capacity := len(k1 + k2 + v1 + v2)
	lru := New(int64(capacity), nil)

	lru.Add(k1, String(v1))
	lru.Add(k2, String(v2))
	lru.Add(k3, String(v3))

	if _, ok := lru.Get(k1); ok || lru.Len() != 2 {
		t.Fatalf("remove oldest value failed")
	}
}

func TestOnEvicted(t *testing.T) {
	k1, k2, k3, k4 := "k1", "k2", "k3", "k4"
	v1, v2, v3, v4 := "v1", "v2", "v3", "v4"
	keys := make([]string, 0)
	// if maxBytes equals to 0
	// lru won't remove oldest
	lru := New(int64(8), func(key string, value Value) {
		t.Log("key", key, " value: ", value)
		keys = append(keys, key)
	})
	t.Log(keys)
	lru.Add(k2, String(v1))
	lru.Add(k1, String(v2))
	lru.Add(k3, String(v3))
	lru.Add(k4, String(v4))
	t.Log("keys", keys)
	expect := []string{k2, k1}

	if !reflect.DeepEqual(keys, expect) {
		t.Fatalf("test on Evicted failed")
	}
}
