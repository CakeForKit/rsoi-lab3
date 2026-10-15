package custom_error

import (
	"net/http"
	"strings"
	"testing"
)

func TestNotFoundError(t *testing.T) {
	err := NotFoundError("ticket")

	if !strings.Contains(err.Error(), "ticket not found") {
		t.Fatalf("unexpected error message: %q", err.Error())
	}
	httpError, ok := err.(HttpError)
	if !ok || httpError.Code() != http.StatusNotFound {
		t.Fatalf("expected HTTP 404 error, got %#v", err)
	}
}
