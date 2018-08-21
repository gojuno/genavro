package types

import (
	"bytes"
	"database/sql/driver"
	"fmt"
	"math/big"
	"strings"
)

const _MASK_SIZE = 128

type BitMask struct {
	*big.Int
}

func NewBitMask() BitMask {
	return BitMask{Int: &big.Int{}}
}

func (b1 BitMask) Match(b2 BitMask) bool {
	if b1.Int == nil || b2.Int == nil {
		if (b1.Int == nil && b2.Int == nil) || (b1.Int.BitLen() == 0 && b2.Int.BitLen() == 0) {
			return true
		}
		return false
	}
	result := NewBitMask()
	result.And(b1.Int, b2.Int)
	return result.Cmp(&big.Int{}) == 1
}

// UnmarshalJSON - encoding/json Unmarshaler interface implementation
func (bitmask *BitMask) UnmarshalJSON(data []byte) error {
	bitmask.Int = &big.Int{}
	s := strings.Replace(string(data), "\"", "", 2)
	if r, ok := bitmask.SetString(s, 10); !ok {
		return fmt.Errorf("failed to unmarshall json %s to bitmask", data)
	} else {
		bitmask.Int = r
		return nil
	}
}

func (bitmask BitMask) IsZero() bool {
	return bitmask.Int == nil || bitmask.Int.BitLen() == 0
}

// MarshalJSON - encoding/json Marshaler interface implementation
func (bitmask BitMask) MarshalJSON() ([]byte, error) {
	if bitmask.Int == nil {
		bitmask.Int = &big.Int{}
	}
	return []byte(fmt.Sprintf("%q", bitmask.Int.String())), nil
}

func (bitmask BitMask) Value() (driver.Value, error) {
	return fillMaskToLength(fmt.Sprintf("%b", bitmask)), nil
}

// Scan implements the sql.Scanner interface for database deserialization.
func (bitmask *BitMask) Scan(value interface{}) error {
	bigint := &big.Int{}
	switch v := value.(type) {
	case string:
		str := strings.Replace(v, "\"", "", 2)
		if _, ok := bigint.SetString(str, 2); !ok {
			return fmt.Errorf("failed to scan bitmask %v", v)
		}
	case []byte:
		str := strings.Replace(string(v), "\"", "", 2)
		if _, ok := bigint.SetString(str, 2); !ok {
			return fmt.Errorf("failed to scan bitmask %q", v)
		}
	default:
		return fmt.Errorf("not supported bitmask type %v", v)
	}
	bitmask.Int = bigint

	return nil
}

func fillMaskToLength(bitmask string) string {
	l := len(bitmask)
	if l == _MASK_SIZE {
		return bitmask
	} else if l < _MASK_SIZE {
		buf := bytes.NewBufferString("")
		for i := l; i < _MASK_SIZE; i += 1 {
			buf.WriteString("0")
		}
		buf.WriteString(bitmask)
		return buf.String()
	}

	panic(fmt.Sprintf("there should never happend. mask default len %d current %d mask %q", _MASK_SIZE, l, bitmask))
}

func GenerateMask(position int) string {
	s := ""
	for i := _MASK_SIZE; i > 0; i -= 1 {
		if i == position {
			s += "1"
		} else {
			s += "0"
		}
	}
	return s
}

func ValidateMaskLen(maskLen int) error {
	if maskLen > _MASK_SIZE {
		return fmt.Errorf("mask size %d too big", maskLen)
	}
	return nil
}
