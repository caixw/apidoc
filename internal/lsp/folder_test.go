// SPDX-License-Identifier: MIT

package lsp

import (
	"io/ioutil"
	"log"
	"testing"

	"github.com/issue9/assert"

	"github.com/caixw/apidoc/v7/core"
)

func TestFolder_messageHandler(t *testing.T) {
	a := assert.New(t)

	s := &server{erro: log.New(ioutil.Discard, "", 0)}
	f := &folder{srv: s}
	f.messageHandler(&core.Message{Type: core.Erro, Message: "abc"})
	a.Empty(f.errors).Empty(f.warns)

	f = &folder{srv: s}
	f.messageHandler(&core.Message{Type: core.Erro, Message: &core.Error{Location: core.Location{URI: "uri"}}})
	a.Empty(f.warns).Equal(1, len(f.errors))
	// 相同的错误，不会再次添加
	f.messageHandler(&core.Message{Type: core.Erro, Message: &core.Error{Location: core.Location{URI: "uri"}}})
	a.Empty(f.warns).Equal(1, len(f.errors))

	f = &folder{srv: s}
	f.messageHandler(&core.Message{Type: core.Warn, Message: &core.Error{Location: core.Location{URI: "uri"}}})
	a.Equal(1, len(f.warns)).Empty(f.errors)
	f.messageHandler(&core.Message{Type: core.Warn, Message: &core.Error{Location: core.Location{URI: "uri"}}})
	a.Equal(1, len(f.warns)).Empty(f.errors)

	f = &folder{srv: s}
	a.PanicString(func() {
		f.messageHandler(&core.Message{Message: &core.Error{}, Type: -100})
	}, "unreached")
}
