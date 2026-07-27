package msisdn

import (
	"fmt"
	"strconv"
	"strings"

	msisdnerrors "github.com/yourusername/go-msisdn/errors"
)

// Parse parses raw into a Phone.
//
// If raw begins with "+" or "00" it is treated as already carrying an
// explicit country calling code and region is ignored. Otherwise region
// must be a supported ISO-3166-1 alpha-2 code (e.g. "KE") and raw is
// interpreted as a national/local number, with any domestic trunk prefix
// (e.g. a leading "0") stripped automatically.
//
// Parse returns an error only for structural problems: an empty input, an
// unrecognized calling code, a missing/unknown region for a non-"+"
// number, or input that doesn't resemble a phone number at all. A number
// that parses structurally but fails validation (wrong length, unknown
// prefix range) is still returned, with Phone.IsValid() reporting false
// and Phone.InvalidReason() explaining why -- callers that want a hard
// error for invalid-but-parseable numbers should check IsValid()
// themselves, or use Validate.
func Parse(raw, region string) (*Phone, error) {
	trimmed := strings.TrimSpace(raw)
	if trimmed == "" {
		return nil, msisdnerrors.New(msisdnerrors.CodeEmpty, msisdnerrors.ErrEmptyNumber, "phone number is empty")
	}

	international := false
	rest := trimmed
	switch {
	case strings.HasPrefix(rest, "+"):
		international = true
		rest = rest[1:]
	case strings.HasPrefix(rest, "00"):
		international = true
		rest = rest[2:]
	}

	digits := filterDigits(rest)
	if digits == "" {
		return nil, msisdnerrors.New(msisdnerrors.CodeInvalidFormat, msisdnerrors.ErrInvalidFormat,
			fmt.Sprintf("could not find any digits in %q", raw))
	}

	var meta *countryMeta
	var national string

	if international {
		m, n, err := splitCallingCode(digits)
		if err != nil {
			return nil, err
		}
		meta, national = m, n
		national = stripLikelyTrunkPrefix(meta, national)
	} else {
		if region == "" {
			return nil, msisdnerrors.New(msisdnerrors.CodeUnknownRegion, msisdnerrors.ErrUnknownRegion,
				"a region is required when the number has no leading '+' or '00'")
		}
		meta = lookupByISO(region)
		if meta == nil {
			return nil, msisdnerrors.New(msisdnerrors.CodeUnknownRegion, msisdnerrors.ErrCountryNotSupported,
				fmt.Sprintf("unsupported region %q", region))
		}
		national = digits
		ccStr := strconv.Itoa(meta.CallingCode)
		switch {
		case strings.HasPrefix(national, ccStr) && !containsInt(meta.NSNLengths, len(national)) &&
			containsInt(meta.NSNLengths, len(national)-len(ccStr)):
			// The input already carries the country calling code but
			// dropped the "+" (e.g. "254712345678" for KE) -- prefer
			// that reading over treating the whole thing as an
			// oversized national number.
			national = national[len(ccStr):]
		case meta.TrunkPrefix != "" && strings.HasPrefix(national, meta.TrunkPrefix):
			national = national[len(meta.TrunkPrefix):]
		}
	}

	p := &Phone{meta: meta, national: national, raw: raw}
	p.valid, p.possible, p.reasonCode, p.reason = validateNational(meta, national)
	return p, nil
}

// splitCallingCode tries to split a run of digits (with the leading "+"
// or "00" already removed) into a known country calling code and the
// remaining national number. It tries 3, then 2, then 1 digit prefixes,
// which is sufficient because go-msisdn only registers a curated,
// non-ambiguous set of calling codes (see countries.go).
func splitCallingCode(digits string) (*countryMeta, string, error) {
	for _, l := range []int{3, 2, 1} {
		if len(digits) <= l {
			continue
		}
		code, err := strconv.Atoi(digits[:l])
		if err != nil {
			continue
		}
		if meta := lookupByCallingCode(code); meta != nil {
			return meta, digits[l:], nil
		}
	}
	return nil, "", msisdnerrors.New(msisdnerrors.CodeInvalidCountryCode, msisdnerrors.ErrInvalidCountryCode,
		fmt.Sprintf("no known country calling code found at the start of %q", digits))
}

// stripLikelyTrunkPrefix handles the common data-entry mistake of typing
// a domestic trunk prefix right after the country code, e.g.
// "+2540712345678" instead of "+254712345678". It only strips when doing
// so produces a length that matches one of the country's valid NSN
// lengths and the un-stripped version does not already match.
func stripLikelyTrunkPrefix(meta *countryMeta, national string) string {
	if meta.TrunkPrefix == "" || len(meta.NSNLengths) == 0 {
		return national
	}
	if containsInt(meta.NSNLengths, len(national)) {
		return national // already a valid length, leave it alone
	}
	if strings.HasPrefix(national, meta.TrunkPrefix) {
		stripped := national[len(meta.TrunkPrefix):]
		if containsInt(meta.NSNLengths, len(stripped)) {
			return stripped
		}
	}
	return national
}

// filterDigits strips every character that is not an ASCII digit.
func filterDigits(s string) string {
	var b strings.Builder
	b.Grow(len(s))
	for _, r := range s {
		if r >= '0' && r <= '9' {
			b.WriteRune(r)
		}
	}
	return b.String()
}

func allDigits(s string) bool {
	if s == "" {
		return false
	}
	for _, r := range s {
		if r < '0' || r > '9' {
			return false
		}
	}
	return true
}

func containsInt(list []int, v int) bool {
	for _, x := range list {
		if x == v {
			return true
		}
	}
	return false
}
