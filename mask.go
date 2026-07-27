package msisdn

// maskConfig holds the tunable parameters for Mask.
type maskConfig struct {
	prefixVisible int
	suffixVisible int
	maskChar      rune
}

// MaskOption customizes Mask's behavior. See WithPrefixVisible,
// WithSuffixVisible, and WithMaskChar.
type MaskOption func(*maskConfig)

// WithPrefixVisible sets how many leading digits stay visible. Default 4.
func WithPrefixVisible(n int) MaskOption {
	return func(c *maskConfig) { c.prefixVisible = n }
}

// WithSuffixVisible sets how many trailing digits stay visible. Default 2.
func WithSuffixVisible(n int) MaskOption {
	return func(c *maskConfig) { c.suffixVisible = n }
}

// WithMaskChar sets the character used to replace hidden digits. Default '*'.
func WithMaskChar(r rune) MaskOption {
	return func(c *maskConfig) { c.maskChar = r }
}

// Mask redacts the middle of a phone number for privacy-friendly logging,
// keeping a configurable number of leading and trailing digits visible.
// It operates on digits only (any "+", spaces, or punctuation in number
// are stripped first, matching Clean) and does not require the number to
// be valid or even parseable -- it's a text transform, not a parser.
//
//	msisdn.Mask("254712345678")                          // "2547******78"
//	msisdn.Mask("254712345678", msisdn.WithPrefixVisible(6)) // "254712****78"
func Mask(number string, opts ...MaskOption) string {
	digits := filterDigits(number)
	cfg := maskConfig{prefixVisible: 4, suffixVisible: 2, maskChar: '*'}
	for _, opt := range opts {
		opt(&cfg)
	}

	n := len(digits)
	prefix, suffix := cfg.prefixVisible, cfg.suffixVisible
	if prefix < 0 {
		prefix = 0
	}
	if suffix < 0 {
		suffix = 0
	}
	if prefix+suffix >= n {
		// Not enough digits to safely reveal both ends without
		// exposing the whole number -- mask everything.
		out := make([]rune, n)
		for i := range out {
			out[i] = cfg.maskChar
		}
		return string(out)
	}

	hidden := n - prefix - suffix
	masked := make([]rune, hidden)
	for i := range masked {
		masked[i] = cfg.maskChar
	}
	return digits[:prefix] + string(masked) + digits[n-suffix:]
}

// Mask redacts the middle of the parsed number's E.164 digits. See the
// package-level Mask function for option documentation.
func (p *Phone) Mask(opts ...MaskOption) string {
	if p == nil {
		return ""
	}
	return Mask(itoa(p.meta.CallingCode)+p.national, opts...)
}
