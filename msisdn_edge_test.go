package msisdn

import "testing"

func TestShallowCountryParsing(t *testing.T) {
	// US/CA use the NANP: no trunk prefix, 10-digit NSN.
	p, err := Parse("2015550123", "US")
	if err != nil {
		t.Fatal(err)
	}
	if !p.IsValid() {
		t.Errorf("expected valid US number, reason=%s", p.InvalidReason())
	}
	if p.E164() != "+12015550123" {
		t.Errorf("E164() = %s", p.E164())
	}
	if p.Country() != "United States" || p.CountryCode() != 1 {
		t.Errorf("got country=%s code=%d", p.Country(), p.CountryCode())
	}
	// Shallow countries have no prefix/type tables, so Type() is Unknown
	// even for a structurally valid number.
	if p.Type() != Unknown {
		t.Errorf("expected Unknown type for shallow country, got %v", p.Type())
	}

	// UK with its trunk-zero convention.
	uk, err := Parse("07911123456", "GB")
	if err != nil || !uk.IsValid() {
		t.Fatalf("UK parse failed: %v valid=%v reason=%s", err, uk.IsValid(), uk.InvalidReason())
	}
	if uk.E164() != "+447911123456" {
		t.Errorf("E164() = %s", uk.E164())
	}
}

func TestUnsupportedRegion(t *testing.T) {
	_, err := Parse("0712345678", "ZZ")
	if err == nil {
		t.Fatal("expected an error for an unsupported region")
	}
}

func TestCaseInsensitiveRegion(t *testing.T) {
	p, err := Parse("0712345678", "ke")
	if err != nil || !p.IsValid() {
		t.Fatalf("lowercase region should work: %v", err)
	}
	if p.ISO() != "KE" {
		t.Errorf("ISO() = %s, want KE", p.ISO())
	}
}

func TestInvalidPrefixDetection(t *testing.T) {
	// 9 digits (correct length for Kenya) but starting with a digit that
	// isn't in any registered mobile/fixed/toll-free/premium range.
	p, err := Parse("999999999", "KE")
	if err != nil {
		t.Fatal(err)
	}
	if p.IsValid() {
		t.Error("expected an unrecognized prefix to be invalid")
	}
	if !p.IsPossible() {
		t.Error("a correct-length number should still be Possible")
	}
}

func TestImpossibleNumber(t *testing.T) {
	// Wildly wrong length: not just invalid, but implausible.
	p, err := Parse("11", "KE")
	if err != nil {
		t.Fatal(err)
	}
	if p.IsPossible() {
		t.Error("a 2-digit 'Kenyan number' should not be Possible")
	}
	if p.IsValid() {
		t.Error("a 2-digit 'Kenyan number' should not be Valid")
	}
}

func TestPhoneStringAndEqualNilSafety(t *testing.T) {
	var nilPhone *Phone
	if nilPhone.String() != "" {
		t.Error("nil Phone.String() should be empty, not panic")
	}
	if nilPhone.IsValid() {
		t.Error("nil Phone.IsValid() should be false")
	}
	if nilPhone.Country() != "" {
		t.Error("nil Phone.Country() should be empty")
	}

	p, _ := Parse("0712345678", "KE")
	if p.Equal(nil) {
		t.Error("a valid phone should not equal nil")
	}
}

func TestSupportedCountriesIncludesRequiredSet(t *testing.T) {
	required := []string{"KE", "UG", "TZ", "RW", "NG", "GH", "ZM"}
	got := map[string]bool{}
	for _, iso := range SupportedCountries() {
		got[iso] = true
	}
	for _, iso := range required {
		if !got[iso] {
			t.Errorf("expected %s to be a supported country", iso)
		}
	}
}

func TestToLocal(t *testing.T) {
	local, err := ToLocal("+254712345678", "")
	if err != nil {
		t.Fatal(err)
	}
	if local != "0712345678" {
		t.Errorf("ToLocal(...) = %s, want 0712345678", local)
	}
}

func TestMultipleQuotesGuardAgainstOffByOneSpacing(t *testing.T) {
	p, err := Parse("0803123 4567", "NG")
	if err != nil || !p.IsValid() {
		t.Fatalf("Nigeria number with stray space failed: %v valid=%v reason=%s", err, p.IsValid(), p.InvalidReason())
	}
	if p.International() != "+234 803 123 4567" {
		t.Errorf("International() = %s", p.International())
	}
}
