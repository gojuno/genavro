package utils

//DEPRECATED: don't use it please
func ContainsStr(a []string, s string) bool {
	for _, e := range a {
		if e == s {
			return true
		}
	}
	return false
}
