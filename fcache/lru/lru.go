package lru

import (
	"container/list"
)

// Cache cache definition
// maxBytes max capacity of Cache
// nBytes current used storage
// ll double linked list
type Cache struct {
	maxBytes  int64
	nBytes    int64
	ll        *list.List
	cache     map[string]*list.Element
	onEvicted func(key string, value Value)
}

// entry data type of double linked list
type entry struct {
	key   string
	value Value
}

type Value interface {
	Len() int
}

// New constructor
func New(maxBytes int64, onEvicted func(string, Value)) *Cache {
	return &Cache{
		maxBytes:  maxBytes,
		ll:        list.New(),
		cache:     make(map[string]*list.Element),
		onEvicted: onEvicted,
	}
}

// Get get value from cache
func (c *Cache) Get(key string) (Value, bool) {
	if ele, ok := c.cache[key]; ok {
		kv := ele.Value.(*entry)
		// assume front is the end of list
		c.ll.MoveToFront(ele)
		return kv.value, true
	}
	return nil, false
}

// removeOldest
func (c *Cache) removeOldest() {
	ele := c.ll.Back()
	if ele != nil {
		c.ll.Remove(ele)
		//	delete ele from cache
		kv := ele.Value.(*entry)
		delete(c.cache, kv.key)
		// TODO ?
		c.nBytes = c.nBytes - int64(kv.value.Len()) - int64(len(kv.key))
		if c.onEvicted != nil {
			c.onEvicted(kv.key, kv.value)
		}
	}
}

// Add add values to cache
func (c *Cache) Add(key string, value Value) {
	// 1. if cache already has key value cover it
	// 1.1 move ele to front
	// 1.2 compute nBytes = current += len(oldValue) - len(newValue)
	// 1.3 replace value *pointer.value=xxx
	// 2. if cache don't has the key value
	// add to cache and move ele to front
	// 3. if maxBytes != 0 && maxBytes less than nBytes remove oldest value

	if ele, ok := c.cache[key]; ok {
		c.ll.MoveToFront(ele)
		kv := ele.Value.(*entry)
		c.nBytes += int64(kv.value.Len()) - int64(value.Len())
		kv.value = value
	} else {
		ele := c.ll.PushFront(&entry{
			key:   key,
			value: value,
		})
		c.cache[key] = ele
		// fix nBytes num
		c.nBytes += int64(len(key)) + int64(value.Len())
	}

	for c.maxBytes != 0 && c.maxBytes < c.nBytes {
		c.removeOldest()
	}
}

func (c *Cache) Len() int {
	return c.ll.Len()
}
