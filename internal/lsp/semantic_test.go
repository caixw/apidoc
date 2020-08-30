// SPDX-License-Identifier: MIT

package lsp

import (
	"testing"

	"github.com/issue9/assert"

	"github.com/caixw/apidoc/v7/core"
	"github.com/caixw/apidoc/v7/core/messagetest"
	"github.com/caixw/apidoc/v7/internal/ast"
)

func TestTokenBuilder_append(t *testing.T) {
	a := assert.New(t)

	b := &tokenBuilder{}
	b.append(core.Range{
		Start: core.Position{Line: 1, Character: 11},
		End:   core.Position{Line: 1, Character: 12},
	}, 1)
	a.Equal(b.tokens[0], []int{1, 11, 1, 1, 0})

	// 空值，不会添加内容
	b.append(core.Range{}, 1)
	a.Equal(1, len(b.tokens))

	// 长度为 0
	b.append(core.Range{
		Start: core.Position{Line: 1, Character: 11},
		End:   core.Position{Line: 1, Character: 11},
	}, 1)
	a.Equal(2, len(b.tokens))

	// 长度为负数
	a.Panic(func() {
		b.append(core.Range{
			Start: core.Position{Line: 1, Character: 11},
			End:   core.Position{Line: 1, Character: 10},
		}, 1)
	})
}

func TestTokenBuilder_build(t *testing.T) {
	a := assert.New(t)

	b := &tokenBuilder{
		tokens: [][]int{
			{1, 2, 5, 0, 0},
			{1, 12, 5, 0, 0},
			{2, 2, 7, 0, 0},
			{2, 2, 5, 0, 0},
			{2, 5, 5, 0, 0},
			{11, 1, 5, 0, 0},
		},
	}
	a.Equal(b.build(), []int{
		1, 2, 5, 0, 0,
		0, 10, 5, 0, 0,
		1, 2, 7, 0, 0,
		0, 0, 5, 0, 0,
		0, 3, 5, 0, 0,
		9, 1, 5, 0, 0,
	})
}

func TestTokenBuilder_sort(t *testing.T) {
	a := assert.New(t)

	b := &tokenBuilder{
		tokens: [][]int{
			{11, 1, 5, 0, 0},
			{1, 2, 5, 0, 0},
			{1, 12, 5, 0, 0},
			{2, 5, 5, 0, 0},
			{2, 2, 7, 0, 0},
			{2, 2, 5, 0, 0},
		},
	}

	b.sort()
	a.Equal(b.tokens, [][]int{
		{1, 2, 5, 0, 0},
		{1, 12, 5, 0, 0},
		{2, 2, 7, 0, 0},
		{2, 2, 5, 0, 0},
		{2, 5, 5, 0, 0},
		{11, 1, 5, 0, 0},
	})
}

func TestSemanticTokens(t *testing.T) {
	a := assert.New(t)

	// NOTE 此处的 apidoc 属性值必须与当前的文档主版本号相同
	b := `<apidoc version="1.1.1" apidoc="6.0.0" created="2020-01-02T13:12:11+08:00">
	<title>标题</title>
	<mimetype>xml</mimetype>
	<api method="GET">
		<path path="/users" />
		<response status="200" />
	</api>

	<api method="POST">
		<path path="/users" />
		<response status="200" type="number" />
	</api>
</apidoc>`
	blk := core.Block{Data: []byte(b), Location: core.Location{URI: "doc.go"}}
	doc := &ast.APIDoc{}
	rslt := messagetest.NewMessageHandler()
	doc.Parse(rslt.Handler, blk)
	rslt.Handler.Stop()
	a.Empty(rslt.Errors)

	a.Equal(semanticTokens(doc, "doc.go", 1, 2, 3), []int{
		0, 1, 6, 1, 0, // apidoc
		0, 7, 7, 2, 0,
		0, 9, 5, 3, 0,
		0, 7, 6, 2, 0,
		0, 8, 5, 3, 0,
		0, 7, 7, 2, 0,
		0, 9, 25, 3, 0,

		1, 2, 5, 1, 0, // <title>
		// {1, 8, 2, 0, 0}, // 元素内容不作解析，直接采用默认的注释颜色
		0, 10, 5, 1, 0, // </title>

		1, 2, 8, 1, 0, // <mimetype>
		// {2, 11, 3, 0, 0}, // 元素内容不作解析，直接采用默认的注释颜色
		0, 14, 8, 1, 0, // </mimetype>

		1, 2, 3, 1, 0, // <api>
		0, 4, 6, 2, 0,
		0, 8, 3, 3, 0,

		1, 3, 4, 1, 0, // path
		0, 5, 4, 2, 0,
		0, 6, 6, 3, 0,

		1, 3, 8, 1, 0, // response
		0, 9, 6, 2, 0,
		0, 8, 3, 3, 0,

		1, 3, 3, 1, 0, // </api>

		2, 2, 3, 1, 0, // api
		0, 4, 6, 2, 0,
		0, 8, 4, 3, 0,

		1, 3, 4, 1, 0, // path
		0, 5, 4, 2, 0,
		0, 6, 6, 3, 0,

		1, 3, 8, 1, 0, // response
		0, 9, 6, 2, 0,
		0, 8, 3, 3, 0,
		0, 5, 4, 2, 0,
		0, 6, 6, 3, 0,

		1, 3, 3, 1, 0, // </api>

		1, 2, 6, 1, 0, // </apidoc>
	})
}
