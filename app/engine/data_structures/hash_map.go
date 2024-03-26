package data_structures

import (
	"sync"
	"time"
)

type HashMap struct {
	store map[string]interface{}
	mx    sync.RWMutex
}

var instance *HashMap

func init() {
	instance = &HashMap{
		store: make(map[string]interface{}),
	}
}

func GetHashMap() *HashMap {
	return instance
}

func (hm *HashMap) Set(key string, val interface{}, expiry int64) error {
	hm.mx.Lock()
	defer hm.mx.Unlock()

	hm.store[key] = val

	if expiry == -1 {
		return nil
	}

	time.AfterFunc(time.Duration(expiry*int64(time.Millisecond)), func() {
		hm.mx.Lock()
		defer hm.mx.Unlock()

		delete(hm.store, key)
	})
	return nil
}

func (hm *HashMap) Get(key string) (interface{}, bool, error) {
	hm.mx.Lock()
	defer hm.mx.Unlock()

	val, ok := hm.store[key]
	if !ok {
		return nil, false, nil
	}

	return val, true, nil
}
