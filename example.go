package msisdn

import (
	"fmt"

	msisdnerrors "github.com/yourusername/go-msisdn/errors"
)

// Example returns a realistic sample E.164 phone number for the given
// ISO-3166-1 alpha-2 country, suitable for tests, fixtures, and demo data.
// It returns an error if the country isn't registered or has no example
// number configured.
//
//	n, _ := msisdn.Example("KE") // "+254712345678"
func Example(iso string) (string, error) {
	meta := lookupByISO(iso)
	if meta == nil {
		return "", msisdnerrors.New(msisdnerrors.CodeUnknownRegion, msisdnerrors.ErrCountryNotSupported,
			fmt.Sprintf("unsupported region %q", iso))
	}
	if meta.ExampleMobile == "" {
		return "", msisdnerrors.New(msisdnerrors.CodeUnknownRegion, msisdnerrors.ErrCountryNotSupported,
			fmt.Sprintf("no example number configured for %q", iso))
	}
	return "+" + itoa(meta.CallingCode) + meta.ExampleMobile, nil
}

// ExamplePhone is like Example but returns a parsed *Phone.
func ExamplePhone(iso string) (*Phone, error) {
	e164, err := Example(iso)
	if err != nil {
		return nil, err
	}
	return Parse(e164, "")
}
