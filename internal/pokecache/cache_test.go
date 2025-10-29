package pokecache

import (
	"fmt"
	"testing"
	"time"
)

// TestAddGet verifies that data can be correctly added to and retrieved from the cache.
func TestAddGet(t *testing.T) {
	const interval = 5 * time.Second
	cases := []struct {
		key string
		val []byte
	}{
		{
			key: "https://example.com/api/1",
			val: []byte("testdata_one"),
		},
		{
			key: "https://example.com/api/2",
			val: []byte("testdata_two"),
		},
	}

	// The interval is irrelevant for this test, but NewCache requires it.
	cache := NewCache(interval)

	for i, c := range cases {
		t.Run(fmt.Sprintf("Test case %v", i), func(t *testing.T) {
			// Add the data
			cache.Add(c.key, c.val)

			// Try to retrieve the data
			val, ok := cache.Get(c.key)

			if !ok {
				t.Errorf("expected to find key %s, but did not", c.key)
				return
			}
			if string(val) != string(c.val) {
				t.Errorf("expected value %s, but got %s", string(c.val), string(val))
				return
			}
		})
	}
}

// TestReapLoop verifies that old entries are automatically removed by the background reaper.
func TestReapLoop(t *testing.T) {
	// Set the base time low so the test runs quickly (e.g., 5 milliseconds)
	const baseTime = 5 * time.Millisecond
	// Wait a bit longer than the reap interval to guarantee cleanup
	const waitTime = baseTime + 5*time.Millisecond

	cache := NewCache(baseTime)
	cache.Add("https://example.com/reaptest", []byte("reapdata"))

	// 1. Check that the key exists immediately after adding
	_, ok := cache.Get("https://example.com/reaptest")
	if !ok {
		t.Errorf("expected to find key before reap")
		return
	}

	// 2. Wait for the reap interval to pass (plus a little buffer)
	time.Sleep(waitTime)

	// 3. Check that the key is gone after the reaper runs
	_, ok = cache.Get("https://example.com/reaptest")
	if ok {
		t.Errorf("expected to not find key after reap")
		return
	}
}
