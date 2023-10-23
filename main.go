package main

import (
	"fmt"
	"math/rand"
	"spinner/pkg/spinner"
	"time"
)

func goSomething(sp *spinner.Spinner, duration time.Duration) {
	// Simulate a long-running task
	time.Sleep(duration)
	sp.Stop()
}

func main() {

	fmt.Println("Some text before the spinners")
	fmt.Println("----------------------------------------")

	rand.Seed(time.Now().UnixNano())
	sm := spinner.NewGroup()

	numTasks := 5
	for i := 1; i <= numTasks; i++ {
		msg := fmt.Sprintf("Task %d", i)
		doneMsg := fmt.Sprintf("Done %d", i)

		randomDuration := time.Duration(rand.Intn(5)+1) * time.Second

		sp := sm.NewSpinner(msg, doneMsg)
		go goSomething(sp, randomDuration)
	}

	// Start the group of spinners
	sm.StartGroup()
	// Wait for all tasks to complete
	sm.WaitForCompletion()

	fmt.Println("Some text after the spinners")
	fmt.Println("----------------------------------------")
}
