package spinner

import (
	"context"
	"fmt"
	"sync"
	"time"

	"github.com/fatih/color"
)

type Spinner struct {
	mu        sync.Mutex
	active    bool
	Msg       string
	DoneMsg   string
	Status    string
	spinChars []string
	stopCh    chan struct{}
	ctx       context.Context
	cancel    context.CancelFunc
	manager   *SpinnerManager
	Name      string
}

func NewSpinner(msg, doneMsg string, ctx context.Context) *Spinner {
	return &Spinner{
		spinChars: []string{"◐", "◓", "◑", "◒"},
		Msg:       msg,
		DoneMsg:   doneMsg,
		ctx:       ctx,
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
			case <-s.ctx.Done():
				s.StopWithStatus("disruption", s.manager.disruptionMessage)
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
	defer s.mu.Unlock()

	if !s.active {
		return
	}

	s.active = false
	gray := color.New(color.FgBlack).Add(color.Faint).SprintFunc()
	s.Status = fmt.Sprintf("✔ %s", gray(s.DoneMsg))

	//To ensure there is no panic if StopWithStatus is called before Start
	if s.stopCh != nil {
		close(s.stopCh)
		s.stopCh = nil
	}

}

func (s *Spinner) StopWithStatus(status, customMsg string) {
	s.mu.Lock()
	defer s.mu.Unlock()

	if !s.active {
		return
	}

	s.active = false
	var message string

	switch status {
	case "success":
		message = fmt.Sprintf("%s %s", color.New(color.FgGreen).Sprint("✔"), customMsg)
	case "failure":
		message = fmt.Sprintf("%s %s", color.New(color.FgRed).Sprint("✘"), customMsg)
	case "disruption":
		message = fmt.Sprintf("%s %s", color.New(color.FgYellow).Sprint("!"), customMsg)
	case "done":
		message = fmt.Sprintf("%s %s", color.New(color.FgBlack).Add(color.Faint).Sprint("✔"), customMsg)
	default:
		message = fmt.Sprintf("%s %s", color.New(color.FgBlack).Add(color.Faint).Sprint("?"), customMsg)
	}

	s.Status = message

	if s.stopCh != nil {
		close(s.stopCh)
		s.stopCh = nil
	}
}
