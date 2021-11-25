// SPDX-License-Identifier: MIT

package lsp

import (
	"io/ioutil"
	"log"
	"testing"

	"github.com/issue9/assert/v2"
	"golang.org/x/text/language"

	"github.com/caixw/apidoc/v7/core"
	"github.com/caixw/apidoc/v7/internal/locale"
	"github.com/caixw/apidoc/v7/internal/lsp/protocol"
)

func TestServer_initialize(t *testing.T) {
	a := assert.New(t, false)
	s := newTestServer(true, log.New(ioutil.Discard, "", 0), log.New(ioutil.Discard, "", 0))
	in := &protocol.InitializeParams{}
	out := &protocol.InitializeResult{}
	a.NotError(s.initialize(false, in, out))
	a.Equal(out.ServerInfo.Name, core.Name)
	a.Equal(s.clientParams, in).Equal(s.serverResult, out)
	a.Equal(s.state, serverInitializing)

	s = newTestServer(true, log.New(ioutil.Discard, "", 0), log.New(ioutil.Discard, "", 0))
	in = &protocol.InitializeParams{
		Capabilities: protocol.ClientCapabilities{Workspace: &protocol.WorkspaceClientCapabilities{
			WorkspaceFolders: true,
		}},
		InitializationOptions: &protocol.InitializationOptions{Locale: "cmn-hant"},
	}
	out = &protocol.InitializeResult{}
	a.NotError(s.initialize(false, in, out))
	a.Equal(locale.Tag(), language.MustParse("cmn-hant"))
	a.True(out.Capabilities.Workspace.WorkspaceFolders.Supported).
		True(out.Capabilities.Workspace.WorkspaceFolders.ChangeNotifications)

	s = newTestServer(true, log.New(ioutil.Discard, "", 0), log.New(ioutil.Discard, "", 0))
	in = &protocol.InitializeParams{
		Capabilities: protocol.ClientCapabilities{TextDocument: protocol.TextDocumentClientCapabilities{
			Hover:        &protocol.HoverCapabilities{},
			FoldingRange: &protocol.FoldingRangeClientCapabilities{},
			Definition:   &protocol.DefinitionClientCapabilities{},
		}},
	}
	out = &protocol.InitializeResult{}
	a.NotError(s.initialize(false, in, out))
	a.False(out.Capabilities.HoverProvider).
		True(out.Capabilities.DefinitionProvider).
		True(out.Capabilities.FoldingRangeProvider)

	s = newTestServer(true, log.New(ioutil.Discard, "", 0), log.New(ioutil.Discard, "", 0))
	in = &protocol.InitializeParams{
		Capabilities: protocol.ClientCapabilities{TextDocument: protocol.TextDocumentClientCapabilities{
			Hover:      &protocol.HoverCapabilities{ContentFormat: []protocol.MarkupKind{protocol.MarkupKindMarkdown}},
			References: &protocol.ReferenceClientCapabilities{},
			Completion: &protocol.CompletionClientCapabilities{},
		}},
	}
	out = &protocol.InitializeResult{}
	a.NotError(s.initialize(false, in, out))
	a.True(out.Capabilities.HoverProvider).
		False(out.Capabilities.DefinitionProvider).
		False(out.Capabilities.FoldingRangeProvider).
		True(out.Capabilities.ReferencesProvider).
		NotNil(out.Capabilities.CompletionProvider)
}
