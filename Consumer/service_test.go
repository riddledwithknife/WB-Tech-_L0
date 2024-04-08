package main

import (
	"strconv"
	"sync"
	"testing"
)

func TestOrdersCache_Set_Get(t *testing.T) {
	t.Run("Work as expected", func(t *testing.T) {
		cache := NewOrdersCache()

		testData := `{"name":"Dan","age":21,"email":"Dan@example.com"}`
		cache.Set("key1", testData)

		retrievedData, ok := cache.Get("key1")
		if !ok {
			t.Error("Expected data not found")
		}

		if retrievedData != testData {
			t.Error("Retrieved data doesn't match")
		}
	})
}

func TestConcurrentCacheAccess(t *testing.T) {
	t.Run("Work as expected", func(t *testing.T) {
		cache := NewOrdersCache()

		const numRoutines = 1000

		var wg sync.WaitGroup
		wg.Add(numRoutines)

		for i := 0; i < numRoutines; i++ {
			go func(id int) {
				defer wg.Done()

				jsonData := `{"id": ` + strconv.Itoa(id) + `}`

				cache.Set("key"+strconv.Itoa(id), jsonData)

				retrievedData, ok := cache.Get("key" + strconv.Itoa(id))
				if !ok {
					t.Errorf("Expected data not found for key %d", id)
					return
				}

				if retrievedData != jsonData {
					t.Errorf("Retrieved data doesn't match for key %d: expected %s, got %s", id, jsonData, retrievedData)
					return
				}
			}(i)
		}

		wg.Wait()
	})
}
