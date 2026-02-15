package main

import (
	"fmt"
	"math/rand"
	"time"
)

// DONE: generator
// DONE: fanout - fanin
// DONE: block sequence
// DOING: replica
// testing

// Replica pattern
// - feeder to load data from a website,
// 	 we fanout to load data with timeout,
//   # situation
//   # 		when some of goroutine is time out
//   # 		we don't want to discard loaded data, so we replicate a fetchers
//   #		to let them race to finish maybe before timeout,
//   #		so we **increase the chance** to get all feed data

func main() {
	a := gen()
	b := gen()
	all := fanIn(a, b)
	for s := range all {
		fmt.Println("recv: ", s)
	}
}

func gen() <-chan int {
	c := make(chan int)
	go func() {
		defer close(c)
		for i := range 10 {
			time.Sleep(time.Duration(rand.Intn(1e3) * int(time.Millisecond)))
			c <- i
		}
	}()
	return c
}

func fanIn(c1, c2 <-chan int) <-chan string {
	c := make(chan string)
	// go func() {
	// 	for c1V := range c1 {
	// 		c <- fmt.Sprintf("c1: %d", c1V)
	// 	}
	// }()
	// go func() {
	// 	for c2V := range c2 {
	// 		c <- fmt.Sprintf("c2: %d", c2V)
	// 	}
	// }()
	go func() {
		for {
			a, ok0 := <-c1
			b, ok1 := <-c2
			if !ok0 || !ok1 {
				close(c)
				return
			}
			c <- fmt.Sprintf("block sequence: %v/%v", a, b)
		}
	}()

	return c
}
