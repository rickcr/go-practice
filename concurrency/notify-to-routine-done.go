package main

import (
	"fmt"
	"time"
)

// <-chan read only channel  (chan<- send only channel)
func MyProcess(mydonechan <-chan bool) {
	for {
		select {
		default:
			fmt.Println("Doing stuff in myProcess()")
		case <-mydonechan:
			fmt.Println("*** Received done signal()")
			return
		}
	}

}
func main() {
	fmt.Println("main()")
	mydonechan := make(chan bool)
	go MyProcess(mydonechan)
	time.Sleep(5 * time.Second)
	mydonechan <- true
	fmt.Println("main() DONE (sent just sent done signal)")
}
