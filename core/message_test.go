// SPDX-License-Identifier: MIT

package core

import (
	"bytes"
	"fmt"
	"testing"
	"time"

	"github.com/issue9/assert/v2"

	"github.com/caixw/apidoc/v7/internal/locale"
)

var _ fmt.Stringer = Erro

func TestType_String(t *testing.T) {
	a := assert.New(t, false)
	a.Equal("ERRO", Erro.String())
	a.Equal("SUCC", Succ.String())
	a.Equal("INFO", Info.String())
	a.Equal("WARN", Warn.String())
	a.Equal("<unknown>", MessageType(-22).String())
}

func TestHandler(t *testing.T) {
	a := assert.New(t, false)

	erro := new(bytes.Buffer)
	warn := new(bytes.Buffer)
	info := new(bytes.Buffer)
	succ := new(bytes.Buffer)
	h := NewMessageHandler(func(msg *Message) {
		switch msg.Type {
		case Erro:
			erro.WriteString("erro")
		case Warn:
			warn.WriteString("warn")
		case Info:
			info.WriteString("info")
		case Succ:
			succ.WriteString("succ")
		default:
			panic("panic")
		}
	})
	a.NotNil(h)

	h.Error((Location{URI: "erro.go"}).NewError(locale.ErrInvalidUTF8Character))
	h.Warning((Location{URI: "warn.go"}).NewError(locale.ErrInvalidUTF8Character))
	h.Info((Location{URI: "info.go"}).NewError(locale.ErrInvalidUTF8Character))
	h.Success((Location{URI: "succ.go"}).NewError(locale.ErrInvalidUTF8Character))

	time.Sleep(1 * time.Second) // 等待 channel 完成
	a.Equal(erro.String(), "erro")
	a.Equal(warn.String(), "warn")
	a.Equal(info.String(), "info")
	a.Equal(succ.String(), "succ")

	h.Stop()
	a.Panic(func() { // 已经关闭 messages
		h.Error((Location{URI: "erro"}).NewError(locale.ErrInvalidUTF8Character))
	})
}

func TestHandler_Stop(t *testing.T) {
	a := assert.New(t, false)
	var exit bool

	h := NewMessageHandler(func(msg *Message) {
		time.Sleep(time.Second)
		exit = true
	})
	a.NotNil(h)

	h.Locale(Erro, locale.ErrInvalidUTF8Character)
	h.Stop() // 此处会阻塞，等待完成
	a.True(exit)
}
