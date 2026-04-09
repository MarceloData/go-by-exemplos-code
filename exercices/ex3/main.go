package main

import (
	"fmt"
	"sync"
)

type List struct {
	nums []int
}

func Double(num int) int {
	return num * 2
}

func main() {
	var wg sync.WaitGroup
	list := List{nums: []int{1, 2, 3, 4, 5}}
	ch := make(chan int, len(list.nums))
	newList := List{nums: []int{}}
	for _, num := range list.nums {
		wg.Add(1)
		go func(n int) {
			defer wg.Done()
			ch <- Double(n)
		}(num)
	}
	wg.Wait()
	close(ch)
	for result := range ch {
		newList.nums = append(newList.nums, result)
	}
	fmt.Println(newList.nums)
}
