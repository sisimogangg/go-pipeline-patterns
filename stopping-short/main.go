package main

import (
	"fmt"
	"sync"
)

func gen(nums ...int) <-chan int {
	out := make(chan int, len(nums))
	for _, n := range nums {
		out <- n
	}
	close(out)
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
		//return
		/* putting this here will cause resource leak because now merge(c1, c2)
		   is stuck trying to send through the second number */

		/*
			We need to arrange for the upstream stages of our pipeline to exit even when
			the downstream stages fail to receive all the inbound values.
				* One way to do this is to change the outbound channels to have a buffer
		*/
	}
}

func merge(cs ...<-chan int) <-chan int {
	var wg sync.WaitGroup
	out := make(chan int, 1)

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
