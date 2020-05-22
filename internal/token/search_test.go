// SPDX-License-Identifier: MIT

package token

import (
	"reflect"
	"testing"

	"github.com/issue9/assert"

	"github.com/caixw/apidoc/v7/core"
	"github.com/caixw/apidoc/v7/core/messagetest"
)

var _ baser = &Base{}

func TestSearchUsage(t *testing.T) {
	a := assert.New(t)

	b := `<apidoc id="55">
	<name>n</name>
	<name>n</name>
</apidoc>`
	p, err := NewParser(core.Block{Data: []byte(b)})
	a.NotError(err).NotNil(p)
	rslt := messagetest.NewMessageHandler()
	obj := &struct {
		BaseTag
		RootName struct{}    `apidoc:"apidoc,meta,usage-root"`
		ID       intAttr     `apidoc:"id,attr,usage-id"`
		Name     []stringTag `apidoc:"name,elem,usage-name"`
	}{}
	Decode(rslt.Handler, p, obj)
	a.Empty(rslt.Errors)
	a.Equal(obj.Start.Line, 0).
		Equal(obj.Start.Character, 0)
	a.True(obj.End.Line > 0).
		True(obj.End.Character > 0)
	v := reflect.ValueOf(obj)

	r := core.Range{
		Start: core.Position{Line: 0, Character: 1},
		End:   core.Position{Line: 0, Character: 2},
	}
	usage, found := SearchUsage(v, r)
	a.True(found).Equal(usage, "usage-root")

	// id
	r = core.Range{
		Start: core.Position{Line: 0, Character: 8},
		End:   core.Position{Line: 0, Character: 8},
	}
	usage, found = SearchUsage(v, r)
	a.True(found).Equal(usage, "usage-id")

	// 55
	r = core.Range{
		Start: core.Position{Line: 0, Character: 11},
		End:   core.Position{Line: 0, Character: 15},
	}
	usage, found = SearchUsage(v, r)
	a.True(found).Equal(usage, "usage-id")

	// apidoc
	r = core.Range{
		Start: core.Position{Line: 1, Character: 0},
		End:   core.Position{Line: 1, Character: 0},
	}
	usage, found = SearchUsage(v, r)
	a.True(found).Equal(usage, "usage-root")

	// name[0]
	r = core.Range{
		Start: core.Position{Line: 1, Character: 1},
		End:   core.Position{Line: 1, Character: 8},
	}
	usage, found = SearchUsage(v, r)
	a.True(found).Equal(usage, "usage-name")

	// name[1]
	r = core.Range{
		Start: core.Position{Line: 2, Character: 2},
		End:   core.Position{Line: 2, Character: 8},
	}
	usage, found = SearchUsage(v, r)
	a.True(found).Equal(usage, "usage-name")

	// name[0]-name[1]
	r = core.Range{
		Start: core.Position{Line: 1, Character: 2},
		End:   core.Position{Line: 2, Character: 8},
	}
	usage, found = SearchUsage(v, r)
	a.True(found).Equal(usage, "usage-root")

	// apidoc
	r = core.Range{
		Start: core.Position{Line: 3, Character: 1},
		End:   core.Position{Line: 3, Character: 3},
	}
	usage, found = SearchUsage(v, r)
	a.True(found).Equal(usage, "usage-root")

	// apidoc 跨元素
	r = core.Range{
		Start: core.Position{Line: 0, Character: 9},
		End:   core.Position{Line: 1, Character: 3},
	}
	usage, found = SearchUsage(v, r)
	a.True(found).Equal(usage, "usage-root")
}
