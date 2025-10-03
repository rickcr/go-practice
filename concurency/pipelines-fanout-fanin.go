package main

import (
	"fmt"
	"math/rand"
	"runtime"
	"sync"
	"time"
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

func fanIn[T any](done <-chan bool, streamsOfPrimes []<-chan T) <-chan T {
	//collect all the prime channels into a single channel
	fannedInChannel := make(chan T)
	var wg sync.WaitGroup

	transfer := func(c <-chan T) {
		defer wg.Done()
		for i := range c {
			select {
			case <-done:
			case fannedInChannel <- i:
			}
		}
	}

	for _, c := range streamsOfPrimes {
		wg.Add(1)
		go transfer(c)
	}

	go func() {
		wg.Wait()
		close(fannedInChannel)
	}()

	return fannedInChannel
}

func main() {
	start := time.Now()
	
	mydonechannel := make(chan bool)
	defer close(mydonechannel) //close the done channel when the function exits
	var randNumFetcher = func() int { return rand.Intn(500_000_000) }
	randNumStream := repeatFunc(mydonechannel, randNumFetcher)

	CupCount := runtime.NumCPU()
	println("cup count is ", CupCount)

	//fan out
	primeFinderChannels := make([]<-chan int, CupCount)
	for i := 0; i < CupCount; i++ {
		primeFinderChannels[i] = primeFinder(mydonechannel, randNumStream)
	}

	//fan in
	fannedInChannel := fanIn(mydonechannel, primeFinderChannels)

	for i := range take(mydonechannel, fannedInChannel, 10) {
		fmt.Println(i)
	}

	fmt.Println(time.Since(start))
}
