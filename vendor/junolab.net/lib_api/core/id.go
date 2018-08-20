package core

import (
	"database/sql/driver"
	"encoding/json"
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/pborman/uuid"
	pkg_errors "github.com/pkg/errors"
)

const (
	// InvalidGroup is a group of users that have an invalid identifer.
	InvalidGroup = -1
)

// ID represents an identifier
type ID string

const (
	internalUUIDPrefix            string = "00000000-0000-0000-0000-000"
	internalUUIDPrefixCustomerBot string = "00000000-0000-0000-0000-0004"
	InternalUUIDPrefixProviderBot string = "00000000-0000-0000-0000-0003"
	EmptyID                       ID     = ""
)

// NewID creates new uuid4 id
func NewID() ID {
	return ID(uuid.New())
}

// ValidateOptional - validating implementation
func (id ID) ValidateOptional() error {
	if id == "" {
		return nil
	}
	if uuid.Parse(id.String()) == nil {
		return fmt.Errorf("id must be valid UUID: %v", id)
	}
	return nil
}

// IsEmpty checks if id is empty
func (id ID) IsEmpty() bool {
	return id == EmptyID
}

func Parse(id string) ID {
	return ID(uuid.Parse(id).String())
}

// NewInternalID generates ID for internal use.
// The ID has the following format: 00000000-0000-0000-0000-000yxxxxxxxx
// y - is securityGroupFlag (0 - all, 1 - providers, 2 - customers)
// xxxxxxxx - some random value generated by uuid.New()
func NewInternalID(securityGroupFlag int) ID {
	id := uuid.New()
	internalUUIDGroupIndex := len(internalUUIDPrefix) + 1
	sgFlagStr := strconv.FormatInt(int64(securityGroupFlag), 16)
	internalID := internalUUIDPrefix + sgFlagStr + id[internalUUIDGroupIndex:]
	return ID(internalID)
}

// WithTime adds time to the end of the uuid.
// xxxxxxxx-xxxx-xxxx-xxxx-[xxxxxxxxxxxx]
// [xxxxxxxxxxxx]::int + time.Unix()
// return xxxxxxxx-xxxx-xxxx-xxxx-uint::hex
func (id ID) WithTime(t time.Time) (ID, error) {
	newID, err := id.addUint64(uint64(t.Unix()))
	return newID, err
}

// AbTest used for selecting given part of identifiers for A/B testing.
func (id ID) AbTest(percents int) bool {
	if percents <= 0 {
		return false
	}
	if percents >= 100 {
		return true
	}
	value, err := id.extractLastPartAsUint()
	if err != nil {
		return false
	}
	return value%100 < uint64(percents)
}

// InABTestRange used for selecting given part of identifiers in [from..to) range for A/B testing.
// 'from' and 'to' should be percent values within [0..100] range each.
// Ranges "crossing" 100% boundary are not supported - i.e. it is not possible to check e.g. for [95..5) range.
func (id ID) InABTestRange(from, to int) bool {
	v := id.ABTestGroup()
	return from <= v && v < to
}

// ABTestGroup returns a group number.
//
// Let identifiers be mapped to the multiplicative group of integers
// modulo 100, this method returns a group number the identifier belongs
// to.
//
// When the identifier is invalid, method returns -1 value.
func (id ID) ABTestGroup() int {
	value, err := id.extractLastPartAsUint()
	if err != nil {
		return InvalidGroup
	}
	return int(value % 100)
}

// ParseID converts string to ID, returning false if string is not a valid UUID.
func ParseID(s string) (id ID, ok bool) {
	uuidVal := uuid.Parse(s)
	if len(uuidVal) == 0 {
		return "", false
	}
	return ID(uuidVal.String()), true
}

// ResetChunkCounter resets last UUID part (as delimited by -) to all zeroes. Can be used to implement counters (i.e. incrementing UUIDs)
func (id *ID) ResetChunkCounter() {
	parts := strings.Split(string(*id), "-")
	if len(parts) == 0 {
		return
	}

	parts[len(parts)-1] = "000000000000"
	*id = ID(strings.Join(parts, "-"))
}

// Inc increments last UUID part (as delimited by -). Can be used to implement counters (i.e. incrementing UUIDs).
func (id *ID) Inc() {
	newID, err := id.addUint64(uint64(1))
	if err != nil {
		return
	}
	*id = newID
}

// validate - validating implementation
func (id ID) validate() error {
	if uuid.Parse(id.String()) == nil {
		return fmt.Errorf("id must be valid UUID: %v", id)
	}
	return nil
}

func (id ID) Validate() error {
	return id.validate()
}

func (id ID) ValidateJSON() error {
	return id.validate()
}

// Eq defines ID equality
func (id ID) Eq(other ID) bool {
	otherUUID := uuid.Parse(other.String())
	currentUUID := uuid.Parse(id.String())
	return uuid.Equal(currentUUID, otherUUID)
}

// String returns string representation of ID
func (id ID) String() string {
	return string(id)
}

// UnmarshalJSON - encoding/json Unmarshaler interface implementation
func (id *ID) UnmarshalJSON(data []byte) error {
	var idStr string
	if err := json.Unmarshal(data, &idStr); err != nil {
		return err
	}

	value := ID(idStr)
	if err := value.ValidateOptional(); err != nil {
		return err
	}
	*id = value
	return nil
}

// MarshalJSON - encoding/json marshaler interface implementation
func (id ID) MarshalJSON() ([]byte, error) {
	if err := id.ValidateOptional(); err != nil {
		return nil, err
	}

	buf := make([]byte, 0, len(id)+2)
	buf = append(buf, '"')
	buf = append(buf, id...)
	buf = append(buf, '"')

	return buf, nil
}

// Value implements driver.Valuer interface
func (id ID) Value() (driver.Value, error) {
	return id.String(), nil
}

// Scan implements the sql.Scanner interface for database deserialization.
func (id *ID) Scan(value interface{}) error {
	var scannedID ID
	switch v := value.(type) {
	case nil:
		scannedID = ""
		//return errors.New("null ID value not supported")
	case string:
		scannedID = ID(v)
	case []byte:
		scannedID = ID(string(v))
	default:
		return fmt.Errorf("not supported ID type: %v", v)
	}

	if err := scannedID.ValidateOptional(); err != nil {
		return err
	}
	*id = scannedID
	return nil
}

// IsInternal checks whether ID is Internal or not.
// The ID is Internal if it starts with 00000000-0000-0000-0000-000
func (id ID) IsInternal() bool {
	return strings.HasPrefix(string(id), internalUUIDPrefix)
}

// IsInternalCustomerBot checks whether given ID is for internal customer bot.
func (id ID) IsInternalCustomerBot() bool {
	return strings.HasPrefix(string(id), internalUUIDPrefixCustomerBot)
}

// ToNullID converts ID to nullable NullID value
func (id *ID) ToNullID() NullID {
	if id == nil || *id == EmptyID {
		return NullID{}
	}
	return NullID{Valid: true, ID: *id}
}

func (id ID) addUint64(delta uint64) (ID, error) {
	idStr := id.String()
	lastPartIndex := strings.LastIndex(idStr, "-")
	lastPartUint, err := strconv.ParseUint(idStr[lastPartIndex+1:], 16, 64)
	if err != nil {
		return id, pkg_errors.Wrapf(err, "failed to parse last part of id %s to uint", id.String())
	}
	lastPartUint += delta
	return ID(idStr[:lastPartIndex+1] + fmt.Sprintf("%012x", lastPartUint)), nil
}

func (id ID) extractLastPartAsUint() (uint64, error) {
	idStr := id.String()
	lastPartIndex := strings.LastIndex(idStr, "-")
	lastPartUint, err := strconv.ParseUint(idStr[lastPartIndex+1:], 16, 64)
	return lastPartUint, err
}

func isHexDigit(c rune) bool {
	switch {
	case '0' <= c && c <= '9':
		return true
	case 'a' <= c && c <= 'f':
		return true
	case 'A' <= c && c <= 'F':
		return true
	default:
		return false
	}
}