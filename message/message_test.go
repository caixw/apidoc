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

var _ HandlerFunc = NewLogHandlerFunc(nil, nil, nil, nil)

func TestHandler(t *testing.T) {
	a := assert.New(t)
	erro := new(bytes.Buffer)
	warn := new(bytes.Buffer)
	info := new(bytes.Buffer)
	succ := new(bytes.Buffer)
	errolog := log.New(erro, "[ERRO]", 0)
	warnlog := log.New(warn, "[WARN]", 0)
	infolog := log.New(info, "[INFO]", 0)
	succlog := log.New(succ, "[SUCC]", 0)

	h := NewHandler(NewLogHandlerFunc(errolog, warnlog, infolog, succlog))
	a.NotError(h)

	h.Error(Erro, NewLocaleError("erro.go", "", 0, locale.ErrRequired))
	h.Error(Warn, NewLocaleError("warn.go", "", 0, locale.ErrRequired))

	time.Sleep(1 * time.Second) // 等待 channel 完成
	a.Contains(erro.String(), "erro.go")
	a.Contains(warn.String(), "warn.go")

	h.Stop()
	a.Panic(func() { // 已经关闭 messages
		h.Error(Erro, NewLocaleError("erro.go", "", 0, locale.ErrRequired))
	})
}

func TestHandler_Stop(t *testing.T) {
	a := assert.New(t)
	var exit bool

	h := NewHandler(func(msg *Message) {
		time.Sleep(time.Second)
		exit = true
	})
	a.NotError(h)

	h.Message(Erro, locale.ErrRequired)
	h.Stop() // 此处会阻塞，等待完成
	a.True(exit)
}
