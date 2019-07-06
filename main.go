package main

import "fmt"

// Squaring numbers example
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
	out := sq(in)

	for n := range out {
		fmt.Println(n)
	}

	/*
			Since sq has the same type for its inbound and outbound channels, we can compose it any number of times

		for n := range sq(sq(gen(2, 3))) {
			fmt.Println(n)
		}
	*/

	//	time.Sleep(100 * time.Millisecond)

}
