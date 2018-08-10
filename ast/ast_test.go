package ast

import (
	"github.com/stretchr/testify/assert"
	"go_interpreter/token"
	"testing"
)

func TestString(t *testing.T) {
	prog := &Program{
		Statements: []Statement{
			&LetStatement{token.Token{token.LET, "let"},
				&Identifier{token.Token{token.IDENT, "v1"}, "v1"},
				&Identifier{token.Token{token.IDENT, "v2"}, "v2"},
			},
		},
	}

	assert.Equal(t, prog.String(), "let v1 = v2;", "Program")
}
