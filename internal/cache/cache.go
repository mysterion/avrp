package cache

import "sync"

var mu sync.Mutex

var mp map[string]string

func init() {
	mu.Lock()
	defer mu.Unlock()

	mp = map[string]string{}
}

func Get(k string) string {
	return mp[k]
}

func Set(k string, v string) {
	mu.Lock()
	defer mu.Unlock()

	mp[k] = v
}
