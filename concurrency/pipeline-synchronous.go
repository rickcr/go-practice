package main

import "fmt"

func main() {
	nums := []int{2, 3, 4, 7, 8}

	dataChannel := sliceToChannel(nums)

	finalChannel := sq(dataChannel)

	//blocks until channel gets a value or is closed
	for num := range finalChannel {
		fmt.Println(num)
	}

}

func sq(inputChannel <-chan int) <-chan int {
	outChannel := make(chan int)
	go func() {
		//blocks until output channel reads a value (in main) or is closed
		for num := range inputChannel {
			fmt.Printf("received num %d to sq and sending sq to outChannel\n", num)
			outChannel <- num * num
		}
		close(outChannel)
	}()
	return outChannel
}

func sliceToChannel(nums []int) <-chan int {
	outChannel := make(chan int)
	go func() {
		//blocks until output channel reads a value (in sq func)  or is closed
		for _, num := range nums {
			fmt.Printf("sent our num %d to outChannel\n", num)
			outChannel <- num
		}
		close(outChannel)
	}()
	return outChannel
}
