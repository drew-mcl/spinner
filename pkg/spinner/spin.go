package spinner

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

func (s *Spinner) StopWithStatus(status string) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.active = false
	s.Status = ""

	var symbol string

	switch status {
	case "success":
		symbol = color.New(color.FgGreen).Sprint("✔")
	case "failure":
		symbol = color.New(color.FgRed).Sprint("✘")
	case "disruption":
		symbol = color.New(color.FgYellow).Sprint("!")
	}

	gray := color.New(color.FgBlack).Add(color.Faint).SprintFunc()

	s.Status = fmt.Sprintf("%s %s", symbol, gray(s.DoneMsg))
}

type SpinnerManager struct {
	spinners []*Spinner
	mu       sync.Mutex
	started  bool
	wg       sync.WaitGroup
}

func NewGroup() *SpinnerManager {
	return &SpinnerManager{}
}

func (sm *SpinnerManager) NewSpinner(msg, doneMsg string) *Spinner {
	sp := NewSpinner(msg, doneMsg)
	sm.mu.Lock()
	sm.spinners = append(sm.spinners, sp)
	sm.mu.Unlock()

	sm.wg.Add(1)

	go func() {
		for {
			sp.mu.Lock()
			if !sp.active {
				sp.mu.Unlock()
				sm.wg.Done()
				return
			}
			sp.mu.Unlock()
			time.Sleep(100 * time.Millisecond)
		}
	}()

	return sp
}

func (sm *SpinnerManager) StartGroup() {
	sm.mu.Lock()
	defer sm.mu.Unlock()

	if sm.started {
		return
	}

	sm.started = true

	// Reserve space for spinners
	for range sm.spinners {
		fmt.Println()
	}

	for _, s := range sm.spinners {
		s.Start()
	}

	go func() {
		for {
			moveCursorUp(len(sm.spinners))
			sm.mu.Lock()
			for _, s := range sm.spinners {
				s.mu.Lock()
				clearCurrentLine()
				if s.Status != "" {
					fmt.Println(s.Status)
				}
				s.mu.Unlock()
			}
			sm.mu.Unlock()
			moveCursorDown(len(sm.spinners))
			time.Sleep(200 * time.Millisecond)
		}
	}()
}

func (sm *SpinnerManager) WaitForCompletion() {
	sm.wg.Wait()
	time.Sleep(200 * time.Millisecond)
}
