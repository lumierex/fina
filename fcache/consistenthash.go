package fcache

import (
	"hash/crc32"
	"sort"
	"strconv"
)

// Hash 将key映射到2^32 - 1 中
type Hash func([]byte) uint32

type Map struct {
	hash     Hash
	replicas int
	keys     []int
	hashMap  map[int]string
}

func New(replicas int, fn Hash) *Map {
	m := &Map{
		hash:     fn,
		replicas: replicas,
		hashMap:  make(map[int]string),
	}
	if m.hash == nil {
		m.hash = crc32.ChecksumIEEE
	}
	return m
}

func (m *Map) Add(keys ...string) {
	for _, key := range keys {
		for i := 0; i < m.replicas; i++ {
			// 虚拟节点都指向key
			hashKey := int(m.hash([]byte(strconv.Itoa(i) + key)))
			m.keys = append(m.keys, hashKey)
			m.hashMap[hashKey] = key
		}
	}
	sort.Ints(m.keys)
}

func (m *Map) Get(key string) string {
	if key == "" {
		return ""
	}
	hash := int(m.hash([]byte(key)))
	idx := sort.Search(len(m.keys), func(i int) bool {
		return m.keys[i] > hash
	})

	// 通过hashKey去找真是key
	return m.hashMap[m.keys[idx%len(m.keys)]]
}
