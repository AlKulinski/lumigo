package timer

import "time"

type ImmediateTicker struct {
	C      <-chan time.Time
	ticker *time.Ticker
}

func (t *ImmediateTicker) Stop() {
	t.ticker.Stop()
}

func NewImmediateTicker(d time.Duration) *ImmediateTicker {
	ch := make(chan time.Time)
	ticker := time.NewTicker(d)

	go func() {
		defer close(ch)
		defer ticker.Stop()

		ch <- time.Now()

		for t := range ticker.C {
			ch <- t
		}
	}()

	return &ImmediateTicker{
		C:      ch,
		ticker: ticker,
	}
}
