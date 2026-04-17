package pokecache

import (
	"testing"
	"time"
)

func TestCache(t *testing.T) {
	cache := NewCache(5 * time.Minute)

	cache.Add("key1", []byte("value1"))
	cache.Add("key2", []byte("value2"))

	val, ok := cache.Get("key1")
	if !ok || string(val) != "value1" {
		t.Errorf("Expected value1, got %v", val)
	}

	val, ok = cache.Get("key2")
	if !ok || string(val) != "value2" {
		t.Errorf("Expected value2, got %v", val)
	}

	val, ok = cache.Get("key3")
	if ok {
		t.Errorf("Expected not found, got %v", val)
	}
}
