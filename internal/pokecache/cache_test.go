package pokecache

import (
	"testing"
	"time"
)

func TestAddGet(t *testing.T) {

	cache := NewCache(5 * time.Second)

	cache.Add("testkey", []byte("testvalue"))

	val, ok := cache.Get("testkey")
	if !ok {
		t.Errorf("Expected to find the key")
	}

	if string(val) != "testvalue" {
		t.Errorf("Expected value doesnt match")
	}
}
