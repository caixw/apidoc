// SPDX-License-Identifier: MIT

package output

import (
	"testing"

	"github.com/issue9/assert"

	"github.com/caixw/apidoc/v6/doc/doctest"
)

func TestRender(t *testing.T) {
	a := assert.New(t)
	doc := doctest.Get()
	o := &Options{
		Path: "./apidoc.xml",
	}
	a.NotError(o.Sanitize())

	a.NotError(Render(doc, o))
}

func TestRender_openapiJSON(t *testing.T) {
	a := assert.New(t)
	doc := doctest.Get()

	o := &Options{
		Type: OpenapiJSON,
		Path: "./openapi.json",
	}
	a.NotError(o.Sanitize())

	a.NotError(Render(doc, o))
}

func TestBuffer(t *testing.T) {
	a := assert.New(t)
	doc := doctest.Get()

	o := &Options{}
	a.NotError(o.Sanitize())
	buf, err := Buffer(doc, o)
	a.NotError(err).NotNil(buf)
}

func TestFilterDoc(t *testing.T) {
	a := assert.New(t)

	d := doctest.Get()
	o := &Options{}
	a.NotError(o.Sanitize())
	filterDoc(d, o)
	a.Equal(3, len(d.Tags))

	d = doctest.Get()
	o = &Options{
		Tags: []string{"t1"},
	}
	a.NotError(o.Sanitize())
	filterDoc(d, o)
	a.Equal(1, len(d.Tags)).
		Equal(2, len(d.Apis))

	d = doctest.Get()
	o = &Options{
		Tags: []string{"t1", "t2"},
	}
	a.NotError(o.Sanitize())
	filterDoc(d, o)
	a.Equal(2, len(d.Tags)).
		Equal(2, len(d.Apis))

	d = doctest.Get()
	o = &Options{
		Tags: []string{"tag1"},
	}
	a.NotError(o.Sanitize())
	filterDoc(d, o)
	a.Equal(1, len(d.Tags)).
		Equal(1, len(d.Apis))

	d = doctest.Get()
	o = &Options{
		Tags: []string{"not-exists"},
	}
	a.NotError(o.Sanitize())
	filterDoc(d, o)
	a.Equal(0, len(d.Tags)).
		Equal(0, len(d.Apis))
}
