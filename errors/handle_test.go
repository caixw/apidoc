// Copyright 2018 by caixw, All rights reserved.
// Use of this source code is governed by a MIT
// license that can be found in the LICENSE file.

package errors

import (
	"bytes"
	"log"
	"testing"
	"time"

	"github.com/issue9/assert"

	"github.com/caixw/apidoc/internal/locale"
)

func TestHandler(t *testing.T) {
	a := assert.New(t)
	erro := new(bytes.Buffer)
	warn := new(bytes.Buffer)
	errolog := log.New(erro, "[ERRO]", 0)
	warnlog := log.New(warn, "[WARN]", 0)

	h := NewHandler(NewHandlerFunc(errolog, warnlog))
	a.NotError(h)

	le := LocaleError{MessageKey: locale.ErrRequired}
	h.SyntaxError(&Error{File: "erro.go", LocaleError: le})
	h.SyntaxWarn(&Error{File: "warn.go", LocaleError: le})

	time.Sleep(1 * time.Second) // 等待 channel 完成
	a.Contains(erro.String(), "erro.go")
	a.Contains(warn.String(), "warn.go")

	h.Stop()
	a.Panic(func() {
		h.SyntaxError(&Error{File: "erro.go", LocaleError: le})
	})
}
