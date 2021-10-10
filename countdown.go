package main

import (
	"fmt"
	"time"
)

type Countdown struct {
	t        time.Time
	duration time.Duration
	running  bool
}

func NewCountdown(d time.Duration) *Countdown {
	return &Countdown{duration: d}
}

func (c *Countdown) Set(d time.Duration) {
	c.duration = d
}

func (c *Countdown) Start() {
	if c.running {
		return
	}
	c.running = true
	c.t = time.Now()
}

func (c *Countdown) Stop() {
	if !c.running {
		return
	}
	c.duration = c.TimeLeft()
	c.running = false
}

func (c *Countdown) TimeLeft() time.Duration {
	if !c.running {
		return c.duration
	}
	endTime := c.t.Add(c.duration)
	return endTime.Sub(time.Now())
}

func (c *Countdown) String() string {
	left := c.TimeLeft()

	// format the time nicely
	min := left / (60 * time.Second)
	sec := left % (60 * time.Second)
	// sub 1 min gets 100 millisecond resolution
	// XXX: configuratble?
	if min >= 1 {
		sec /= time.Second
		return fmt.Sprintf("%02d:%02d", min, sec)
	} else {
		sec /= 100 * time.Millisecond
		return fmt.Sprintf("%02d:%03d", min, sec)
	}
}

