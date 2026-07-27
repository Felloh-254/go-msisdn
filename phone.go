// Package msisdn is a developer-friendly phone number toolkit for backend
// systems, fintech, CRM, telecom, and iGaming platforms: parsing,
// validation, normalization, formatting, type/operator detection, masking,
// and batch helpers, built around a single Phone value type.
package msisdn

import (
	"strconv"

	msisdnerrors "github.com/yourusername/go-msisdn/errors"
	"github.com/yourusername/go-msisdn/operators"
)

// Phone represents a parsed phone number together with everything
// go-msisdn was able to determine about it: its country, national
// significant number, validity, type, and (where available) operator.
//
// Phone is immutable and safe for concurrent use. The zero value is not
// usable; construct a Phone via Parse, ParseMany, or by decoding JSON /
// scanning from a database column.
type Phone struct {
	meta       *countryMeta
	national   string // digits only, no trunk prefix, no country code
	raw        string
	valid      bool
	possible   bool
	reasonCode msisdnerrors.Code
	reason     string
}

// Country returns the country's display name, e.g. "Kenya". Returns ""
// for a zero-value Phone.
func (p *Phone) Country() string {
	if p == nil || p.meta == nil {
		return ""
	}
	return p.meta.Name
}

// ISO returns the ISO-3166-1 alpha-2 region code, e.g. "KE".
func (p *Phone) ISO() string {
	if p == nil || p.meta == nil {
		return ""
	}
	return p.meta.ISO
}

// CountryCode returns the E.164 country calling code, e.g. 254.
func (p *Phone) CountryCode() int {
	if p == nil || p.meta == nil {
		return 0
	}
	return p.meta.CallingCode
}

// NationalNumber returns the national significant number (no trunk
// prefix, no country code) as an unsigned integer, e.g. 712345678. Use
// NationalNumberString if you need to preserve leading zeros (rare, but
// possible for some countries' number ranges).
func (p *Phone) NationalNumber() uint64 {
	if p == nil {
		return 0
	}
	n, _ := strconv.ParseUint(p.national, 10, 64)
	return n
}

// NationalNumberString returns the national significant number as a
// digit string, preserving any leading zeros.
func (p *Phone) NationalNumberString() string {
	if p == nil {
		return ""
	}
	return p.national
}

// Raw returns the exact string that was originally passed to Parse.
func (p *Phone) Raw() string {
	if p == nil {
		return ""
	}
	return p.raw
}

// IsValid reports whether the number is fully valid: known country,
// correct length, and (for deeply-supported countries) a recognized
// number-type prefix range.
func (p *Phone) IsValid() bool {
	return p != nil && p.valid
}

// IsPossible reports whether the number could plausibly be dialable --
// i.e. it is not off by an implausible margin in length -- even if it
// isn't fully Valid (for example, correct length but an unrecognized
// prefix range). Every Valid number is also Possible.
func (p *Phone) IsPossible() bool {
	return p != nil && p.possible
}

// InvalidReason returns a human-readable explanation of why the number is
// invalid, or "" if it is valid.
func (p *Phone) InvalidReason() string {
	if p == nil {
		return ""
	}
	return p.reason
}

// Type classifies the number (mobile, fixed line, toll free, ...). It
// returns Unknown for countries go-msisdn only shallowly supports, or if
// the number's prefix doesn't fall into a known range.
func (p *Phone) Type() NumberType {
	if p == nil || p.meta == nil || !p.meta.Deep {
		return Unknown
	}
	return classify(p.meta, p.national)
}

// IsMobile reports whether Type() is Mobile or FixedLineOrMobile.
func (p *Phone) IsMobile() bool {
	t := p.Type()
	return t == Mobile || t == FixedLineOrMobile
}

// IsFixedLine reports whether Type() is FixedLine or FixedLineOrMobile.
func (p *Phone) IsFixedLine() bool {
	t := p.Type()
	return t == FixedLine || t == FixedLineOrMobile
}

// IsTollFree reports whether Type() is TollFree.
func (p *Phone) IsTollFree() bool { return p.Type() == TollFree }

// IsPremiumRate reports whether Type() is PremiumRate.
func (p *Phone) IsPremiumRate() bool { return p.Type() == PremiumRate }

// IsVoIP reports whether Type() is VoIP.
func (p *Phone) IsVoIP() bool { return p.Type() == VoIP }

// IsPager reports whether Type() is Pager.
func (p *Phone) IsPager() bool { return p.Type() == Pager }

// Operator returns the detected mobile network operator name, e.g.
// "Safaricom", or "" if unknown (either because the country has no
// operator table registered, or no prefix rule matched).
func (p *Phone) Operator() string {
	if p == nil || p.meta == nil {
		return ""
	}
	name, _ := operators.Lookup(p.meta.ISO, p.national)
	return name
}

// String implements fmt.Stringer, returning the E.164 form.
func (p *Phone) String() string {
	if p == nil || p.meta == nil {
		return ""
	}
	return p.E164()
}

// Equal reports whether two Phone values refer to the same number,
// compared by E.164 form.
func (p *Phone) Equal(other *Phone) bool {
	if p == nil || other == nil {
		return p == other
	}
	return p.E164() == other.E164()
}
