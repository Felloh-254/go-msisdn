package operators

func init() {
	RegisterCountry("RW", []Rule{
		{
			Operator: "MTN Rwanda",
			Prefixes: []string{"78"},
		},
		{
			Operator: "Airtel Rwanda",
			Prefixes: []string{"72", "73"},
		},
	})
}
