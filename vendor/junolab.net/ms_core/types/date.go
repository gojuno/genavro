package types

import (
	"fmt"
	"time"
)

const (
	DATE_FORMAT = "2006-01-02"
)

type Date struct {
	Time `db:"-"`
}

//Generates date value from various input types:
//string, time, int(timestamp), Date, Time
func NewDate(i interface{}) (d Date, err error) {
	var t Time
	var tt time.Time

	switch v := i.(type) {
	case Date:
		t = v.Time
	case string:
		tt, err = time.Parse(DATE_FORMAT, v)
		t = Time{tt, DATE_FORMAT}
	default:
		t, err = NewTime(v)
	}

	if err != nil {
		return
	}

	t.Fmt = DATE_FORMAT

	return Date{t}, nil
}

// UnmarshalJSON implements the json.Unmarshaler interface.
func (d *Date) UnmarshalJSON(dateBytes []byte) (err error) {
	strDate, err := Unquote(dateBytes)
	if err != nil {
		return
	}

	t, err := time.Parse(DATE_FORMAT, strDate)
	if err != nil {
		return
	}

	oldFmt := d.Fmt
	if oldFmt == "" {
		oldFmt = DATE_FORMAT
	}

	(*d).Time = Time{t, oldFmt}

	return
}

//Scan implements the sql.Scanner interface for database deserialization.
func (d *Date) Scan(value interface{}) (err error) {
	oldFmt := d.Fmt
	if oldFmt == "" {
		oldFmt = DATE_FORMAT
	}

	switch v := value.(type) {
	case nil:
		*d, err = NewDate(time.Time{})
	case time.Time:
		*d, err = NewDate(v)
	default:
		var s string
		s, err = Unquote(v)
		if err != nil {
			return
		}
		*d, err = NewDate(s)
	}

	d.Fmt = oldFmt

	return err
}

// UnmarshalText implements the encoding.TextUnmarshaler interface for XML
// deserialization.
func (d *Date) UnmarshalText(text []byte) error {
	oldFmt := d.Fmt
	if oldFmt == "" {
		oldFmt = DATE_FORMAT
	}

	str := string(text)

	v, err := NewDate(str)
	*d = v
	if err != nil {
		return fmt.Errorf("Error decoding string %q: %s", str, err)
	}

	d.Fmt = oldFmt

	return nil
}
