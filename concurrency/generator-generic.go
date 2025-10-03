package main

import (
	"fmt"
	"math/rand"
)

// Generator - generic function can return a channel of any type and a fucntion that returns that type
// the done channel is used to signal the end of the stream and can be of any type
func repeatFunc[T any, F any](mydonechannel <-chan F, myfunc func() T) <-chan T {
	mystreamchannel := make(chan T)
	go func() {
		defer close(mystreamchannel)
		for {
			select {
			case <-mystreamchannel:
				return
			case mystreamchannel <- myfunc():
			}
		}
	}()
	return mystreamchannel
}

func main() {
	mydonechannel := make(chan bool)
	defer close(mydonechannel) //close the done channel when the function exits
	//or if not annon randNum := rand.Intn(500_000)
	myRepeatChannel := repeatFunc(mydonechannel, func() int {
		return rand.Intn(500_000)
	})
	for i := range myRepeatChannel {
		fmt.Println(i)
	}
}
