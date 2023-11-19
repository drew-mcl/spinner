package spinner

import (
	"context"
	"fmt"
	"spinner/pkg/term"
	"sync"
	"time"
)

type SpinnerManager struct {
	spinners          []*Spinner
	mu                sync.Mutex
	started           bool
	wg                sync.WaitGroup
	doneCh            chan bool // To signal when all spinners are done
	quit              chan bool // To signal the updating goroutine to stop
	ctx               context.Context
	cancel            context.CancelFunc
	disruptionMessage string
}

func NewGroup() *SpinnerManager {
	ctx, cancel := context.WithCancel(context.Background())
	return &SpinnerManager{
		quit:   make(chan bool),
		doneCh: make(chan bool),
		ctx:    ctx,
		cancel: cancel,
	}
}

func (sm *SpinnerManager) NewSpinner(msg, doneMsg string) *Spinner {
	sp := &Spinner{
		spinChars: []string{"◐", "◓", "◑", "◒"},
		Msg:       msg,
		DoneMsg:   doneMsg,
		ctx:       sm.ctx,
		manager:   sm,
	}
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

func (sm *SpinnerManager) DisruptAllSpinners(disruptionMessage string) {
	sm.mu.Lock()
	sm.disruptionMessage = disruptionMessage
	sm.mu.Unlock()

	sm.cancel() // Cancels the manager's context
}

func (sm *SpinnerManager) StartGroup() {
	sm.mu.Lock()
	defer sm.mu.Unlock()

	if sm.started {
		return
	}

	sm.started = true
	term.HideCursor() // Added hideCursor from the second code

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
					term.MoveCursorUp(len(sm.spinners))
				}
				for _, s := range sm.spinners {
					term.ClearCurrentLine()
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
	term.ShowCursor()
}

func (sm *SpinnerManager) DisruptAllNotCompleted(customMsg string) {
	sm.mu.Lock()
	defer sm.mu.Unlock()

	for _, spinner := range sm.spinners {
		if spinner.active {
			spinner.StopWithStatus("disruption", customMsg)
		}
	}
}
