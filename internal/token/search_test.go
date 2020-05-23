// SPDX-License-Identifier: MIT

package token

import (
	"reflect"
	"testing"

	"github.com/issue9/assert"

	"github.com/caixw/apidoc/v7/core"
	"github.com/caixw/apidoc/v7/core/messagetest"
)

var _ tiper = &Base{}

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

	pos := core.Position{Line: 0, Character: 1}
	tip := SearchUsage(v, pos)
	a.NotNil(tip).Equal(tip.Usage, "usage-root")

	// id
	pos = core.Position{Line: 0, Character: 8}
	tip = SearchUsage(v, pos)
	a.NotNil(tip).Equal(tip.Usage, "usage-id")

	// 55
	pos = core.Position{Line: 0, Character: 11}
	tip = SearchUsage(v, pos)
	a.NotNil(tip).Equal(tip.Usage, "usage-id")

	// apidoc
	pos = core.Position{Line: 1, Character: 0}
	tip = SearchUsage(v, pos)
	a.NotNil(tip).Equal(tip.Usage, "usage-root")

	// name[0]
	pos = core.Position{Line: 1, Character: 1}
	tip = SearchUsage(v, pos)
	a.NotNil(tip).Equal(tip.Usage, "usage-name")

	// name[1]
	pos = core.Position{Line: 2, Character: 2}
	tip = SearchUsage(v, pos)
	a.NotNil(tip).Equal(tip.Usage, "usage-name")

	// apidoc
	pos = core.Position{Line: 3, Character: 1}
	tip = SearchUsage(v, pos)
	a.NotNil(tip).Equal(tip.Usage, "usage-root")
}
