// Command basic demonstrates the core go-msisdn API end to end. Run it
// with:
//
//	go run ./examples/basic
package main

import (
	"encoding/json"
	"fmt"

	"github.com/yourusername/go-msisdn"
)

func main() {
	// --- Parsing -----------------------------------------------------
	phone, err := msisdn.Parse("0712345678", "KE")
	if err != nil {
		panic(err)
	}

	fmt.Println("Country:       ", phone.Country())
	fmt.Println("ISO:           ", phone.ISO())
	fmt.Println("Country code:  ", phone.CountryCode())
	fmt.Println("National no.:  ", phone.NationalNumber())
	fmt.Println("E.164:         ", phone.E164())
	fmt.Println("International: ", phone.International())
	fmt.Println("National:      ", phone.National())
	fmt.Println("Valid:         ", phone.IsValid())
	fmt.Println("Type:          ", phone.Type())
	fmt.Println("Operator:      ", phone.Operator())
	fmt.Println("Masked:        ", phone.Mask())

	// --- Validation ----------------------------------------------------
	result := msisdn.Validate("0712345", "KE")
	fmt.Printf("\nValidate(short number) -> Valid=%v Reason=%q Code=%s\n",
		result.Valid, result.Reason, result.Code)

	// --- Normalization & comparison ------------------------------------
	norm, _ := msisdn.Normalize("+254 712 345 678", "KE")
	fmt.Println("\nNormalized:    ", norm)
	fmt.Println("Equal check:   ", msisdn.Equal("0712345678", "+254712345678", "KE"))

	// --- Masking --------------------------------------------------------
	fmt.Println("\nMask default:  ", msisdn.Mask("254712345678"))
	fmt.Println("Mask 6/2:      ", msisdn.Mask("254712345678", msisdn.WithPrefixVisible(6)))

	// --- Batch processing -------------------------------------------------
	numbers := []string{"0712345678", "0722000000", "invalid", "0771234567"}
	fmt.Println("\nBatch validation:")
	for _, res := range msisdn.ValidateMany(numbers, "KE") {
		fmt.Printf("  %-14s valid=%-5v reason=%s\n", "", res.Valid, res.Reason)
	}

	// --- Deduplication ----------------------------------------------------
	dupes := []string{"0712345678", "+254712345678", "254712345678"}
	fmt.Println("\nDeduped:       ", msisdn.Dedupe(dupes, "KE"))

	// --- Example numbers for tests ----------------------------------------
	example, _ := msisdn.Example("NG")
	fmt.Println("\nNigeria example:", example)

	// --- JSON --------------------------------------------------------------
	type User struct {
		Name  string       `json:"name"`
		Phone msisdn.Phone `json:"phone"`
	}
	u := User{Name: "Wanjiru", Phone: *phone}
	b, _ := json.Marshal(u)
	fmt.Println("\nJSON:          ", string(b))
}
