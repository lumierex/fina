package lru

import (
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
