package parser

import "github.com/pkg/errors"

type MatchFunc func([]Token) (bool, int)
type ReplaceFunc func([]Token) ([]Token, error)

func ReplaceTokens(tokens []Token, chain []MatchFunc, replace ReplaceFunc) ([]Token, error) {
	result := make([]Token, 0, len(tokens))

	var ci int // chain index
	var si int // start index of matching tokens
	var ti int // current token index

	gotMatch := func() error {
		var err error
		// we got a match for all tokens[si:ti]
		oldSize := len(tokens)

		tokensToAdd := tokens
		diff := len(tokensToAdd) - oldSize
		tokensToAdd = tokensToAdd[si : ti+diff]
		if replace != nil {
			tokensToAdd, err = replace(tokensToAdd)
			if err != nil {
				return errors.WithStack(err)
			}
		}
		result = append(result, tokensToAdd...)
		return nil
	}

	for ti < len(tokens) {
		if ok, advance := chain[ci](tokens[ti:]); ok {
			ti += advance
			ci++
			if ci == len(chain) {
				// got a match
				if err := gotMatch(); err != nil {
					return result, errors.WithStack(err)
				}
				ci = 0
				si = ti
			}
			continue
		}
		ti = si
		result = append(result, tokens[ti])
		ci = 0
		ti++
		si = ti
	}

	if ci == 0 {
		return result, nil
	}
	// tokens are done, but the chain is not over yet
	// check if the chain is still solvable
	for ; ci < len(chain); ci++ {
		if ok, _ := chain[ci](nil); !ok {
			return append(result, tokens[si:]...), nil
		}
	}

	// got a match
	if err := gotMatch(); err != nil {
		return result, errors.WithStack(err)
	}
	return result, nil
}
