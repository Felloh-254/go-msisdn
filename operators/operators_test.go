package operators

import "testing"

func TestLookup(t *testing.T) {
	cases := []struct {
		name     string
		iso      string
		national string
		want     string
		wantOK   bool
	}{
		{"Kenya Safaricom", "KE", "712345678", "Safaricom", true},
		{"Kenya Airtel", "KE", "733123456", "Airtel Kenya", true},
		{"Kenya Telkom", "KE", "770123456", "Telkom Kenya", true},
		{"Kenya lowercase iso", "ke", "712345678", "Safaricom", true},
		{"Uganda MTN", "UG", "771234567", "MTN Uganda", true},
		{"Tanzania Vodacom", "TZ", "754123456", "Vodacom Tanzania", true},
		{"Rwanda MTN", "RW", "781234567", "MTN Rwanda", true},
		{"Nigeria MTN", "NG", "8031234567", "MTN Nigeria", true},
		{"Ghana MTN", "GH", "241234567", "MTN Ghana", true},
		{"Zambia MTN", "ZM", "961234567", "MTN Zambia", true},
		{"Unknown country", "XX", "712345678", "", false},
		{"Unknown prefix", "KE", "999999999", "", false},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			got, ok := Lookup(tc.iso, tc.national)
			if ok != tc.wantOK {
				t.Fatalf("Lookup(%q, %q) ok = %v, want %v", tc.iso, tc.national, ok, tc.wantOK)
			}
			if got != tc.want {
				t.Fatalf("Lookup(%q, %q) = %q, want %q", tc.iso, tc.national, got, tc.want)
			}
		})
	}
}

func TestLongestPrefixWins(t *testing.T) {
	RegisterCountry("ZZ", []Rule{
		{Operator: "Generic7", Prefixes: []string{"7"}},
		{Operator: "Specific712", Prefixes: []string{"712"}},
	})

	got, ok := Lookup("ZZ", "712345678")
	if !ok || got != "Specific712" {
		t.Fatalf("expected longest-prefix match Specific712, got %q (ok=%v)", got, ok)
	}

	got, ok = Lookup("ZZ", "700000000")
	if !ok || got != "Generic7" {
		t.Fatalf("expected fallback Generic7, got %q (ok=%v)", got, ok)
	}
}

func TestSupportedCountries(t *testing.T) {
	countries := SupportedCountries()
	want := map[string]bool{"KE": true, "UG": true, "TZ": true, "RW": true, "NG": true, "GH": true, "ZM": true}
	found := map[string]bool{}
	for _, c := range countries {
		found[c] = true
	}
	for iso := range want {
		if !found[iso] {
			t.Errorf("expected %s to be registered as supported", iso)
		}
	}
}
