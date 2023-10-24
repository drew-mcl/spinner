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
	stopCh    chan struct{}
}

func NewSpinner(msg, doneMsg string) *Spinner {
	return &Spinner{
		spinChars: []string{"◐", "◓", "◑", "◒"},
		Msg:       msg,
		DoneMsg:   doneMsg,
	}
}

func (s *Spinner) Start() {
	s.mu.Lock() // Locking to avoid race conditions
	s.active = true
	s.stopCh = make(chan struct{})
	s.mu.Unlock()

	blue := color.New(color.FgBlue).SprintFunc()

	go func() {
		for {
			select {
			case <-s.stopCh:
				return
			default:
				s.mu.Lock()
				for i := 0; i < len(s.spinChars) && s.active; i++ { // added s.active check
					s.Status = fmt.Sprintf("%s %s", blue(s.spinChars[i]), s.Msg)
					time.Sleep(200 * time.Millisecond)
				}
				s.mu.Unlock()
			}
		}
	}()
}

func (s *Spinner) Stop() {
	s.mu.Lock() // Locking to ensure thread safety
	s.active = false
	gray := color.New(color.FgBlack).Add(color.Faint).SprintFunc()
	s.Status = fmt.Sprintf("✔ %s", gray(s.DoneMsg))
	s.mu.Unlock()

	close(s.stopCh) // Signal the spinner to stop
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
	doneCh   chan bool // To signal when all spinners are done
	quit     chan bool // To signal the updating goroutine to stop

}

func NewGroup() *SpinnerManager {
	return &SpinnerManager{
		quit:   make(chan bool),
		doneCh: make(chan bool),
	}
}

func (sm *SpinnerManager) NewSpinner(msg, doneMsg string) *Spinner {
	sp := NewSpinner(msg, doneMsg)
	sm.mu.Lock()
	sm.spinners = append(sm.spinners, sp)
	sm.mu.Unlock()

	go func() {
		<-sp.stopCh
		for _, spinner := range sm.spinners {
			if spinner.active {
				return // If any spinner is still active, return
			}
		}
		sm.doneCh <- true // If all spinners are done, signal doneCh
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
	hideCursor() // Added hideCursor from the second code

	for _, s := range sm.spinners {
		s.Start()
	}

	firstDraw := true // Added the firstDraw logic from the second code
	go func() {
		for {
			select {
			case <-sm.quit: // Listen to quit signal
				return
			default:
				sm.mu.Lock()
				if !firstDraw {
					moveCursorUp(len(sm.spinners))
				}
				for _, s := range sm.spinners {
					clearCurrentLine()
					if s.Status != "" {
						fmt.Println(s.Status)
					} else {
						// Print a placeholder for inactive spinners
						fmt.Println()
					}
				}
				sm.mu.Unlock()

				if firstDraw {
					firstDraw = false
				}
				time.Sleep(100 * time.Millisecond)
			}
		}
	}()
}

func (sm *SpinnerManager) StopGroup() {
	time.Sleep(300 * time.Millisecond) // Sleep to allow the last update
	sm.quit <- true                    // Signal to stop updating spinners
	sm.ResetTerminal()
}

func (sm *SpinnerManager) ResetTerminal() {
	showCursor() // Restore the cursor visibility
}
