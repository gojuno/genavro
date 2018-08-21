package types

import (
	"fmt"

	"github.com/shopspring/decimal"
	lib_decimal "junolab.net/lib_api/decimal"
)

// DEPRECATED: use lib_api/decimal.Decimal
type Decimal struct {
	decimal.Decimal `db:"-"`
	precision       int
	precisionIsSet  bool
}

// DEPRECATED: use lib_api/decimal.Decimal
func NewDecimal(i interface{}) (d Decimal, err error) {
	switch v := i.(type) {
	case string:
		d.Decimal, err = decimal.NewFromString(v)
	case float32:
		d.Decimal = decimal.NewFromFloat(float64(v))
	case float64:
		d.Decimal = decimal.NewFromFloat(v)
	case int:
		d.Decimal = decimal.NewFromFloat(float64(v))
	case int64:
		d.Decimal = decimal.NewFromFloat(float64(v))
	case decimal.Decimal:
		d.Decimal = v
	case lib_decimal.Decimal:
		d.Decimal, err = decimal.NewFromString(v.String())
	default:
		err = fmt.Errorf("Can't make decimal from: %#v", v)
	}

	return
}

// Creates an instance of Decimal from float64 with requested precision p
func newDecimalWithPrecision(f float64, p int) (d Decimal) {
	d, _ = NewDecimal(f)
	d.SetPrecision(p)

	return d
}

// ToLibDecimal converts value to lib_api/decimal.Decimal
func (d *Decimal) ToLibDecimal() lib_decimal.Decimal {
	return lib_decimal.NewFromDecimal(d.Decimal)
}

//SetPrecision sets/updates precision of the decimal instance
func (d *Decimal) SetPrecision(p int) *Decimal {
	d.precision = p
	d.precisionIsSet = true

	return d
}

//Returns a new instance of Decimal with requested precision p
func (d Decimal) WithPrecision(p int) Decimal {
	nd := d
	nd.SetPrecision(p)

	return nd
}

func (d Decimal) String() string {
	return d.StringFixed(int32(d.precision))
}

// MarshalJSON implements the json.Marshaler interface.
func (d Decimal) MarshalJSON() ([]byte, error) {
	if !d.precisionIsSet {
		return nil, fmt.Errorf("Decimal precision must be set before marshalling: %v, precision: %d", d, d.precision)
	}

	return []byte(d.String()), nil
}

// UnmarshalJSON implements the json.Unmarshaler interface.
func (d *Decimal) UnmarshalJSON(data []byte) error {
	if err := d.Decimal.UnmarshalJSON(data); err != nil {
		return err
	}

	d.SetPrecision(-int(d.Decimal.Exponent()))

	return nil
}

func (r Decimal) Validate() error {
	return nil
}

// Scan implements the sql.Scanner interface for database deserialization.
func (d *Decimal) Scan(value interface{}) (err error) {
	if !d.precisionIsSet {
		return fmt.Errorf("Decimal precision must be set before reading from DB")
	}

	//scanner will keep precision of the Decimal that d points to
	oldPrecision := d.precision

	switch v := value.(type) {
	case []byte:
		var s string
		s, err = Unquote(v)
		*d, err = NewDecimal(s)
	default:
		*d, err = NewDecimal(v)
	}

	if err != nil {
		return err
	}

	d.SetPrecision(oldPrecision)

	return nil
}

func (d Decimal) IsZero() bool {
	return d.Decimal.Equals(decimal.Zero)
}

func (d Decimal) LessThan(v float64) bool {
	return d.Cmp(decimal.NewFromFloat(v)) == -1
}

func (d Decimal) MoreThan(v float64) bool {
	return d.Cmp(decimal.NewFromFloat(v)) == 1
}
