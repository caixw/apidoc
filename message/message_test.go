// SPDX-License-Identifier: MIT

package message

import (
	"bytes"
	"testing"
	"time"

	"github.com/issue9/assert"

	"github.com/caixw/apidoc/v5/internal/locale"
)

func TestHandler(t *testing.T) {
	a := assert.New(t)

	erro := new(bytes.Buffer)
	warn := new(bytes.Buffer)
	h := NewHandler(func(msg *Message) {
		switch msg.Type {
		case Erro:
			erro.WriteString("erro")
		case Warn:
			warn.WriteString("warn")
		default:
			panic("panic")
		}
	})
	a.NotError(h)

	h.Error(Erro, NewLocaleError("erro.go", "", 0, locale.ErrRequired))
	h.Error(Warn, NewLocaleError("warn.go", "", 0, locale.ErrRequired))

	time.Sleep(1 * time.Second) // 等待 channel 完成
	a.Equal(erro.String(), "erro")
	a.Equal(warn.String(), "warn")

	h.Stop()
	a.Panic(func() { // 已经关闭 messages
		h.Error(Erro, NewLocaleError("erro", "", 0, locale.ErrRequired))
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
