package parser

func Equal(tkn Token) func(min, max int) MatchFunc {
	return func(min, max int) MatchFunc {
		return func(tokens []Token) (bool, int) {
			max := max
			if max <= 0 || max > len(tokens) {
				max = len(tokens)
			}
			gotHits := 0
			for i := 0; i < max; i++ {
				if tkn != tokens[i] {
					break
				}
				gotHits++
			}
			if gotHits >= min {
				return true, gotHits
			}
			return false, 0
		}
	}
}

func EqualType[T any](min, max int) MatchFunc {
	return func(tokens []Token) (bool, int) {
		max := max
		if max <= 0 || max > len(tokens) {
			max = len(tokens)
		}
		gotHits := 0
		for i := 0; i < max; i++ {
			if _, ok := tokens[i].(T); !ok {
				break
			}
			gotHits++
		}

		if gotHits >= min {
			return true, gotHits
		}
		return false, 0
	}
}

func Or(matchFuncs ...MatchFunc) MatchFunc {
	return func(tokens []Token) (bool, int) {
		for _, f := range matchFuncs {
			if ok, n := f(tokens); ok {
				return ok, n
			}
		}
		return false, 0
	}
}
