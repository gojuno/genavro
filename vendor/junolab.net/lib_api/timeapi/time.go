package timeapi

import (
	"time"
	"unsafe"
)

// Time represents timestamp in milliseconds from 1 January 1970
// We agreed to use this type for duration between ms to make testing easier.
type Time int64

// Now return the current time
func Now() Time {
	return NewTime(time.Now().UTC())
}

// NewTime creates a time from time.Time
func NewTime(t time.Time) Time {
	nanos := t.UnixNano()
	if nanos <= 0 {
		return Time(0)
	}
	return Time(nanos / int64(time.Millisecond))
}

// Add appends duration to the current time
func (t Time) Add(d Duration) Time {
	return t + Time(d)
}

// Sub returns a duration of time window
func (t Time) Sub(t2 Time) Duration {
	return Duration(t - t2)
}

// Actual returns now() if the local time is undefined
func (t Time) Actual() Time {
	if t.IsZero() {
		return Now()
	}
	return t
}

// GoTime returns time.Time corresponding to the given time
// It returns empty GoTime in case of empty Time, due to it converts empty Time to non zero GoTime
func (t Time) GoTime() time.Time {
	if t.IsZero() {
		return time.Time{}
	}
	return time.Unix(int64(t)/1000, (int64(t)%1000)*int64(time.Millisecond)).UTC()
}

// UnixMilliseconds returns timestamp in UTC milliseconds.
func (t Time) UnixMilliseconds() int64 {
	return int64(t)
}

// IsZero checks Time on zero
func (t Time) IsZero() bool {
	return t <= 0
}

// String returns string format of Time
func (t Time) String() string {
	return t.GoTime().Format(time.RFC3339)
}

// TimestampsAsInts return overview timestamps as a slice of int64
func TimestampsAsInts(t []Time) []int64 {
	return *(*[]int64)(unsafe.Pointer(&t))
}
