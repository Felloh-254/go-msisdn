package msisdn

import (
	"database/sql/driver"
	"fmt"
)

// Value implements driver.Valuer, storing the number as its E.164 string
// so it round-trips cleanly through PostgreSQL, MySQL, and SQLite text/
// varchar columns.
func (p Phone) Value() (driver.Value, error) {
	if p.meta == nil {
		return nil, nil
	}
	return p.E164(), nil
}

// Scan implements sql.Scanner, so a Phone (or *Phone) field on a struct
// can be populated directly from a database column via database/sql. The
// column value is parsed with Parse using "" as the region, so it must
// already be in E.164 form ("+254712345678") -- exactly what Value
// produces, which is what makes the pair round-trip safe.
func (p *Phone) Scan(src interface{}) error {
	if src == nil {
		*p = Phone{}
		return nil
	}

	var s string
	switch v := src.(type) {
	case string:
		s = v
	case []byte:
		s = string(v)
	default:
		return fmt.Errorf("msisdn: cannot scan %T into Phone", src)
	}

	if s == "" {
		*p = Phone{}
		return nil
	}

	parsed, err := Parse(s, "")
	if err != nil {
		return fmt.Errorf("msisdn: scanning Phone: %w", err)
	}
	*p = *parsed
	return nil
}
