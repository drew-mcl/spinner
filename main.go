package main

import (
	"fmt"
	"sync"
	"time"

	"github.com/fatih/color"
)

func clearCurrentLine() {
	fmt.Print("\033[K")
}

func moveCursorUp(lines int) {
	fmt.Printf("\033[%dA", lines)
}

func moveCursorDown(lines int) {
	fmt.Printf("\033[%dB", lines)
}

func hideCursor() {
	fmt.Print("\033[?25l")
}

func showCursor() {
	fmt.Print("\033[?25h")
}

type Spinner struct {
	mu        sync.Mutex
	active    bool
	Msg       string
	DoneMsg   string
	Status    string
	spinChars []string
}

func NewSpinner(msg, doneMsg string) *Spinner {
	return &Spinner{
		spinChars: []string{"◐", "◓", "◑", "◒"},
		Msg:       msg,
		DoneMsg:   doneMsg,
	}
}

func (s *Spinner) Start() {
	s.mu.Lock()
	s.active = true
	s.mu.Unlock()

	blue := color.New(color.FgBlue).SprintFunc()

	go func() {
		for {
			for i := 0; i < len(s.spinChars); i++ {
				s.mu.Lock()
				if !s.active {
					s.mu.Unlock()
					return
				}
				s.Status = fmt.Sprintf("%s %s", blue(s.spinChars[i]), s.Msg)
				s.mu.Unlock()
				time.Sleep(200 * time.Millisecond)
			}
		}
	}()
}

func (s *Spinner) Stop() {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.active = false
	gray := color.New(color.FgBlack).Add(color.Faint).SprintFunc()
	s.Status = fmt.Sprintf("✔ %s", gray(s.DoneMsg))
}

type SpinnerManager struct {
	spinners []*Spinner
	mu       sync.Mutex
}

func NewSpinnerManager() *SpinnerManager {
	return &SpinnerManager{}
}

func (sm *SpinnerManager) AddSpinner(s *Spinner) {
	sm.mu.Lock()
	defer sm.mu.Unlock()
	sm.spinners = append(sm.spinners, s)
}

func (sm *SpinnerManager) Render() {
	sm.mu.Lock()
	defer sm.mu.Unlock()

	// Move cursor to initial position
	moveCursorUp(len(sm.spinners))

	for _, s := range sm.spinners {
		s.mu.Lock()
		clearCurrentLine()
		fmt.Println(s.Status)
		s.mu.Unlock()
	}
}

func main() {
	hideCursor()
	defer showCursor()

	sm := NewSpinnerManager()

	spinner1 := NewSpinner("Loading...", "Done")
	spinner2 := NewSpinner("Processing...", "Failed")
	spinner3 := NewSpinner("Warning...", "Disrupted")

	sm.AddSpinner(spinner1)
	sm.AddSpinner(spinner2)
	sm.AddSpinner(spinner3)

	fmt.Println("Whatever comes above")
	fmt.Println("--------------------")

	spinner1.Start()
	spinner2.Start()
	spinner3.Start()

	// Move the cursor down to make space for the spinners
	moveCursorDown(len(sm.spinners))

	go func() {
		for {
			sm.Render()
			time.Sleep(200 * time.Millisecond)
		}
	}()

	time.Sleep(3 * time.Second)
	spinner1.Stop()

	time.Sleep(2 * time.Second)
	spinner2.Stop()

	time.Sleep(1 * time.Second)
	spinner3.Stop()

	time.Sleep(2 * time.Second) // Just to keep the final state visible for a bit

	fmt.Println("--------------------")
	fmt.Println("Whatever comes after")
}
