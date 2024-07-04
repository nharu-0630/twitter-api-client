package api

import (
	"net/http"
	"strconv"
	"time"
)

type RateLimit struct {
	limit     int
	remaining int
	reset     int
}

func (r *RateLimit) Wait() {
	if r.remaining == 0 {
		t := time.Unix(int64(r.reset), 0)
		time.Sleep(time.Until(t))
	}
}

func (r *RateLimit) Update(header http.Header) {
	if limit := header.Get("x-rate-limit-limit"); limit != "" {
		intLimit, err := strconv.Atoi(limit)
		if err == nil {
			r.limit = intLimit
		}
	}
	if remaining := header.Get("x-rate-limit-remaining"); remaining != "" {
		intRemaining, err := strconv.Atoi(remaining)
		if err == nil {
			r.remaining = intRemaining
		}
	}
	if reset := header.Get("x-rate-limit-reset"); reset != "" {
		intReset, err := strconv.Atoi(reset)
		if err == nil {
			r.reset = intReset
		}
	}
}
