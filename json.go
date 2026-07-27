package msisdn

import (
	"encoding/json"
	"fmt"
)

// MarshalJSON implements json.Marshaler, encoding the Phone as its E.164
// string, e.g. "+254712345678". A zero-value Phone marshals to null.
func (p Phone) MarshalJSON() ([]byte, error) {
	if p.meta == nil {
		return []byte("null"), nil
	}
	return json.Marshal(p.E164())
}

// UnmarshalJSON implements json.Unmarshaler. It accepts a JSON string
// containing an E.164 number (e.g. "+254712345678"); a plain national
// number without a region hint cannot be unambiguously resolved from JSON
// alone and will produce an error. null decodes to the zero Phone.
//
//	type User struct {
//	    Phone msisdn.Phone `json:"phone"`
//	}
func (p *Phone) UnmarshalJSON(data []byte) error {
	if string(data) == "null" {
		*p = Phone{}
		return nil
	}

	var s string
	if err := json.Unmarshal(data, &s); err != nil {
		return fmt.Errorf("msisdn: Phone must be a JSON string: %w", err)
	}
	if s == "" {
		*p = Phone{}
		return nil
	}

	parsed, err := Parse(s, "")
	if err != nil {
		return fmt.Errorf("msisdn: unmarshaling Phone: %w", err)
	}
	*p = *parsed
	return nil
}
