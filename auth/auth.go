package auth

import "errors"

var (
	ErrUnauthorized = errors.New("unauthorized")
)

type Checker interface {
	Check(username string, password string) bool
}

// // Validator is used to validate tokens
// type Validator interface {
// 	// Validate returns the username and an error
// 	Validate(value string) (string, error)
// }
