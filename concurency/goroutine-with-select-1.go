package main

import (
	"fmt"
	"time"
)

func rick_just_loop(c chan int, quit chan string) {
	x := 0
	for {
		fmt.Println("rick_just_loop waiting for quit ")
		time.Sleep(1 * time.Second)
		select {
		case c <- x:
			//fmt.Println("case c, we incremented x")
			x = x + 1
		case <-quit:
			fmt.Println("quit case received")
			return
			//default:
			//	continue
		}
	}
}

func main() {
	c := make(chan int)
	quitchannel := make(chan string) //just showing can pass whatever

	go func() {
		//Thread doing stuff, notifies some other work, that
		//it can stop
		for i := 0; i < 10; i++ {
			fmt.Println("in loop ", i)
			fmt.Println("print channel c int: ", <-c)
			time.Sleep(1 * time.Second)
		}
		fmt.Println("quitting so sending stop to channel!")
		quitchannel <- "stop"
	}()

	rick_just_loop(c, quitchannel)
}
