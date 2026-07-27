package msisdn

// This file is pure data. It registers every country go-msisdn knows
// about at package init time. Two tiers exist:
//
//   - "Deep" countries (Kenya, Uganda, Tanzania, Rwanda, Nigeria, Ghana,
//     Zambia) have full national-significant-number length validation and
//     per-number-type prefix ranges, matching the countries called out in
//     the project brief.
//   - "Shallow" countries have a name, calling code, and (where the rule
//     is simple and well known) a trunk prefix / length, so Parse,
//     Country/ISO/CountryCode and E.164 formatting work everywhere, even
//     though IsValid()/Type() fall back to best-effort heuristics.
//
// Extending coverage is purely additive: append to either table below, or
// call the unexported registerCountry from a new file in this package.

func init() {
	registerDeepCountries()
	registerShallowCountries()
}

func registerDeepCountries() {
	registerCountry(countryMeta{
		ISO: "KE", Name: "Kenya", CallingCode: 254, TrunkPrefix: "0",
		NSNLengths: []int{9},
		MobilePrefixes: []string{
			"70", "71", "72", "74", "79",
			"110", "111", "112", "113", "114", "115",
			"73", "75", "78",
			"100", "101", "102", "103", "105", "106", "107", "108", "109",
			"76", "77",
		},
		FixedLinePrefixes:   []string{"20", "40", "41", "43", "51", "53", "55", "57", "61", "62", "65"},
		TollFreePrefixes:    []string{"800"},
		PremiumRatePrefixes: []string{"900"},
		ExampleMobile:       "712345678",
		Deep:                true,
	})

	registerCountry(countryMeta{
		ISO: "UG", Name: "Uganda", CallingCode: 256, TrunkPrefix: "0",
		NSNLengths: []int{9},
		MobilePrefixes: []string{
			"77", "78", "76", "39", "70", "74", "75", "79", "71",
		},
		FixedLinePrefixes: []string{"20", "31", "41", "43", "46", "48"},
		TollFreePrefixes:  []string{"800"},
		ExampleMobile:     "712345678",
		Deep:              true,
	})

	registerCountry(countryMeta{
		ISO: "TZ", Name: "Tanzania", CallingCode: 255, TrunkPrefix: "0",
		NSNLengths: []int{9},
		MobilePrefixes: []string{
			"74", "75", "76", "68", "69", "78", "71", "65", "67", "62", "73",
		},
		FixedLinePrefixes: []string{"22", "23", "24", "25", "26", "27", "28"},
		TollFreePrefixes:  []string{"800"},
		ExampleMobile:     "754123456",
		Deep:              true,
	})

	registerCountry(countryMeta{
		ISO: "RW", Name: "Rwanda", CallingCode: 250, TrunkPrefix: "0",
		NSNLengths:        []int{9},
		MobilePrefixes:    []string{"78", "72", "73", "79"},
		FixedLinePrefixes: []string{"25"},
		ExampleMobile:     "781234567",
		Deep:              true,
	})

	registerCountry(countryMeta{
		ISO: "NG", Name: "Nigeria", CallingCode: 234, TrunkPrefix: "0",
		NSNLengths: []int{10},
		MobilePrefixes: []string{
			"803", "806", "813", "816", "810", "814", "903", "906", "916",
			"802", "808", "812", "901", "902", "907", "912",
			"805", "807", "815", "811", "905",
			"809", "817", "818", "908", "909",
		},
		TollFreePrefixes: []string{"700"},
		ExampleMobile:    "8031234567",
		Deep:             true,
	})

	registerCountry(countryMeta{
		ISO: "GH", Name: "Ghana", CallingCode: 233, TrunkPrefix: "0",
		NSNLengths: []int{9},
		MobilePrefixes: []string{
			"24", "54", "55", "59", "20", "50", "26", "27", "56", "57",
		},
		FixedLinePrefixes: []string{"30", "31", "32", "33", "34"},
		ExampleMobile:     "241234567",
		Deep:              true,
	})

	registerCountry(countryMeta{
		ISO: "ZM", Name: "Zambia", CallingCode: 260, TrunkPrefix: "0",
		NSNLengths:        []int{9},
		MobilePrefixes:    []string{"96", "76", "97", "77", "95", "75"},
		FixedLinePrefixes: []string{"21"},
		ExampleMobile:     "961234567",
		Deep:              true,
	})
}

// shallowCountry is a compact declaration for countries go-msisdn supports
// at the "name + calling code (+ optional length)" level.
type shallowCountry struct {
	iso, name           string
	callingCode         int
	trunkPrefix         string
	nsnLength           int // 0 = unknown/not enforced
	exampleNationalStub string
}

func registerShallowCountries() {
	list := []shallowCountry{
		{"US", "United States", 1, "", 10, "2015550123"},
		{"CA", "Canada", 1, "", 10, "5145550123"},
		{"GB", "United Kingdom", 44, "0", 10, "7911123456"},
		{"IE", "Ireland", 353, "0", 9, "851234567"},
		{"FR", "France", 33, "0", 9, "612345678"},
		{"DE", "Germany", 49, "0", 10, "1512345678"},
		{"ES", "Spain", 34, "", 9, "612345678"},
		{"PT", "Portugal", 351, "", 9, "912345678"},
		{"IT", "Italy", 39, "", 10, "3123456789"},
		{"NL", "Netherlands", 31, "0", 9, "612345678"},
		{"BE", "Belgium", 32, "0", 9, "470123456"},
		{"CH", "Switzerland", 41, "0", 9, "781234567"},
		{"SE", "Sweden", 46, "0", 9, "701234567"},
		{"NO", "Norway", 47, "", 8, "40612345"},
		{"DK", "Denmark", 45, "", 8, "20123456"},
		{"FI", "Finland", 358, "0", 9, "412345678"},
		{"PL", "Poland", 48, "", 9, "512345678"},
		{"AT", "Austria", 43, "0", 10, "6641234567"},
		{"GR", "Greece", 30, "", 10, "6912345678"},
		{"ZA", "South Africa", 27, "0", 9, "821234567"},
		{"EG", "Egypt", 20, "0", 10, "1001234567"},
		{"MA", "Morocco", 212, "0", 9, "612345678"},
		{"ET", "Ethiopia", 251, "0", 9, "911234567"},
		{"DZ", "Algeria", 213, "0", 9, "551234567"},
		{"TN", "Tunisia", 216, "", 8, "20123456"},
		{"CI", "Côte d'Ivoire", 225, "", 10, "0712345678"},
		{"SN", "Senegal", 221, "", 9, "701234567"},
		{"CM", "Cameroon", 237, "", 9, "671234567"},
		{"IN", "India", 91, "0", 10, "9123456789"},
		{"PK", "Pakistan", 92, "0", 10, "3001234567"},
		{"BD", "Bangladesh", 880, "0", 10, "1712345678"},
		{"CN", "China", 86, "0", 11, "13123456789"},
		{"JP", "Japan", 81, "0", 10, "9012345678"},
		{"KR", "South Korea", 82, "0", 10, "1012345678"},
		{"ID", "Indonesia", 62, "0", 11, "81234567890"},
		{"PH", "Philippines", 63, "0", 10, "9171234567"},
		{"VN", "Vietnam", 84, "0", 9, "912345678"},
		{"TH", "Thailand", 66, "0", 9, "812345678"},
		{"MY", "Malaysia", 60, "0", 9, "123456789"},
		{"SG", "Singapore", 65, "", 8, "81234567"},
		{"AU", "Australia", 61, "0", 9, "412345678"},
		{"NZ", "New Zealand", 64, "0", 9, "211234567"},
		{"BR", "Brazil", 55, "0", 11, "11987654321"},
		{"MX", "Mexico", 52, "", 10, "5512345678"},
		{"AR", "Argentina", 54, "0", 10, "91123456789"},
		{"CO", "Colombia", 57, "", 10, "3001234567"},
		{"PE", "Peru", 51, "", 9, "912345678"},
		{"CL", "Chile", 56, "", 9, "912345678"},
		{"RU", "Russia", 7, "8", 10, "9123456789"},
		{"TR", "Turkey", 90, "0", 10, "5321234567"},
		{"SA", "Saudi Arabia", 966, "0", 9, "512345678"},
		{"AE", "United Arab Emirates", 971, "0", 9, "501234567"},
		{"IL", "Israel", 972, "0", 9, "521234567"},
		{"QA", "Qatar", 974, "", 8, "33123456"},
	}
	for _, c := range list {
		m := countryMeta{
			ISO: c.iso, Name: c.name, CallingCode: c.callingCode,
			TrunkPrefix: c.trunkPrefix, Deep: false,
			ExampleMobile: c.exampleNationalStub,
		}
		if c.nsnLength > 0 {
			m.NSNLengths = []int{c.nsnLength}
		}
		registerCountry(m)
	}
}
