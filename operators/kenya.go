package operators

func init() {
	RegisterCountry("KE", []Rule{
		{
			Operator: "Safaricom",
			Prefixes: []string{
				"70", "71", "72", "74", "79",
				"110", "111", "112", "113", "114", "115",
			},
		},
		{
			Operator: "Airtel Kenya",
			Prefixes: []string{
				"73", "75", "78",
				"100", "101", "102", "103", "105", "106", "107", "108", "109",
			},
		},
		{
			Operator: "Telkom Kenya",
			Prefixes: []string{"77"},
		},
		{
			Operator: "Equitel",
			Prefixes: []string{"76"},
		},
	})
}
