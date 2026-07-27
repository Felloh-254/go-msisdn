package msisdn

// countryMeta describes everything go-msisdn needs to know about a single
// country/region in order to parse, validate, format and type-detect
// numbers for it, without depending on an external libphonenumber binding.
//
// This is intentionally a plain data table (see countries.go for the
// values). Adding or correcting a country is a data change here, never a
// change to parser/validator/formatter logic.
type countryMeta struct {
	// ISO is the ISO-3166-1 alpha-2 region code, e.g. "KE".
	ISO string
	// Name is the human-readable country name, e.g. "Kenya".
	Name string
	// CallingCode is the E.164 country calling code, e.g. 254.
	CallingCode int
	// TrunkPrefix is the digit(s) dialled domestically before the
	// national significant number, e.g. "0". Empty if the country has
	// none (e.g. the USA/Canada under the NANP).
	TrunkPrefix string
	// NSNLengths lists the valid lengths (in digits) of the national
	// significant number (i.e. without country code or trunk prefix).
	// If empty, length is not strictly validated (used for countries
	// where go-msisdn only has shallow/general support).
	NSNLengths []int
	// MobilePrefixes are NSN prefixes that indicate a mobile number.
	MobilePrefixes []string
	// FixedLinePrefixes are NSN prefixes that indicate a fixed-line
	// number. May overlap with MobilePrefixes for "fixed line or
	// mobile" ranges, handled via FixedOrMobilePrefixes instead.
	FixedLinePrefixes []string
	// FixedOrMobilePrefixes are NSN prefixes where the range is shared
	// between mobile and fixed-line allocations and cannot be
	// distinguished from the number alone.
	FixedOrMobilePrefixes []string
	// TollFreePrefixes are NSN prefixes for toll-free numbers.
	TollFreePrefixes []string
	// PremiumRatePrefixes are NSN prefixes for premium-rate numbers.
	PremiumRatePrefixes []string
	// VoIPPrefixes are NSN prefixes for VoIP numbers.
	VoIPPrefixes []string
	// PagerPrefixes are NSN prefixes for pager numbers.
	PagerPrefixes []string
	// ExampleMobile is a realistic (but not necessarily allocated) NSN
	// used by Example(), without trunk prefix or country code.
	ExampleMobile string
	// Deep indicates the country has detailed prefix/length metadata
	// (the 7 countries called out in the project brief, plus any the
	// caller registers). Countries with Deep == false still get name,
	// calling code and best-effort E.164 formatting, but IsValid/Type
	// checks are best-effort only.
	Deep bool
}

// countryRegistry indexes countries by ISO code.
var countryRegistry = map[string]*countryMeta{}

// callingCodeIndex indexes ISO codes by calling code. A calling code can
// be shared by multiple regions (e.g. +1 for the US, Canada, and NANP
// members); the first registered region is treated as the primary one.
var callingCodeIndex = map[int][]string{}

func registerCountry(m countryMeta) {
	c := m
	countryRegistry[c.ISO] = &c
	callingCodeIndex[c.CallingCode] = append(callingCodeIndex[c.CallingCode], c.ISO)
}

// lookupByISO returns metadata for an ISO-3166-1 alpha-2 code (case
// insensitive), or nil if unknown.
func lookupByISO(iso string) *countryMeta {
	return countryRegistry[normalizeISO(iso)]
}

// lookupByCallingCode returns the primary metadata for an E.164 calling
// code, or nil if unknown. Longer, more specific calling codes (e.g. some
// Caribbean +1 area codes are out of scope) are matched first.
func lookupByCallingCode(code int) *countryMeta {
	isos, ok := callingCodeIndex[code]
	if !ok || len(isos) == 0 {
		return nil
	}
	return countryRegistry[isos[0]]
}

func normalizeISO(iso string) string {
	if len(iso) != 2 {
		return iso
	}
	b := []byte(iso)
	if b[0] >= 'a' && b[0] <= 'z' {
		b[0] -= 32
	}
	if b[1] >= 'a' && b[1] <= 'z' {
		b[1] -= 32
	}
	return string(b)
}

// SupportedCountries returns the ISO codes of every registered country,
// including both "deep" (fully validated) and "shallow" (name/calling
// code only) entries.
func SupportedCountries() []string {
	out := make([]string, 0, len(countryRegistry))
	for iso := range countryRegistry {
		out = append(out, iso)
	}
	return out
}
