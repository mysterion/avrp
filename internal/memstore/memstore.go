package memstore

import "sync"

var mu sync.Mutex

var mp map[string]interface{}

func init() {
	mu.Lock()
	defer mu.Unlock()

	mp = map[string]interface{}{}
}

func MpGet(k string) interface{} {
	return mp[k]
}

func MpSet(k string, v interface{}) {
	mu.Lock()
	defer mu.Unlock()

	mp[k] = v
}
