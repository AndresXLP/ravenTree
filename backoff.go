package ravenTree

import (
	"time"
)

type strategy int

const (
	Default strategy = iota
	Lineal
	Exponential
)

// BackoffStrategy defines the strategy for implementing backoff delays in retry operations.
// Fields:
//   - BackoffDelay: Specifies the duration to wait before the next retry.
//   - strategy: Determines the type of backoff (Default, Lineal, or Exponential).
type BackoffStrategy struct {
	BackoffDelay time.Duration
	strategy
}

// BackoffDefault initializes a BackoffStrategy with a default backoff delay.
func BackoffDefault() *BackoffStrategy {
	return &BackoffStrategy{
		BackoffDelay: 0,
		strategy:     Default,
	}
}

// BackoffLineal initializes a BackoffStrategy with a linear backoff delay starting at 1 second.
//
// The delay will increase linearly with each retry.
func BackoffLineal() *BackoffStrategy {
	return &BackoffStrategy{
		BackoffDelay: 1 * time.Second,
		strategy:     Lineal,
	}
}

// BackoffExponential initializes a BackoffStrategy with an exponential backoff delay starting at 1 second.
//
// The delay doubles with each retry.
func BackoffExponential() *BackoffStrategy {
	return &BackoffStrategy{
		BackoffDelay: 1 * time.Second,
		strategy:     Exponential,
	}
}

// Next applies the specified backoff strategy to wait before the next retry, adjusting the delay
// based on the strategy (Default, Lineal, or Exponential).
//
// The delay will not exceed maxDelay.
func (b *BackoffStrategy) Next(maxDelay time.Duration) {
	switch b.strategy {
	case Default:
		time.Sleep(b.BackoffDelay)

	case Lineal:
		time.Sleep(b.BackoffDelay)
		b.BackoffDelay += 1 * time.Second

		if b.BackoffDelay > maxDelay {
			b.BackoffDelay = maxDelay
		}

	case Exponential:
		time.Sleep(b.BackoffDelay)
		b.BackoffDelay *= 2

		if b.BackoffDelay > maxDelay {
			b.BackoffDelay = maxDelay
		}
	}
}
