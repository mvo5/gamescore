package main

import (
	"time"
	
	. "gopkg.in/check.v1"
)

var _ = Suite(&countdownSuite{})

type countdownSuite struct{}

func (s *countdownSuite) TestStartStop(c *C) {
	ct := NewCountdown(20*time.Second)
	c.Check(ct.TimeLeft(), Equals, time.Duration(20*time.Second))
	// not running yet
	time.Sleep(100*time.Millisecond)
	c.Check(ct.TimeLeft(), Equals, time.Duration(20*time.Second))

	ct.Start()
	time.Sleep(100*time.Millisecond)
	c.Check(ct.TimeLeft().Truncate(time.Second), Equals, time.Duration(19*time.Second))

	ct.Stop()
	time.Sleep(100*time.Millisecond)
	c.Check(ct.TimeLeft().Truncate(time.Second), Equals, time.Duration(19*time.Second))
}

func (s *countdownSuite) TestFormatTime(c *C) {
	countdown := NewCountdown(320 * time.Second)
	c.Check(countdown.String(), Equals, "05:20")

	countdown.Set(50 * time.Second)
	c.Check(countdown.String(), Equals, "00:500")

	countdown.Set( 5 * time.Second)
	c.Check(countdown.String(), Equals, "00:050")

}
