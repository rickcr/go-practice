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
			case <-mydonechannel:
				return
			case mystreamchannel <- myfunc():
			}
			//blocks until mystreamchannel gets a value or is closed, so it won't resume until then
		}
	}()
	return mystreamchannel
}

func take[T any, K any](mydonechannel <-chan K, streamOfRepeatFuncVals <-chan T, n int) <-chan T {
	takenchannel := make(chan T)
	go func() {
		defer close(takenchannel)
		for i := 0; i < n; i++ {
			select {
			case <-mydonechannel:
				return
				//taking the one value from stream (<-stream) and then writing it to the takechannel
				//if you only did <-stream you'd be writing the same value over and over again'
			case takenchannel <- <-streamOfRepeatFuncVals:
				//blocks until takenchannel gets a value or is closed, so it won't resume until taken channel is closed
			}
		}
	}()

	return takenchannel
}

func primeFinder(done <-chan bool, randNumStream <-chan int) <-chan int {
	isPrime := func(num int) bool {
		for i := 2; i < num; i++ {
			if num%i == 0 {
				return false
			}
		}
		return true
	}
	primes := make(chan int)
	go func() {
		defer close(primes)
		for num := range randNumStream {
			select {
			case <-done:
				return
			case randomInt := <-randNumStream:
				if isPrime(randomInt) {
					primes <- randomInt
				}
			}
			if isPrime(num) {
				primes <- num
			}
		}
	}()
	return primes

}

func main() {
	mydonechannel := make(chan bool)
	defer close(mydonechannel) //close the done channel when the function exits
	var randNumFetcher = func() int { return rand.Intn(500_000_000) }
	randNumStream := repeatFunc(mydonechannel, randNumFetcher)
	primeStream := primeFinder(mydonechannel, randNumStream)

	myTakenChannel := take(mydonechannel, primeStream, 10)

	for i := range myTakenChannel {
		fmt.Println(i)
	}

}
