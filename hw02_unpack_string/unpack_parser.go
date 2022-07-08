package hw02unpackstring

import (
	"strings"
	"unicode"
)

const (
	escapeRune = '\\'
	emptyRune  = 0
)

type parser struct {
	currentRuneToMultiply rune
	isEscapeMode          bool
}

func (p *parser) handleNextRune(r rune) (string, error) {
	returnValue := ""
	var err error

	switch {
	case p.isEscapeMode:
		defer func() { p.isEscapeMode = false }()
		if p.isEscapable(r) {
			p.currentRuneToMultiply = r
		} else {
			err = ErrInvalidString
		}
	case p.currentRuneToMultiply != emptyRune:
		multiplier := 1
		if unicode.IsDigit(r) {
			multiplier = int(r - '0')
		} else {
			defer func() { p.handleNextRune(r) }()
		}

		returnValue = p.release(multiplier)
	case p.isEscapeRune(r):
		p.isEscapeMode = true
	case !unicode.IsDigit(r):
		p.currentRuneToMultiply = r
	default:
		err = ErrInvalidString // сразу встретили цифру, возвращаем ошибку
	}

	return returnValue, err
}

func (p *parser) release(multiplier int) string {
	defer func() {
		p.currentRuneToMultiply = emptyRune
		p.isEscapeMode = false
	}()

	return strings.Repeat(string(p.currentRuneToMultiply), multiplier)
}

func (p parser) isEscapable(r rune) bool {
	return p.isEscapeRune(r) || unicode.IsDigit(r)
}

func (p parser) isEscapeRune(r rune) bool {
	return r == escapeRune
}
