package cache_test

import (
	"testing"

	"github.com/mysterion/avrp/internal/cache"
)

func TestCache(t *testing.T) {
	cache.Set("hello", "world")
	s := cache.Get("hello")
	if s != "world" {
		t.Fatal(s)
	}

	w := cache.Get("world")
	if w != "" {
		t.Fatal(w)

	}
}
