// Copyright 2018 by caixw, All rights reserved.
// Use of this source code is governed by a MIT
// license that can be found in the LICENSE file.

package doc

import (
	"io/ioutil"
	"log"

	"github.com/caixw/apidoc/doc/lexer"
	"github.com/caixw/apidoc/errors"
	"github.com/caixw/apidoc/internal/input"
)

func newLexer(data string) *lexer.Lexer {
	erro := log.New(ioutil.Discard, "[ERRO]", 0)
	warn := log.New(ioutil.Discard, "[WARN]", 0)
	h := errors.NewHandler(errors.NewLogHandlerFunc(erro, warn))
	return lexer.New(input.Block{Data: []byte(data)}, h)
}

func newTag(data string) *lexer.Tag {
	return newLexer(data).Tag()
}
