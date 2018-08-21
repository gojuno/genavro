package decimal

import (
	"database/sql/driver"
	"fmt"

	"github.com/shopspring/decimal"
)

var (
	MinusOne = New(-1, 0)
	Zero     = Decimal{val: decimal.Zero}
	One      = New(1, 0)
	Ten      = New(10, 0)
	Hundred  = New(100, 0)
)

type Decimal struct {
	val decimal.Decimal
}

// NewFromString creates new Decimal from string
func NewFromString(value string) (Decimal, error) {
	d, err := decimal.NewFromString(value)
	return Decimal{val: d}, err
}

// New returns a new fixed-point decimal, value * 10 ^ exp
func New(value int64, exp int32) Decimal {
	return Decimal{val: decimal.New(value, exp)}
}

// NewFromFloat creates new Decimal from float
func NewFromFloat(value float64) Decimal {
	return Decimal{val: decimal.NewFromFloat(value)}
}

// NewFromFloatWithExponent returns value rounded to 10 ^ exp (exp would typically be negative),
// e.g. NewFromFloatWithExponent(10.333, -2) = Decimal(10.33)
func NewFromFloatWithExponent(value float64, exp int32) Decimal {
	return Decimal{val: decimal.NewFromFloatWithExponent(value, exp)}
}

// NewFromFloatWithPrecision returns value rounded to 10 ^ (-precision) (precision would typically be positive),
// e.g. NewFromFloatWithPrecision(10.333, 2) = Decimal(10.33)
func NewFromFloatWithPrecision(value float64, precision int32) Decimal {
	return NewFromFloatWithExponent(value, -precision)
}

// NewFromDecimal creates Decimal from original shopspring decimal type
func NewFromDecimal(d decimal.Decimal) Decimal {
	return Decimal{val: d}
}

func (d Decimal) Abs() Decimal {
	return Decimal{val: d.val.Abs()}
}

func (d Decimal) Add(d2 Decimal) Decimal {
	return Decimal{val: d.val.Add(d2.val)}
}

func (d Decimal) Sub(d2 Decimal) Decimal {
	return Decimal{val: d.val.Sub(d2.val)}
}

func (d Decimal) Mul(d2 Decimal) Decimal {
	return Decimal{val: d.val.Mul(d2.val)}
}

func (d Decimal) Div(d2 Decimal) Decimal {
	return Decimal{val: d.val.Div(d2.val)}
}

// Cmp compares the numbers represented by d and d2 and returns:
//
//     -1 if d <  d2
//      0 if d == d2
//     +1 if d >  d2
//
func (d Decimal) Cmp(d2 Decimal) int {
	return d.val.Cmp(d2.val)
}

func (d Decimal) Equal(d2 Decimal) bool {
	return d.val.Equals(d2.val) // in recent lib version Equals deprecated in favor of Equal
}

func (d Decimal) GreaterThan(d2 Decimal) bool {
	return d.Cmp(d2) == 1
}

func (d Decimal) GreaterThanOrEqual(d2 Decimal) bool {
	cmp := d.Cmp(d2)
	return cmp == 1 || cmp == 0
}

func (d Decimal) LessThan(d2 Decimal) bool {
	return d.Cmp(d2) == -1
}

func (d Decimal) LessThanOrEqual(d2 Decimal) bool {
	cmp := d.Cmp(d2)
	return cmp == -1 || cmp == 0
}

func (d Decimal) IsPositive() bool {
	return d.Cmp(Zero) > 0
}

func (d Decimal) IsNegative() bool {
	return d.Cmp(Zero) < 0
}

func (d Decimal) IsZero() bool {
	return d.Equal(Zero)
}

// IntPart returns the integer component of the decimal.
func (d Decimal) IntPart() int64 {
	return d.val.IntPart()
}

// Exponent returns the exponent, or scale component of the decimal.
// Typically negative, e.g. for 0.01 it would be -2
func (d Decimal) Exponent() int32 {
	return d.val.Exponent()
}

// Float64 returns the nearest float64 value for d and a bool indicating
// whether f represents d exactly.
// For more details, see the documentation for big.Rat.Float64
func (d Decimal) Float64() (f float64, exact bool) {
	return d.val.Float64()
}

// Precision returns the decimal precision.
// Typically positive, e.g. for 0.01 it would be 2
func (d Decimal) Precision() int32 {
	return -d.Exponent()
}

func (d Decimal) Round(places int32) Decimal {
	return Decimal{val: d.val.Round(places)}
}

func (d Decimal) Floor() Decimal {
	return Decimal{val: d.val.Floor()}
}

func (d Decimal) Ceil() Decimal {
	return Decimal{val: d.val.Ceil()}
}

func (d Decimal) Truncate(precision int32) Decimal {
	return Decimal{val: d.val.Truncate(precision)}
}

func (d Decimal) Negate() Decimal {
	return d.Mul(MinusOne)
}

func (d Decimal) String() string {
	return d.val.String()
}

func (d Decimal) StringFixed(places int32) string {
	return d.val.StringFixed(places)
}

// MarshalJSON implements the json.Marshaler interface.
func (d Decimal) MarshalJSON() ([]byte, error) {
	return d.val.MarshalJSON() // this marshals value as quoted string
}

// UnmarshalJSON implements the json.Unmarshaler interface.
func (d *Decimal) UnmarshalJSON(decimalBytes []byte) error {
	return d.val.UnmarshalJSON(decimalBytes)
}

// MarshalText implements the encoding.TextMarshaler interface.
func (d Decimal) MarshalText() ([]byte, error) {
	return d.val.MarshalText()
}

// UnmarshalText implements the encoding.TextUnmarshaler.
func (d *Decimal) UnmarshalText(text []byte) error {
	return d.val.UnmarshalText(text)
}

// Value implements the driver.Valuer interface for database serialization.
func (d Decimal) Value() (driver.Value, error) {
	return d.val.Value()
}

// Scan implements the sql.Scanner interface for database deserialization.
func (d *Decimal) Scan(value interface{}) error {
	switch v := value.(type) {
	case string:
		return d.val.Scan([]byte(v))
	default:
		return d.val.Scan(v)
	}
}

// Decode implements the junolab.net/ms/core/config.Decodable interface
func (d *Decimal) Decode(data interface{}) error {
	val, err := NewFromString(fmt.Sprintf("%v", data))
	if err != nil {
		return err
	}
	*d = val
	return nil
}

func (d Decimal) ValidatePositive() error {
	if !d.GreaterThan(Zero) {
		return fmt.Errorf("value %v is not positive", d)
	}
	return nil
}

func (d Decimal) ValidateNegative() error {
	if !d.LessThan(Zero) {
		return fmt.Errorf("value %v is not negative", d)
	}
	return nil
}

func (d Decimal) ValidateNonPositive() error {
	if d.GreaterThan(Zero) {
		return fmt.Errorf("value %v is positive", d)
	}
	return nil
}

func (d Decimal) ValidateNonNegative() error {
	if d.LessThan(Zero) {
		return fmt.Errorf("value %v is negative", d)
	}
	return nil
}

func (d *Decimal) ValueOrZero() Decimal {
	if d == nil {
		return Zero
	}
	return *d
}

func Min(first Decimal, rest ...Decimal) Decimal {
	ans := first
	for _, item := range rest {
		if item.Cmp(ans) < 0 {
			ans = item
		}
	}
	return ans
}

func Max(first Decimal, rest ...Decimal) Decimal {
	ans := first
	for _, item := range rest {
		if item.Cmp(ans) > 0 {
			ans = item
		}
	}
	return ans
}
