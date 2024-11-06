package ravenTree

import (
	"time"
)

type strategy int

const (
	// Default sets BackoffDelay to 0 seconds.
	// This strategy does not impose any additional delay between retries.
	Default strategy = iota

	// Lineal sets BackoffDelay to 1 second and increases the delay linearly with each retry.
	// For example, the delays will be 1s, 2s, 3s, etc., until MaxDelay is reached.
	Lineal

	// Exponential sets BackoffDelay to 1 second and doubles the delay with each retry.
	// For example, the delays will be 1s, 2s, 4s, 8s, etc., until MaxDelay is reached.
	Exponential
)

const (
	tenSecond = 10 * time.Second
	oneSecond = 1 * time.Second
)

// Backoff defines the strategy for implementing backoff delays in retry operations.
// Fields:
//   - BackoffDelay: Specifies the duration to wait before the next retry.
//   - MaxDelay: Specifies the maximum duration to
//   - Strategy: Determines the type of backoff (Default, Lineal, or Exponential).
type Backoff struct {
	BackoffDelay time.Duration
	MaxDelay     time.Duration
	Strategy     strategy
}

type BackoffOptions func(*Backoff)

// NewBackoff creates a new Backoff with optional parameters.
//
// If no options are provided, it defaults to:
// - BackoffDelay: 0 seconds
// - MaxDelay: 10 seconds
// - strategy: Default
//
// If MaxDelay is set to a value less than BackoffDelay, MaxDelay will be updated
// to match BackoffDelay to ensure valid configuration.
func NewBackoff(opts ...BackoffOptions) *Backoff {
	backoff := &Backoff{
		BackoffDelay: time.Duration(0),
		MaxDelay:     tenSecond,
		Strategy:     Default,
	}

	for _, opt := range opts {
		opt(backoff)
	}

	if backoff.MaxDelay < backoff.BackoffDelay {
		backoff.MaxDelay = backoff.BackoffDelay
	}

	return backoff
}

// WithStrategy sets the backoff strategy type (Default, Lineal, Exponential).
//
// If s is not a valid strategy, the function will default to the Default strategy.
func WithStrategy(s strategy) BackoffOptions {
	return func(b *Backoff) {
		switch s {
		case Lineal, Exponential:
			b.Strategy = s
			b.BackoffDelay = oneSecond
		case Default:
			b.Strategy = Default
		default:
			b.Strategy = Default
		}
	}
}

// WithBackoffDelay sets an initial backoff delay.
func WithBackoffDelay(delay time.Duration) BackoffOptions {
	return func(b *Backoff) {
		b.BackoffDelay = delay
	}
}

// WithMaxDelay set the maximum backoff delay.
func WithMaxDelay(delay time.Duration) BackoffOptions {
	return func(b *Backoff) {
		b.MaxDelay = delay
	}
}

// Next applies the specified backoff strategy to wait before the next retry, adjusting the delay
// based on the strategy (Default, Lineal, or Exponential).
//
// The delay will not exceed maxDelay.
func (b *Backoff) Next() {
	switch b.Strategy {
	case Default:
		time.Sleep(b.BackoffDelay)

	case Lineal:
		time.Sleep(b.BackoffDelay)
		b.BackoffDelay += oneSecond

	case Exponential:
		time.Sleep(b.BackoffDelay)
		b.BackoffDelay *= 2
	}

	if b.BackoffDelay >= b.MaxDelay {
		b.BackoffDelay = b.MaxDelay
	}
}
