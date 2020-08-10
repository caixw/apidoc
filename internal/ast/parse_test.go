// SPDX-License-Identifier: MIT

package ast

import (
	"io"
	"net/http"
	"testing"

	"github.com/issue9/assert"

	"github.com/caixw/apidoc/v7/core"
	"github.com/caixw/apidoc/v7/core/messagetest"
	"github.com/caixw/apidoc/v7/internal/xmlenc"
)

func newParser(a *assert.Assertion, data string, uri core.URI) (*xmlenc.Parser, *messagetest.Result) {
	rslt := messagetest.NewMessageHandler()
	p, err := xmlenc.NewParser(rslt.Handler, core.Block{Data: []byte(data), Location: core.Location{URI: uri}})
	a.NotError(err).NotNil(p)
	return p, rslt
}

func TestAPIDoc_ParseBlocks(t *testing.T) {
	a := assert.New(t)

	rslt := messagetest.NewMessageHandler()
	doc := &APIDoc{}
	doc.ParseBlocks(rslt.Handler, func(blocks chan core.Block) {
		blocks <- core.Block{Data: []byte(`<api method="GET"><path path="/p1" /></api>`)}
		blocks <- core.Block{Data: []byte(`<api method="POST"><path path="/p1" /></api>`)}
		blocks <- core.Block{Data: []byte(`ErrNoDocFormat`)}
	})
	rslt.Handler.Stop()
	a.Empty(rslt.Errors)

	// 带错误返回
	rslt = messagetest.NewMessageHandler()
	doc = &APIDoc{}
	doc.ParseBlocks(rslt.Handler, func(blocks chan core.Block) {
		blocks <- core.Block{Data: []byte(`<api method="GET"><path path="/p1" /></api>`)}
		blocks <- core.Block{Data: []byte(`<api method="POST"><path path="/p1" /></api>`)}
		blocks <- core.Block{Data: []byte(`<api method="GET" />`)} // 少 path
	})
	rslt.Handler.Stop()
	a.NotEmpty(rslt.Errors)
}

func TestAPIDoc_Parse(t *testing.T) {
	a := assert.New(t)

	rslt := messagetest.NewMessageHandler()
	d := &APIDoc{}
	d.Parse(rslt.Handler, core.Block{Data: []byte("<api>")})
	rslt.Handler.Stop()
	a.Empty(rslt.Errors)

	// 直接结束标签
	rslt = messagetest.NewMessageHandler()
	d.Parse(rslt.Handler, core.Block{Data: []byte("</api>")})
	rslt.Handler.Stop()
	a.NotEmpty(rslt.Errors)

	// 多个 apidoc 标签
	rslt = messagetest.NewMessageHandler()
	d.Title = &Element{Content: Content{Value: "title"}}
	d.Parse(rslt.Handler, core.Block{Data: []byte("<apidoc />")})
	rslt.Handler.Stop()
	a.Error(rslt.Errors[0])

	// 未知标签
	rslt = messagetest.NewMessageHandler()
	d.Parse(rslt.Handler, core.Block{Data: []byte("<tag />")})
	rslt.Handler.Stop()
	a.Empty(rslt.Errors)

	// 先解析 api，再解析 apidoc，不会覆盖 apidoc.APIs
	rslt = messagetest.NewMessageHandler()
	d = &APIDoc{
		APIs: []*API{
			{
				Method: &MethodAttribute{Value: xmlenc.String{Value: http.MethodGet}},
			},
		},
	}
	d.Parse(rslt.Handler, core.Block{Data: []byte(`<apidoc><title>title</title><mimetype>application/json</mimetype></apidoc>`)})
	rslt.Handler.Stop()
	a.Empty(rslt.Errors).
		Equal(1, len(d.APIs))
}

func TestGetTagName(t *testing.T) {
	a := assert.New(t)

	p, rslt := newParser(a, "  <root>xx</root>", "")
	root, err := getTagName(p)
	a.NotError(err).Equal(root, "root")
	rslt.Handler.Stop()
	a.Empty(rslt.Errors)

	p, rslt = newParser(a, "<!-- xx -->  <root>xx</root>", "")
	root, err = getTagName(p)
	a.NotError(err).Equal(root, "root")
	rslt.Handler.Stop()
	a.Empty(rslt.Errors)

	p, rslt = newParser(a, "<!-- xx -->   <root>xx</root>", "")
	root, err = getTagName(p)
	a.NotError(err).Equal(root, "root")
	rslt.Handler.Stop()
	a.Empty(rslt.Errors)

	// 无效格式
	p, rslt = newParser(a, "<!-- xx   <root>xx</root>", "")
	root, err = getTagName(p)
	a.Error(err).Equal(root, "")
	rslt.Handler.Stop()
	a.Empty(rslt.Errors)

	// 无效格式
	p, rslt = newParser(a, "</root>", "")
	root, err = getTagName(p)
	a.Error(err).Equal(root, "")
	rslt.Handler.Stop()
	a.Empty(rslt.Errors)

	// io.EOF
	p, rslt = newParser(a, "<!-- xx -->", "")
	root, err = getTagName(p)
	a.Equal(err, io.EOF).Equal(root, "")
	rslt.Handler.Stop()
	a.Empty(rslt.Errors)
}

func TestAPIDoc_sortAPIs(t *testing.T) {
	a := assert.New(t)

	doc := &APIDoc{
		APIs: []*API{
			{
				Path:   &Path{Path: &Attribute{Value: xmlenc.String{Value: "/p1/p3"}}},
				Method: &MethodAttribute{Value: xmlenc.String{Value: http.MethodPost}},
			},
			{
				Path:   &Path{Path: &Attribute{Value: xmlenc.String{Value: "/p1/p3"}}},
				Method: &MethodAttribute{Value: xmlenc.String{Value: http.MethodGet}},
			},
			{
				Path:   &Path{Path: &Attribute{Value: xmlenc.String{Value: "/p1/p2"}}},
				Method: &MethodAttribute{Value: xmlenc.String{Value: http.MethodGet}},
			},
			{
				Path:   &Path{Path: &Attribute{Value: xmlenc.String{Value: "/p1/p3"}}},
				Method: &MethodAttribute{Value: xmlenc.String{Value: http.MethodPut}},
			},
		},
	}
	doc.sortAPIs()

	api := doc.APIs[0]
	a.Equal(api.Path.Path.V(), "/p1/p2")

	api = doc.APIs[1]
	a.Equal(api.Path.Path.V(), "/p1/p3").Equal(api.Method.V(), http.MethodGet)

	api = doc.APIs[2]
	a.Equal(api.Path.Path.V(), "/p1/p3").Equal(api.Method.V(), http.MethodPost)

	api = doc.APIs[3]
	a.Equal(api.Path.Path.V(), "/p1/p3").Equal(api.Method.V(), http.MethodPut)
}
