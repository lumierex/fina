package fcache

import (
	"strconv"
	"testing"
)

func TestMap_Get(t *testing.T) {

	hash := New(3, func(key []byte) uint32 {
		i, _ := strconv.Atoi(string(key))
		return uint32(i)
	})
	hash.Add("6", "4", "2")
}
