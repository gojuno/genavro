package astparser

import "errors"

type Config struct {
	InputDir      string
	ExcludeRegexp string
	IncludeRegexp string
}

func (c *Config) validate() error {
	if c.IncludeRegexp != "" && c.ExcludeRegexp != "" {
		return errors.New("both include and exclude regexps are set")
	}

	return nil
}

func (c *Config) prepare() error {
	if err := c.validate(); err != nil {
		return err
	}

	if c.InputDir == "" {
		c.InputDir = "./"
	}

	return nil
}
