package errors

import (
	stderrors "errors"
	"testing"
)

func TestValidationErrorUnwrap(t *testing.T) {
	ve := New(CodeInvalidLength, ErrInvalidLength, "too short")

	if !stderrors.Is(ve, ErrInvalidLength) {
		t.Error("expected errors.Is to match the wrapped sentinel")
	}
	if stderrors.Is(ve, ErrInvalidPrefix) {
		t.Error("did not expect errors.Is to match an unrelated sentinel")
	}
	if ve.Error() != "too short" {
		t.Errorf("Error() = %q, want %q", ve.Error(), "too short")
	}
	if ve.Code != CodeInvalidLength {
		t.Errorf("Code = %s, want %s", ve.Code, CodeInvalidLength)
	}
}

func TestSentinelErrorsAreDistinct(t *testing.T) {
	sentinels := []error{
		ErrEmptyNumber, ErrUnknownRegion, ErrInvalidCountryCode,
		ErrInvalidLength, ErrInvalidPrefix, ErrInvalidFormat,
		ErrImpossibleNumber, ErrCountryNotSupported,
	}
	for i, a := range sentinels {
		for j, b := range sentinels {
			if i != j && stderrors.Is(a, b) {
				t.Errorf("sentinel %d unexpectedly matches sentinel %d", i, j)
			}
		}
	}
}
