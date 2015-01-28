package timer

import (
	"time"
)

// one dispatcher per goroutine
type Dispatcher struct {
	ChanTimer chan *Timer
}

type Timer struct {
	t  *time.Timer
	cb func()
}

func NewDispatcher(l int) *Dispatcher {
	disp := new(Dispatcher)
	disp.ChanTimer = make(chan *Timer, l)
	return disp
}

func (disp *Dispatcher) AfterFunc(d time.Duration, cb func()) *Timer {
	if cb == nil {
		cb = func() {}
	}

	t := new(Timer)
	t.cb = cb
	t.t = time.AfterFunc(d, func() {
		disp.ChanTimer <- t
	})
	return t
}

func (t *Timer) Stop() {
	t.t.Stop()
	t.cb = nil
}

func (t *Timer) IsRunning() bool {
	return t.cb != nil
}

func (t *Timer) Cb() {
	if t.cb != nil {
		t.cb()
		t.cb = nil
	}
}