// SPDX-License-Identifier: MIT

package ast

import (
	"testing"

	"github.com/issue9/assert"

	"github.com/caixw/apidoc/v7/core"
	"github.com/caixw/apidoc/v7/core/messagetest"
	"github.com/caixw/apidoc/v7/internal/token"
)

func TestAPIDoc_ParseBlocks(t *testing.T) {
	a := assert.New(t)

	erro, _, h := messagetest.MessageHandler()
	doc := &APIDoc{}
	doc.ParseBlocks(h, func(blocks chan core.Block) {
		blocks <- core.Block{Data: []byte(`<api method="GET"><path path="/p1" /></api>`)}
		blocks <- core.Block{Data: []byte(`<api method="POST"><path path="/p1" /></api>`)}
		blocks <- core.Block{Data: []byte(`ErrNoDocFormat`)}
	})
	h.Stop()
	a.Empty(erro.String())

	// 带错误返回
	erro, _, h = messagetest.MessageHandler()
	doc = &APIDoc{}
	doc.ParseBlocks(h, func(blocks chan core.Block) {
		blocks <- core.Block{Data: []byte(`<api method="GET"><path path="/p1" /></api>`)}
		blocks <- core.Block{Data: []byte(`<api method="POST"><path path="/p1" /></api>`)}
		blocks <- core.Block{Data: []byte(`<api method="GET" />`)} // 少 path
	})
	h.Stop()
	a.NotEmpty(erro.String())
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
