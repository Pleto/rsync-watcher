package main

import (
	"sync"
	"time"
)

// EventDebouncer helps delay an action until a certain amount of time has passed
// without any new triggers.
// If it gets triggered again before the time is up, the timer resets,
// and the action is postponed.
type EventDebouncer struct {
	interval time.Duration
	action   func()
	timer    *time.Timer
	mu       sync.Mutex
}

func NewEventDebouncer(interval time.Duration, action func()) *EventDebouncer {
	return &EventDebouncer{
		interval: interval,
		action:   action,
	}
}

func (d *EventDebouncer) Trigger() {
	d.mu.Lock()
	defer d.mu.Unlock()

	if d.timer != nil {
		d.timer.Reset(d.interval)
	} else {
		d.timer = time.AfterFunc(d.interval, func() {
			d.mu.Lock()
			defer d.mu.Unlock()
			d.action()
			d.timer = nil
		})
	}
}
