package main

import (
	"fmt"
	"spinner/pkg/spinner" // replace with the actual path
	"time"
)

func main() {
	fmt.Println("Some stuff before the spinner group")
	fmt.Println("--------------------")

	sm := spinner.NewGroup()
	s1 := sm.NewSpinner("Loading...", "Done")
	s2 := sm.NewSpinner("Processing...", "Failed")
	s3 := sm.NewSpinner("Warning...", "Disrupted")

	sm.StartGroup()

	time.Sleep(3 * time.Second)
	s1.Stop()

	time.Sleep(2 * time.Second)
	s2.Stop()

	time.Sleep(1 * time.Second)
	s3.Stop()

	fmt.Println("--------------------")
	fmt.Println("Some stuff after the spinner group")
}
