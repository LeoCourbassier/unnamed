package main

import (
	"sync"
	"time"
)

// NewDebouncer returns a debounced function that takes another functions as its argument.
// This function will be called when the debounced function stops being called
// for the given duration.
// The debounced function can be invoked with different functions, if needed,
// the last one will win.
func NewDebouncer(after time.Duration) func(f func()) {
	d := &Debouncer{after: after}

	return func(f func()) {
		d.add(f)
	}
}

type Debouncer struct {
	mu    sync.Mutex
	after time.Duration
	timer *time.Timer
}

func (d *Debouncer) add(f func()) {
	d.mu.Lock()
	defer d.mu.Unlock()

	if d.timer != nil {
		d.timer.Stop()
	}
	d.timer = time.AfterFunc(d.after, f)
}
