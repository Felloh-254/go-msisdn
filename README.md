# go-msisdn

A developer-friendly phone number toolkit for Go — parsing, validation, normalization,
formatting, type detection, and telecom operator lookup, built around a single, clean
`Phone` type. Designed for backend systems, fintech, CRM, telecom, and iGaming platforms
where you need more than a regex to know whether `+254 (712)-345-678` is a real,
dialable Safaricom number.

```go
phone, err := msisdn.Parse("0712345678", "KE")
// phone.E164()      -> "+254712345678"
// phone.Country()   -> "Kenya"
// phone.Operator()  -> "Safaricom"
// phone.IsMobile()  -> true
```

---

## Why go-msisdn

Most Go "phone validators" are a regex and a prayer. go-msisdn instead gives you:

- A **single `Phone` value type** that carries country, national number, validity,
  number type, and operator — instead of scattering that logic across your codebase.
- **Structured validation results**, not just `true`/`false` — you get _why_ a number
  is invalid (wrong length, unrecognized prefix, unknown country code, ...).
- **Configurable, data-driven operator detection** — adding a country or operator is a
  data change, never a code change.
- **First-class JSON and `database/sql` support** — `Phone` fields on your structs just
  work with Postgres/MySQL/SQLite and with `encoding/json`.
- **Batch helpers** for validating, normalizing, and parsing lists of numbers, which is
  most of what you actually do with phone numbers in production.

### A note on libphonenumber

The original design for this library called for building on top of Google's
libphonenumber via a Go port (e.g. `nyaruka/phonenumbers`). This build was produced in
a network-sandboxed environment that could not reach the Go module proxy or
`golang.org`, so it ships instead with a **self-contained, zero-dependency** parsing and
validation engine with detailed metadata for its "deep" countries (see below) and
lighter-weight support for ~50 more.

The engine sits behind the same `Phone` API a libphonenumber-backed implementation
would expose, and the parsing/validation/formatting logic is isolated in a small number
of files (`parse.go`, `validate.go`, `format.go`, `metadata.go`, `countries.go`) so that
swapping in a real libphonenumber binding later — for full global coverage — is a
contained, backwards-compatible change and does not touch the public API.

---

## Installation

```bash
go get https://github.com/Felloh-254/go-msisdn.git
```

Requires Go 1.21+. Zero external dependencies.

---

## Quick Start

```go
package main

import (
	"fmt"

	"github.com/Felloh-254/go-msisdn"
)

func main() {
	phone, err := msisdn.Parse("+254712345678", "")
	if err != nil {
		panic(err)
	}
	fmt.Println(phone.E164()) // +254712345678
}
```

---

## API Examples

### Parsing

```go
phone, err := msisdn.Parse("0712345678", "KE")
if err != nil {
	// err is only returned for structural problems: empty input, unknown
	// calling code, or a missing/unsupported region for a non-"+" number.
}

phone.Country()        // "Kenya"
phone.ISO()             // "KE"
phone.CountryCode()     // 254
phone.NationalNumber()  // 712345678
phone.E164()            // "+254712345678"
phone.International()   // "+254 712 345678"
phone.National()        // "0712 345678"
phone.IsValid()          // true
phone.Type()             // msisdn.Mobile
phone.Operator()         // "Safaricom"
```

A number that _parses_ but doesn't _validate_ (wrong length, unrecognized prefix) is
still returned — `Parse` doesn't error out on that, because "is this number valid" and
"could I even make sense of this input" are different questions:

```go
phone, _ := msisdn.Parse("07123", "KE") // too short
phone.IsValid()        // false
phone.InvalidReason()  // "invalid length for Kenya: got 5 digits, expected [9]"
```

### Validation

```go
result := msisdn.Validate("07123", "KE")
// result.Valid    -> false
// result.Possible -> true (right ballpark, just short)
// result.Reason   -> "invalid length for Kenya: got 5 digits, expected [9]"
// result.Code     -> "INVALID_LENGTH"

phone, _ := msisdn.Parse("0712345678", "KE")
phone.IsValid()    // true
phone.IsPossible() // true
```

`Validate` never returns a Go `error` — a structurally unparseable input (empty string,
missing region, unknown calling code) is reported as an invalid `ValidationResult`
instead, so it's safe to use directly for form/API input validation.

### Normalization

```go
msisdn.Normalize("0712345678", "KE")    // "254712345678", nil
msisdn.Normalize("712345678", "KE")     // "254712345678", nil
msisdn.Normalize("+254712345678", "")   // "254712345678", nil
msisdn.Normalize("254712345678", "KE")  // "254712345678", nil
```

### Formatting

```go
msisdn.Format("0712345678", "KE", msisdn.E164)         // "+254712345678"
msisdn.Format("0712345678", "KE", msisdn.National)      // "0712 345678"
msisdn.Format("0712345678", "KE", msisdn.International) // "+254 712 345678"
msisdn.Format("0712345678", "KE", msisdn.RFC3966)       // "tel:+254712345678"

// Or on an already-parsed Phone:
phone.Format(msisdn.National) // "0712 345678"
```

> **Naming note:** the design brief asked for `Format(number, E164)` /
> `Format(number, NATIONAL)`. Go doesn't allow a type and a top-level function to share
> an identifier, so the _type_ is `msisdn.Style` while the constants (`E164`,
> `National`, `International`, `RFC3966`) and the top-level `Format` function keep
> exactly the requested call shape: `msisdn.Format(number, region, msisdn.National)`.

### Cleaning

```go
msisdn.Clean("+254 (712)-345-678") // "254712345678"
```

### Comparison

```go
msisdn.Equal("0712345678", "254712345678", "KE")  // true
msisdn.Equal("0712345678", "+254712345678", "KE")  // true
msisdn.Equal("0712345678", "0722345678", "KE")     // false

phone1.Equal(phone2) // compare two already-parsed *Phone values
```

### Masking

```go
msisdn.Mask("254712345678") // "2547******78"

msisdn.Mask("254712345678",
	msisdn.WithPrefixVisible(6),
	msisdn.WithSuffixVisible(2),
) // "254712****78"

phone.Mask() // mask an already-parsed Phone's E.164 form
```

### Country metadata

```go
phone.Country()     // "Kenya"
phone.ISO()          // "KE"
phone.CountryCode()  // 254
```

### Number type detection

```go
phone.Type()             // msisdn.Mobile
phone.IsMobile()          // true
phone.IsFixedLine()       // false
phone.IsTollFree()        // false
phone.IsPremiumRate()     // false
phone.IsVoIP()            // false
phone.IsPager()           // false
```

Recognized types: `Mobile`, `FixedLine`, `FixedLineOrMobile`, `TollFree`,
`PremiumRate`, `VoIP`, `Pager`, `Unknown`.

### Operator detection

```go
phone, _ := msisdn.Parse("0712345678", "KE")
phone.Operator() // "Safaricom"
```

Operator lookup is entirely data-driven (see [`operators/`](./operators)) — a
longest-prefix match against a per-country table registered via
`operators.RegisterCountry`. Adding a country or a new operator prefix is a small,
self-contained data file, not a change to any lookup logic:

```go
operators.RegisterCountry("XX", []operators.Rule{
	{Operator: "ExampleTel", Prefixes: []string{"70", "71"}},
})
```

### Local conversion

```go
msisdn.ToLocal("+254712345678", "") // "0712345678", nil
phone.Local()                        // "0712345678"
```

### Deduplication

```go
msisdn.Dedupe([]string{
	"0712345678",
	"+254712345678",
	"254712345678",
}, "KE")
// []string{"+254712345678"}
```

### Batch processing

```go
numbers := []string{"0712345678", "not-a-number", "0771234567"}

msisdn.ParseMany(numbers, "KE")     // []ParseResult
msisdn.ValidateMany(numbers, "KE")  // []ValidationResult
msisdn.NormalizeMany(numbers, "KE") // []NormalizeResult
```

Every result is keyed to its input by index, and per-item failures never abort the
batch.

### Example numbers (for tests & fixtures)

```go
msisdn.Example("KE") // "+254712345678", nil
msisdn.ExamplePhone("NG") // *Phone, nil
```

### Database support

```go
type User struct {
	ID    int
	Phone msisdn.Phone
}

_, err := db.Exec(`INSERT INTO users (id, phone) VALUES ($1, $2)`, user.ID, user.Phone)

var u User
err = db.QueryRow(`SELECT id, phone FROM users WHERE id = $1`, id).Scan(&u.ID, &u.Phone)
```

`Phone` implements `driver.Valuer` (stores as E.164 text) and `sql.Scanner` (parses
E.164 text back into a `Phone`), so it works naturally as a column type against
PostgreSQL, MySQL, and SQLite `text`/`varchar` columns.

### JSON support

```go
type User struct {
	Name  string      `json:"name"`
	Phone msisdn.Phone `json:"phone"`
}

data := []byte(`{"name":"Wanjiru","phone":"+254712345678"}`)
var u User
json.Unmarshal(data, &u)
u.Phone.E164() // "+254712345678"

json.Marshal(u) // {"name":"Wanjiru","phone":"+254712345678"}
```

---

## Supported Countries

### Deep support (full length + number-type + operator detection)

| Country  | ISO | Calling code | Operators                                       |
| -------- | --- | ------------ | ----------------------------------------------- |
| Kenya    | KE  | 254          | Safaricom, Airtel Kenya, Telkom Kenya, Equitel  |
| Uganda   | UG  | 256          | MTN Uganda, Airtel Uganda, Africell Uganda, UTL |
| Tanzania | TZ  | 255          | Vodacom, Airtel, Yas (Tigo), Halotel, TTCL      |
| Rwanda   | RW  | 250          | MTN Rwanda, Airtel Rwanda                       |
| Nigeria  | NG  | 234          | MTN, Airtel, Globacom, 9mobile                  |
| Ghana    | GH  | 233          | MTN Ghana, Vodafone Ghana, AirtelTigo           |
| Zambia   | ZM  | 260          | MTN Zambia, Airtel Zambia, Zamtel               |

### Shallow support (name, calling code, E.164/national formatting)

United States, Canada, United Kingdom, Ireland, France, Germany, Spain, Portugal, Italy,
Netherlands, Belgium, Switzerland, Sweden, Norway, Denmark, Finland, Poland, Austria,
Greece, South Africa, Egypt, Morocco, Ethiopia, Algeria, Tunisia, Côte d'Ivoire, Senegal,
Cameroon, India, Pakistan, Bangladesh, China, Japan, South Korea, Indonesia,
Philippines, Vietnam, Thailand, Malaysia, Singapore, Australia, New Zealand, Brazil,
Mexico, Argentina, Colombia, Peru, Chile, Russia, Turkey, Saudi Arabia, United Arab
Emirates, Israel, Qatar.

`IsValid()`/`Type()` are best-effort for shallow countries: length is checked where the
rule is simple and well known, but number-type/operator prefix ranges aren't populated.
Adding deep support for any of these (or a new country entirely) means adding an entry
to `countries.go` — see [Extending](#extending-the-library) below.

> **Known limitation:** fixed-line prefix ranges for the deep countries are
> illustrative, not exhaustive. Mobile ranges (the primary use case for OTP/fintech/
> iGaming workloads) are the focus of this build.

---

## Supported Features

- [x] Parsing (`Parse`, `ParseMany`)
- [x] Validation with structured reasons (`Validate`, `ValidateMany`, `IsValid`, `IsPossible`)
- [x] Normalization (`Normalize`, `NormalizeMany`)
- [x] Multi-format output (`Format`, `E164`, `National`, `International`, `RFC3966`)
- [x] Cleaning (`Clean`)
- [x] Comparison (`Equal`)
- [x] Configurable masking (`Mask`, `WithPrefixVisible`, `WithSuffixVisible`, `WithMaskChar`)
- [x] Country metadata (`Country`, `ISO`, `CountryCode`)
- [x] Number type detection (`Type`, `IsMobile`, `IsFixedLine`, `IsTollFree`, `IsPremiumRate`, `IsVoIP`, `IsPager`)
- [x] Data-driven operator detection (`Operator`, `operators.RegisterCountry`)
- [x] Local conversion (`ToLocal`, `Phone.Local`)
- [x] Deduplication (`Dedupe`)
- [x] Batch processing (`ParseMany`, `ValidateMany`, `NormalizeMany`)
- [x] Example number generator (`Example`, `ExamplePhone`)
- [x] `database/sql` `Scanner`/`Valuer`
- [x] `encoding/json` `Marshaler`/`Unmarshaler`

---

## Architecture

```
go-msisdn/
├── go.mod
├── LICENSE
├── README.md
├── *.go                 # root "msisdn" package: Phone type, Parse, Validate,
│                         # Normalize, Format, Mask, batch helpers, JSON/DB support
├── errors/               # sentinel errors + structured ValidationError
│   └── errors.go
├── operators/             # data-driven operator (MNO) prefix lookup
│   ├── operators.go        # registry + longest-prefix-match Lookup
│   ├── kenya.go, uganda.go, tanzania.go, rwanda.go,
│   │   nigeria.go, ghana.go, zambia.go     # per-country data, each an init()
│   └── operators_test.go
├── examples/
│   └── basic/main.go       # runnable end-to-end example (`go run ./examples/basic`)
├── msisdn_test.go          # table-driven core tests
└── msisdn_edge_test.go     # edge cases
```

### Why this layout, and not `cmd/`, `internal/`, `pkg/`

The brief's suggested layout (`cmd/`, `internal/`, `pkg/`, plus one subpackage per
concern: `parser/`, `validator/`, `formatter/`, `types/`) is a common enterprise-Java-
style convention, but it's explicitly discouraged for libraries by the Go community
(there's no `cmd/` because this is a library, not a binary; wrapping everything in
`pkg/` adds a directory level with no semantic value; and splitting `parser`/
`validator`/`formatter`/`types` into separate packages just to keep a single cohesive
concept — "a phone number" — apart invites import cycles, since a parser needs the
`Phone` type, the formatter needs the parser's output, and the validator needs both).

Instead:

- The **root package (`msisdn`)** owns the `Phone` type and everything that operates
  directly on it (parsing, validation, formatting, normalization, masking, batching,
  JSON, DB). This is idiomatic for a focused Go library — see how `time`, `net/url`, or
  `net/mail` are structured: one cohesive package, split across files by concern, not
  by artificial package boundaries.
- **`operators/`** is a genuinely separate concern (it doesn't need to know about
  parsing or formatting — just "given a country and a prefix, what operator?") and is
  explicitly meant to be user-extensible, so it's a real subpackage with its own
  registry.
- **`errors/`** is a real subpackage because sentinel errors are meant to be a stable,
  independently-importable contract (`errors.Is(err, msisdnerrors.ErrInvalidLength)`)
  that shouldn't force importing the whole library.
- **`examples/`** holds a runnable example program rather than living under `cmd/`
  (which implies "this repo produces a CLI binary", which it doesn't).
- Tests live next to the code they test (`msisdn_test.go`, `operators_test.go`), which
  is the standard Go convention — there's no separate `tests/` directory, since Go's
  tooling (`go test ./...`) doesn't need one and a parallel test tree would just drift
  out of sync with the code.

### Public API design

- **One entry type, `Phone`.** Everything you'd want to know about a number hangs off
  one value, returned by `Parse`. No separate "metadata" object to keep in sync.
- **Package-level functions for one-shot use** (`Normalize`, `Clean`, `Equal`, `Mask`,
  `Format`, `Validate`, `Dedupe`) so simple call sites don't need to hold onto a `Phone`
  at all — mirroring how `strings` and `path/filepath` are used.
- **Errors vs. invalidity are different things.** `Parse` returns a Go `error` only when
  it structurally cannot make sense of the input (empty string, unrecognized country
  code, missing region). A number that parses but fails business-rule validation
  (wrong length, unknown prefix) is returned successfully, with `IsValid()`/
  `IsPossible()`/`InvalidReason()` telling you why. This mirrors how you'd want to
  handle user-submitted numbers in a signup form: you don't want a Go `error` (and the
  awkward `if err != nil` branching that implies) for "this looks like a phone number
  but it's the wrong length" — you want a value you can inspect and show a message for.
- **Structured error codes** (`errors.Code`, e.g. `"INVALID_LENGTH"`) sit alongside
  human-readable reasons, so API responses can switch on a stable code while logs get a
  readable message.
- **Functional options for `Mask`** (`WithPrefixVisible`, ...) rather than a config
  struct, since masking has a good, unsurprising default and options are the idiomatic
  way to make optional knobs discoverable via autocomplete without an explosion of
  `MaskN` function variants.
- **Nil-safe methods on `*Phone`.** Every getter is safe to call on a nil `*Phone`
  (returns zero values), so partially-initialized state (e.g. a struct field that
  hasn't been parsed yet) doesn't panic when read.

---

## Testing

```bash
go test ./...
go test ./... -cover
go test ./... -v
```

Tests are table-driven throughout and cover: Kenyan/Ugandan/Tanzanian/Rwandan/Nigerian/
Ghanaian/Zambian numbers, invalid and impossible numbers, every input format (E.164,
"00" international prefix, national with/without trunk zero, punctuated), country and
operator detection, number-type classification, JSON round-tripping, `database/sql`
`Scan`/`Value`, nil-safety, and shallow-country fallback behavior.

---

## Extending the Library

**Add an operator** (existing country):

```go
operators.RegisterCountry("KE", append(operators.SupportedCountries() /* ... */))
// or simply add a new Rule to operators/kenya.go
```

**Add a new country with full support**: add a `countryMeta` entry to
`registerDeepCountries()` in `countries.go` (calling code, trunk prefix, valid NSN
lengths, per-type prefix ranges) and, if you want operator detection, a new file under
`operators/` calling `operators.RegisterCountry` in an `init()`.

**Add shallow support for a new country**: append a row to the `list` in
`registerShallowCountries()` in `countries.go`.

No existing code needs to change in either case — the registries are additive.

---

## Contributing

1. Fork the repo and create a feature branch.
2. Keep the "deep vs. shallow" split in mind: if you're adding real validation depth
   for a country, prefer accuracy over exhaustiveness — cite a source for prefix ranges
   in your PR description where possible.
3. Add table-driven tests for anything you change; `go test ./... -cover` should not
   regress.
4. Run `gofmt -l .` and `go vet ./...` before opening a PR.
5. Open a PR describing the change and, for new countries/operators, the source of the
   numbering data.

---

## License

MIT — see [LICENSE](./LICENSE).
