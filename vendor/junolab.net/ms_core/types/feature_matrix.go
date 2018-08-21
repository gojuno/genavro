package types

import "junolab.net/lib_api/core"

type (
	FeatureMatrix map[AppID]map[Feature]FeatureEntry

	ABTestFunc func(core.ID) bool

	AppIDTree map[AppID]AppID

	FeatureEntry struct {
		AppVersion AppVersion
		VersionCmp int
		ABTestFunc ABTestFunc
	}
)

const (
	GreaterOrEqualVersion = iota
	LessVersion
	EqualVersion
	NotEqualVersion
)

func (fm FeatureMatrix) hasFeature(feature Feature, userID core.ID, appID AppID, appVersion AppVersion, appIDTree AppIDTree) bool {
	featuresMap, ok := fm[appID]
	if ok {
		entry, ok := featuresMap[feature]
		if ok {
			if appID.Validate() != nil || appVersion.Validate() != nil || entry.AppVersion.Validate() != nil {
				return false
			}
			switch entry.VersionCmp {
			case NotEqualVersion:
				if entry.AppVersion == appVersion {
					return false
				}
			case EqualVersion:
				if entry.AppVersion != appVersion {
					return false
				}
			case GreaterOrEqualVersion:
				if appVersion.Less(entry.AppVersion) {
					return false
				}
			case LessVersion:
				if entry.AppVersion.Less(appVersion) || entry.AppVersion == appVersion {
					return false
				}
			default:
				return false
			}
			if entry.ABTestFunc != nil && !entry.ABTestFunc(userID) {
				return false
			}
			return true
		}
	}

	appID, ok = appIDTree[appID]
	if ok {
		// try to find in parent app id.
		return fm.hasFeature(feature, userID, appID, appVersion, appIDTree)
	}
	return false
}

// ExtractUserFeatures gets all supported features for user.
func (fm FeatureMatrix) ExtractUserFeatures(userID core.ID, appID AppID, appVersion AppVersion, appIDTree AppIDTree) Features {
	features := Features{}
	featuresMap, ok := fm[appID]
	if ok {
		for f := range featuresMap {
			if fm.hasFeature(f, userID, appID, appVersion, appIDTree) {
				features = append(features, f)
			}
		}
	}
	appID, ok = appIDTree[appID]
	if ok {
		// try to find in parent app id.
		return append(features, fm.ExtractUserFeatures(userID, appID, appVersion, appIDTree)...)
	}
	return features
}

func (fm FeatureMatrix) IsFeatureSupported(feature Feature, userID core.ID, appID AppID, appVersion AppVersion, appIDTree AppIDTree) bool {
	return fm.hasFeature(feature, userID, appID, appVersion, appIDTree)
}
