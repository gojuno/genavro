package types

import "fmt"

type Phone struct {
	Num         string `json:"num"`
	CountryCode string `json:"country_code"`
}

func (p Phone) IsValid() bool {
	return p.Num != "" && p.CountryCode != ""
}

func (p Phone) Validate() error {
	if !p.IsValid() {
		return fmt.Errorf("invalid phone")
	}

	if err := CheckValidRunes("+- ()0123456789", p.CountryCode); err != nil {
		return fmt.Errorf("invalid country code: %v", err)
	}

	if err := CheckValidRunes("- ()0123456789", p.Num); err != nil {
		return fmt.Errorf("invalid number: %v", err)
	}

	return nil
}
