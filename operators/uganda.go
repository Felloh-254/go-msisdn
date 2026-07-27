package operators

func init() {
	RegisterCountry("UG", []Rule{
		{
			Operator: "MTN Uganda",
			Prefixes: []string{"77", "78", "76", "39"},
		},
		{
			Operator: "Airtel Uganda",
			Prefixes: []string{"70", "74", "75"},
		},
		{
			Operator: "Africell Uganda",
			Prefixes: []string{"79"},
		},
		{
			Operator: "Uganda Telecom",
			Prefixes: []string{"71"},
		},
	})
}
