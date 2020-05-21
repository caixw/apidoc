// SPDX-License-Identifier: MIT

package ast

import (
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

	d := &APIDoc{}
	err := d.Parse(core.Block{Data: []byte("<api>")})
	a.Equal(err, ErrNoDocFormat)

	// 直接结束标签
	err = d.Parse(core.Block{Data: []byte("</api>")})
	a.Equal(err, ErrNoDocFormat)

	// 多个 apidoc 标签
	d.Title = &Element{Content: Content{Value: "title"}}
	err = d.Parse(core.Block{Data: []byte("<apidoc />")})
	a.Error(err)

	// 未知标签
	err = d.Parse(core.Block{Data: []byte("<tag />")})
	a.Equal(err, ErrNoDocFormat)

	// 先解析 api，再解析 apidoc，不会覆盖 apidoc.APIs
	d = &APIDoc{
		Apis: []*API{
			{
				Method: &MethodAttribute{Value: token.String{Value: http.MethodGet}},
			},
		},
	}
	err = d.Parse(core.Block{Data: []byte(`<apidoc><title>title</title><mimetype>application/json</mimetype></apidoc>`)})
	a.NotError(err)
	a.Equal(1, len(d.Apis))
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
	a.Equal(err, ErrNoDocFormat).Equal(root, "")

	// io.EOF
	p, err = token.NewParser(core.Block{Data: []byte("<!-- xx -->")})
	a.NotError(err).NotNil(p)
	root, err = getTagName(p)
	a.Equal(err, ErrNoDocFormat).Equal(root, "")
}
