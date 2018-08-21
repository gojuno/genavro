package types

import (
	"github.com/hashicorp/go-version"
)

type (
	AppVersion string
)

func (s AppVersion) String() string {
	return string(s)
}

func (s AppVersion) Less(appVersion AppVersion) bool {
	versions, err := parseVersions(s, appVersion)
	if err != nil {
		return false
	}
	return versions[0].LessThan(versions[1])
}

func (s AppVersion) Equal(appVersion AppVersion) bool {
	versions, err := parseVersions(s, appVersion)
	if err != nil {
		return false
	}
	return versions[0].Equal(versions[1])
}

func (s AppVersion) Validate() error {
	_, err := version.NewVersion(s.String())
	return err
}

func parseVersions(appVersions ...AppVersion) ([]*version.Version, error) {
	versions := make([]*version.Version, len(appVersions))
	for i, appVer := range appVersions {
		ver, err := version.NewVersion(appVer.String())
		if err != nil {
			return nil, err
		}
		versions[i] = ver
	}
	return versions, nil
}
