package types

import (
	"database/sql/driver"
	"errors"
	"fmt"

	"github.com/shopspring/decimal"
	lib_decimal "junolab.net/lib_api/decimal"
)

// DEPRECATED: use lib_api/decimal.Decimal if you need just amount, or lib_api/money.Money, if you need amount+currency
type Money struct {
	decimal.Decimal `db:"-"`
}

var ZeroMoney Money

// Generates money value from various input types
// DEPRECATED: use lib_api/decimal.Decimal if you need just amount, or lib_api/money.Money, if you need amount+currency
func NewMoney(i interface{}) (Money, error) {
	var err error
	var d decimal.Decimal

	switch v := i.(type) {
	case string:
		d, err = decimal.NewFromString(v)
	case float64:
		d = decimal.NewFromFloat(v)
	case int64:
		d = decimal.New(v, 0)
	case int:
		d = decimal.New(int64(v), 0)
	case decimal.Decimal:
		d = v
	case Decimal:
		d = v.Decimal
	case Money:
		d = v.Decimal
	case lib_decimal.Decimal:
		d, err = decimal.NewFromString(v.String())
	default:
		err = fmt.Errorf("Can't convert %+v to types.Money", i)
	}

	return Money{d}, err
}

// DEPRECATED: use lib_api/decimal.Decimal if you need just amount, or lib_api/money.Money, if you need amount+currency
func NewMoneyFromDecimal(d decimal.Decimal) Money {
	return Money{d}
}

// DEPRECATED: use lib_api/decimal.Decimal if you need just amount, or lib_api/money.Money, if you need amount+currency
func NewMoneyFromString(s string) (Money, error) {
	// TODO: this is pretty much basic implementation - should be rewritten when currency is added
	decValue, err := decimal.NewFromString(s)
	return Money{decValue}, err
}

// Generates money value from float64
// DEPRECATED: use lib_api/decimal.Decimal if you need just amount, or lib_api/money.Money, if you need amount+currency
func NewMoneyFromFloat(i float64) Money {
	return Money{decimal.NewFromFloat(i)}
} //TODO: rewrite it through base NewMoney() method

// DEPRECATED: use lib_api/decimal.Decimal if you need just amount, or lib_api/money.Money, if you need amount+currency
func NewMoneyFromFloatWithExponent(i float64, exp int32) Money {
	return Money{decimal.NewFromFloatWithExponent(i, exp)}
}

func (m Money) ToLibDecimal() lib_decimal.Decimal {
	return lib_decimal.NewFromDecimal(m.Decimal)
}

func (m Money) String() string {
	return m.StringFixed(2)
}

// MarshalJSON implements the json.Marshaler interface.
func (m Money) MarshalJSON() ([]byte, error) {
	return []byte(m.StringFixed(2)), nil
}

// UnmarshalJSON implements the json.Unmarshaler interface.
func (m *Money) UnmarshalJSON(decimalBytes []byte) error {
	d := &decimal.Decimal{}
	if err := d.UnmarshalJSON(decimalBytes); err != nil {
		return err
	}

	*m = Money{*d}

	return m.Validate()
}

func (r Money) Validate() error {
	//TODO: we'll need more sophisticated validations once
	//we have many currencies with many conversion/rounding rules
	//right now we just check for 2 or less digits after decimal point
	if r.Decimal.Exponent() < -2 {
		return errors.New("money value can't have more than 2 digits after decimal point")
	}
	return nil
}

func (r Money) ValidateNotNegative() error {
	if r.LessThan(0.0) {
		return fmt.Errorf("money value is less than 0.0: %v", r)
	}
	return r.Validate()
}

// Scan implements the sql.Scanner interface for database deserialization.
func (m *Money) Scan(value interface{}) error {
	str, err := Unquote(value)
	if err != nil {
		return err
	}
	*m, err = NewMoneyFromString(str)

	return err
}

// Value implements the driver.Valuer interface for database serialization.
func (m Money) Value() (driver.Value, error) {
	return m.String(), nil
}

func (m1 Money) Add(m2 Money) Money {
	return NewMoneyFromDecimal(m1.Decimal.Add(m2.Decimal))
}

func (m *Money) Inc(i interface{}) error {
	im, err := NewMoney(i)
	if err == nil {
		m.Decimal = m.Decimal.Add(im.Decimal)
	}

	return err
}

func (m *Money) Dec(i interface{}) error {
	im, err := NewMoney(i)
	if err == nil {
		m.Decimal = m.Decimal.Sub(im.Decimal)
	}
	return err
}

func (m *Money) Mul(i interface{}) error {
	im, err := NewMoney(i)
	if err == nil {
		m.Decimal = m.Decimal.Mul(im.Decimal)
	}
	return err
}

func (m *Money) IncFloat(i float64) *Money {
	m.Decimal = m.Decimal.Add(decimal.NewFromFloat(i))
	return m
}

func (m *Money) DecFloat(i float64) *Money {
	m.Decimal = m.Decimal.Sub(decimal.NewFromFloat(i))
	return m
}

func (m *Money) MulFloat(i float64) Money {
	return NewMoneyFromDecimal(m.Decimal.Mul(decimal.NewFromFloat(i)))
}

func (m Money) IsZero() bool {
	return m.Decimal.Equals(decimal.Zero)
}

func (m Money) IsPositive() bool {
	return m.Decimal.Cmp(decimal.Zero) == 1
}

func (m Money) Float64() float64 {
	f, _ := m.Decimal.Round(2).Float64()
	return f
}

func (m Money) Equals(m2 Money) bool {
	return m.Decimal.Equals(m2.Decimal)
}

func (m Money) PercentOf(total Money) decimal.Decimal {
	if total.IsPositive() {
		return m.Div(total.Decimal).Mul(decimal.New(100, 0))
	}
	return decimal.Zero
}

func (m Money) Negate() Money {
	m.Mul(-1)
	return m
}

func (m Money) LessThan(v float64) bool {
	return m.Cmp(decimal.NewFromFloat(v)) == -1
}

func (m Money) MoreThan(v float64) bool {
	return m.Cmp(decimal.NewFromFloat(v)) == 1
}
