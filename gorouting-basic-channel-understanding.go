package main

import (
	"fmt"
	"time"
)

func A_Process(c chan string) {
	fmt.Println("A process go routine kicked off")
	//mimic slow
	time.Sleep(4000 * time.Millisecond)
	fmt.Println("A process done")
	
	//if you commented this out, it would cause a deadlock since the main method is waiting 
	//for the channel to be closed
	//or for the channel to receive a message
	c <- "stuff"
}

func main() {
	fmt.Println("main()")
	myChannel := make(chan string)
	go A_Process(myChannel)

	//THIS WILL BLOCK until A Process returns a message
	msg := <-myChannel

	fmt.Println("Main method ended. Channel msg is: ", msg)

}
