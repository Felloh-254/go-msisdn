package msisdn

import "strconv"

// Normalize parses number (using region as the default country for
// numbers without a leading "+"/"00") and returns it as a bare digit
// string of calling-code + national number, with no "+", spaces, or other
// punctuation -- e.g. "254712345678". This is the canonical MSISDN form
// used as a storage/lookup key throughout telecom and fintech systems.
func Normalize(number, region string) (string, error) {
	p, err := Parse(number, region)
	if err != nil {
		return "", err
	}
	return strconv.Itoa(p.meta.CallingCode) + p.national, nil
}

// Clean strips every character that is not an ASCII digit, including any
// leading "+". It does not parse or validate the result -- it's a pure
// text-cleaning utility, useful as a first pass before Parse, or for
// sanitizing free-text input.
func Clean(number string) string {
	return filterDigits(number)
}

// Equal reports whether two phone number strings refer to the same
// number, regardless of how each is formatted (local, E.164, spaced,
// punctuated, ...). An optional defaultRegion is used for either input
// that doesn't carry an explicit "+"/"00" country code, exactly as with
// Parse; if omitted, such inputs are compared as invalid (since they
// can't be unambiguously resolved to a country).
func Equal(number1, number2 string, defaultRegion ...string) bool {
	region := ""
	if len(defaultRegion) > 0 {
		region = defaultRegion[0]
	}
	p1, err1 := Parse(number1, region)
	if err1 != nil {
		return false
	}
	p2, err2 := Parse(number2, region)
	if err2 != nil {
		return false
	}
	return p1.Equal(p2)
}

// ToLocal converts a number back into the domestic/local dialling form
// (trunk prefix + national number, no country code), e.g.
// "254712345678" -> "0712345678". It's the inverse of prefixing a local
// number with a country code.
func ToLocal(number, region string) (string, error) {
	p, err := Parse(number, region)
	if err != nil {
		return "", err
	}
	return p.Local(), nil
}

// Dedupe normalizes every number in numbers (using defaultRegion for any
// that lack an explicit country code) and returns the unique E.164 forms,
// preserving the order of first appearance. Numbers that fail to parse
// are silently skipped -- use ValidateMany first if you need to know
// which inputs were dropped and why.
func Dedupe(numbers []string, defaultRegion string) []string {
	seen := make(map[string]bool, len(numbers))
	out := make([]string, 0, len(numbers))
	for _, raw := range numbers {
		p, err := Parse(raw, defaultRegion)
		if err != nil {
			continue
		}
		e164 := p.E164()
		if !seen[e164] {
			seen[e164] = true
			out = append(out, e164)
		}
	}
	return out
}
