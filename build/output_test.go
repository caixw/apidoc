// SPDX-License-Identifier: MIT

package build

import (
	"testing"

	"github.com/issue9/assert"

	"github.com/caixw/apidoc/v6/internal/ast/asttest"
)

func TestStylesheet(t *testing.T) {
	a := assert.New(t)
	a.NotEmpty(stylesheetURL)
}

func TestOptions_contains(t *testing.T) {
	a := assert.New(t)

	o := &Output{}
	a.True(o.contains("tag"))
	a.True(o.contains(""))

	o.Tags = []string{"t1", "t2"}
	a.True(o.contains("t1"))
	a.False(o.contains("not-exists"))
	a.False(o.contains(""))
}

func TestOutput_Sanitize(t *testing.T) {
	a := assert.New(t)

	var o *Output
	a.Error(o.Sanitize())

	// 默认的 Type
	o = &Output{}
	a.NotError(o.Sanitize())
	a.Equal(o.marshal, marshaler(apidocMarshaler))

	o = &Output{Type: "invalid-type"}
	a.Error(o.Sanitize())

	o = &Output{Type: ApidocXML}
	o.Path = "./testdir/apidoc.json"
	a.NotError(o.Sanitize())
	a.Equal(o.Style, stylesheetURL).
		Equal(2, len(o.procInst)).
		Contains(o.procInst[1], stylesheetURL)
}

func TestOptions_buffer(t *testing.T) {
	a := assert.New(t)

	doc := asttest.Get()
	o := &Output{
		Type: OpenapiJSON,
		Path: "./openapi.json",
	}
	a.NotError(o.Sanitize())
	a.NotError(o.buffer(doc))

	doc = asttest.Get()
	o = &Output{}
	a.NotError(o.Sanitize())
	buf, err := o.buffer(doc)
	a.NotError(err).NotNil(buf)
}

func TestFilterDoc(t *testing.T) {
	a := assert.New(t)

	d := asttest.Get()
	o := &Output{}
	a.NotError(o.Sanitize())
	filterDoc(d, o)
	a.Equal(3, len(d.Tags))

	d = asttest.Get()
	o = &Output{
		Tags: []string{"t1"},
	}
	a.NotError(o.Sanitize())
	filterDoc(d, o)
	a.Equal(1, len(d.Tags)).
		Equal(2, len(d.Apis))

	d = asttest.Get()
	o = &Output{
		Tags: []string{"t1", "t2"},
	}
	a.NotError(o.Sanitize())
	filterDoc(d, o)
	a.Equal(2, len(d.Tags)).
		Equal(2, len(d.Apis))

	d = asttest.Get()
	o = &Output{
		Tags: []string{"tag1"},
	}
	a.NotError(o.Sanitize())
	filterDoc(d, o)
	a.Equal(1, len(d.Tags)).
		Equal(1, len(d.Apis))

	d = asttest.Get()
	o = &Output{
		Tags: []string{"not-exists"},
	}
	a.NotError(o.Sanitize())
	filterDoc(d, o)
	a.Equal(0, len(d.Tags)).
		Equal(0, len(d.Apis))
}
