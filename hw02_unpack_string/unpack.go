package hw02unpackstring

import (
	"errors"
	"strconv"
	"strings"
	"unicode"
)

var ErrInvalidString = errors.New("invalid string")

func Unpack(str string) (string, error) {
	var (
		result      strings.Builder
		prevChar    rune
		lastStrIndx int = len(str) - 1
	)

	for i, char := range str {
		isFirstIteration := i == 0
		if (isFirstIteration && unicode.IsDigit(char)) || (unicode.IsDigit(prevChar) && unicode.IsDigit(char)) {
			return "", ErrInvalidString
		}

		if unicode.IsDigit(char) {
			howMuchRepeat, err := strconv.Atoi(string(char))
			if err != nil {
				return "", err
			}

			result.WriteString(strings.Repeat(string(prevChar), howMuchRepeat))

		} else if !isFirstIteration && !unicode.IsDigit(prevChar) && !unicode.IsDigit(char) {
			result.WriteString(string(prevChar))
		}

		prevChar = char

		if lastStrIndx == i && !unicode.IsDigit(char) {
			result.WriteString(string(char))
		}
	}

	return result.String(), nil
}
