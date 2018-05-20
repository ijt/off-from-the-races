package main

import (
	"fmt"
	"time"
)

func main() {
	sq := func(x int) int { return x * x }
	inc := func(x int) int { return x + 1 }
	done := make(chan interface{})
	go func() {
		time.Sleep(time.Millisecond)
		close(done)
	}()
	for x := range mapPipeline(done, gen(2, 3), sq, inc) {
		fmt.Println(x)
	}
}

func mapPipeline(done <-chan interface{}, in <-chan int, fs ...func(int) int) <-chan int {
	if len(fs) == 0 {
		return in
	}
	in = mapPipeline(done, in, fs[:len(fs)-1]...)
	f := fs[len(fs)-1]
	out := make(chan int)
	go func() {
		defer close(out)
		for x := range in {
			select {
			case <-done:
				return
			default:
				out <- f(x)
			}
		}
	}()
	return out
}

func gen(nums ...int) <-chan int {
	out := make(chan int)
	go func() {
		for _, n := range nums {
			out <- n
		}
		close(out)
	}()
	return out
}
