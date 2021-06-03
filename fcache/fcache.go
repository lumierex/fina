package fcache

import (
	"fmt"
	"log"
	"sync"
)

type Getter interface {
	Get(key string) ([]byte, error)
}

type GetterFunc func(key string) ([]byte, error)

func (f GetterFunc) Get(key string) ([]byte, error) {
	return f(key)
}

var (
	mu     sync.RWMutex
	groups = make(map[string]*Group)
)

// Group namespace for different cache
// getter cache not hit exec callback(getter)
type Group struct {
	name      string
	getter    Getter
	mainCache cache
}

func NewGroup(name string, cacheBytes int64, getter Getter) *Group {
	mu.Lock()
	defer mu.Unlock()

	g := &Group{
		name:   name,
		getter: getter,
		mainCache: cache{
			cacheBytes: cacheBytes,
		},
	}
	groups[name] = g
	return g
}

// GetGroup get cache group
func GetGroup(name string) *Group {
	mu.RLock()
	g := groups[name]
	mu.RUnlock()
	return g
}

func (group *Group) Get(key string) (ByteView, error) {

	if key == "" {
		return ByteView{}, fmt.Errorf("key is required")
	}

	if value, ok := group.mainCache.get(key); ok {
		log.Println("fcache hit: ", key)
		return value, nil
	}

	// if group cache not hit the key
	// then exec callback to load from locally
	// or load from remote server
	return group.load(key)

}

func (group *Group) load(key string) (ByteView, error) {
	return group.loadLocally(key)
}

func (group *Group) loadLocally(key string) (ByteView, error) {
	byteValue, err := group.getter.Get(key)
	if err != nil {
		return ByteView{}, err
	}
	value := ByteView{b: cloneBytes(byteValue)}
	group.populate(key, value)
	return value, nil

}

func (group *Group) populate(key string, v ByteView) {
	group.mainCache.add(key, v)
}
