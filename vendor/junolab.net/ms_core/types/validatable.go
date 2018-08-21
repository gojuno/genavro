package types

// Validatable is used for requests defining custom validation rules
type Validatable interface {
	Validate() error
}

type ValidatableJSON interface {
	ValidateJSON() error
}
