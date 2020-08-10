// SPDX-License-Identifier: MIT

package ast

import (
	"reflect"
	"testing"

	"github.com/issue9/assert"

	"github.com/caixw/apidoc/v7/core"
	"github.com/caixw/apidoc/v7/core/messagetest"
)

func TestSearch(t *testing.T) {
	a := assert.New(t)
	data := `<apidoc version="1.0.0">
		<title>title</title>
		<tag name="t1" title="tag1" />
		<tag name="t2" title="tag2" />
		<api method="GET">
			<tag>t1</tag>
			<path path="/users" />
			<response status="200" />
		</api>
		<api method="POST">
			<tag>t1</tag>
			<tag>t2</tag>
			<path path="/users" />
			<response status="200" />
		</api>
	</apidoc>`

	rslt := messagetest.NewMessageHandler()
	doc := &APIDoc{}
	doc.Parse(rslt.Handler, core.Block{Data: []byte(data), Location: core.Location{URI: "doc.go"}})
	rslt.Handler.Stop()

	r := doc.Search("doc.go", core.Position{}, nil)
	a.NotNil(r).Equal(r.R(), core.Range{End: core.Position{Line: 15, Character: 10}})

	r = doc.Search("doc.go", core.Position{Character: 100}, nil)
	a.NotNil(r).Equal(r.R(), core.Range{End: core.Position{Line: 15, Character: 10}})

	// 超出范围
	r = doc.Search("doc.go", core.Position{Line: 100, Character: 100}, nil)
	a.Nil(r)

	// title.title
	r = doc.Search("doc.go", core.Position{Line: 1, Character: 10}, nil)
	a.NotNil(r).Equal(r.R(), core.Range{
		Start: core.Position{Line: 1, Character: 9},
		End:   core.Position{Line: 1, Character: 14},
	})

	// title.title，URI 不匹配
	r = doc.Search("not.exists", core.Position{Line: 1, Character: 10}, nil)
	a.Nil(r)

	// tags[0].name
	r = doc.Search("doc.go", core.Position{Line: 2, Character: 9}, nil)
	a.NotNil(r).Equal(r.R(), core.Range{
		Start: core.Position{Line: 2, Character: 7},
		End:   core.Position{Line: 2, Character: 16},
	})

	// 两个数组元素的中间
	r = doc.Search("doc.go", core.Position{Line: 2, Character: 40}, nil)
	a.NotNil(r).Equal(r.R(), core.Range{
		End: core.Position{Line: 15, Character: 10},
	})

	// tags[1].title
	r = doc.Search("doc.go", core.Position{Line: 3, Character: 17}, nil)
	a.NotNil(r).Equal(r.R(), core.Range{
		Start: core.Position{Line: 3, Character: 17},
		End:   core.Position{Line: 3, Character: 29},
	})

	// tags[0]，因为 referenceType 限定，只能搜索到 ast.Tag 实例
	referencerType := reflect.TypeOf((*Referencer)(nil)).Elem()
	r = search(reflect.ValueOf(doc), core.Position{Line: 3, Character: 17}, referencerType)
	a.NotNil(r).Equal(r.R(), core.Range{
		Start: core.Position{Line: 3, Character: 2},
		End:   core.Position{Line: 3, Character: 32},
	})
	refs := r.(Referencer).References()
	a.Equal(1, len(refs))
	a.NotNil(refs[0].Target).
		Equal(refs[0].Location, core.Location{
			URI: "doc.go",
			Range: core.Range{
				Start: core.Position{Line: 11, Character: 3},
				End:   core.Position{Line: 11, Character: 16},
			},
		})

	// api[0]
	r = doc.Search("doc.go", core.Position{Line: 5, Character: 3}, nil)
	a.NotNil(r).Equal(r.R(), core.Range{
		Start: core.Position{Line: 5, Character: 3},
		End:   core.Position{Line: 5, Character: 16},
	})

	// api[0]，URI 不匹配
	r = doc.Search("not-exists", core.Position{Line: 5, Character: 3}, nil)
	a.Nil(r)

	// api[0]，不匹配 api，匹配至整个 apidoc
	doc.APIs[0].URI = "api.go"
	r = doc.Search("doc.go", core.Position{Line: 5, Character: 3}, nil)
	a.NotNil(r).Equal(r.R(), core.Range{
		End: core.Position{Line: 15, Character: 10},
	})
}
