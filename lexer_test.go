package parser

import (
	"bytes"
	"io"
	"testing"
	"testing/iotest"

	"github.com/stretchr/testify/require"
)

func Test_Lex(t *testing.T) {
	tests := []struct {
		name    string
		args    string
		want    []Token
		wantErr bool
	}{
		{
			name:    "Empty",
			args:    "",
			want:    nil,
			wantErr: false,
		},
		{
			name: "Simple Text",
			args: "Hello World!",
			want: []Token{
				RuneToken{Rune: 'H'},
				RuneToken{Rune: 'e'},
				RuneToken{Rune: 'l'},
				RuneToken{Rune: 'l'},
				RuneToken{Rune: 'o'},
				RuneToken{Rune: ' '},
				RuneToken{Rune: 'W'},
				RuneToken{Rune: 'o'},
				RuneToken{Rune: 'r'},
				RuneToken{Rune: 'l'},
				RuneToken{Rune: 'd'},
				RuneToken{Rune: '!'},
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			got, err := Lex(bytes.NewReader([]byte(tt.args)))
			if (err != nil) != tt.wantErr {
				t.Errorf("Lex() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			require.Equal(t, tt.want, got)
		})
	}
}

func Test_Lex_UnexpectedEOF(t *testing.T) {
	tokens, err := Lex(iotest.ErrReader(io.ErrUnexpectedEOF))
	require.Nil(t, tokens)
	require.Error(t, err)
}
