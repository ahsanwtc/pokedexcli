package cache

import (
	"fmt"
	"testing"
	"time"
)

func TestGet(t *testing.T) {
	const interval = 5 * time.Second
	cases := []struct{
		key string
		value []byte
	}{
		{ key: "abc", value: []byte("value-abc") },
		{ key: "def", value: []byte("value-def") },
	}

	for i, c := range cases {
		t.Run(fmt.Sprintf("should get a cache hit %v", i), func(t *testing.T) {
			_cache := NewCache(interval)
			_cache.Set(c.key, c.value)

			testData, ok := _cache.Get(c.key)
			if !ok {
				t.Errorf("expected to find key: %s", c.key)
				return
			}

			if string(c.value) != string(testData) {
				t.Errorf("expected to find value: %s", string(c.value))
				return
			}
		})
	}
}

func TestReapLoop(t *testing.T) {
	const interval = 3 * time.Second
	wait := interval + 5 * time.Millisecond
	_cache := NewCache(interval)
	_cache.Set("abc", []byte("testdata"))

	if _, ok := _cache.Get("abc"); !ok {
		t.Error("expected to find the key")
		return
	}

	time.Sleep(wait)

	if _, ok := _cache.Get("abc"); ok {
		t.Error("expected NOT to find the key")
		return
	}
}