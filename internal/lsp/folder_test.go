// SPDX-License-Identifier: MIT

package lsp

import (
	"io/ioutil"
	"log"
	"testing"

	"github.com/issue9/assert/v2"

	"github.com/caixw/apidoc/v7/core"
	"github.com/caixw/apidoc/v7/internal/locale"
	"github.com/caixw/apidoc/v7/internal/lsp/protocol"
)

func TestFolder_messageHandler(t *testing.T) {
	a := assert.New(t, false)

	s := newTestServer(true, log.New(ioutil.Discard, "", 0), log.New(ioutil.Discard, "", 0))
	f := &folder{srv: s, diagnostics: map[core.URI]*protocol.PublishDiagnosticsParams{}}
	f.messageHandler(&core.Message{Type: core.Erro, Message: "abc"})
	a.Empty(f.diagnostics)

	err := locale.NewError(locale.ErrInvalidUTF8Character)

	f = &folder{srv: s, diagnostics: map[core.URI]*protocol.PublishDiagnosticsParams{}}
	f.messageHandler(&core.Message{Type: core.Erro, Message: &core.Error{Location: core.Location{URI: "uri"}, Err: err}})
	a.Equal(1, len(f.diagnostics))
	// 相同的错误，不会再次添加
	f.messageHandler(&core.Message{Type: core.Erro, Message: &core.Error{Location: core.Location{URI: "uri"}, Err: err}})
	a.Equal(1, len(f.diagnostics))

	f = &folder{srv: s, diagnostics: map[core.URI]*protocol.PublishDiagnosticsParams{}}
	f.messageHandler(&core.Message{Type: core.Warn, Message: &core.Error{Location: core.Location{URI: "uri"}, Err: err}})
	a.Equal(1, len(f.diagnostics))
	f.messageHandler(&core.Message{Type: core.Warn, Message: &core.Error{Location: core.Location{URI: "uri"}, Err: err}})
	a.Equal(1, len(f.diagnostics))

	f = &folder{srv: s, diagnostics: map[core.URI]*protocol.PublishDiagnosticsParams{}}
	a.PanicString(func() {
		f.messageHandler(&core.Message{Message: &core.Error{}, Type: -100})
	}, "unreached")
}
