package main

import (
	"fmt"
	"math/rand/v2"
	"sync"
	"time"
)

// TaskResult represents the data returned by a goroutine
type TaskResult struct {
	TaskID      int
	Duration    time.Duration
	CompletedAt time.Time
}

// simulateTask simulates a task with a random duration
func simulateTask(taskID int, ch chan<- TaskResult, wg *sync.WaitGroup) {
	defer wg.Done() // Signal completion of the goroutine

	// Create a new random source for this goroutine
	source := rand.NewPCG(uint64(time.Now().UnixNano()), uint64(taskID))
	r := rand.New(source)

	// Simulate task duration between 500ms and 3000ms
	duration := time.Duration(r.IntN(2500)+500) * time.Millisecond
	time.Sleep(duration)

	// Send result to the channel
	result := TaskResult{
		TaskID:      taskID,
		Duration:    duration,
		CompletedAt: time.Now(),
	}
	fmt.Printf("Task %d sending result at %v\n", taskID, time.Now().Format("15:04:05.000"))
	ch <- result
}

// processResult simulates heavy processing of a result and prints it
func processResult(result TaskResult, wg *sync.WaitGroup) {
	defer wg.Done()

	// Simulate heavy processing (e.g., 200-800ms of work)
	source := rand.NewPCG(uint64(time.Now().UnixNano()), uint64(result.TaskID))
	r := rand.New(source)
	processingTime := time.Duration(r.IntN(600)+200) * time.Millisecond
	time.Sleep(processingTime)

	// Print the result
	fmt.Printf("Task %d completed in %v at %v (processed in %v)\n",
		result.TaskID,
		result.Duration,
		result.CompletedAt.Format("15:04:05.000"),
		processingTime)
}

// processResults reads results from the channel and processes them concurrently
func processResults(ch <-chan TaskResult, processWG *sync.WaitGroup) {
	for result := range ch {
		// Simulate a slow receiver
		time.Sleep(500 * time.Millisecond) // Artificial delay
		// Increment processWG for each result being processed
		processWG.Add(1)
		// Process each result in a separate goroutine
		go processResult(result, processWG)
	}
	// After the channel is closed, wait for all processing to complete
	processWG.Wait()
}

func main() {
	fmt.Println("=== Concurrent Task Example with Real-Time Notifications ===")

	// Number of tasks
	numTasks := 5

	// Create a buffered channel to receive task results
	resultChan := make(chan TaskResult, numTasks)
	/* because our for result runs quickly, we could get by without a buffered channel, but if
	it were an expaensive long running processing in the loop, then

	Blocking: When a task goroutine sends to the unbuffered resultChan, it blocks until the
	result-processing goroutine reads the result. The artificial 500ms delay in the for ...
	range loop means that each task goroutine may wait up to 500ms before its ch <- result operation completes.

	Output: You’ll see the “Task X sending result” messages spaced out by at least 500ms due to
	the receiver’s delay, indicating that task goroutines are blocking. With a buffered channel,
	these messages would appear closer to the actual completion times of the tasks, as sends would
	not block (up to the buffer capacity).
	*/

	// Create WaitGroups: one for tasks, one for result processing
	var taskWG, processWG sync.WaitGroup

	// Launch a goroutine to collect results and process them concurrently
	//if annoymous function
	/*
	go func() {
		for result := range resultChan {
			// Simulate a slow receiver
			//time.Sleep(500 * time.Millisecond) // Artificial delay
			//with slow receiver, we'll be ok since we have a buffered channel so you'll
			//see all the Task %d sending calls complete quickly, BUT if we had an unbuffered
			//and this sleep above wass long, you'd see that go routine block waiting for this process to
			//complete first before it can send the next result.
			
			// Increment processWG for each result being processed
			processWG.Add(1)
			// Process each result in a separate goroutine if processing is slow
			go processResult(result, &processWG)
		}
		// After the channel is closed, wait for all processing to complete
		processWG.Wait()
	}()
	*/

	// Launch the result processing goroutine
	go processResults(resultChan, &processWG)
	
	// Launch task goroutines
	fmt.Println("Launching", numTasks, "tasks...")
	start := time.Now()

	for i := 0; i < numTasks; i++ {
		taskWG.Add(1)
		go simulateTask(i, resultChan, &taskWG)
	}

	// Wait for all task goroutines to complete
	taskWG.Wait()
	fmt.Println("All tasks have completed.")

	// Close the channel since no more results will be sent
	close(resultChan)

	// Wait for all result processing to finish
	processWG.Wait()
	fmt.Println("All result processing has completed.")

	// Print total execution time
	fmt.Printf("\nProcessed %d tasks in %v\n", numTasks, time.Since(start))
}
