package api

import (
	"time"
)

func NewRateLimit(limit int, interval time.Duration) *RateLimit {
	return &RateLimit{limit: limit, interval: interval, remaining: limit, reset: time.Now()}
}

func (r *RateLimit) Wait() {
	if r.remaining == 0 {
		time.Sleep(time.Until(r.reset))
		r.Reset()
	}
}

func (r *RateLimit) Call() {
	if r.limit == r.remaining {
		r.Reset()
	}
	r.Wait()
	r.remaining--
}

func (r *RateLimit) GuestCall() bool {
	if r.remaining == 0 {
		return false
	}
	r.remaining--
	return true
}

func (r *RateLimit) Reset() {
	r.remaining = r.limit
	r.reset = time.Now().Add(r.interval)
}
