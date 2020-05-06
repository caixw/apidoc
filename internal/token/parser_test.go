// SPDX-License-Identifier: MIT

package token

import (
	"errors"
	"io"
	"testing"

	"github.com/issue9/assert"

	"github.com/caixw/apidoc/v7/core"
)

func TestParser_Token(t *testing.T) {
	a := assert.New(t)
	start := core.Position{
		Line:      11,
		Character: 22,
	}
	uri := core.URI("file:///path")
	data := []*struct {
		input string
		elems []interface{}
		err   *core.SyntaxError
	}{
		{},
		{
			input: `<?xml version="1.0" encoding="utf-8"?>
* <apidoc version="2.0">
*	<title>标题</title>
*	<desc type="html"><![CDATA[<h1>h1</h1>]]></desc>
* </apidoc>
<!-- comment -->  `, // 尾部包含空格
			elems: []interface{}{
				&Instruction{
					Range: core.Range{
						Start: core.Position{Line: 11, Character: 22},
						End:   core.Position{Line: 11, Character: 60},
					},
					Name: String{
						Range: core.Range{
							Start: core.Position{Line: 11, Character: 24},
							End:   core.Position{Line: 11, Character: 27},
						},
						Value: "xml",
					},
					Attributes: []*Attribute{
						{
							Range: core.Range{
								Start: core.Position{Line: 11, Character: 28},
								End:   core.Position{Line: 11, Character: 41},
							},
							Name: String{
								Range: core.Range{
									Start: core.Position{Line: 11, Character: 28},
									End:   core.Position{Line: 11, Character: 35},
								},
								Value: "version",
							},
							Value: String{
								Range: core.Range{
									Start: core.Position{Line: 11, Character: 37},
									End:   core.Position{Line: 11, Character: 40},
								},
								Value: "1.0",
							},
						},
						{
							Range: core.Range{
								Start: core.Position{Line: 11, Character: 42},
								End:   core.Position{Line: 11, Character: 58},
							},
							Name: String{
								Range: core.Range{
									Start: core.Position{Line: 11, Character: 42},
									End:   core.Position{Line: 11, Character: 50},
								},
								Value: "encoding",
							},
							Value: String{
								Range: core.Range{
									Start: core.Position{Line: 11, Character: 52},
									End:   core.Position{Line: 11, Character: 57},
								},
								Value: "utf-8",
							},
						},
					}, // end Instruction.Attributes
				}, // end Instruction
				&String{
					Range: core.Range{
						Start: core.Position{Line: 11, Character: 60},
						End:   core.Position{Line: 12, Character: 2},
					},
					Value: "\n* ",
				},
				&StartElement{
					Range: core.Range{
						Start: core.Position{Line: 12, Character: 2},
						End:   core.Position{Line: 12, Character: 24},
					},
					Name: String{
						Range: core.Range{
							Start: core.Position{Line: 12, Character: 3},
							End:   core.Position{Line: 12, Character: 9},
						},
						Value: "apidoc",
					},
					Attributes: []*Attribute{
						{
							Range: core.Range{
								Start: core.Position{Line: 12, Character: 10},
								End:   core.Position{Line: 12, Character: 23},
							},
							Name: String{
								Range: core.Range{
									Start: core.Position{Line: 12, Character: 10},
									End:   core.Position{Line: 12, Character: 17},
								},
								Value: "version",
							},
							Value: String{
								Range: core.Range{
									Start: core.Position{Line: 12, Character: 19},
									End:   core.Position{Line: 12, Character: 22},
								},
								Value: "2.0",
							},
						},
					},
				}, // end StartElement

				&String{
					Range: core.Range{
						Start: core.Position{Line: 12, Character: 24},
						End:   core.Position{Line: 13, Character: 2},
					},
					Value: "\n*\t",
				},

				&StartElement{
					Range: core.Range{
						Start: core.Position{Line: 13, Character: 2},
						End:   core.Position{Line: 13, Character: 9},
					},
					Name: String{
						Value: "title",
						Range: core.Range{
							Start: core.Position{Line: 13, Character: 3},
							End:   core.Position{Line: 13, Character: 8},
						},
					},
				},

				&String{
					Range: core.Range{
						Start: core.Position{Line: 13, Character: 9},
						End:   core.Position{Line: 13, Character: 11},
					},
					Value: "标题",
				},

				&EndElement{
					Range: core.Range{
						Start: core.Position{Line: 13, Character: 11},
						End:   core.Position{Line: 13, Character: 19},
					},
					Name: String{
						Value: "title",
						Range: core.Range{
							Start: core.Position{Line: 13, Character: 13},
							End:   core.Position{Line: 13, Character: 18},
						},
					},
				},

				&String{
					Range: core.Range{
						Start: core.Position{Line: 13, Character: 19},
						End:   core.Position{Line: 14, Character: 2},
					},
					Value: "\n*\t",
				},

				&StartElement{
					Range: core.Range{
						Start: core.Position{Line: 14, Character: 2},
						End:   core.Position{Line: 14, Character: 20},
					},
					Name: String{
						Value: "desc",
						Range: core.Range{
							Start: core.Position{Line: 14, Character: 3},
							End:   core.Position{Line: 14, Character: 7},
						},
					},
					Attributes: []*Attribute{
						{
							Range: core.Range{
								Start: core.Position{Line: 14, Character: 8},
								End:   core.Position{Line: 14, Character: 19},
							},
							Name: String{
								Range: core.Range{
									Start: core.Position{Line: 14, Character: 8},
									End:   core.Position{Line: 14, Character: 12},
								},
								Value: "type",
							},
							Value: String{
								Range: core.Range{
									Start: core.Position{Line: 14, Character: 14},
									End:   core.Position{Line: 14, Character: 18},
								},
								Value: "html",
							},
						},
					},
				}, // end StartElement

				&CData{
					Range: core.Range{
						Start: core.Position{Line: 14, Character: 20},
						End:   core.Position{Line: 14, Character: 43},
					},
					Value: String{
						Range: core.Range{
							Start: core.Position{Line: 14, Character: 29},
							End:   core.Position{Line: 14, Character: 40},
						},
						Value: "<h1>h1</h1>",
					},
				},

				&EndElement{
					Range: core.Range{
						Start: core.Position{Line: 14, Character: 43},
						End:   core.Position{Line: 14, Character: 50},
					},
					Name: String{
						Range: core.Range{
							Start: core.Position{Line: 14, Character: 45},
							End:   core.Position{Line: 14, Character: 49},
						},
						Value: "desc",
					},
				},

				&String{
					Range: core.Range{
						Start: core.Position{Line: 14, Character: 50},
						End:   core.Position{Line: 15, Character: 2},
					},
					Value: "\n* ",
				},

				&EndElement{
					Range: core.Range{
						Start: core.Position{Line: 15, Character: 2},
						End:   core.Position{Line: 15, Character: 11},
					},
					Name: String{
						Range: core.Range{
							Start: core.Position{Line: 15, Character: 4},
							End:   core.Position{Line: 15, Character: 10},
						},
						Value: "apidoc",
					},
				},

				&String{
					Range: core.Range{
						Start: core.Position{Line: 15, Character: 11},
						End:   core.Position{Line: 16, Character: 0},
					},
					Value: "\n",
				},

				&Comment{
					Range: core.Range{
						Start: core.Position{Line: 16, Character: 0},
						End:   core.Position{Line: 16, Character: 16},
					},
					Value: String{
						Range: core.Range{
							Start: core.Position{Line: 16, Character: 4},
							End:   core.Position{Line: 16, Character: 13},
						},
						Value: " comment ",
					},
				},
				&String{
					Range: core.Range{
						Start: core.Position{Line: 16, Character: 16},
						End:   core.Position{Line: 16, Character: 18},
					},
					Value: "  ",
				},
				nil, nil,
			}, // end elems
		},

		{
			input: `<apidoc version="2.0" /> 
  `, // 尾部包含空格
			elems: []interface{}{
				&StartElement{
					Close: true,
					Range: core.Range{
						Start: core.Position{Line: 11, Character: 22},
						End:   core.Position{Line: 11, Character: 46},
					},
					Name: String{
						Range: core.Range{
							Start: core.Position{Line: 11, Character: 23},
							End:   core.Position{Line: 11, Character: 29},
						},
						Value: "apidoc",
					},
					Attributes: []*Attribute{
						{
							Range: core.Range{
								Start: core.Position{Line: 11, Character: 30},
								End:   core.Position{Line: 11, Character: 43},
							},
							Name: String{
								Range: core.Range{
									Start: core.Position{Line: 11, Character: 30},
									End:   core.Position{Line: 11, Character: 37},
								},
								Value: "version",
							},
							Value: String{
								Range: core.Range{
									Start: core.Position{Line: 11, Character: 39},
									End:   core.Position{Line: 11, Character: 42},
								},
								Value: "2.0",
							},
						},
					},
				}, // end StartElement
				&String{
					Value: " \n  ",
					Range: core.Range{
						Start: core.Position{Line: 11, Character: 46},
						End:   core.Position{Line: 12, Character: 2},
					},
				},
				nil,
			},
		},

		{
			input: `<apidoc version="2.0">123
	</apidoc>
  `, // 尾部包含空格
			elems: []interface{}{
				&StartElement{
					Range: core.Range{
						Start: core.Position{Line: 11, Character: 22},
						End:   core.Position{Line: 11, Character: 44},
					},
					Name: String{
						Range: core.Range{
							Start: core.Position{Line: 11, Character: 23},
							End:   core.Position{Line: 11, Character: 29},
						},
						Value: "apidoc",
					},
					Attributes: []*Attribute{
						{
							Range: core.Range{
								Start: core.Position{Line: 11, Character: 30},
								End:   core.Position{Line: 11, Character: 43},
							},
							Name: String{
								Range: core.Range{
									Start: core.Position{Line: 11, Character: 30},
									End:   core.Position{Line: 11, Character: 37},
								},
								Value: "version",
							},
							Value: String{
								Range: core.Range{
									Start: core.Position{Line: 11, Character: 39},
									End:   core.Position{Line: 11, Character: 42},
								},
								Value: "2.0",
							},
						},
					},
				}, // end StartElement
				&String{
					Value: "123\n\t",
					Range: core.Range{
						Start: core.Position{Line: 11, Character: 44},
						End:   core.Position{Line: 12, Character: 1},
					},
				},
				&EndElement{
					Range: core.Range{
						Start: core.Position{Line: 12, Character: 1},
						End:   core.Position{Line: 12, Character: 10},
					},
					Name: String{
						Range: core.Range{
							Start: core.Position{Line: 12, Character: 3},
							End:   core.Position{Line: 12, Character: 9},
						},
						Value: "apidoc",
					},
				},
				&String{
					Value: "\n  ",
					Range: core.Range{
						Start: core.Position{Line: 12, Character: 10},
						End:   core.Position{Line: 13, Character: 2},
					},
				},
				nil, nil,
			},
		},

		{ // 嵌套自闭合对象
			input: `<apidoc><apidoc /></apidoc> `,
			elems: []interface{}{
				&StartElement{
					Range: core.Range{
						Start: core.Position{Line: 11, Character: 22},
						End:   core.Position{Line: 11, Character: 30},
					},
					Name: String{
						Range: core.Range{
							Start: core.Position{Line: 11, Character: 23},
							End:   core.Position{Line: 11, Character: 29},
						},
						Value: "apidoc",
					},
				}, // end StartElement
				&StartElement{
					Close: true,
					Range: core.Range{
						Start: core.Position{Line: 11, Character: 30},
						End:   core.Position{Line: 11, Character: 40},
					},
					Name: String{
						Range: core.Range{
							Start: core.Position{Line: 11, Character: 31},
							End:   core.Position{Line: 11, Character: 37},
						},
						Value: "apidoc",
					},
				},
				&EndElement{
					Range: core.Range{
						Start: core.Position{Line: 11, Character: 40},
						End:   core.Position{Line: 11, Character: 49},
					},
					Name: String{
						Range: core.Range{
							Start: core.Position{Line: 11, Character: 42},
							End:   core.Position{Line: 11, Character: 48},
						},
						Value: "apidoc",
					},
				},
			},
		},
	}

	for _, item := range data {
		p, err := NewParser(core.Block{
			Location: core.Location{
				URI:   uri,
				Range: core.Range{Start: start},
			},
			Data: []byte(item.input),
		})
		a.NotError(err, "error %s at %s", err, item.input).
			NotNil(p, "nil at %s", item.input)

		for i, elem := range item.elems {
			e, err := p.Token()
			if elem == nil {
				a.Equal(err, io.EOF, "%s no io.EOF at %s:%d", err, item.input, i)
				a.Nil(e, "not nil at %s:%d", item.input, i)
			} else {
				a.NotError(err, "error %s at %d", err, i)
				a.Equal(e, elem, "not equal at %s:%d\nv1=%#v\nv2=%#v", item.input, i, e, elem)
			}
		}
	}
}

func TestParser_parseStartElement(t *testing.T) {
	a := assert.New(t)
	start := core.Position{
		Line:      11,
		Character: 22,
	}
	uri := core.URI("file:///path")

	data := []*struct {
		input string
		elem  *StartElement
		err   *core.SyntaxError
	}{
		{
			err: &core.SyntaxError{
				Location: core.Location{
					URI: uri,
					Range: core.Range{
						Start: core.Position{Line: 11, Character: 22},
						End:   core.Position{Line: 11, Character: 22},
					},
				},
			},
		},
		{ // 没有结束标签
			input: `tag version="1.0"`,
			err: &core.SyntaxError{
				Location: core.Location{
					URI: uri,
					Range: core.Range{
						Start: core.Position{Line: 11, Character: 39},
						End:   core.Position{Line: 11, Character: 39},
					},
				},
			},
		},
		{ // 没有标签名
			input: ">",
			err: &core.SyntaxError{
				Location: core.Location{
					URI: uri,
					Range: core.Range{
						Start: core.Position{Line: 11, Character: 22},
						End:   core.Position{Line: 11, Character: 22},
					},
				},
			},
		},
		{
			input: `tag>`,
			elem: &StartElement{
				Range: core.Range{
					Start: core.Position{Line: 11, Character: 22},
					End:   core.Position{Line: 11, Character: 26},
				},
				Name: String{
					Range: core.Range{
						Start: core.Position{Line: 11, Character: 22},
						End:   core.Position{Line: 11, Character: 25},
					},
					Value: "tag",
				},
			},
		},
		{
			input: `tag/>`,
			elem: &StartElement{
				Range: core.Range{
					Start: core.Position{Line: 11, Character: 22},
					End:   core.Position{Line: 11, Character: 27},
				},
				Name: String{
					Range: core.Range{
						Start: core.Position{Line: 11, Character: 22},
						End:   core.Position{Line: 11, Character: 25},
					},
					Value: "tag",
				},
				Close: true,
			},
		},

		{
			input: `tag ver="1.0" enc="utf8"/>`,
			elem: &StartElement{
				Range: core.Range{
					Start: core.Position{Line: 11, Character: 22},
					End:   core.Position{Line: 11, Character: 48},
				},
				Name: String{
					Range: core.Range{
						Start: core.Position{Line: 11, Character: 22},
						End:   core.Position{Line: 11, Character: 25},
					},
					Value: "tag",
				},
				Attributes: []*Attribute{
					{
						Range: core.Range{
							Start: core.Position{Line: 11, Character: 26},
							End:   core.Position{Line: 11, Character: 35},
						},
						Name: String{
							Range: core.Range{
								Start: core.Position{Line: 11, Character: 26},
								End:   core.Position{Line: 11, Character: 29},
							},
							Value: "ver",
						},
						Value: String{
							Range: core.Range{
								Start: core.Position{Line: 11, Character: 31},
								End:   core.Position{Line: 11, Character: 34},
							},
							Value: "1.0",
						},
					},
					{
						Range: core.Range{
							Start: core.Position{Line: 11, Character: 36},
							End:   core.Position{Line: 11, Character: 46},
						},
						Name: String{
							Range: core.Range{
								Start: core.Position{Line: 11, Character: 36},
								End:   core.Position{Line: 11, Character: 39},
							},
							Value: "enc",
						},
						Value: String{
							Range: core.Range{
								Start: core.Position{Line: 11, Character: 41},
								End:   core.Position{Line: 11, Character: 45},
							},
							Value: "utf8",
						},
					},
				}, // end Attributes
				Close: true,
			},
		},
	}

	for _, item := range data {
		p, err := NewParser(core.Block{
			Location: core.Location{
				URI:   uri,
				Range: core.Range{Start: start},
			},
			Data: []byte(item.input),
		})
		a.NotError(err, "error %s at %s", err, item.input).
			NotNil(p, "nil at %s", item.input)

		elem, err := p.parseStartElement(p.Position())
		if item.err != nil {
			serr, ok := err.(*core.SyntaxError)
			a.True(ok, "false at %s", item.input).
				Equal(serr.Location, item.err.Location, "not equal at %s\nv1=%+v\nv2=%+v", item.input, serr.Location, item.err.Location)
		} else {
			a.NotError(err, "error %s at %s", err, item.input).
				Equal(elem, item.elem, "not equal at %s\nv1=%+v\nv2=%+v", item.input, elem, item.elem)
		}
	}
}

func TestParser_parseEndElement(t *testing.T) {
	a := assert.New(t)
	start := core.Position{
		Line:      11,
		Character: 22,
	}
	uri := core.URI("file:///path")

	data := []*struct {
		input string
		elem  *EndElement
		err   *core.SyntaxError
	}{
		{
			err: &core.SyntaxError{
				Location: core.Location{
					URI: uri,
					Range: core.Range{
						Start: core.Position{Line: 11, Character: 22},
						End:   core.Position{Line: 11, Character: 22},
					},
				},
			},
		},
		{
			input: ">",
			err: &core.SyntaxError{
				Location: core.Location{
					URI: uri,
					Range: core.Range{
						Start: core.Position{Line: 11, Character: 22},
						End:   core.Position{Line: 11, Character: 22},
					},
				},
			},
		},
		{
			input: `tag>`,
			elem: &EndElement{
				Range: core.Range{
					Start: core.Position{Line: 11, Character: 22},
					End:   core.Position{Line: 11, Character: 26},
				},
				Name: String{
					Range: core.Range{
						Start: core.Position{Line: 11, Character: 22},
						End:   core.Position{Line: 11, Character: 25},
					},
					Value: "tag",
				},
			},
		},
	}

	for _, item := range data {
		p, err := NewParser(core.Block{
			Location: core.Location{
				URI:   uri,
				Range: core.Range{Start: start},
			},
			Data: []byte(item.input),
		})
		a.NotError(err, "error %s at %s", err, item.input).
			NotNil(p, "nil at %s", item.input)

		elem, err := p.parseEndElement(p.Position())
		if item.err != nil {
			serr, ok := err.(*core.SyntaxError)
			a.True(ok, "false at %s", item.input).
				Equal(serr.Location, item.err.Location, "not equal at %s\nv1=%+v\nv2=%+v", item.input, serr.Location, item.err.Location)
		} else {
			a.NotError(err, "error %s at %s", err, item.input).
				Equal(elem, item.elem, "not equal at %s\nv1=%+v\nv2=%+v", item.input, elem, item.elem)
		}
	}
}

func TestParser_parseCData(t *testing.T) {
	a := assert.New(t)
	start := core.Position{
		Line:      11,
		Character: 22,
	}
	uri := core.URI("file:///path")

	data := []*struct {
		input string
		cdata *CData
		err   *core.SyntaxError
	}{
		{
			err: &core.SyntaxError{
				Location: core.Location{
					URI: uri,
					Range: core.Range{
						Start: core.Position{Line: 11, Character: 22},
						End:   core.Position{Line: 11, Character: 22},
					},
				},
			},
		},
		{
			input: `<h1></h1>`,
			err: &core.SyntaxError{
				Location: core.Location{
					URI: uri,
					Range: core.Range{
						Start: core.Position{Line: 11, Character: 22},
						End:   core.Position{Line: 11, Character: 22},
					},
				},
			},
		},
		{
			input: "<h1>\nxxx]]>",
			cdata: &CData{
				Range: core.Range{
					Start: core.Position{Line: 11, Character: 22},
					End:   core.Position{Line: 12, Character: 6},
				},
				Value: String{
					Range: core.Range{
						Start: core.Position{Line: 11, Character: 22},
						End:   core.Position{Line: 12, Character: 3},
					},
					Value: "<h1>\nxxx",
				},
			},
		},
		{ // cdata 转义
			input: "<h1>]]]]><![CDATA[>\nxxx]]>",
			cdata: &CData{
				Range: core.Range{
					Start: core.Position{Line: 11, Character: 22},
					End:   core.Position{Line: 12, Character: 6},
				},
				Value: String{
					Range: core.Range{
						Start: core.Position{Line: 11, Character: 22},
						End:   core.Position{Line: 12, Character: 3},
					},
					Value: "<h1>]]]]><![CDATA[>\nxxx",
				},
			},
		},

		{
			input: "<h1>]]]]><![CDATA[>\n12]]]]><![CDATA[>34\nxxx]]>",
			cdata: &CData{
				Range: core.Range{
					Start: core.Position{Line: 11, Character: 22},
					End:   core.Position{Line: 13, Character: 6},
				},
				Value: String{
					Range: core.Range{
						Start: core.Position{Line: 11, Character: 22},
						End:   core.Position{Line: 13, Character: 3},
					},
					Value: "<h1>]]]]><![CDATA[>\n12]]]]><![CDATA[>34\nxxx",
				},
			},
		},
	}

	for _, item := range data {
		p, err := NewParser(core.Block{
			Location: core.Location{
				URI:   uri,
				Range: core.Range{Start: start},
			},
			Data: []byte(item.input),
		})
		a.NotError(err, "error %s at %s", err, item.input).
			NotNil(p, "nil at %s", item.input)

		cdata, err := p.parseCData(p.Position())
		if item.err != nil {
			serr, ok := err.(*core.SyntaxError)
			a.True(ok, "false at %s", item.input).
				Equal(serr.Location, item.err.Location, "not equal at %s\nv1=%+v\nv2=%+v", item.input, serr.Location, item.err.Location)
		} else {
			a.NotError(err, "error %s at %s", err, item.input).
				Equal(cdata, item.cdata, "not equal at %s\nv1=%+v\nv2=%+v", item.input, cdata, item.cdata)
		}
	}
}

func TestParser_parseInstruction(t *testing.T) {
	a := assert.New(t)
	start := core.Position{
		Line:      11,
		Character: 22,
	}
	uri := core.URI("file:///path")

	data := []*struct {
		input string
		pi    *Instruction
		err   *core.SyntaxError
	}{
		{
			err: &core.SyntaxError{
				Location: core.Location{
					URI: uri,
					Range: core.Range{
						Start: core.Position{Line: 11, Character: 22},
						End:   core.Position{Line: 11, Character: 22},
					},
				},
			},
		},
		{ // 缺少结束符号 ?>
			input: `xml version="1.0"`,
			err: &core.SyntaxError{
				Location: core.Location{
					URI: uri,
					Range: core.Range{
						Start: core.Position{Line: 11, Character: 39},
						End:   core.Position{Line: 11, Character: 39},
					},
				},
			},
		},
		{ // version 被当作标签名，之后找不到 ?>
			input: `version="1.0"`,
			err: &core.SyntaxError{
				Location: core.Location{
					URI: uri,
					Range: core.Range{
						Start: core.Position{Line: 11, Character: 29},
						End:   core.Position{Line: 11, Character: 29},
					},
				},
			},
		},
		{
			input: `xml?>version="1 "?>`,
			pi: &Instruction{
				Range: core.Range{
					Start: core.Position{Line: 11, Character: 22},
					End:   core.Position{Line: 11, Character: 27},
				},
				Name: String{
					Value: "xml",
					Range: core.Range{
						Start: core.Position{Line: 11, Character: 22},
						End:   core.Position{Line: 11, Character: 25},
					},
				},
			},
		},
		{ // 属性值中包含空格
			input: `xml version="1 0" ?>`,
			pi: &Instruction{
				Range: core.Range{
					Start: core.Position{Line: 11, Character: 22},
					End:   core.Position{Line: 11, Character: 42},
				},
				Name: String{
					Value: "xml",
					Range: core.Range{
						Start: core.Position{Line: 11, Character: 22},
						End:   core.Position{Line: 11, Character: 25},
					},
				},
				Attributes: []*Attribute{
					{
						Range: core.Range{
							Start: core.Position{Line: 11, Character: 26},
							End:   core.Position{Line: 11, Character: 39},
						},
						Name: String{
							Range: core.Range{
								Start: core.Position{Line: 11, Character: 26},
								End:   core.Position{Line: 11, Character: 33},
							},
							Value: "version",
						},
						Value: String{
							Range: core.Range{
								Start: core.Position{Line: 11, Character: 35},
								End:   core.Position{Line: 11, Character: 38},
							},
							Value: "1 0",
						},
					},
				}, // end Attributes
			},
		},
		{
			input: `xml version="1.0" encoding="utf-8"?>`,
			pi: &Instruction{
				Range: core.Range{
					Start: core.Position{Line: 11, Character: 22},
					End:   core.Position{Line: 11, Character: 58},
				},
				Name: String{
					Value: "xml",
					Range: core.Range{
						Start: core.Position{Line: 11, Character: 22},
						End:   core.Position{Line: 11, Character: 25},
					},
				},
				Attributes: []*Attribute{
					{
						Range: core.Range{
							Start: core.Position{Line: 11, Character: 26},
							End:   core.Position{Line: 11, Character: 39},
						},
						Name: String{
							Range: core.Range{
								Start: core.Position{Line: 11, Character: 26},
								End:   core.Position{Line: 11, Character: 33},
							},
							Value: "version",
						},
						Value: String{
							Range: core.Range{
								Start: core.Position{Line: 11, Character: 35},
								End:   core.Position{Line: 11, Character: 38},
							},
							Value: "1.0",
						},
					},
					{
						Range: core.Range{
							Start: core.Position{Line: 11, Character: 40},
							End:   core.Position{Line: 11, Character: 56},
						},
						Name: String{
							Range: core.Range{
								Start: core.Position{Line: 11, Character: 40},
								End:   core.Position{Line: 11, Character: 48},
							},
							Value: "encoding",
						},
						Value: String{
							Range: core.Range{
								Start: core.Position{Line: 11, Character: 50},
								End:   core.Position{Line: 11, Character: 55},
							},
							Value: "utf-8",
						},
					},
				}, // end Attributes
			},
		},
	}

	for _, item := range data {
		p, err := NewParser(core.Block{
			Location: core.Location{
				URI:   uri,
				Range: core.Range{Start: start},
			},
			Data: []byte(item.input),
		})
		a.NotError(err, "error %s at %s", err, item.input).
			NotNil(p, "nil at %s", item.input)

		pi, err := p.parseInstruction(p.Position())
		if item.err != nil {
			serr, ok := err.(*core.SyntaxError)
			a.True(ok, "false at %s", item.input).
				Equal(serr.Location, item.err.Location, "not equal at %s\nv1=%+v\nv2=%+v", item.input, serr.Location, item.err.Location)
		} else {
			a.NotError(err, "error %s at %s", err, item.input).
				Equal(pi, item.pi, "not equal at %s\nv1=%+v\nv2=%+v", item.input, pi, item.pi)
		}
	}
}

func TestParser_parseAttributes(t *testing.T) {
	a := assert.New(t)
	start := core.Position{
		Line:      11,
		Character: 22,
	}
	uri := core.URI("file:///path")

	data := []*struct {
		input string
		attrs []*Attribute
		err   *core.SyntaxError
	}{
		{},
		{
			input: `name="value"`,
			attrs: []*Attribute{
				{
					Range: core.Range{
						Start: core.Position{Line: 11, Character: 22},
						End:   core.Position{Line: 11, Character: 34},
					},
					Name: String{
						Value: "name",
						Range: core.Range{
							Start: core.Position{Line: 11, Character: 22},
							End:   core.Position{Line: 11, Character: 26},
						},
					},
					Value: String{
						Value: "value",
						Range: core.Range{
							Start: core.Position{Line: 11, Character: 28},
							End:   core.Position{Line: 11, Character: 33},
						},
					},
				},
			},
		},

		{
			input: `name="value"
	name="value"`,
			attrs: []*Attribute{
				{
					Range: core.Range{
						Start: core.Position{Line: 11, Character: 22},
						End:   core.Position{Line: 11, Character: 34},
					},
					Name: String{
						Value: "name",
						Range: core.Range{
							Start: core.Position{Line: 11, Character: 22},
							End:   core.Position{Line: 11, Character: 26},
						},
					},
					Value: String{
						Value: "value",
						Range: core.Range{
							Start: core.Position{Line: 11, Character: 28},
							End:   core.Position{Line: 11, Character: 33},
						},
					},
				},
				{
					Range: core.Range{
						Start: core.Position{Line: 12, Character: 1},
						End:   core.Position{Line: 12, Character: 13},
					},
					Name: String{
						Value: "name",
						Range: core.Range{
							Start: core.Position{Line: 12, Character: 1},
							End:   core.Position{Line: 12, Character: 5},
						},
					},
					Value: String{
						Value: "value",
						Range: core.Range{
							Start: core.Position{Line: 12, Character: 7},
							End:   core.Position{Line: 12, Character: 12},
						},
					},
				},
			},
		},

		{
			input: `name="" xx=`,
			err: &core.SyntaxError{
				Location: core.Location{
					URI: uri,
					Range: core.Range{
						Start: core.Position{Line: 11, Character: 33},
						End:   core.Position{Line: 11, Character: 33},
					},
				},
			},
		},
	}

	for _, item := range data {
		p, err := NewParser(core.Block{
			Location: core.Location{
				URI:   uri,
				Range: core.Range{Start: start},
			},
			Data: []byte(item.input),
		})
		a.NotError(err, "error %s at %s", err, item.input).
			NotNil(p, "nil at %s", item.input)

		attrs, err := p.parseAttributes()
		if item.err != nil {
			serr, ok := err.(*core.SyntaxError)
			a.True(ok).
				Equal(serr.Location, item.err.Location, "not equal at %s\nv1=%+v\nv2=%+v", item.input, serr.Location, item.err.Location)
			break
		}
		a.NotError(err, "error %s at %s", err, item.input).
			Equal(len(attrs), len(item.attrs), "not equal at %s,v1=%+v,v2=%+v", item.input, len(attrs), len(item.attrs))

		for i, attr := range attrs {
			a.Equal(attr, item.attrs[i], "not equal at %s:%d\nv1=%+v\nv2=%+v", item.input, i, attr, item.attrs[i])
		}
	}
}

func TestParser_parseAttribute(t *testing.T) {
	a := assert.New(t)
	start := core.Position{
		Line:      11,
		Character: 22,
	}
	uri := core.URI("file:///path")

	data := []*struct {
		input string
		attr  *Attribute
		err   *core.SyntaxError
	}{
		{},
		{
			input: `name="value"`,
			attr: &Attribute{
				Range: core.Range{
					Start: core.Position{Line: 11, Character: 22},
					End:   core.Position{Line: 11, Character: 34},
				},
				Name: String{
					Value: "name",
					Range: core.Range{
						Start: core.Position{Line: 11, Character: 22},
						End:   core.Position{Line: 11, Character: 26},
					},
				},
				Value: String{
					Value: "value",
					Range: core.Range{
						Start: core.Position{Line: 11, Character: 28},
						End:   core.Position{Line: 11, Character: 33},
					},
				},
			},
		},
		{ // 属性值包含 =
			input: `name="val=e"`,
			attr: &Attribute{
				Range: core.Range{
					Start: core.Position{Line: 11, Character: 22},
					End:   core.Position{Line: 11, Character: 34},
				},
				Name: String{
					Value: "name",
					Range: core.Range{
						Start: core.Position{Line: 11, Character: 22},
						End:   core.Position{Line: 11, Character: 26},
					},
				},
				Value: String{
					Value: "val=e",
					Range: core.Range{
						Start: core.Position{Line: 11, Character: 28},
						End:   core.Position{Line: 11, Character: 33},
					},
				},
			},
		},
		{ // 属性值包含 >
			input: `name="val>e"`,
			attr: &Attribute{
				Range: core.Range{
					Start: core.Position{Line: 11, Character: 22},
					End:   core.Position{Line: 11, Character: 34},
				},
				Name: String{
					Value: "name",
					Range: core.Range{
						Start: core.Position{Line: 11, Character: 22},
						End:   core.Position{Line: 11, Character: 26},
					},
				},
				Value: String{
					Value: "val>e",
					Range: core.Range{
						Start: core.Position{Line: 11, Character: 28},
						End:   core.Position{Line: 11, Character: 33},
					},
				},
			},
		},
		{
			input: "\tname\t=\n\"value\"",
			attr: &Attribute{
				Range: core.Range{
					Start: core.Position{Line: 11, Character: 23},
					End:   core.Position{Line: 12, Character: 7},
				},
				Name: String{
					Value: "name",
					Range: core.Range{
						Start: core.Position{Line: 11, Character: 23},
						End:   core.Position{Line: 11, Character: 27},
					},
				},
				Value: String{
					Value: "value",
					Range: core.Range{
						Start: core.Position{Line: 12, Character: 1},
						End:   core.Position{Line: 12, Character: 6},
					},
				},
			},
		},
		{ // 空的属性值
			input: "\tname\t=\n\"\"",
			attr: &Attribute{
				Range: core.Range{
					Start: core.Position{Line: 11, Character: 23},
					End:   core.Position{Line: 12, Character: 2},
				},
				Name: String{
					Value: "name",
					Range: core.Range{
						Start: core.Position{Line: 11, Character: 23},
						End:   core.Position{Line: 11, Character: 27},
					},
				},
				Value: String{
					Value: "",
					Range: core.Range{
						Start: core.Position{Line: 12, Character: 1},
						End:   core.Position{Line: 12, Character: 1},
					},
				},
			},
		},

		{
			input: `name `,
			err: &core.SyntaxError{
				Location: core.Location{
					URI: uri,
					Range: core.Range{
						Start: core.Position{Line: 11, Character: 27},
						End:   core.Position{Line: 11, Character: 27},
					},
				},
			},
		},
		{
			input: `name=`,
			err: &core.SyntaxError{
				Location: core.Location{
					URI: uri,
					Range: core.Range{
						Start: core.Position{Line: 11, Character: 27},
						End:   core.Position{Line: 11, Character: 27},
					},
				},
			},
		},
		{
			input: `name="`,
			err: &core.SyntaxError{
				Location: core.Location{
					URI: uri,
					Range: core.Range{
						Start: core.Position{Line: 11, Character: 28},
						End:   core.Position{Line: 11, Character: 28},
					},
				},
			},
		},
		{
			input: `name=="val"`,
			err: &core.SyntaxError{
				Location: core.Location{
					URI: uri,
					Range: core.Range{
						Start: core.Position{Line: 11, Character: 27},
						End:   core.Position{Line: 11, Character: 27},
					},
				},
			},
		},
	}

	for _, item := range data {
		p, err := NewParser(core.Block{
			Location: core.Location{
				URI:   uri,
				Range: core.Range{Start: start},
			},
			Data: []byte(item.input),
		})
		a.NotError(err, "error %s at %s", err, item.input).
			NotNil(p, "nil at %s", item.input)

		attr, err := p.parseAttribute()
		if item.err != nil {
			serr, ok := err.(*core.SyntaxError)
			a.True(ok, "false at %s", item.input).
				Equal(serr.Location, item.err.Location, "not equal at %s\nv1=%+v\nv2=%+v", item.input, serr.Location, item.err.Location)
		} else {
			a.NotError(err, "error %s at %s", err, item.input).
				Equal(attr, item.attr, "not equal at %s\nv1=%+v\nv2=%+v", item.input, attr, item.attr)
		}
	}
}

func TestParser_WithError(t *testing.T) {
	a := assert.New(t)

	err1 := errors.New("err1")
	p, err := NewParser(core.Block{})
	a.NotError(err).NotNil(p)

	err = p.WithError(core.Position{}, core.Position{}, "field1", err1)
	serr, ok := err.(*core.SyntaxError)
	a.True(ok).Equal(serr.Err, err1)

	err2 := core.NewSyntaxErrorWithError(core.Location{}, "", err1)
	err = p.WithError(core.Position{}, core.Position{}, "field1", err2)
	serr, ok = err.(*core.SyntaxError)
	a.True(ok).Equal(serr.Err, err1)

	err3 := core.NewSyntaxErrorWithError(core.Location{}, "", err2)
	err = p.WithError(core.Position{}, core.Position{}, "field1", err3)
	serr, ok = err.(*core.SyntaxError)
	a.True(ok).Equal(serr.Err, err1)

	err4 := core.NewSyntaxErrorWithError(core.Location{}, "", err3)
	err = p.WithError(core.Position{}, core.Position{}, "field1", err4)
	serr, ok = err.(*core.SyntaxError)
	a.True(ok).Equal(serr.Err, err1)
}