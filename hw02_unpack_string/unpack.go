package hw02unpackstring

import (
	"errors"
	"strings"
)

var ErrInvalidString = errors.New("некорректная строка")

func Unpack(packedString string) (string, error) {
	if len(packedString) == 0 {
		return "", nil
	}

	resultStringBuilder := strings.Builder{}
	p := new(parser)
	packedRunes := []rune(packedString)
	// добавим в конец пустую руну, чтобы парсер высвободил, если что-то накопилось
	packedRunes = append(packedRunes, emptyRune)
	for _, r := range packedRunes {
		unpackedPart, err := p.handleNextRune(r)
		if err != nil {
			return "", err
		}

		resultStringBuilder.WriteString(unpackedPart)
	}

	return resultStringBuilder.String(), nil
}
