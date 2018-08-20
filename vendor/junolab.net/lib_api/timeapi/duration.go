package timeapi

import (
	"time"
)

// Duration represents duration in milliseconds
// We agreed to use this type for duration between ms to make testing easier.
type Duration int64

// NewDuration creates Duration from time.Duration
func NewDuration(d time.Duration) Duration {
	return Duration(d.Nanoseconds() / int64(time.Millisecond))
}

// NewDurationFromSeconds creates Duration from seconds (int64)
func NewDurationFromSeconds(s int64) Duration {
	return Duration(s * 1000)
}

// GoDuration returns time.Duration
func (d Duration) GoDuration() time.Duration {
	return time.Duration(d) * time.Millisecond
}

// Seconds returns duration in seconds as int64
func (d Duration) Seconds() int64 {
	return int64(d / 1000)
}

// IsZero checks Duration on zero
func (d Duration) IsZero() bool {
	return d == 0
}
