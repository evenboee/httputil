package auth

type MapChecker map[string]string

func (m MapChecker) Check(username string, password string) bool {
	if pass, ok := m[username]; ok {
		return pass == password
	}
	return false
}
