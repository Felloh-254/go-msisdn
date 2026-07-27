package msisdn

// NumberType classifies what kind of line a phone number belongs to.
type NumberType int

const (
	// Unknown means go-msisdn could not determine the number type, most
	// often because the country only has shallow support.
	Unknown NumberType = iota
	// Mobile is a mobile/cellular number.
	Mobile
	// FixedLine is a landline number.
	FixedLine
	// FixedLineOrMobile is used when the numbering range is shared
	// between mobile and fixed-line allocations and cannot be told
	// apart from the digits alone.
	FixedLineOrMobile
	// TollFree is a toll-free (freephone) number.
	TollFree
	// PremiumRate is a premium-rate number.
	PremiumRate
	// VoIP is a voice-over-IP number.
	VoIP
	// Pager is a pager number.
	Pager
)

// String implements fmt.Stringer.
func (t NumberType) String() string {
	switch t {
	case Mobile:
		return "MOBILE"
	case FixedLine:
		return "FIXED_LINE"
	case FixedLineOrMobile:
		return "FIXED_LINE_OR_MOBILE"
	case TollFree:
		return "TOLL_FREE"
	case PremiumRate:
		return "PREMIUM_RATE"
	case VoIP:
		return "VOIP"
	case Pager:
		return "PAGER"
	default:
		return "UNKNOWN"
	}
}

// MarshalJSON implements json.Marshaler so NumberType serializes as its
// string name (e.g. "MOBILE") rather than a bare integer.
func (t NumberType) MarshalJSON() ([]byte, error) {
	return []byte(`"` + t.String() + `"`), nil
}
