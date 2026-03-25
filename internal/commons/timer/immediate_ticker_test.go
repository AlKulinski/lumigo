package timer

import (
	"context"
	"testing"
	"time"
)

const tickInterval = 50 * time.Millisecond

// First tick must arrive well before the first scheduled interval fires.
func TestImmediateTicker_firesImmediately(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	ticker := NewImmediateTicker(tickInterval, ctx)

	select {
	case <-ticker.C:
		// ok
	case <-time.After(tickInterval / 2):
		t.Error("first tick did not fire immediately; waited half the interval with no tick")
	}
}

// Second tick must arrive approximately one interval after construction.
func TestImmediateTicker_secondTickAfterInterval(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	ticker := NewImmediateTicker(tickInterval, ctx)

	start := time.Now()
	<-ticker.C // immediate tick

	select {
	case <-ticker.C:
		elapsed := time.Since(start)
		if elapsed < tickInterval-5*time.Millisecond {
			t.Errorf("second tick arrived too early: %v (interval=%v)", elapsed, tickInterval)
		}
	case <-time.After(3 * tickInterval):
		t.Error("second tick never arrived")
	}
}

// Ticks must keep arriving at each interval.
func TestImmediateTicker_multipleTicks(t *testing.T) {
	const n = 4
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	ticker := NewImmediateTicker(tickInterval, ctx)

	timeout := time.After(time.Duration(n+1) * tickInterval * 2)
	for i := range n {
		select {
		case _, ok := <-ticker.C:
			if !ok {
				t.Fatalf("channel closed unexpectedly after %d ticks", i)
			}
		case <-timeout:
			t.Fatalf("timed out waiting for tick %d/%d", i+1, n)
		}
	}
}

// Cancelling the context must stop ticks and close the channel.
func TestImmediateTicker_contextCancel_closesChannel(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())

	ticker := NewImmediateTicker(tickInterval, ctx)
	<-ticker.C // drain the immediate tick
	cancel()

	// Channel must close after cancellation.
	select {
	case _, ok := <-ticker.C:
		if ok {
			t.Error("received a tick after context cancellation")
		}
		// channel closed — expected
	case <-time.After(3 * tickInterval):
		t.Error("channel was not closed after context cancellation")
	}
}

// Stop() must cancel the internal context, halt ticks, and close the channel.
func TestImmediateTicker_stop_haltsTicks(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	ticker := NewImmediateTicker(tickInterval, ctx)
	<-ticker.C // drain the immediate tick
	ticker.Stop()

	select {
	case _, ok := <-ticker.C:
		if ok {
			t.Error("received a tick after Stop()")
		}
		// channel closed — expected
	case <-time.After(3 * tickInterval):
		t.Error("channel was not closed after Stop()")
	}
}

// Ticks must carry timestamps that are non-zero and monotonically increasing.
func TestImmediateTicker_tickTimestamps(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	ticker := NewImmediateTicker(tickInterval, ctx)

	var prev time.Time
	timeout := time.After(5 * tickInterval)
	for range 3 {
		select {
		case ts := <-ticker.C:
			if ts.IsZero() {
				t.Error("received zero timestamp")
			}
			if !prev.IsZero() && !ts.After(prev) {
				t.Errorf("tick timestamp %v is not after previous %v", ts, prev)
			}
			prev = ts
		case <-timeout:
			t.Fatal("timed out waiting for tick")
		}
	}
}

