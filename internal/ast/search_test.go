// SPDX-License-Identifier: MIT

package ast

import (
	"reflect"
	"testing"

	"github.com/issue9/assert"

	"github.com/caixw/apidoc/v7/core"
	"github.com/caixw/apidoc/v7/internal/xmlenc"
)

func TestSearch(t *testing.T) {
	a := assert.New(t)

	b := xmlenc.Base{
		Range: core.Range{
			Start: core.Position{Character: 1},
			End:   core.Position{Character: 10},
		},
	}
	r := search(reflect.ValueOf(b), core.Position{Character: 5}, nil)
	a.NotNil(r)

	r = search(reflect.ValueOf(&b), core.Position{Character: 5}, nil)
	a.NotNil(r)

	pb := &b
	r = search(reflect.ValueOf(&pb), core.Position{Character: 5}, nil)
	a.NotNil(r)

	// 超出范围
	r = search(reflect.ValueOf(&b), core.Position{Character: 55}, nil)
	a.Nil(r)

	doc := &APIDoc{ // 0,0 - 10,10
		BaseTag: xmlenc.BaseTag{
			Base: xmlenc.Base{
				Range: core.Range{
					End: core.Position{Character: 10, Line: 10},
				},
			},
		},
		Version: &VersionAttribute{ // 1,20 - 2,15
			BaseAttribute: xmlenc.BaseAttribute{
				Base: xmlenc.Base{
					Range: core.Range{
						Start: core.Position{Character: 20, Line: 1},
						End:   core.Position{Character: 15, Line: 2},
					},
				},
			},
			Value: xmlenc.String{ // 1,22 - 2,10
				Range: core.Range{
					Start: core.Position{Character: 22, Line: 1},
					End:   core.Position{Character: 10, Line: 2},
				},
			},
		},
		Tags: []*Tag{
			{ // 3,0 - 3,10
				BaseTag: xmlenc.BaseTag{
					Base: xmlenc.Base{
						Range: core.Range{
							Start: core.Position{Character: 0, Line: 3},
							End:   core.Position{Character: 10, Line: 3},
						},
					},
				},
				Title: &Attribute{ // 3,1 - 3,8
					BaseAttribute: xmlenc.BaseAttribute{
						Base: xmlenc.Base{
							Range: core.Range{
								Start: core.Position{Character: 1, Line: 3},
								End:   core.Position{Character: 8, Line: 3},
							},
						},
					},
				},
			},
			{ // 4,0 - 4,10
				BaseTag: xmlenc.BaseTag{
					Base: xmlenc.Base{
						Range: core.Range{
							Start: core.Position{Character: 0, Line: 4},
							End:   core.Position{Character: 10, Line: 4},
						},
					},
				},
			},
		},
	}

	r = search(reflect.ValueOf(doc), core.Position{}, nil)
	a.NotNil(r).Equal(r.R(), core.Range{End: core.Position{Character: 10, Line: 10}})

	r = search(reflect.ValueOf(doc), core.Position{Character: 100}, nil)
	a.NotNil(r).Equal(r.R(), core.Range{End: core.Position{Character: 10, Line: 10}})

	// 超出范围
	r = search(reflect.ValueOf(doc), core.Position{Line: 100, Character: 100}, nil)
	a.Nil(r)

	// Version.Value
	r = search(reflect.ValueOf(doc), core.Position{Line: 1, Character: 100}, nil)
	a.NotNil(r).Equal(r.R(), core.Range{
		Start: core.Position{Character: 22, Line: 1},
		End:   core.Position{Character: 10, Line: 2},
	})

	// Version
	r = search(reflect.ValueOf(doc), core.Position{Line: 2, Character: 14}, nil)
	a.NotNil(r).Equal(r.R(), core.Range{
		Start: core.Position{Character: 20, Line: 1},
		End:   core.Position{Character: 15, Line: 2},
	})

	// tags[1]
	r = search(reflect.ValueOf(doc), core.Position{Line: 4, Character: 9}, nil)
	a.NotNil(r).Equal(r.R(), core.Range{
		Start: core.Position{Character: 0, Line: 4},
		End:   core.Position{Character: 10, Line: 4},
	})

	// 两个数组元素的中间
	r = search(reflect.ValueOf(doc), core.Position{Line: 4, Character: 11}, nil)
	a.NotNil(r).Equal(r.R(), core.Range{
		End: core.Position{Character: 10, Line: 10},
	})

	// tags[0].title
	r = search(reflect.ValueOf(doc), core.Position{Line: 3, Character: 2}, nil)
	a.NotNil(r).Equal(r.R(), core.Range{
		Start: core.Position{Character: 1, Line: 3},
		End:   core.Position{Character: 8, Line: 3},
	})

	// tags[0]，因为 referenceType 限定，只能搜索到 ast.Tag 实例
	referencerType := reflect.TypeOf((*Referencer)(nil)).Elem()
	r = search(reflect.ValueOf(doc), core.Position{Line: 3, Character: 2}, referencerType)
	a.NotNil(r).Equal(r.R(), core.Range{
		Start: core.Position{Character: 0, Line: 3},
		End:   core.Position{Character: 10, Line: 3},
	})
}
