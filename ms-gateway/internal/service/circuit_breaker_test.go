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

func TestCircuitBreakerSuccessResetsFailureCount(t *testing.T) {
	breaker := &circuitBreaker{}

	breaker.failure()
	breaker.failure()
	breaker.success()

	for range circuitFailureThreshold - 1 {
		breaker.failure()
	}
	if !breaker.allow() {
		t.Fatal("breaker retained failures from before a successful response")
	}
}

func TestCircuitBreakerClosesAfterSuccessfulProbe(t *testing.T) {
	breaker := &circuitBreaker{openedAt: time.Now().Add(-circuitRetryTimeout)}

	if !breaker.allow() {
		t.Fatal("breaker did not allow a probe after retry timeout")
	}
	breaker.success()
	if !breaker.allow() {
		t.Fatal("breaker remained open after successful probe")
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
