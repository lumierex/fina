package fcache

import (
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
