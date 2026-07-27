package msisdn

import (
	"fmt"
	"strings"
)

// Style identifies an output representation for a phone number.
//
// Note on naming: the initial design brief called this type "Format" and
// asked for a top-level Format(number, STYLE) function. Go does not allow
// a type and a function to share an identifier in the same package, so
// the type is named Style here and the top-level function keeps the name
// Format -- exactly matching the requested call shape, msisdn.Format(n,
// msisdn.National), while staying compilable.
type Style int

const (
	// E164 is the "+254712345678" form: a plus sign, calling code, and
	// national significant number, with no other characters.
	E164 Style = iota
	// National is the domestic dialling form, e.g. "0712 345678".
	National
	// International is the E.164 digits with human-friendly spacing and
	// a leading "+", e.g. "+254 712 345678".
	International
	// RFC3966 is the "tel:+254712345678" URI form.
	RFC3966
)

func (s Style) String() string {
	switch s {
	case E164:
		return "E164"
	case National:
		return "NATIONAL"
	case International:
		return "INTERNATIONAL"
	case RFC3966:
		return "RFC3966"
	default:
		return "UNKNOWN"
	}
}

// Format parses number (optionally using region as the default country
// when number has no leading "+") and renders it using the requested
// Style. It's a convenience wrapper around Parse + Phone.Format for
// callers who don't need to keep the parsed Phone around.
//
//	msisdn.Format("0712345678", "KE", msisdn.E164) // "+254712345678"
func Format(number, region string, style Style) (string, error) {
	p, err := Parse(number, region)
	if err != nil {
		return "", err
	}
	return p.Format(style), nil
}

// Format renders the phone number using the requested Style.
func (p *Phone) Format(style Style) string {
	switch style {
	case National:
		return p.National()
	case International:
		return p.International()
	case RFC3966:
		return p.RFC3966()
	default:
		return p.E164()
	}
}

// E164 renders the number as "+<callingcode><nationalnumber>", e.g.
// "+254712345678". This is the canonical, comparison-safe form.
func (p *Phone) E164() string {
	return "+" + itoa(p.meta.CallingCode) + p.national
}

// International renders the number as "+<callingcode> <spaced national>",
// e.g. "+254 712 345678".
func (p *Phone) International() string {
	return "+" + itoa(p.meta.CallingCode) + " " + spaceNational(p.meta.ISO, p.national)
}

// National renders the number in domestic dialling form, e.g.
// "0712 345678". If the country has no trunk prefix (e.g. NANP countries)
// this is the same as the spaced national significant number.
func (p *Phone) National() string {
	return p.meta.TrunkPrefix + spaceNational(p.meta.ISO, p.national)
}

// RFC3966 renders the number as a "tel:" URI, e.g. "tel:+254712345678".
func (p *Phone) RFC3966() string {
	return "tel:" + p.E164()
}

// Local returns the domestic dialling form without cosmetic spacing, e.g.
// "0712345678". This is what the project brief calls "local conversion":
// turning an international number back into the form a subscriber would
// dial domestically.
func (p *Phone) Local() string {
	return p.meta.TrunkPrefix + p.national
}

// spaceNational applies light, country-aware grouping to a national
// significant number for display purposes. It is deliberately simple: a
// 3-N split for the deep African countries covered by this library, and a
// generic 3s-from-the-left split elsewhere. This is cosmetic only --
// E164() and Equal() never depend on it.
func spaceNational(iso, national string) string {
	switch iso {
	case "KE", "UG", "TZ", "RW", "GH", "ZM":
		if len(national) == 9 {
			return national[:3] + " " + national[3:]
		}
	case "NG":
		if len(national) == 10 {
			return national[:3] + " " + national[3:6] + " " + national[6:]
		}
	}
	// Generic fallback: group in 3s from the left, keep the remainder.
	var b strings.Builder
	for i, r := range national {
		if i != 0 && i%3 == 0 {
			b.WriteByte(' ')
		}
		b.WriteRune(r)
	}
	return b.String()
}

func itoa(n int) string {
	return fmt.Sprintf("%d", n)
}
