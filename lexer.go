package parser

import (
	"bufio"
	"io"
	"unicode/utf8"

	"github.com/pkg/errors"
)

type Token any

type RuneToken struct {
	Rune rune
}

func Lex(r io.Reader) ([]Token, error) {
	var tokens []Token
	reader := bufio.NewReader(r)

	for {
		r, _, err := reader.ReadRune()
		if r != 0 && r != utf8.RuneError {
			tokens = append(tokens, RuneToken{Rune: r})
		}
		if err != nil {
			if errors.Is(err, io.EOF) {
				break
			}
			return tokens, errors.WithStack(err)
		}
	}
	return tokens, nil
}
