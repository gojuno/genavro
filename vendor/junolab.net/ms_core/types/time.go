package types

import (
	"database/sql/driver"
	"fmt"
	"time"
)

const (
	TIME_FORMAT     = time.RFC3339
	USA_TIME_FORMAT = "Mon, _2 Jan 2006 03:04PM"
)

type Time struct {
	time.Time `db:"-"`
	Fmt       string `db:"-"`
}

var dbTimeFormats []string = []string{
	"2006-01-02 15:04:05", //MySQL default datetime representation
	TIME_FORMAT,
}

var ZeroTime Time = Time{
	Time: time.Time{},
	Fmt:  TIME_FORMAT,
}

//Generates time value from various input types: string, time, int64(treated as unix timestamp)
func NewTime(i interface{}) (d Time, err error) {
	var t time.Time
	switch v := i.(type) {
	case string:
		t, err = time.Parse(TIME_FORMAT, v)
		if err != nil {
			return d, err
		}
	case []byte:
		t, err = time.Parse(TIME_FORMAT, string(v))
		if err != nil {
			return d, err
		}
	case int:
		t = time.Unix(int64(v), 0).UTC()
	case int64:
		t = time.Unix(v, 0).UTC()
	case time.Time:
		t = v
	case Time:
		t = v.Time
	}

	return Time{t, TIME_FORMAT}, nil
}

func NewTimeNow() Time {
	return Time{time.Now(), TIME_FORMAT}
}

func NewTimeFromTime(t time.Time) Time {
	return Time{t, TIME_FORMAT}
}

func (t Time) Ms() int64 {
	return t.Time.UnixNano() / 1e6
}

// MarshalJSON implements the json.Marshaler interface.
func (d Time) MarshalJSON() ([]byte, error) {
	return []byte(fmt.Sprintf("%q", d.String())), nil
}

// UnmarshalJSON implements the json.Unmarshaler interface.
func (d *Time) UnmarshalJSON(dateBytes []byte) (err error) {
	strDate, err := Unquote(dateBytes)
	if err != nil {
		return
	}

	d.Time, err = time.Parse(TIME_FORMAT, strDate)
	if d.Fmt == "" {
		d.Fmt = TIME_FORMAT
	}

	return
}

// Scan implements the sql.Scanner interface for database deserialization.
func (t *Time) Scan(value interface{}) (err error) {
	oldFmt := t.Fmt
	if oldFmt == "" {
		oldFmt = TIME_FORMAT
	}

	switch v := value.(type) {
	case nil:
		*t, err = NewTime(time.Time{})
	case time.Time:
		*t, err = NewTime(v)
	default:
		var s string
		s, err = Unquote(v)
		if err != nil {
			return err
		}

		var tt time.Time
		for _, format := range dbTimeFormats {
			tt, err = time.Parse(format, s)
			if err == nil {
				*t, err = NewTime(tt)
				break
			}
		}
	}

	t.Fmt = oldFmt
	t.Time = t.Time.In(time.UTC) //workaround for:  https://github.com/lib/pq/issues/329

	return err
}

func (t Time) String() string {
	var fmt = t.Fmt
	if "" == t.Fmt {
		fmt = TIME_FORMAT
	}

	return t.Time.Format(fmt)
}

// Value implements the driver.Valuer interface for database serialization.
func (t Time) Value() (driver.Value, error) {
	if t.IsZero() {
		return nil, nil
	}

	utc := t.Time.In(time.UTC)

	return utc, nil
}

// UnmarshalText implements the encoding.TextUnmarshaler interface for XML
// deserialization.
func (t *Time) UnmarshalText(text []byte) error {
	oldFmt := t.Fmt
	if t.Fmt == "" {
		t.Fmt = TIME_FORMAT
	}

	str := string(text)

	v, err := NewTime(str)
	*t = v
	if err != nil {
		return fmt.Errorf("Error decoding string %q: %s", str, err)
	}

	t.Fmt = oldFmt

	return nil
}

// MarshalText implements the encoding.TextMarshaler interface for XML
// serialization.
func (t Time) MarshalText() (text []byte, err error) {
	return []byte(t.String()), nil
}

// Between accepts time.Time and types.Time types and checks if
// t is between start and end of time interval
func (t Time) Between(start interface{}, end interface{}) bool {
	var tStart, tEnd time.Time

	switch s := start.(type) {
	case time.Time:
		tStart = s
	case Time:
		tStart = s.Time
	}

	switch e := end.(type) {
	case time.Time:
		tEnd = e
	case Time:
		tEnd = e.Time
	}

	fmt.Println(fmt.Sprintf("%s between %s and %s", t.Time.String(), tStart.String(), tEnd.String()))

	return t.Time.After(tStart) && t.Time.Before(tEnd)
}
