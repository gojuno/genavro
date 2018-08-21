package types

import (
	"database/sql/driver"
	"fmt"
	"strings"
)

type (
	AppID       string
	Application string
)

const (
	AppIdVersionDelimiter = ":"
)

var (
	ErrInvalidAppIDFormat      = fmt.Errorf("invalid app id format")
	ErrInvalidAppVersionFormat = fmt.Errorf("invalid app version format")
)

func (s AppID) String() string {
	return string(s)
}

func (s AppID) WithAppVersion(appVersion AppVersion) Application {
	if s.String() == "" && appVersion.String() == "" {
		return Application("")
	}
	return Application(strings.Join([]string{s.String(), appVersion.String()}, AppIdVersionDelimiter))
}

func (s AppID) Validate() error {
	if len(s) == 0 {
		return ErrInvalidAppIDFormat
	}
	return nil
}

func (s Application) GetAppID() (AppID, error) {
	subs := strings.Split(s.String(), AppIdVersionDelimiter)
	if len(subs) == 2 {
		return AppID(subs[0]), nil
	}
	return "", ErrInvalidAppIDFormat
}

func (s Application) GetAppVersion() (AppVersion, error) {
	subs := strings.Split(s.String(), AppIdVersionDelimiter)
	if len(subs) == 2 {
		return AppVersion(subs[1]), nil
	}
	return "", ErrInvalidAppVersionFormat
}

func (s Application) String() string {
	return string(s)
}

func (s *Application) Scan(value interface{}) error {
	switch v := value.(type) {
	case nil:
		*s = ""
	case string:
		*s = Application(v)
	case []byte:
		*s = Application(string(v))
	default:
		return fmt.Errorf("failed to scan application, %v", value)
	}
	return nil
}

// Value implements the driver Valuer interface.
func (s Application) Value() (driver.Value, error) {
	return s.String(), nil
}
