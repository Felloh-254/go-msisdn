// Package operators provides a configurable, data-driven mobile network
// operator (MNO) lookup keyed by ISO-3166 country code and national
// significant number prefix.
//
// It intentionally contains no per-country conditional logic. Each
// supported country registers a table of prefix->operator rules via
// RegisterCountry (see kenya.go, uganda.go, etc.), and Lookup performs a
// longest-prefix match against that table. Adding a new country is a data
// change, not a code change.
package operators

import "sort"

// Rule maps a set of national-number prefixes to an operator name.
type Rule struct {
	// Operator is the display name of the network operator, e.g. "Safaricom".
	Operator string
	// Prefixes are digit strings matched against the start of the
	// national significant number (without country code or trunk
	// prefix), e.g. "70", "701".
	Prefixes []string
}

var registry = map[string][]Rule{}

// RegisterCountry registers (or replaces) the operator prefix table for an
// ISO-3166-1 alpha-2 country code. It is safe to call from package init()
// functions and is typically used that way by the built-in country files
// in this package, but is exported so consumers can register additional
// countries or override the defaults at runtime.
func RegisterCountry(iso string, rules []Rule) {
	registry[normalizeISO(iso)] = rules
}

// SupportedCountries returns the ISO codes for which operator data has
// been registered, sorted alphabetically.
func SupportedCountries() []string {
	out := make([]string, 0, len(registry))
	for iso := range registry {
		out = append(out, iso)
	}
	sort.Strings(out)
	return out
}

// Lookup returns the operator name for a national significant number
// (digits only, no country calling code, no trunk prefix) in the given
// ISO country. It performs a longest-prefix match so that more specific
// rules (e.g. "701") take precedence over shorter ones (e.g. "70").
//
// The second return value is false if no operator could be determined,
// either because the country has no registered data or no prefix matched.
func Lookup(iso, nationalNumber string) (string, bool) {
	rules, ok := registry[normalizeISO(iso)]
	if !ok {
		return "", false
	}

	best := ""
	bestLen := -1
	for _, rule := range rules {
		for _, prefix := range rule.Prefixes {
			if len(prefix) > bestLen && hasPrefix(nationalNumber, prefix) {
				best = rule.Operator
				bestLen = len(prefix)
			}
		}
	}
	if bestLen == -1 {
		return "", false
	}
	return best, true
}

func hasPrefix(s, prefix string) bool {
	return len(s) >= len(prefix) && s[:len(prefix)] == prefix
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
