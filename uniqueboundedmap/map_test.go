package uniqueboundedmap

import (
	"fmt"
	"sync"
	"testing"
)

func TestMap(T *testing.T) {
	lMap := NewLockedMap[string, int](10)

	// Simulating concurrent access
	wg := sync.WaitGroup{}
	for i := 0; i < 15; i++ {
		wg.Add(1)
		go func(i int) {
			defer wg.Done()
			key := fmt.Sprintf("key%d", i)
			lMap.Put(key, i)
		}(i)
	}

	wg.Wait()

	// Deleting some items
	for i := 0; i < 5; i++ {
		wg.Add(1)
		go func(i int) {
			defer wg.Done()
			key := fmt.Sprintf("key%d", i)
			lMap.Delete(key)
		}(i)
	}

	wg.Wait()

	// Output the remaining items in the map
	lMap.mutex.Lock()
	fmt.Println("Remaining items in the map:")
	for k, v := range lMap.m {
		fmt.Printf("%s: %d\n", k, v)
	}
	lMap.mutex.Unlock()
}
