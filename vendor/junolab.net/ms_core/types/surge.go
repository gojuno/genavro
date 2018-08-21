package types

import (
	"database/sql/driver"
	"fmt"

	"github.com/shopspring/decimal"
)

const surgePrecision = 2

var (
	// DefaultSurge represents the base of surges
	DefaultSurge = NewSurgeFromFloat(1)

	zero       = NewSurgeFromFloat(0.)
	decimal100 = decimal.NewFromFloat(100.0)
	coeff      = decimal.NewFromFloat(1.0)
	minSurge   = decimal.NewFromFloat(1.0)
	maxSurge   = decimal.NewFromFloat(9.99)
)

// Surge regulate price multiplier which depends on supply/demand in the real world.
type Surge struct {
	d decimal.Decimal
}

// NewSurgeFromDecimal creates a surge from a decimal with precision round
func NewSurgeFromDecimal(d decimal.Decimal) Surge {
	return Surge{d.Round(surgePrecision)}
}

// NewSurgeFromFloat creates a surge from float using precision in exponent
func NewSurgeFromFloat(f float64) Surge {
	return NewSurgeFromDecimal(floatToSurgeDecimal(f))
}

// NewSurgeFromString create a surge from string using decimal reader
func NewSurgeFromString(surgeString string) (Surge, error) {
	d, err := decimal.NewFromString(surgeString)
	if err != nil {
		return Surge{}, err
	}
	return NewSurgeFromDecimal(d), nil
}

// Decimal returns a surge in decimal.Decimal type
func (s Surge) Decimal() decimal.Decimal {
	return s.d
}

// ToCoefficient returns a surge minus 1 (surge - 1)
func (s Surge) ToCoefficient() decimal.Decimal {
	return s.d.Sub(coeff)
}

// ToPercentInt returns an integer part of a surge coefficient (surge - 1) * 100
func (s Surge) ToPercentInt() int64 {
	return s.ToCoefficient().Mul(decimal100).Round(0).IntPart()
}

// IsSurgeExist returns true if an integer part is positive and non-zero
func (s Surge) IsSurgeExist() bool {
	return s.ToPercentInt() > 0
}

// IsZero checks whether surge is zero or not
func (s Surge) IsZero() bool {
	return s.Equals(zero)
}

// Equals returns whether the numbers represented by s and s2 are equal.
func (s Surge) Equals(s2 Surge) bool {
	return s.d.Equals(s2.d)
}

// Less returns true if s < s2
func (s Surge) Less(s2 Surge) bool {
	return s.d.Cmp(s2.d) < 0
}

// Greater returns true if s > s2
func (s Surge) Greater(s2 Surge) bool {
	return s.d.Cmp(s2.d) > 0
}

// ApplyTo multiplies a float by surge
func (s Surge) ApplyTo(f float64) decimal.Decimal {
	return s.d.Mul(floatToSurgeDecimal(f))
}

// IsDefault returns true if a surge equals 1.0
func (s Surge) IsDefault() bool {
	return s.Equals(DefaultSurge)
}

// Validate checks surge precision and max/min bound
func (s Surge) Validate() error {
	if s.d.Exponent() > surgePrecision {
		return fmt.Errorf("surge precision should be less or equal to %d, given %d", surgePrecision, s.d.Exponent())
	}
	if s.d.Cmp(minSurge) == -1 {
		return fmt.Errorf("surge is less than %v: %v", minSurge, s.d)
	}
	if s.d.Cmp(maxSurge) == 1 {
		return fmt.Errorf("surge is more than %v: %v", maxSurge, s.d)
	}
	return nil
}

// MarshalJSON implements the json.Marshaler interface.
func (s Surge) MarshalJSON() ([]byte, error) {
	return []byte(s.String()), nil
}

// UnmarshalJSON implements the json.Unmarshaler interface.
func (s *Surge) UnmarshalJSON(data []byte) error {
	var d decimal.Decimal
	if err := d.UnmarshalJSON(data); err != nil {
		return err
	}
	*s = NewSurgeFromDecimal(d)
	return s.Validate()
}

func (s Surge) String() string {
	return s.d.StringFixed(surgePrecision)
}

// Scan implements the sql.Scanner interface for database deserialization.
func (s *Surge) Scan(value interface{}) error {
	str, err := Unquote(value)
	if err != nil {
		return err
	}
	*s, err = NewSurgeFromString(str)
	return err
}

// Value implements the driver.Valuer interface for database serialization.
func (s Surge) Value() (driver.Value, error) {
	return s.String(), nil
}

func floatToSurgeDecimal(f float64) decimal.Decimal {
	return decimal.NewFromFloatWithExponent(f, -surgePrecision)
}
