package chimer

import (
	"time"
)

type Chime int64

func (w Chime) String() string {
	switch w {
	case Hour:
		return "Hour"
	case QuarterPast:
		return "QuarterPast"
	case HalfPast:
		return "HalfPast"
	case QuarterTo:
		return "QuarterTo"
	case None:
		return "None"
	default:
		panic("unknown value")
	}
}

const (
	None Chime = iota
	Hour
	QuarterPast
	HalfPast
	QuarterTo
)

func secondsFromMidnight(t time.Time) int {
	hour := t.Hour()
	minute := t.Minute()
	seconds := t.Second()

	return (hour*60+minute)*60 + seconds
}

func GetChime(t time.Time, tolerance time.Duration) (int, Chime) {
	t = t.Local()
	hour := t.Hour()
	seconds := secondsFromMidnight(t)

	if hour > 12 {
		hour -= 12
	}

	if hourItIs, increment := isHour(seconds, int(tolerance.Seconds())); hourItIs {
		return hour + increment, Hour
	}

	if isHalfHour(seconds, int(tolerance.Seconds())) {
		return hour, HalfPast
	}

	if isQuarterHour(seconds, int(tolerance.Seconds())) {
		if t.Minute() > 40 {
			return hour, QuarterTo
		}

		return hour, QuarterPast
	}

	return hour, None
}

const quarterlySeconds = 15 * 60
const halfHourlySeconds = 2 * quarterlySeconds
const hourSeconds = 4 * quarterlySeconds

func isHour(s, t int) (bool, int) {
	return isMultipleOfWithTolerance(hourSeconds, s, t)
}

func isQuarterHour(s, t int) bool {
	ret, _ := isMultipleOfWithTolerance(quarterlySeconds, s, t)
	return ret
}

func isHalfHour(s, t int) bool {
	ret, _ := isMultipleOfWithTolerance(halfHourlySeconds, s, t)
	return ret
}

func isMultipleOfWithTolerance(m, s, t int) (bool, int) {
	r := s % m

	// if we are on the dot
	if r == 0 {
		return true, 0
	}

	// if there is no tolerance
	if t == 0 {
		return false, 0
	}

	// a negative numbers for tolerance?!
	if t < 0 {
		t = -t
	}

	if r >= (m - t) {
		return true, 1
	}

	if r <= t {
		return true, 0
	}

	return false, 0
}
