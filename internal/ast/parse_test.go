// SPDX-License-Identifier: MIT

package ast

import (
	"io"
	"net/http"
	"testing"

	"github.com/issue9/assert"

	"github.com/caixw/apidoc/v7/core"
	"github.com/caixw/apidoc/v7/core/messagetest"
	"github.com/caixw/apidoc/v7/internal/token"
)

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
				Method: &MethodAttribute{Value: token.String{Value: http.MethodGet}},
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

	p, err := token.NewParser(core.Block{Data: []byte("  <root>xx</root>")})
	a.NotError(err).NotNil(p)
	root, err := getTagName(p)
	a.NotError(err).Equal(root, "root")

	p, err = token.NewParser(core.Block{Data: []byte("<!-- xx -->  <root>xx</root>")})
	a.NotError(err).NotNil(p)
	root, err = getTagName(p)
	a.NotError(err).Equal(root, "root")

	p, err = token.NewParser(core.Block{Data: []byte("<!-- xx -->   <root>xx</root>")})
	a.NotError(err).NotNil(p)
	root, err = getTagName(p)
	a.NotError(err).Equal(root, "root")

	// 无效格式
	p, err = token.NewParser(core.Block{Data: []byte("<!-- xx   <root>xx</root>")})
	a.NotError(err).NotNil(p)
	root, err = getTagName(p)
	a.Error(err).Equal(root, "")

	// 无效格式
	p, err = token.NewParser(core.Block{Data: []byte("</root>")})
	a.NotError(err).NotNil(p)
	root, err = getTagName(p)
	a.Error(err).Equal(root, "")

	// io.EOF
	p, err = token.NewParser(core.Block{Data: []byte("<!-- xx -->")})
	a.NotError(err).NotNil(p)
	root, err = getTagName(p)
	a.Equal(err, io.EOF).Equal(root, "")
}
