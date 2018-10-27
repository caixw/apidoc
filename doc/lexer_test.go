// Copyright 2018 by caixw, All rights reserved.
// Use of this source code is governed by a MIT
// license that can be found in the LICENSE file.

package doc

import (
	"github.com/caixw/apidoc/doc/lexer"
	"github.com/caixw/apidoc/input"
)

func newLexer(data string) *lexer.Lexer {
	return lexer.New(input.Block{Data: []byte(data)})
}

func newTag(data string) *lexer.Tag {
	return &lexer.Tag{
		Data: []byte(data),
	}
}
