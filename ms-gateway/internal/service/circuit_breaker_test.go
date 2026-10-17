package service

import (
	"testing"
	"time"
)

func TestCircuitBreakerOpensAfterFailureThreshold(t *testing.T) {
	breaker := &circuitBreaker{}

	for range circuitFailureThreshold - 1 {
		breaker.failure()
		if !breaker.allow() {
			t.Fatal("breaker opened before failure threshold")
		}
	}

	breaker.failure()
	if breaker.allow() {
		t.Fatal("breaker allowed a request after failure threshold")
	}
}

func TestCircuitBreakerClosedSuccessKeepsFailureHistory(t *testing.T) {
	breaker := &circuitBreaker{}

	breaker.failure()
	breaker.failure()
	breaker.success()
	breaker.failure()

	if breaker.allow() {
		t.Fatal("breaker discarded failures after a successful request in Closed")
	}
}

func TestCircuitBreakerDoesNotOpenForOldFailures(t *testing.T) {
	fallenAt := time.Now()
	breaker := &circuitBreaker{
		failures: [circuitFailureThreshold]time.Time{
			fallenAt.Add(-circuitFailureWindow - 2*time.Millisecond),
			fallenAt.Add(-circuitFailureWindow - time.Millisecond),
			{},
		},
		next:  2,
		count: circuitFailureThreshold - 1,
	}

	breaker.failure()
	if !breaker.openedAt.IsZero() {
		t.Fatal("breaker opened for failures outside the time window")
	}
}

func TestCircuitBreakerClosesAfterSuccessfulProbe(t *testing.T) {
	breaker := &circuitBreaker{
		failures: [circuitFailureThreshold]time.Time{
			time.Now(), time.Now(), time.Now(),
		},
		count:    circuitFailureThreshold,
		openedAt: time.Now().Add(-circuitRetryTimeout),
	}

	if !breaker.allow() {
		t.Fatal("breaker did not allow a probe after retry timeout")
	}
	breaker.success()
	if !breaker.allow() {
		t.Fatal("breaker remained open after successful probe")
	}
	if breaker.count != 0 {
		t.Fatal("breaker did not clear failure history after a successful probe")
	}
}

func TestCircuitBreakerReopensAfterFailedProbe(t *testing.T) {
	breaker := &circuitBreaker{openedAt: time.Now().Add(-circuitRetryTimeout)}

	if !breaker.allow() {
		t.Fatal("breaker did not allow a probe after retry timeout")
	}
	breaker.failure()
	if breaker.allow() {
		t.Fatal("breaker allowed a request after failed probe")
	}
}
