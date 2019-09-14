// SPDX-License-Identifier: MIT

package message

import (
	"bytes"
	"log"
	"testing"
	"time"

	"github.com/issue9/assert"

	"github.com/caixw/apidoc/v5/internal/locale"
)

var _ HandlerFunc = NewLogHandlerFunc(nil, nil)

func TestHandler(t *testing.T) {
	a := assert.New(t)
	erro := new(bytes.Buffer)
	warn := new(bytes.Buffer)
	errolog := log.New(erro, "[ERRO]", 0)
	warnlog := log.New(warn, "[WARN]", 0)

	h := NewHandler(NewLogHandlerFunc(errolog, warnlog))
	a.NotError(h)

	le := locale.NewLocale(locale.ErrRequired)
	h.Error(&SyntaxError{File: "erro.go", locale: le})
	h.Warn(&SyntaxError{File: "warn.go", locale: le})

	time.Sleep(1 * time.Second) // 等待 channel 完成
	a.Contains(erro.String(), "erro.go")
	a.Contains(warn.String(), "warn.go")

	h.Stop()
	a.Panic(func() {
		h.Error(&SyntaxError{File: "erro.go", locale: le})
	})
}
