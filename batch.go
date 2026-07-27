package msisdn

// ParseResult pairs a batch input with its parse outcome, so results can
// always be matched back to the original input by index -- important
// because Parse itself can fail per-item without aborting the batch.
type ParseResult struct {
	Input string `json:"input"`
	Phone *Phone `json:"phone,omitempty"`
	Error error  `json:"-"`
	// ErrorMessage mirrors Error as a string for JSON consumers.
	ErrorMessage string `json:"error,omitempty"`
}

// ParseMany parses every number in numbers against defaultRegion, never
// stopping at the first failure. Check each result's Error/ErrorMessage
// to see which inputs failed and why.
func ParseMany(numbers []string, defaultRegion string) []ParseResult {
	results := make([]ParseResult, len(numbers))
	for i, raw := range numbers {
		p, err := Parse(raw, defaultRegion)
		r := ParseResult{Input: raw, Phone: p, Error: err}
		if err != nil {
			r.ErrorMessage = err.Error()
		}
		results[i] = r
	}
	return results
}

// ValidateMany runs Validate over every number in numbers against
// defaultRegion.
func ValidateMany(numbers []string, defaultRegion string) []ValidationResult {
	results := make([]ValidationResult, len(numbers))
	for i, raw := range numbers {
		results[i] = Validate(raw, defaultRegion)
	}
	return results
}

// NormalizeResult pairs a batch input with its normalized form (or error).
type NormalizeResult struct {
	Input        string `json:"input"`
	Normalized   string `json:"normalized,omitempty"`
	Error        error  `json:"-"`
	ErrorMessage string `json:"error,omitempty"`
}

// NormalizeMany runs Normalize over every number in numbers against
// defaultRegion.
func NormalizeMany(numbers []string, defaultRegion string) []NormalizeResult {
	results := make([]NormalizeResult, len(numbers))
	for i, raw := range numbers {
		n, err := Normalize(raw, defaultRegion)
		r := NormalizeResult{Input: raw, Normalized: n, Error: err}
		if err != nil {
			r.ErrorMessage = err.Error()
		}
		results[i] = r
	}
	return results
}
