package ravenTree_test

import (
	"testing"
	"time"

	tree "github.com/AndresXLP/ravenTree"
)

func TestNewBackoff_WithDefaultValues(t *testing.T) {
	maxBackoffDelay := 10 * time.Second
	backoff := tree.NewBackoff(
		tree.WithStrategy(tree.Default))

	if backoff.BackoffDelay != 0 {
		t.Errorf("Expected BackoffDelay to be 0, got %v", backoff.BackoffDelay)
	}

	if backoff.MaxDelay != maxBackoffDelay {
		t.Errorf("Expected MaxDelay to be %v, got %v", maxBackoffDelay, backoff.MaxDelay)
	}

	if backoff.Strategy != tree.Default {
		t.Errorf("Expected strategy to be Default, got %v", backoff.Strategy)
	}

	since := time.Now()

	backoff.Next()

	elapsed := time.Since(since)

	tolerance := 50 * time.Microsecond
	if elapsed > tolerance {
		t.Errorf("Expected elapsed to be within %v, got %v", tolerance, elapsed)
	}
}

func TestNewBackoff_WithBackoffOptionsLinealStrategy(t *testing.T) {
	customBackoffDelay := 3 * time.Second
	customMaxDelay := 5 * time.Second
	backoff := tree.NewBackoff(
		tree.WithStrategy(tree.Lineal),
		tree.WithBackoffDelay(customBackoffDelay),
		tree.WithMaxDelay(customMaxDelay),
	)

	if backoff.BackoffDelay != customBackoffDelay {
		t.Errorf("Expected BackoffDelay to be %v, got %v", customBackoffDelay, backoff.BackoffDelay)
	}

	if backoff.MaxDelay != customMaxDelay {
		t.Errorf("Expected MaxDelay to be %v, got %v", customMaxDelay, backoff.MaxDelay)
	}

	if backoff.Strategy != tree.Lineal {
		t.Errorf("Expected Strategy to be %v, got %v", tree.Lineal, backoff.Strategy)
	}

	since := time.Now()

	backoff.Next()

	elapsed := time.Since(since)

	tolerance := customBackoffDelay + 50*time.Millisecond
	if elapsed > tolerance {
		t.Errorf("Expected elapsed to be %v, got %v", tolerance, elapsed)
	}
}

func TestNewBackoff_WithBackoffOptionsExponentialStrategy(t *testing.T) {
	backoffDelayExpected := 1 * time.Second
	backoff := tree.NewBackoff(
		tree.WithStrategy(tree.Exponential),
	)

	if backoff.BackoffDelay != backoffDelayExpected {
		t.Errorf("Expected BackoffDelay to be %v, got %v", backoffDelayExpected, backoff.BackoffDelay)
	}

	if backoff.Strategy != tree.Exponential {
		t.Errorf("Expected Strategy to be %v, got %v", tree.Exponential, backoff.Strategy)
	}

	since := time.Now()

	backoff.Next()

	elapsed := time.Since(since)

	tolerance := backoffDelayExpected + 50*time.Millisecond
	if elapsed > tolerance {
		t.Errorf("Expected elapsed to be %v, got %v", tolerance, elapsed)
	}
}

func TestNewBackoff_WithBackoffOptionsMaxDelayLessThanBackoffDelay(t *testing.T) {
	overrideDelay := 5 * time.Second
	backoff := tree.NewBackoff(
		tree.WithStrategy(tree.Lineal),
		tree.WithBackoffDelay(overrideDelay),
		tree.WithMaxDelay(1*time.Second),
	)

	if backoff.MaxDelay != overrideDelay {
		t.Errorf("Expected MaxDelay to be %v, got %v", overrideDelay, backoff.MaxDelay)
	}

	since := time.Now()

	backoff.Next()

	elapsed := time.Since(since)

	tolerance := overrideDelay + 50*time.Millisecond
	if elapsed > tolerance {
		t.Errorf("Expected elapsed to be %v, got %v", tolerance, elapsed)
	}
}

func TestNewBackoff_WithBackoffOptionsInvalidStrategy(t *testing.T) {
	backoff := tree.NewBackoff(
		tree.WithStrategy(5),
	)

	if backoff.Strategy != tree.Default {
		t.Errorf("Expected Strategy to be %v, got %v", tree.Default, backoff.Strategy)
	}
}
