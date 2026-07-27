package operators

func init() {
	RegisterCountry("ZM", []Rule{
		{
			Operator: "MTN Zambia",
			Prefixes: []string{"96", "76"},
		},
		{
			Operator: "Airtel Zambia",
			Prefixes: []string{"97", "77"},
		},
		{
			Operator: "Zamtel",
			Prefixes: []string{"95", "75"},
		},
	})
}
