// SPDX-License-Identifier: MIT

package lsp

import (
	"io/ioutil"
	"log"
	"testing"

	"github.com/issue9/assert"

	"github.com/caixw/apidoc/v7/core"
	"github.com/caixw/apidoc/v7/core/messagetest"
	"github.com/caixw/apidoc/v7/internal/ast"
	"github.com/caixw/apidoc/v7/internal/lsp/protocol"
	"github.com/caixw/apidoc/v7/internal/xmlenc"
)

func TestDeleteURI(t *testing.T) {
	a := assert.New(t)

	d := &ast.APIDoc{}
	d.APIDoc = &ast.APIDocVersionAttribute{Value: xmlenc.String{Value: "1.0.0"}}
	d.URI = core.URI("uri1")
	d.APIs = []*ast.API{
		{ //1
			BaseTag: xmlenc.BaseTag{Base: xmlenc.Base{Location: core.Location{URI: "uri1"}}},
		},
		{ //2
			BaseTag: xmlenc.BaseTag{Base: xmlenc.Base{Location: core.Location{URI: "uri2"}}},
		},
		{ //3
			BaseTag: xmlenc.BaseTag{Base: xmlenc.Base{Location: core.Location{URI: "uri3"}}},
		},
		{ //4
		},
	}

	a.True(deleteURI(d, "uri3"))
	a.Equal(3, len(d.APIs)).NotNil(d.APIDoc)

	// 同时会删除 1,4
	a.True(deleteURI(d, "uri1"))
	a.Equal(1, len(d.APIs)).Nil(d.APIDoc)

	a.True(deleteURI(d, "uri2"))
	a.Equal(0, len(d.APIs)).Nil(d.APIDoc)

	a.False(deleteURI(d, "uri2"))
}

func TestServer_textDocumentFoldingRange(t *testing.T) {
	a := assert.New(t)

	const referenceDefinitionDoc = `<apidoc version="1.1.1">
	<title>标题</title>
	<mimetype>xml</mimetype>
	<tag name="t1" title="tag1" />
	<tag name="t2" title="tag2" />
	<api method="GET">
		<tag>t1</tag>
		<path path="/users" />
		<response status="200" />
	</api>
</apidoc>`

	uri := core.URI("file:///root/doc.go")
	blk := core.Block{Data: []byte(referenceDefinitionDoc), Location: core.Location{URI: uri}}
	rslt := messagetest.NewMessageHandler()
	doc := &ast.APIDoc{}
	doc.Parse(rslt.Handler, blk)
	rslt.Handler.Stop()
	a.Empty(rslt.Errors)

	s := newTestServer(true, log.New(ioutil.Discard, "", 0), log.New(ioutil.Discard, "", 0))
	s.folders = append(s.folders, &folder{srv: s, doc: doc})

	a.NotNil(s)
	result := make([]protocol.FoldingRange, 0, 10)
	err := s.textDocumentFoldingRange(false, &protocol.FoldingRangeParams{}, &result) // 未指定 URI
	a.NotError(err).Empty(result)

	err = s.textDocumentFoldingRange(false, &protocol.FoldingRangeParams{
		TextDocument: protocol.TextDocumentIdentifier{URI: uri},
	}, &result)
	a.NotError(err)
	a.Equal(result, []protocol.FoldingRange{
		{
			StartLine: 0,
			EndLine:   10,
			Kind:      "comment",
		},
		{
			StartLine: 5,
			EndLine:   9,
			Kind:      "comment",
		},
	})
}
