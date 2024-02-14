package auth

type MapChecker map[string]string

func (m MapChecker) Check(username string, password string) bool {
	if pass, ok := m[username]; ok {
		return pass == password
	}
	return false
}

type checker struct {
	f func(string, string) bool
}

func (c checker) Check(username, password string) bool {
	return c.f(username, password)
}

func CheckerFrom(f func(string, string) bool) Checker {
	return checker{f: f}
}
