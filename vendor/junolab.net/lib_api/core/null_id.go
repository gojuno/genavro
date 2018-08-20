package core

import (
	"database/sql/driver"

	"github.com/pkg/errors"
)

// Nullable ID for storing to database
type NullID struct {
	Valid bool
	ID    ID
}

// NewNullID creates new NullID from ID
func NewNullID(id ID) NullID {
	return NullID{ID: id, Valid: len(id) > 0}
}

func (ni NullID) String() string {
	if !ni.Valid {
		return "nil"
	}
	return ni.ID.String()
}

func (ni NullID) ToPtr() *ID {
	if ni.Valid {
		return &ni.ID
	}
	return nil
}

func (ni NullID) Value() (driver.Value, error) {
	if ni.Valid {
		return ni.ID.Value()
	}
	return nil, nil
}

func (ni *NullID) Scan(src interface{}) error {
	switch v := src.(type) {
	case nil:
		*ni = NullID{}
		return nil
	default:
		err := ni.ID.Scan(v)
		if err != nil {
			return err
		}
		ni.Valid = true
		return nil
	}
}

// UnmarshalJSON - encoding/json Unmarshaler interface implementation
func (ni *NullID) UnmarshalJSON(data []byte) error {
	if len(data) == 0 || string(data) == "null" {
		ni.ID, ni.Valid = "", false
		return nil
	}

	if err := ni.ID.UnmarshalJSON(data); err != nil {
		return errors.Wrap(err, "failed to unmarshal NullID.ID value")
	}

	ni.Valid = true
	return nil
}

// MarshalJSON - encoding/json marshaler interface implementation
func (ni NullID) MarshalJSON() ([]byte, error) {
	if !ni.Valid {
		return []byte("null"), nil
	}

	return ni.ID.MarshalJSON()
}
