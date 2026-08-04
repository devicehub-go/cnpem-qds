package watchdog

import (
	"context"
	"errors"
	"time"
)

var (
	ErrTimeout = errors.New("alive check does not occurs on deadline")
)

type Watchdog struct {
	ctx     context.Context
	timeout time.Duration
	reset   chan struct{}
	err     chan error
}

// Creates a new instance of watchdog
func New(ctx context.Context, timeout time.Duration) *Watchdog {
	return &Watchdog{
		ctx:     ctx,
		timeout: timeout,
		reset:   make(chan struct{}, 1),
		err:     make(chan error, 1),
	}
}

// Starts to monitor the alive check of system
func (w *Watchdog) Start() {
	go w.run()
}

func (w *Watchdog) run() {
	timer := time.NewTimer(w.timeout)
	defer timer.Stop()
	defer close(w.err)

	for {
		select {
		case <-w.ctx.Done():
			return
		case <-timer.C:
			w.err <- ErrTimeout
			return
		case <-w.reset:
			if !timer.Stop() {
				select {
				case <-timer.C:
				default:
				}
			}
			timer.Reset(w.timeout)
		}
	}
}

// Sets the alive status
func (w *Watchdog) Kick() {
	select {
	case w.reset <- struct{}{}:
	default:
	}
}

// Returns a channel that indicates if a timeout occurs
func (w *Watchdog) Err() <-chan error {
	return w.err
}
