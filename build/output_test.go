// SPDX-License-Identifier: MIT

package build

import (
	"testing"

	"github.com/issue9/assert/v3"

	"github.com/caixw/apidoc/v7/core"
	"github.com/caixw/apidoc/v7/internal/ast/asttest"
	"github.com/caixw/apidoc/v7/internal/docs"
)

func TestOptions_contains(t *testing.T) {
	a := assert.New(t, false)

	o := &Output{}
	a.True(o.contains("tag"))
	a.True(o.contains(""))

	o.Tags = []string{"t1", "t2"}
	a.True(o.contains("t1"))
	a.False(o.contains("not-exists"))
	a.False(o.contains(""))
}

func TestOutput_Sanitize(t *testing.T) {
	a := assert.New(t, false)

	// 默认的 Type
	o := &Output{}
	a.NotError(o.sanitize())
	a.Equal(o.Type, APIDocXML).NotNil(o.marshal)

	o = &Output{Type: "invalid-type"}
	a.Error(o.sanitize())

	o = &Output{Type: APIDocXML}
	o.Path = "./testdir/apidoc.json"
	a.NotError(o.sanitize())
	a.Equal(o.Style, docs.StylesheetURL(core.OfficialURL)).
		Equal(2, len(o.procInst)).
		Contains(o.procInst[1], docs.StylesheetURL(core.OfficialURL))

	o.Version = "1.0.0"
	a.NotError(o.sanitize())
	o.Version = "1"
	a.Error(o.sanitize())
}

func TestOptions_buffer(t *testing.T) {
	a := assert.New(t, false)

	doc := asttest.Get()
	o := &Output{
		Type: OpenapiJSON,
		Path: "./openapi.json",
	}
	a.NotError(o.sanitize())
	_, err := o.buffer(doc)
	a.NotError(err)

	doc = asttest.Get()
	o = &Output{}
	a.NotError(o.sanitize())
	buf, err := o.buffer(doc)
	a.NotError(err).NotNil(buf)
}

func TestFilterDoc(t *testing.T) {
	a := assert.New(t, false)

	d := asttest.Get()
	o := &Output{}
	a.NotError(o.sanitize())
	filterDoc(d, o)
	a.Equal(3, len(d.Tags))

	d = asttest.Get()
	o = &Output{
		Tags: []string{"t1"},
	}
	a.NotError(o.sanitize())
	filterDoc(d, o)
	a.Equal(1, len(d.Tags)).
		Equal(2, len(d.APIs))

	d = asttest.Get()
	o = &Output{
		Tags: []string{"t1", "t2"},
	}
	a.NotError(o.sanitize())
	filterDoc(d, o)
	a.Equal(2, len(d.Tags)).
		Equal(2, len(d.APIs))

	d = asttest.Get()
	o = &Output{
		Tags: []string{"tag1"},
	}
	a.NotError(o.sanitize())
	filterDoc(d, o)
	a.Equal(1, len(d.Tags)).
		Equal(1, len(d.APIs))

	d = asttest.Get()
	o = &Output{
		Tags: []string{"not-exists"},
	}
	a.NotError(o.sanitize())
	filterDoc(d, o)
	a.Equal(0, len(d.Tags)).
		Equal(0, len(d.APIs))
}
