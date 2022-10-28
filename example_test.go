package parser_test

import (
	"bytes"
	"fmt"
	"strings"

	"github.com/Eun/go-parser"
)

func ExampleReplaceTokens() {
	// this example parses the git commit message

	tokens, err := parser.Lex(bytes.NewReader([]byte(`
# Please enter the commit message for your changes. Lines starting
# with '#' will be ignored, and an empty message aborts the commit.
1. Implemented awesome feature
2. Fixed a nasty bug
# not sure about this:
3. Optimization
`)))
	if err != nil {
		panic(err)
	}

	// 1. replace '\n' RuneTokens with a NewLineToken
	// this makes matching easier

	chain := []parser.MatchFunc{
		parser.Equal(parser.RuneToken{Rune: '\n'})(1, 1),
	}

	type NewLineToken struct{}

	tokens, err = parser.ReplaceTokens(tokens, chain, func(tokens []parser.Token) ([]parser.Token, error) {
		return []parser.Token{NewLineToken{}}, nil
	})
	if err != nil {
		panic(err)
	}

	// 2. replace all lines start with '#'
	chain = []parser.MatchFunc{
		parser.Equal(parser.RuneToken{Rune: '#'})(1, 0),
		parser.EqualType[parser.RuneToken](0, 0),
		parser.Equal(NewLineToken{})(1, 1),
	}

	type Comment struct {
		Text string
	}

	tokens, err = parser.ReplaceTokens(tokens, chain, func(tokens []parser.Token) ([]parser.Token, error) {
		// we could filter out the first '#' here
		// but let's skip this for simplicity
		var sb strings.Builder
		for _, token := range tokens {
			if t, ok := token.(parser.RuneToken); ok {
				sb.WriteRune(t.Rune)
			}
		}
		return []parser.Token{
			Comment{Text: sb.String()},
		}, nil
	})
	if err != nil {
		panic(err)
	}

	// 3. Everything that is left is the commit
	var sb strings.Builder
	for _, token := range tokens {
		switch t := token.(type) {
		case parser.RuneToken:
			sb.WriteRune(t.Rune)
		case NewLineToken:
			sb.WriteRune('\n')
		}
	}
	fmt.Println(strings.TrimSpace(sb.String()))
}
