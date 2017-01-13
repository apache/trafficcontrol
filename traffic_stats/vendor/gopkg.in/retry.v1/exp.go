package retry

import (
	"time"
)

// Exponential represents an exponential backoff retry strategy.
// To limit the number of attempts or their overall duration, wrap
// this in LimitCount or LimitDuration.
type Exponential struct {
	// Initial holds the initial delay.
	Initial time.Duration
	// Factor holds the factor that the delay time will be multiplied
	// by on each iteration.
	Factor float64
	// MaxDelay holds the maximum delay between the start
	// of attempts. If this is zero, there is no maximum delay.
	MaxDelay time.Duration
}

type exponentialTimer struct {
	strategy Exponential
	start    time.Time
	end      time.Time
	delay    time.Duration
}

// NewTimer implements Strategy.NewTimer.
func (r Exponential) NewTimer(now time.Time) Timer {
	return &exponentialTimer{
		strategy: r,
		start:    now,
		delay:    r.Initial,
	}
}

// NextSleep implements Timer.NextSleep.
func (a *exponentialTimer) NextSleep(now time.Time) (time.Duration, bool) {
	sleep := a.delay - now.Sub(a.start)
	if sleep <= 0 {
		sleep = 0
	}
	// Set the start of the next try.
	a.start = now.Add(sleep)
	a.delay = time.Duration(float64(a.delay) * a.strategy.Factor)
	if a.strategy.MaxDelay > 0 && a.delay > a.strategy.MaxDelay {
		a.delay = a.strategy.MaxDelay
	}
	return sleep, true
}
