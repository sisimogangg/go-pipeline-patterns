package main

import (
	"fmt"
	"sync"
)

/*
Fan-out:
Multiple functions can read from the same channel until that channel is closed

Fan-in:
A function can read from multiple inputs and proceed until all are closed by multiplexing [merging]
the input channels onto a single channel that's closed when all the inputs are closed

*/

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

func sq(ch <-chan int) <-chan int {
	out := make(chan int)
	go func() {
		for n := range ch {
			out <- n * n
		}
		close(out)
	}()

	return out
}

func main() {
	in := gen(2, 3)

	// Fanning-out
	// Distributing the sq worj accross two goroutines that both read from in.
	c1 := sq(in)
	c2 := sq(in)

	//fmt.Println(<-c1)
	//fmt.Println(<-c2)
	for n := range merge(c1, c2) {
		fmt.Println(n)
	}
}

func merge(cs ...<-chan int) <-chan int {
	var wg sync.WaitGroup
	out := make(chan int)

	// Start an output goroutine for each input channel in cs.  output
	// copies values from c to out until c is closed, then calls wg.Done.
	output := func(c <-chan int) {
		for n := range c {
			out <- n
		}
		wg.Done()
	}
	wg.Add(len(cs))
	for _, c := range cs {
		go output(c)
	}

	// Start a goroutine to close out once all the output goroutines are
	// done.  This must start after the wg.Add call.
	go func() {
		wg.Wait()
		close(out)
	}()
	return out
}
