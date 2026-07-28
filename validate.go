package msisdn

import (
	"fmt"

	msisdnerrors "github.com/Felloh-254/go-msisdn/errors"
)

// ValidationResult is the structured outcome of Validate: it tells you
// not just whether a number is valid, but why not.
type ValidationResult struct {
	// Valid is true if the number is fully valid.
	Valid bool `json:"valid"`
	// Possible is true if the number is at least plausible (right
	// ballpark length), even if not fully Valid.
	Possible bool `json:"possible"`
	// Reason is a human-readable explanation, empty when Valid is true.
	Reason string `json:"reason,omitempty"`
	// Code is a stable, machine-readable reason code, empty when Valid
	// is true. See the errors package for possible values.
	Code msisdnerrors.Code `json:"code,omitempty"`
	// Phone is the parsed number, or nil if parsing itself failed
	// (e.g. empty input, unrecognized calling code, missing region).
	Phone *Phone `json:"-"`
}

// Validate parses number (see Parse for how region is used) and reports
// detailed validation information. Unlike Parse, Validate never returns a
// Go error -- structural parse failures are reported as an invalid
// ValidationResult instead, so this is the simplest entry point for
// "is this number OK, and if not, why?" checks such as form validation.
func Validate(number, region string) ValidationResult {
	p, err := Parse(number, region)
	if err != nil {
		var ve *msisdnerrors.ValidationError
		if as(err, &ve) {
			return ValidationResult{Valid: false, Possible: false, Reason: ve.Reason, Code: ve.Code}
		}
		return ValidationResult{Valid: false, Possible: false, Reason: err.Error()}
	}
	return ValidationResult{
		Valid:    p.valid,
		Possible: p.possible,
		Reason:   p.reason,
		Code:     p.reasonCode,
		Phone:    p,
	}
}

// as is a tiny local wrapper around errors.As so this file only needs one
// stdlib import site; kept unexported since it's only used here.
func as(err error, target **msisdnerrors.ValidationError) bool {
	ve, ok := err.(*msisdnerrors.ValidationError)
	if ok {
		*target = ve
	}
	return ok
}

// validateNational determines whether a national significant number is
// valid/possible for meta, returning a reason code and message when not.
func validateNational(meta *countryMeta, national string) (valid, possible bool, code msisdnerrors.Code, reason string) {
	if !allDigits(national) {
		return false, false, msisdnerrors.CodeInvalidFormat, "national number must contain only digits"
	}

	if len(meta.NSNLengths) > 0 && !containsInt(meta.NSNLengths, len(national)) {
		possible := isLengthPlausible(meta.NSNLengths, len(national))
		code := msisdnerrors.CodeInvalidLength
		if !possible {
			code = msisdnerrors.CodeImpossible
		}
		return false, possible, code, fmt.Sprintf(
			"invalid length for %s: got %d digits, expected %v",
			meta.Name, len(national), meta.NSNLengths,
		)
	}

	if meta.Deep && !matchesAnyKnownRange(meta, national) {
		return false, true, msisdnerrors.CodeInvalidPrefix, fmt.Sprintf(
			"%q does not match any known number range for %s", national, meta.Name,
		)
	}

	return true, true, "", ""
}

// isLengthPlausible reports whether got is close enough to any allowed
// length to be merely "wrong" rather than outright "impossible".
func isLengthPlausible(allowed []int, got int) bool {
	for _, want := range allowed {
		diff := want - got
		if diff < 0 {
			diff = -diff
		}
		if diff <= 2 {
			return true
		}
	}
	return false
}

// classify returns the NumberType for national within meta, using a
// longest-prefix-match across all registered ranges (mobile, fixed line,
// toll free, premium rate, VoIP, pager). meta.Deep is assumed true.
func classify(meta *countryMeta, national string) NumberType {
	type candidate struct {
		t      NumberType
		prefix string
	}
	var best *candidate
	consider := func(t NumberType, prefixes []string) {
		for _, pre := range prefixes {
			if hasPrefix(national, pre) && (best == nil || len(pre) > len(best.prefix)) {
				best = &candidate{t: t, prefix: pre}
			}
		}
	}
	consider(Mobile, meta.MobilePrefixes)
	consider(FixedLine, meta.FixedLinePrefixes)
	consider(FixedLineOrMobile, meta.FixedOrMobilePrefixes)
	consider(TollFree, meta.TollFreePrefixes)
	consider(PremiumRate, meta.PremiumRatePrefixes)
	consider(VoIP, meta.VoIPPrefixes)
	consider(Pager, meta.PagerPrefixes)

	if best == nil {
		return Unknown
	}
	return best.t
}

func matchesAnyKnownRange(meta *countryMeta, national string) bool {
	return classify(meta, national) != Unknown
}

func hasPrefix(s, prefix string) bool {
	return len(s) >= len(prefix) && s[:len(prefix)] == prefix
}
