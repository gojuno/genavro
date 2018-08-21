package types

import (
	"database/sql/driver"
	"errors"
	"fmt"
	"net/mail"
	"strings"
	"unicode/utf8"
)

type Email string

// IsInternal checks whether Email is internal or not.
// Email is internal it ends with one of provided internal domains
func (email Email) IsInternal(internalDomains []string) bool {
	for _, d := range internalDomains {
		if strings.HasSuffix(string(email), "@"+d) {
			return true
		}
	}
	return false
}

func (email Email) Eq(other Email) bool {
	return email == other
}

func (email Email) String() string {
	return string(email)
}

// Value implements the driver Valuer interface.
func (email Email) Value() (driver.Value, error) {
	return email.String(), nil
}

// Scan implements the sql.Scanner interface for database deserialization.
func (email *Email) Scan(value interface{}) error {
	switch v := value.(type) {
	case nil:
		*email = ""
	case string:
		*email = Email(v)
	case []byte:
		*email = Email(string(v))
	default:
		return errors.New(fmt.Sprintf("Can't convert: %v to Email", v))
	}

	return nil
}

func (r Email) Validate() error {
	if r.String() == "" {
		return nil
	}

	return r.ValidateNonEmpty()
}

func (r Email) ValidateNonEmpty() error {
	if utf8.RuneCountInString(r.String()) < 5 {
		return fmt.Errorf("min len of email is 5: %v", r.String())
	}
	if utf8.RuneCountInString(r.String()) > 128 {
		return fmt.Errorf("max len of email is 128: %v", r.String())
	}

	email := r.String()
	if _, err := mail.ParseAddress(email); err != nil {
		return fmt.Errorf("email '%s' is not valid: %v", email, err)
	}
	return nil
}

func (email Email) Sanitize() Email {
	return Email(strings.ToLower(strings.TrimSpace(string(email))))
}
