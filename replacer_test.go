package parser

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/require"
)

type TextToken struct {
	Text string
}

func replaceRuneTokensWithTextToken(tokens []Token) ([]Token, error) {
	var sb strings.Builder
	for _, token := range tokens {
		switch t := token.(type) {
		case RuneToken:
			sb.WriteRune(t.Rune)
		case TextToken:
			sb.WriteString(t.Text)
		}
	}

	return []Token{
		TextToken{Text: sb.String()},
	}, nil
}

func Test_Equal(t *testing.T) {
	tests := []struct {
		name        string
		args        []Token
		matchFuncs  []MatchFunc
		replaceFunc ReplaceFunc
		want        []Token
		wantErr     bool
	}{
		{
			name: "Exact Match",
			args: []Token{
				RuneToken{'1'},
				RuneToken{'2'},
				RuneToken{'3'},
				RuneToken{'4'},
			},
			matchFuncs: []MatchFunc{
				Equal(RuneToken{'1'})(1, 1),
				Equal(RuneToken{'2'})(1, 1),
				Equal(RuneToken{'3'})(1, 1),
				Equal(RuneToken{'4'})(1, 1),
			},
			replaceFunc: replaceRuneTokensWithTextToken,
			want: []Token{
				TextToken{Text: "1234"},
			},
			wantErr: false,
		},
		{
			name: "Double Match",
			args: []Token{
				RuneToken{'1'},
				RuneToken{'2'},
				RuneToken{'3'},
				RuneToken{'4'},
				RuneToken{'1'},
				RuneToken{'2'},
				RuneToken{'3'},
				RuneToken{'4'},
			},
			matchFuncs: []MatchFunc{
				Equal(RuneToken{'1'})(1, 1),
				Equal(RuneToken{'2'})(1, 1),
				Equal(RuneToken{'3'})(1, 1),
				Equal(RuneToken{'4'})(1, 1),
			},
			replaceFunc: replaceRuneTokensWithTextToken,
			want: []Token{
				TextToken{Text: "1234"},
				TextToken{Text: "1234"},
			},
			wantErr: false,
		},
		{
			name: "Longer Chain than Input",
			args: []Token{
				RuneToken{'1'},
			},
			matchFuncs: []MatchFunc{
				Equal(RuneToken{'1'})(1, 1),
				Equal(RuneToken{'2'})(0, 1),
			},
			replaceFunc: replaceRuneTokensWithTextToken,
			want: []Token{
				TextToken{Text: "1"},
			},
			wantErr: false,
		},
		{
			name: "Match with Other Tokens Between",
			args: []Token{
				RuneToken{'1'},
				RuneToken{'2'},
				RuneToken{'3'},
				RuneToken{'4'},
				RuneToken{'5'},
				RuneToken{'1'},
				RuneToken{'2'},
				RuneToken{'3'},
				RuneToken{'4'},
			},
			matchFuncs: []MatchFunc{
				Equal(RuneToken{'1'})(1, 1),
				Equal(RuneToken{'2'})(1, 1),
				Equal(RuneToken{'3'})(1, 1),
				Equal(RuneToken{'4'})(1, 1),
			},
			replaceFunc: replaceRuneTokensWithTextToken,
			want: []Token{
				TextToken{Text: "1234"},
				RuneToken{Rune: '5'},
				TextToken{Text: "1234"},
			},
			wantErr: false,
		},
		{
			name: "Optional - Token exists",
			args: []Token{
				RuneToken{'1'},
				RuneToken{'2'},
				RuneToken{'3'},
				RuneToken{'4'},
			},
			matchFuncs: []MatchFunc{
				Equal(RuneToken{'1'})(1, 1),
				Equal(RuneToken{'2'})(0, 1),
				Equal(RuneToken{'3'})(1, 1),
				Equal(RuneToken{'4'})(1, 1),
			},
			replaceFunc: replaceRuneTokensWithTextToken,
			want: []Token{
				TextToken{Text: "1234"},
			},
			wantErr: false,
		},
		{
			name: "Optional - Token does not exist",
			args: []Token{
				RuneToken{'1'},
				RuneToken{'3'},
				RuneToken{'4'},
			},
			matchFuncs: []MatchFunc{
				Equal(RuneToken{'1'})(1, 1),
				Equal(RuneToken{'2'})(0, 1),
				Equal(RuneToken{'3'})(1, 1),
				Equal(RuneToken{'4'})(1, 1),
			},
			replaceFunc: replaceRuneTokensWithTextToken,
			want: []Token{
				TextToken{Text: "134"},
			},
			wantErr: false,
		},
		{
			name: "No max limit - One occurrence",
			args: []Token{
				RuneToken{'1'},
				RuneToken{'2'},
				RuneToken{'3'},
				RuneToken{'4'},
			},
			matchFuncs: []MatchFunc{
				Equal(RuneToken{'1'})(1, 1),
				Equal(RuneToken{'2'})(1, 0),
				Equal(RuneToken{'3'})(1, 1),
				Equal(RuneToken{'4'})(1, 1),
			},
			replaceFunc: replaceRuneTokensWithTextToken,
			want: []Token{
				TextToken{Text: "1234"},
			},
			wantErr: false,
		},
		{
			name: "No max limit - Two occurrences",
			args: []Token{
				RuneToken{'1'},
				RuneToken{'2'},
				RuneToken{'2'},
				RuneToken{'3'},
				RuneToken{'4'},
			},
			matchFuncs: []MatchFunc{
				Equal(RuneToken{'1'})(1, 1),
				Equal(RuneToken{'2'})(1, 0),
				Equal(RuneToken{'3'})(1, 1),
				Equal(RuneToken{'4'})(1, 1),
			},
			replaceFunc: replaceRuneTokensWithTextToken,
			want: []Token{
				TextToken{Text: "12234"},
			},
			wantErr: false,
		},
		{
			name: "No max limit - Tree occurrences",
			args: []Token{
				RuneToken{'1'},
				RuneToken{'2'},
				RuneToken{'2'},
				RuneToken{'2'},
				RuneToken{'3'},
				RuneToken{'4'},
			},
			matchFuncs: []MatchFunc{
				Equal(RuneToken{'1'})(1, 1),
				Equal(RuneToken{'2'})(1, 3),
				Equal(RuneToken{'3'})(1, 1),
				Equal(RuneToken{'4'})(1, 1),
			},
			replaceFunc: replaceRuneTokensWithTextToken,
			want: []Token{
				TextToken{Text: "122234"},
			},
			wantErr: false,
		},
		{
			name: "No Match - Wrong Token Value",
			args: []Token{
				RuneToken{'1'},
				RuneToken{'2'},
				RuneToken{'3'},
				RuneToken{'4'},
			},
			matchFuncs: []MatchFunc{
				Equal(RuneToken{'1'})(1, 1),
				Equal(RuneToken{'5'})(1, 1),
				Equal(RuneToken{'3'})(1, 1),
				Equal(RuneToken{'4'})(1, 1),
			},
			replaceFunc: replaceRuneTokensWithTextToken,
			want: []Token{
				RuneToken{'1'},
				RuneToken{'2'},
				RuneToken{'3'},
				RuneToken{'4'},
			},
			wantErr: false,
		},
		{
			name: "No Match - Max Limit",
			args: []Token{
				RuneToken{'1'},
				RuneToken{'2'},
				RuneToken{'2'},
				RuneToken{'3'},
				RuneToken{'4'},
			},
			matchFuncs: []MatchFunc{
				Equal(RuneToken{'1'})(1, 1),
				Equal(RuneToken{'2'})(1, 1),
				Equal(RuneToken{'3'})(1, 1),
				Equal(RuneToken{'4'})(1, 1),
			},
			replaceFunc: replaceRuneTokensWithTextToken,
			want: []Token{
				RuneToken{'1'},
				RuneToken{'2'},
				RuneToken{'2'},
				RuneToken{'3'},
				RuneToken{'4'},
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			got, err := ReplaceTokens(tt.args, tt.matchFuncs, tt.replaceFunc)
			if (err != nil) != tt.wantErr {
				t.Errorf("ReplaceTokens() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			require.Equal(t, tt.want, got)
		})
	}
}

func Test_EqualType(t *testing.T) {
	tests := []struct {
		name        string
		args        []Token
		matchFuncs  []MatchFunc
		replaceFunc ReplaceFunc
		want        []Token
		wantErr     bool
	}{
		{
			name: "Exact Match",
			args: []Token{
				RuneToken{'1'},
				TextToken{"2"},
				RuneToken{'3'},
				RuneToken{'4'},
			},
			matchFuncs: []MatchFunc{
				EqualType[RuneToken](1, 1),
				EqualType[TextToken](1, 1),
				EqualType[RuneToken](1, 1),
				EqualType[RuneToken](1, 1),
			},
			replaceFunc: replaceRuneTokensWithTextToken,
			want: []Token{
				TextToken{Text: "1234"},
			},
			wantErr: false,
		},
		{
			name: "Optional - Token exists",
			args: []Token{
				RuneToken{'1'},
				TextToken{"2"},
				RuneToken{'3'},
				RuneToken{'4'},
			},
			matchFuncs: []MatchFunc{
				EqualType[RuneToken](1, 1),
				EqualType[TextToken](0, 1),
				EqualType[RuneToken](1, 1),
				EqualType[RuneToken](1, 1),
			},
			replaceFunc: replaceRuneTokensWithTextToken,
			want: []Token{
				TextToken{Text: "1234"},
			},
			wantErr: false,
		},
		{
			name: "Optional - Token does not exist",
			args: []Token{
				RuneToken{'1'},
				RuneToken{'3'},
				RuneToken{'4'},
			},
			matchFuncs: []MatchFunc{
				EqualType[RuneToken](1, 1),
				EqualType[TextToken](0, 1),
				EqualType[RuneToken](1, 1),
				EqualType[RuneToken](1, 1),
			},
			replaceFunc: replaceRuneTokensWithTextToken,
			want: []Token{
				TextToken{Text: "134"},
			},
			wantErr: false,
		},
		{
			name: "No max limit - One occurrence",
			args: []Token{
				RuneToken{'1'},
				TextToken{"2"},
				RuneToken{'3'},
				RuneToken{'4'},
			},
			matchFuncs: []MatchFunc{
				EqualType[RuneToken](1, 1),
				EqualType[TextToken](1, 0),
				EqualType[RuneToken](1, 1),
				EqualType[RuneToken](1, 1),
			},
			replaceFunc: replaceRuneTokensWithTextToken,
			want: []Token{
				TextToken{Text: "1234"},
			},
			wantErr: false,
		},
		{
			name: "No max limit - Two occurrences",
			args: []Token{
				RuneToken{'1'},
				TextToken{"2"},
				TextToken{"2"},
				RuneToken{'3'},
				RuneToken{'4'},
			},
			matchFuncs: []MatchFunc{
				EqualType[RuneToken](1, 1),
				EqualType[TextToken](1, 0),
				EqualType[RuneToken](1, 1),
				EqualType[RuneToken](1, 1),
			},
			replaceFunc: replaceRuneTokensWithTextToken,
			want: []Token{
				TextToken{Text: "12234"},
			},
			wantErr: false,
		},
		{
			name: "No max limit - Tree occurrences",
			args: []Token{
				RuneToken{'1'},
				TextToken{"2"},
				TextToken{"2"},
				TextToken{"2"},
				RuneToken{'3'},
				RuneToken{'4'},
			},
			matchFuncs: []MatchFunc{
				EqualType[RuneToken](1, 1),
				EqualType[TextToken](1, 3),
				EqualType[RuneToken](1, 1),
				EqualType[RuneToken](1, 1),
			},
			replaceFunc: replaceRuneTokensWithTextToken,
			want: []Token{
				TextToken{Text: "122234"},
			},
			wantErr: false,
		},
		{
			name: "No Match - Wrong Token",
			args: []Token{
				RuneToken{'1'},
				RuneToken{'2'},
				RuneToken{'3'},
				RuneToken{'4'},
			},
			matchFuncs: []MatchFunc{
				EqualType[RuneToken](1, 1),
				EqualType[TextToken](1, 1),
				EqualType[RuneToken](1, 1),
				EqualType[RuneToken](1, 1),
			},
			replaceFunc: replaceRuneTokensWithTextToken,
			want: []Token{
				RuneToken{'1'},
				RuneToken{'2'},
				RuneToken{'3'},
				RuneToken{'4'},
			},
			wantErr: false,
		},
		{
			name: "No Match - Max Limit",
			args: []Token{
				RuneToken{'1'},
				TextToken{"2"},
				TextToken{"2"},
				RuneToken{'3'},
				RuneToken{'4'},
			},
			matchFuncs: []MatchFunc{
				EqualType[RuneToken](1, 1),
				EqualType[TextToken](1, 1),
				EqualType[RuneToken](1, 1),
				EqualType[RuneToken](1, 1),
			},
			replaceFunc: replaceRuneTokensWithTextToken,
			want: []Token{
				RuneToken{'1'},
				TextToken{"2"},
				TextToken{"2"},
				RuneToken{'3'},
				RuneToken{'4'},
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			got, err := ReplaceTokens(tt.args, tt.matchFuncs, tt.replaceFunc)
			if (err != nil) != tt.wantErr {
				t.Errorf("ReplaceTokens() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			require.Equal(t, tt.want, got)
		})
	}
}

func Test_Or(t *testing.T) {
	tests := []struct {
		name        string
		args        []Token
		matchFuncs  []MatchFunc
		replaceFunc ReplaceFunc
		want        []Token
		wantErr     bool
	}{
		{
			name: "Match first option",
			args: []Token{
				RuneToken{'1'},
				RuneToken{'2'},
				RuneToken{'3'},
				RuneToken{'4'},
			},
			matchFuncs: []MatchFunc{
				Equal(RuneToken{'1'})(1, 1),
				Or(
					Equal(RuneToken{'2'})(1, 1),
					Equal(RuneToken{'5'})(1, 1),
				),
				Equal(RuneToken{'3'})(1, 1),
				Equal(RuneToken{'4'})(1, 1),
			},
			replaceFunc: replaceRuneTokensWithTextToken,
			want: []Token{
				TextToken{Text: "1234"},
			},
			wantErr: false,
		},
		{
			name: "Match second option",
			args: []Token{
				RuneToken{'1'},
				RuneToken{'5'},
				RuneToken{'3'},
				RuneToken{'4'},
			},
			matchFuncs: []MatchFunc{
				Equal(RuneToken{'1'})(1, 1),
				Or(
					Equal(RuneToken{'2'})(1, 1),
					Equal(RuneToken{'5'})(1, 1),
				),
				Equal(RuneToken{'3'})(1, 1),
				Equal(RuneToken{'4'})(1, 1),
			},
			replaceFunc: replaceRuneTokensWithTextToken,
			want: []Token{
				TextToken{Text: "1534"},
			},
			wantErr: false,
		},
		{
			name: "No Match",
			args: []Token{
				RuneToken{'1'},
				RuneToken{'2'},
				RuneToken{'3'},
				RuneToken{'4'},
			},
			matchFuncs: []MatchFunc{
				Equal(RuneToken{'1'})(1, 1),
				Or(
					Equal(RuneToken{'5'})(1, 1),
					Equal(RuneToken{'6'})(1, 1),
				),
				Equal(RuneToken{'3'})(1, 1),
				Equal(RuneToken{'4'})(1, 1),
			},
			replaceFunc: replaceRuneTokensWithTextToken,
			want: []Token{
				RuneToken{'1'},
				RuneToken{'2'},
				RuneToken{'3'},
				RuneToken{'4'},
			},
			wantErr: false,
		},
		{
			name: "First Options has no Max",
			args: []Token{
				RuneToken{'1'},
				RuneToken{'5'},
				RuneToken{'5'},
				RuneToken{'3'},
				RuneToken{'4'},
			},
			matchFuncs: []MatchFunc{
				Equal(RuneToken{'1'})(1, 1),
				Or(
					Equal(RuneToken{'5'})(1, 0),
					Equal(RuneToken{'2'})(1, 1),
				),
				Equal(RuneToken{'3'})(1, 1),
				Equal(RuneToken{'4'})(1, 1),
			},
			replaceFunc: replaceRuneTokensWithTextToken,
			want: []Token{
				TextToken{Text: "15534"},
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			got, err := ReplaceTokens(tt.args, tt.matchFuncs, tt.replaceFunc)
			if (err != nil) != tt.wantErr {
				t.Errorf("ReplaceTokens() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			require.Equal(t, tt.want, got)
		})
	}
}
