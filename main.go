package main

import (
	"fmt"
	"math/rand"
	"spinner/pkg/spinner" // replace with the correct import path
	"sync"
	"time"
)

func main() {

	fmt.Println("This is an example")

	for i := 0; i < 3; i++ {

		sm := spinner.NewGroup()

		sp1 := sm.NewSpinner("Task 1", "Done 1")
		sp2 := sm.NewSpinner("Task 2", "Done 2")
		sp3 := sm.NewSpinner("Task 3", "Done 3")
		sp4 := sm.NewSpinner("Task 4", "Done 4")

		sm.StartGroup()

		var wg sync.WaitGroup

		wg.Add(4) // As there are 4 tasks

		// Simulate tasks with random durations
		go func() {
			defer wg.Done()
			time.Sleep(time.Second * time.Duration(rand.Intn(5)))
			sp1.Stop()
		}()

		go func() {
			defer wg.Done()
			time.Sleep(time.Second * time.Duration(rand.Intn(5)))
			sp2.StopWithStatus("success")
		}()

		go func() {
			defer wg.Done()
			time.Sleep(time.Second * time.Duration(rand.Intn(5)))
			sp3.StopWithStatus("failure")
		}()

		go func() {
			defer wg.Done()
			time.Sleep(time.Second * time.Duration(rand.Intn(5)))
			sp4.StopWithStatus("disruption")
		}()

		// Wait for all tasks to complete
		wg.Wait()

		// Stop the spinner group display
		sm.StopGroup()

		fmt.Println("Some terminal output")
		fmt.Println("Now we go again")
	}
}
