package types

import (
	"fmt"
	"time"
)

// TimingInfo is a member of MsgMeta. Internal members of TimingInfo can be changed via different go-routines,
// so it's not thread-safe in common case. We assume that changing Delta won't cause any issue because it is safe
// for our logic and Delta type is time.Duration. We encode the latest timing info and pass it further
// to another micro-service.
type TimingInfo struct {
	// Initial timeout.
	Timeout time.Duration `json:"timeout"`
	// Aggregates execution time in local context.
	LocalDelta time.Duration `json:"-"`
	// Aggregates execution time in previous context.
	Delta time.Duration `json:"delta"`
	// Timing start time in local context.
	TimingStart time.Time `json:"-"`
}

func NewTimingInfo() *TimingInfo {
	return &TimingInfo{}
}

func (t *TimingInfo) SetTimeout(timeout time.Duration) {
	t.Timeout = timeout
}

func (t *TimingInfo) SetTimeoutDefault() {
	t.Timeout = 2 * time.Second
}

// CheckPoint registers timing check point.
func (t *TimingInfo) CheckPoint() error {
	// Does not collect info if timeout is 0.
	if t.Timeout == 0 {
		return fmt.Errorf("No timeout")
	}

	// Initialize at first checkpoint.
	if t.TimingStart.IsZero() {
		t.TimingStart = time.Now()
		return nil
	}

	t.LocalDelta = t.GetElapsedTime()
	return nil
}

// FixLocalDelta moves local delta to global delta. TimingInfo becomes invalid for further using.
// We can encode it and have to destroy it in local context.
func (t *TimingInfo) FixLocalDelta() {
	// Does not collect info if timeout is 0.
	if t.Timeout == 0 {
		return
	}
	t.Delta += t.LocalDelta
	t.LocalDelta = 0
}

// GetElapsedTime returns nanoseconds elapsed from request creation.
func (t *TimingInfo) GetElapsedTime() time.Duration {
	return time.Since(t.TimingStart)
}

func (t *TimingInfo) IsDeadline() bool {
	return t.Timeout != 0 && t.GetElapsedTime()+t.Delta >= t.Timeout
}

func (t *TimingInfo) DurationUntilDeadline() (time.Duration, bool) {
	if t.Timeout == 0 {
		return 0, false
	}
	elapsed := t.GetElapsedTime() + t.Delta
	if elapsed >= t.Timeout {
		return 0, false
	}
	return t.Timeout - elapsed, true
}

func (t *TimingInfo) IsNextTimeoutValid(timeout time.Duration) bool {
	if t.Timeout == 0 {
		return true
	}
	return timeout <= t.Timeout
}

func (t *TimingInfo) String() string {
	return fmt.Sprintf("Elapsed: %v from %v", t.GetElapsedTime()+t.Delta, t.Timeout)
}
