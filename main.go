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
	s1.StopWithStatus("success")

	time.Sleep(2 * time.Second)
	s2.StopWithStatus("failure")

	time.Sleep(1 * time.Second)
	s3.StopWithStatus("disruption")

	//Required to get the last symbol to render, will move out once made async and non-blocking
	time.Sleep(200 * time.Millisecond)

	fmt.Println("--------------------")
	fmt.Println("Some stuff after the spinner group")
}
