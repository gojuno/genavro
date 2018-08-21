package types

import (
	"database/sql/driver"
	"encoding/json"
	"fmt"

	"github.com/pkg/errors"
)

var EmptyRawMessage = RawMessage([]byte("{}"))

type RawMessage json.RawMessage

// Value implements the sql.Valuer interface for database serialization.
func (m *RawMessage) Value() (driver.Value, error) {
	if m == nil {
		return nil, nil
	}
	return string(*m), nil
}

// Scan implements the sql.Scanner interface for database deserialization.
func (m *RawMessage) Scan(value interface{}) error {
	switch v := value.(type) {
	case nil:
	case string:
		t := RawMessage(v)
		*m = t
	case []byte:
		// important to make bytes copy - this could be a reference to some internal database/sql or driver buffer reused for subsequent rows
		t := RawMessage(string(v))
		*m = t
	default:
		return fmt.Errorf("Not supported RawMessage type: %v", v)
	}
	return nil
}

// UnmarshalJSON - encoding/json Unmarshaler interface implementation
func (m *RawMessage) UnmarshalJSON(data []byte) error {
	t := RawMessage(data)
	*m = t
	return nil
}

// MarshalJSON - encoding/json marshaler interface implementation
func (m RawMessage) MarshalJSON() ([]byte, error) {
	return m, nil
}

func (m *RawMessage) Bytes() []byte {
	return *m
}

// Decode Unmarshal JSON content to struct
func (m *RawMessage) Decode(ptr interface{}) error {
	if m == nil {
		return fmt.Errorf("nil raw message")
	}

	err := json.Unmarshal(m.Bytes(), ptr)
	if err != nil {
		return err
	}
	return nil
}

// NewRawMessage Create RawMessage from arbitrary value
func NewRawMessage(some interface{}) (*RawMessage, error) {
	m := new(RawMessage)
	bytes, err := json.Marshal(some)

	if err != nil {
		return nil, fmt.Errorf("failed to marshal value %v : %v", some, err)
	}

	err = m.UnmarshalJSON(bytes)
	if err != nil {
		return nil, err
	}

	return m, nil
}

func (m RawMessage) String() string {
	return fmt.Sprintf("%s", string(m))
}

// Validate checks if RawMessage contains valid JSON.
// Unfortunately, encoding/json package does not export validation method (json.checkValid).
// Call json.Unmarshal with "nil" as param to do only JSON validation, but avoid decoding due to performance optimization.
//
// BenchmarkRawMessage_Validate_Nil-4               1000000              1649 ns/op             328 B/op          4 allocs/op
// BenchmarkRawMessage_Validate_Interface-4          200000              6952 ns/op            1622 B/op         41 allocs/op
func (m RawMessage) Validate() error {
	if err := json.Unmarshal([]byte(m), nil); err != nil {
		// Check error for json.InvalidUnmarshalError with specific value
		if jsonErr, ok := err.(*json.InvalidUnmarshalError); ok && jsonErr.Type == nil {
			return nil
		}
		return errors.Wrap(err, "RawMessage contains invalid JSON")
	}
	return nil
}
