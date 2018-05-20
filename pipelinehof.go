package main

import (
	"fmt"
)

func main() {
	sq := func(x int) int { return x * x }
	inc := func(x int) int { return x + 1 }
	for x := range mapPipeline(gen(2, 3), sq, inc) {
		fmt.Println(x)
	}
}

func mapPipeline(in <-chan int, fs ...func(int) int) <-chan int {
	if len(fs) == 0 {
		return in
	}
	in = mapPipeline(in, fs[:len(fs)-1]...)
	f := fs[len(fs)-1]
	out := make(chan int)
	go func() {
		for x := range in {
			out <- f(x)
		}
		close(out)
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
