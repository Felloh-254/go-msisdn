package operators

func init() {
	RegisterCountry("GH", []Rule{
		{
			Operator: "MTN Ghana",
			Prefixes: []string{"24", "54", "55", "59"},
		},
		{
			Operator: "Vodafone Ghana",
			Prefixes: []string{"20", "50"},
		},
		{
			Operator: "AirtelTigo",
			Prefixes: []string{"26", "27", "56", "57"},
		},
	})
}
