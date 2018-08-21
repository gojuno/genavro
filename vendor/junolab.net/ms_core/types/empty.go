package types

type Empty struct{}

func (e Empty) Validate() error {
	return nil
}

func (e Empty) ValidateJSON() error {
	return nil
}
