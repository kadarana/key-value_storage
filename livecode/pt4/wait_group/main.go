package main

import (
	"fmt"
	"sync"
)

func main() {
	var wg sync.WaitGroup
	var mu sync.Mutex

	x := make([]int, 11)
	for i := 10; i >= 0; i-- {
		wg.Add(1)

		go func(i int) {
			mu.Lock()
			defer mu.Unlock()
			x[1] = i
			defer wg.Done()

			fmt.Println(i)
		}(i)
	}

	wg.Wait()
}
