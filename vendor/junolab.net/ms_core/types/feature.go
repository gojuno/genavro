package types

type (
	Feature  string
	Features []Feature
)

func (f Feature) String() string {
	return string(f)
}

func (fl Features) IsFeatureSupported(feature Feature) bool {
	for _, f := range fl {
		if f == feature {
			return true
		}
	}
	return false
}

func (fl Features) Strings() []string {
	s := make([]string, len(fl))
	for i, f := range fl {
		s[i] = f.String()
	}
	return s
}
