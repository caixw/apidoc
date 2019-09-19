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

var _ HandlerFunc = NewLogHandlerFunc(nil, nil, nil)

func TestHandler(t *testing.T) {
	a := assert.New(t)
	erro := new(bytes.Buffer)
	warn := new(bytes.Buffer)
	info := new(bytes.Buffer)
	errolog := log.New(erro, "[ERRO]", 0)
	warnlog := log.New(warn, "[WARN]", 0)
	infolog := log.New(info, "[INFO]", 0)

	h := NewHandler(NewLogHandlerFunc(errolog, warnlog, infolog))
	a.NotError(h)

	h.Error(Erro, NewError("erro.go", "", 0, locale.ErrRequired))
	h.Error(Warn, NewError("warn.go", "", 0, locale.ErrRequired))

	time.Sleep(1 * time.Second) // 等待 channel 完成
	a.Contains(erro.String(), "erro.go")
	a.Contains(warn.String(), "warn.go")

	h.Stop()
	a.Panic(func() { // 已经关闭 messages
		h.Error(Erro, NewError("erro.go", "", 0, locale.ErrRequired))
	})
}
