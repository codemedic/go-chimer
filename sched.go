package chimer

import (
	"time"
)

type Scheduler struct {
	stop      chan struct{}
	frequency time.Duration
}

// NewScheduler creates a scheduler with the frequency of f duration. Value of f is restricted to 1s <= f <= 1m.
func NewScheduler(f time.Duration) *Scheduler {
	if f < time.Second {
		f = time.Second
	} else if f > time.Minute {
		f = time.Minute
	}

	return &Scheduler{
		stop:      make(chan struct{}),
		frequency: f,
	}
}

// Stop stops the scheduler.
func (s *Scheduler) Stop() {
	close(s.stop)
}

// EveryQuarterHour calls f, with current time t and tolerance, every 0th, 15th, 30th and 45th minutes of the hour. The
// value for tolerance is based on specified frequency for Scheduler.
func (s *Scheduler) EveryQuarterHour(f func(t time.Time, tolerance time.Duration)) {
	tolerance := s.frequency / 2

	var last time.Time
	for {
		t := time.Now().Local()
		nextMinuteWait := time.Minute - (time.Duration(t.Second())*time.Second + time.Duration(t.Nanosecond()))

		wait := s.frequency
		if nextMinuteWait < s.frequency {
			wait = nextMinuteWait
		}

		timer := time.NewTimer(wait)
		select {
		case t = <-timer.C:
			switch t.Minute() {
			case 0, 15, 30, 45:
				if last.IsZero() || t.Sub(last) > time.Minute {
					last = t
					go f(t, tolerance)
				}
			}
		case <-s.stop:
			timer.Stop()
			return
		}
	}
}
