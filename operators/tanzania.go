package operators

func init() {
	RegisterCountry("TZ", []Rule{
		{
			Operator: "Vodacom Tanzania",
			Prefixes: []string{"74", "75", "76"},
		},
		{
			Operator: "Airtel Tanzania",
			Prefixes: []string{"68", "69", "78"},
		},
		{
			Operator: "Yas (Tigo) Tanzania",
			Prefixes: []string{"71", "65", "67"},
		},
		{
			Operator: "Halotel",
			Prefixes: []string{"62"},
		},
		{
			Operator: "TTCL",
			Prefixes: []string{"73"},
		},
	})
}
