// SPDX-License-Identifier: MIT

package core

import (
	"bytes"
	"fmt"
	"testing"
	"time"

	"github.com/issue9/assert"

	"github.com/caixw/apidoc/v7/internal/locale"
)

var _ fmt.Stringer = Erro

func TestType_String(t *testing.T) {
	a := assert.New(t)
	a.Equal("ERRO", Erro.String())
	a.Equal("SUCC", Succ.String())
	a.Equal("INFO", Info.String())
	a.Equal("WARN", Warn.String())
	a.Equal("<unknown>", MessageType(-22).String())
}

func TestHandler(t *testing.T) {
	a := assert.New(t)

	erro := new(bytes.Buffer)
	warn := new(bytes.Buffer)
	h := NewMessageHandler(func(msg *Message) {
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

	h.Error(Erro, NewLocaleError(Location{URI: "erro.go"}, "", locale.ErrRequired))
	h.Error(Warn, NewLocaleError(Location{URI: "warn.go"}, "", locale.ErrRequired))

	time.Sleep(1 * time.Second) // 等待 channel 完成
	a.Equal(erro.String(), "erro")
	a.Equal(warn.String(), "warn")

	h.Stop()
	a.Panic(func() { // 已经关闭 messages
		h.Error(Erro, NewLocaleError(Location{URI: "erro"}, "", locale.ErrRequired))
	})
}

func TestHandler_Stop(t *testing.T) {
	a := assert.New(t)
	var exit bool

	h := NewMessageHandler(func(msg *Message) {
		time.Sleep(time.Second)
		exit = true
	})
	a.NotError(h)

	h.Message(Erro, locale.ErrRequired)
	h.Stop() // 此处会阻塞，等待完成
	a.True(exit)
}
