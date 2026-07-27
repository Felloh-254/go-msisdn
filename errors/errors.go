// Package errors defines the sentinel errors and structured error codes
// used throughout go-msisdn. Consumers can use errors.Is / errors.As with
// the standard library "errors" package against the values in this file.
package errors

import "errors"

// Code identifies the category of a phone number problem in a way that is
// safe to switch on, log, or expose in an API response (unlike a raw error
// string, which is not guaranteed to be stable).
type Code string

const (
	// CodeEmpty is returned when an empty string was supplied.
	CodeEmpty Code = "EMPTY_NUMBER"
	// CodeUnknownRegion is returned when the supplied region/ISO code is
	// not recognized and the number does not carry an explicit "+".
	CodeUnknownRegion Code = "UNKNOWN_REGION"
	// CodeInvalidCountryCode is returned when the number's calling code
	// does not match any known country.
	CodeInvalidCountryCode Code = "INVALID_COUNTRY_CODE"
	// CodeInvalidLength is returned when the national significant number
	// does not match any valid length for its country.
	CodeInvalidLength Code = "INVALID_LENGTH"
	// CodeInvalidPrefix is returned when the number's leading digits do
	// not correspond to any known number-type range for its country.
	CodeInvalidPrefix Code = "INVALID_PREFIX"
	// CodeInvalidFormat is returned when the raw input contains
	// characters or structure that cannot be interpreted as a phone
	// number at all.
	CodeInvalidFormat Code = "INVALID_FORMAT"
	// CodeImpossible is returned when the number is too short/long to
	// ever be a valid number for the region, regardless of prefix.
	CodeImpossible Code = "IMPOSSIBLE_NUMBER"
)

// Sentinel errors. Use errors.Is(err, msisdnerrors.ErrInvalidLength) etc.
var (
	ErrEmptyNumber         = errors.New("msisdn: empty number")
	ErrUnknownRegion       = errors.New("msisdn: unknown or missing region")
	ErrInvalidCountryCode  = errors.New("msisdn: invalid country code")
	ErrInvalidLength       = errors.New("msisdn: invalid number length")
	ErrInvalidPrefix       = errors.New("msisdn: invalid number prefix")
	ErrInvalidFormat       = errors.New("msisdn: invalid number format")
	ErrImpossibleNumber    = errors.New("msisdn: impossible number")
	ErrCountryNotSupported = errors.New("msisdn: country not supported")
)

// ValidationError is a structured error describing why a phone number is
// not valid. It wraps a sentinel error so callers can use errors.Is, while
// still exposing a machine-readable Code and a human-readable Reason.
type ValidationError struct {
	Code   Code
	Reason string
	err    error
}

func (e *ValidationError) Error() string {
	return e.Reason
}

// Unwrap allows errors.Is(err, ErrInvalidLength) to succeed.
func (e *ValidationError) Unwrap() error {
	return e.err
}

// New builds a ValidationError for the given code/sentinel/reason.
func New(code Code, sentinel error, reason string) *ValidationError {
	return &ValidationError{Code: code, Reason: reason, err: sentinel}
}
