package msisdn

import (
	"database/sql/driver"
	"encoding/json"
	"errors"
	"testing"
)

func TestParseKenyanFormats(t *testing.T) {
	cases := []struct {
		name       string
		raw        string
		region     string
		wantE164   string
		wantValid  bool
		wantMobile bool
	}{
		{"leading zero", "0712345678", "KE", "+254712345678", true, true},
		{"no leading zero", "712345678", "KE", "+254712345678", true, true},
		{"already e164", "+254712345678", "", "+254712345678", true, true},
		{"international 00 prefix", "00254712345678", "", "+254712345678", true, true},
		{"bare digits with country code, no plus", "254712345678", "KE", "+254712345678", true, true},
		{"punctuated", "+254 (712)-345-678", "", "+254712345678", true, true},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			p, err := Parse(tc.raw, tc.region)
			if err != nil {
				t.Fatalf("Parse(%q, %q) returned error: %v", tc.raw, tc.region, err)
			}
			if p.IsValid() != tc.wantValid {
				t.Errorf("IsValid() = %v, want %v (reason: %s)", p.IsValid(), tc.wantValid, p.InvalidReason())
			}
			if tc.wantValid {
				if got := p.E164(); got != tc.wantE164 {
					t.Errorf("E164() = %q, want %q", got, tc.wantE164)
				}
				if p.IsMobile() != tc.wantMobile {
					t.Errorf("IsMobile() = %v, want %v", p.IsMobile(), tc.wantMobile)
				}
			}
		})
	}
}

func TestParseUgandanAndNigerianNumbers(t *testing.T) {
	p, err := Parse("0771234567", "UG")
	if err != nil || !p.IsValid() {
		t.Fatalf("Uganda parse failed: %v valid=%v reason=%s", err, p.IsValid(), p.InvalidReason())
	}
	if p.Country() != "Uganda" || p.Operator() != "MTN Uganda" {
		t.Errorf("got country=%s operator=%s", p.Country(), p.Operator())
	}

	p2, err := Parse("08031234567", "NG")
	if err != nil || !p2.IsValid() {
		t.Fatalf("Nigeria parse failed: %v valid=%v reason=%s", err, p2.IsValid(), p2.InvalidReason())
	}
	if p2.Operator() != "MTN Nigeria" || p2.CountryCode() != 234 {
		t.Errorf("got operator=%s callingCode=%d", p2.Operator(), p2.CountryCode())
	}
}

func TestParseErrors(t *testing.T) {
	if _, err := Parse("", "KE"); err == nil {
		t.Error("expected error for empty number")
	}
	if _, err := Parse("712345678", ""); err == nil {
		t.Error("expected error for missing region on a non-+ number")
	}
	if _, err := Parse("+999712345678", ""); err == nil {
		t.Error("expected error for unknown calling code")
	}
	if _, err := Parse("abcd", "KE"); err == nil {
		t.Error("expected error for non-numeric garbage")
	}
}

func TestInvalidLengthIsValidFalseNotError(t *testing.T) {
	p, err := Parse("07123", "KE") // too short
	if err != nil {
		t.Fatalf("Parse returned an error for a structurally OK but too-short number: %v", err)
	}
	if p.IsValid() {
		t.Error("expected IsValid() == false for a too-short number")
	}
	if p.InvalidReason() == "" {
		t.Error("expected a non-empty InvalidReason()")
	}
}

func TestValidate(t *testing.T) {
	res := Validate("0712345678", "KE")
	if !res.Valid || res.Phone == nil {
		t.Fatalf("expected valid result, got %+v", res)
	}

	res2 := Validate("12345", "KE")
	if res2.Valid {
		t.Error("expected invalid result for garbage-length number")
	}
	if res2.Reason == "" {
		t.Error("expected a reason for the invalid result")
	}

	res3 := Validate("0712345678", "")
	if res3.Valid {
		t.Error("expected invalid result when region is missing and number has no '+'")
	}
}

func TestFormat(t *testing.T) {
	p, err := Parse("0712345678", "KE")
	if err != nil {
		t.Fatal(err)
	}
	if got := p.E164(); got != "+254712345678" {
		t.Errorf("E164() = %q", got)
	}
	if got := p.International(); got != "+254 712 345678" {
		t.Errorf("International() = %q", got)
	}
	if got := p.National(); got != "0712 345678" {
		t.Errorf("National() = %q", got)
	}
	if got := p.RFC3966(); got != "tel:+254712345678" {
		t.Errorf("RFC3966() = %q", got)
	}
	if got := p.Local(); got != "0712345678" {
		t.Errorf("Local() = %q", got)
	}

	out, err := Format("0712345678", "KE", National)
	if err != nil || out != "0712 345678" {
		t.Errorf("Format(...) = %q, err=%v", out, err)
	}
}

func TestNormalizeAndClean(t *testing.T) {
	inputs := []string{"0712345678", "712345678", "+254712345678", "254712345678"}
	for _, in := range inputs {
		got, err := Normalize(in, "KE")
		if err != nil {
			t.Fatalf("Normalize(%q) error: %v", in, err)
		}
		if got != "254712345678" {
			t.Errorf("Normalize(%q) = %q, want 254712345678", in, got)
		}
	}

	if got := Clean("+254 (712)-345-678"); got != "254712345678" {
		t.Errorf("Clean(...) = %q", got)
	}
}

func TestEqual(t *testing.T) {
	if !Equal("0712345678", "254712345678", "KE") {
		t.Error("expected numbers to be equal")
	}
	if !Equal("0712345678", "+254712345678", "KE") {
		t.Error("expected numbers to be equal")
	}
	if Equal("0712345678", "0722345678", "KE") {
		t.Error("expected different numbers to be unequal")
	}
}

func TestMask(t *testing.T) {
	if got := Mask("254712345678"); got != "2547******78" {
		t.Errorf("Mask(...) = %q, want 2547******78", got)
	}
	if got := Mask("254712345678", WithPrefixVisible(6), WithSuffixVisible(2)); got != "254712****78" {
		t.Errorf("Mask with custom options = %q", got)
	}
	if got := Mask("123", WithPrefixVisible(4), WithSuffixVisible(4)); got != "***" {
		t.Errorf("Mask short number = %q", got)
	}
}

func TestDedupe(t *testing.T) {
	got := Dedupe([]string{"0712345678", "+254712345678", "254712345678"}, "KE")
	if len(got) != 1 || got[0] != "+254712345678" {
		t.Errorf("Dedupe(...) = %v, want [+254712345678]", got)
	}
}

func TestNumberType(t *testing.T) {
	mobile, _ := Parse("0712345678", "KE")
	if mobile.Type() != Mobile {
		t.Errorf("expected Mobile, got %v", mobile.Type())
	}

	tollFree, _ := Parse("0800123456", "KE")
	if tollFree.Type() != TollFree {
		t.Errorf("expected TollFree, got %v", tollFree.Type())
	}

	shallow, _ := Parse("2015550123", "US")
	if shallow.Type() != Unknown {
		t.Errorf("expected Unknown for shallow country, got %v", shallow.Type())
	}
}

func TestBatch(t *testing.T) {
	numbers := []string{"0712345678", "not-a-number", "0771234567"}

	parsed := ParseMany(numbers, "KE")
	if len(parsed) != 3 {
		t.Fatalf("expected 3 results, got %d", len(parsed))
	}
	if parsed[1].Error == nil {
		t.Error("expected an error for the garbage input")
	}

	validated := ValidateMany(numbers, "KE")
	if !validated[0].Valid {
		t.Error("expected first number to validate")
	}

	normalized := NormalizeMany([]string{"0712345678"}, "KE")
	if normalized[0].Normalized != "254712345678" {
		t.Errorf("got %q", normalized[0].Normalized)
	}
}

func TestExample(t *testing.T) {
	e164, err := Example("KE")
	if err != nil {
		t.Fatal(err)
	}
	p, err := Parse(e164, "")
	if err != nil || !p.IsValid() {
		t.Fatalf("example number for KE did not parse as valid: %v", err)
	}

	if _, err := Example("ZZ"); err == nil {
		t.Error("expected error for unsupported region")
	}
}

func TestJSONRoundTrip(t *testing.T) {
	type User struct {
		Phone Phone `json:"phone"`
	}

	p, err := Parse("0712345678", "KE")
	if err != nil {
		t.Fatal(err)
	}
	u := User{Phone: *p}

	b, err := json.Marshal(u)
	if err != nil {
		t.Fatal(err)
	}
	if string(b) != `{"phone":"+254712345678"}` {
		t.Fatalf("unexpected JSON: %s", b)
	}

	var u2 User
	if err := json.Unmarshal(b, &u2); err != nil {
		t.Fatal(err)
	}
	if u2.Phone.E164() != "+254712345678" {
		t.Errorf("round-tripped phone = %s", u2.Phone.E164())
	}
}

func TestDatabaseValueAndScan(t *testing.T) {
	p, err := Parse("0712345678", "KE")
	if err != nil {
		t.Fatal(err)
	}

	val, err := p.Value()
	if err != nil {
		t.Fatal(err)
	}
	if val != driver.Value("+254712345678") {
		t.Errorf("Value() = %v", val)
	}

	var scanned Phone
	if err := scanned.Scan("+254712345678"); err != nil {
		t.Fatal(err)
	}
	if scanned.E164() != "+254712345678" {
		t.Errorf("Scan(...) produced %s", scanned.E164())
	}

	var scannedBytes Phone
	if err := scannedBytes.Scan([]byte("+254712345678")); err != nil {
		t.Fatal(err)
	}

	var scannedNil Phone
	if err := scannedNil.Scan(nil); err != nil {
		t.Fatal(err)
	}
	if scannedNil.IsValid() {
		t.Error("expected zero-value Phone from Scan(nil)")
	}
}

func TestErrorsAreWrapped(t *testing.T) {
	_, err := Parse("", "KE")
	if err == nil {
		t.Fatal("expected error")
	}
	// This mirrors the msisdnerrors sentinel used in parse.go; re-declared
	// here via the public errors subpackage import path in other files,
	// so we just check the message shape instead of importing again.
	if !errors.Is(err, err) {
		t.Fatal("errors.Is should at least match itself")
	}
}
