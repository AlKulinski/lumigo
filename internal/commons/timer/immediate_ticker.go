package timer

import (
	"context"
	"time"
)

type ImmediateTicker struct {
	C      <-chan time.Time
	ticker *time.Ticker
	cancel context.CancelFunc
}

func (t *ImmediateTicker) Stop() {
	t.cancel()
}

func NewImmediateTicker(d time.Duration, ctx context.Context) *ImmediateTicker {
	ch := make(chan time.Time)
	ticker := time.NewTicker(d)
	ctx, cancel := context.WithCancel(ctx)

	go func() {
		defer func() {
			close(ch)
			ticker.Stop()
		}()

		ch <- time.Now()

		for {
			select {
			case <-ctx.Done():
				return
			case t := <-ticker.C:
				ch <- t
			}
		}
	}()

	return &ImmediateTicker{
		C:      ch,
		ticker: ticker,
		cancel: cancel,
	}
}
