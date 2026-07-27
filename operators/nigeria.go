package operators

func init() {
	RegisterCountry("NG", []Rule{
		{
			Operator: "MTN Nigeria",
			Prefixes: []string{
				"803", "806", "813", "816", "810", "814", "903", "906", "916",
			},
		},
		{
			Operator: "Airtel Nigeria",
			Prefixes: []string{
				"802", "808", "812", "901", "902", "907", "912",
			},
		},
		{
			Operator: "Globacom",
			Prefixes: []string{
				"805", "807", "815", "811", "905",
			},
		},
		{
			Operator: "9mobile",
			Prefixes: []string{
				"809", "817", "818", "908", "909",
			},
		},
	})
}
