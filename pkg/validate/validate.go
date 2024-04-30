package validate

import (
	"errors"
	"net/mail"
	"regexp"
	"unicode/utf8"
)

const (
	minNameLen     = 1
	maxNameLen     = 50
	minPasswordLen = 1
	maxPasswordLen = 100
)

var (
	ErrInvalidName     = errors.New("invalid name")
	ErrInvalidEmail    = errors.New("invalid email")
	ErrInvalidPassword = errors.New("invalid password")
)

func Name(n string) error {
	length := utf8.RuneCountInString(n)
	re := regexp.MustCompile(`^[a-zA-Z-' ]*$`)

	if !(minNameLen < length && length < maxNameLen && re.MatchString(n)) {
		return ErrInvalidName
	}

	return nil
}

func Email(e string) error {
	_, err := mail.ParseAddress(e)
	if err != nil {
		return ErrInvalidEmail
	}

	return nil
}

func Password(p string) error {
	length := utf8.RuneCountInString(p)
	if minPasswordLen > length || length > maxPasswordLen {
		return ErrInvalidPassword
	}

	return nil
}
